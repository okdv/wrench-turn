package db

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/okdv/wrench-turn/models"
)

var DB *sql.DB

// type Like
// Used by QueryBuilder
type Like struct {
	Fields []string `json:"fields"`
	Match  string   `json:"match"`
	Or     bool     `json:"or"`
}

// CreateDatabase
func CreateDatabase(filename string, sqlString string) error {
	// open new db connection
	db, err := sql.Open("sqlite3", "./"+filename)
	if err != nil {
		return err
	}
	defer db.Close()

	// execute sql
	_, err = db.Exec(sqlString)
	if err != nil {
		return err
	}
	log.Print("Databse created successfully")
	return nil
}

// ConnectDatabase
// Use sqlite pkg to establish a connection
func ConnectDatabase(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./"+filename)
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}

// QueryBuilder
// Take basic parts of SQL query, construct into usable query
// May need ORM in the future if it gets too complex, but should work fine for basic queries used thus far
func QueryBuilder(query string, joins *[]string, wheres *[]string, likes *[]Like, groupBy *string, sort *string) string {
	// start with main query (e.g. SELECT ... FROM ... AS ...)
	var q string = query
	// start where str to be appended to query after joins, wheres, likes
	var whereStr string = ""
	var sortStr string = ""
	var groupStr string = ""
	// loop through joins (e.g. JOIN ... AS ... ON), append to query
	if joins != nil {
		for _, join := range *joins {
			q = q + " " + join
		}
	}
	// loop through wheres (e.g. table.col=x), append to query
	if wheres != nil {
		for i, where := range *wheres {
			if i == 0 {
				whereStr = whereStr + " WHERE"
			}
			// if not first where, append with AND
			if i > 0 {
				whereStr = whereStr + " AND"
			}
			whereStr = whereStr + " " + where
		}
	}
	// loop through likes (e.g. type Like), process it into string, append to query
	if likes != nil {
		for i, like := range *likes {
			// if first like, establish whereStr as WHERE (e.g. WHERE )
			if i == 0 && len(whereStr) == 0 {
				whereStr = whereStr + " WHERE"
				// if not first like, append to whereStr with AND (e.g. WHERE (table.col LIKE %a%' OR table.col2 LIKE %a%') AND)
			} else if len(whereStr) > 0 {
				whereStr = whereStr + " AND"
			}
			// add to whereStr with new section in () (e.g. WHERE (table.col LIKE %a%' OR table.col2 LIKE %a%') AND (table2.col LIKE %a%' AND table2.col LIKE '%a%'))
			whereStr = whereStr + " ("
			for i, field := range like.Fields {
				// if not first like field in section
				if i > 0 {
					// if OR, append OR (e.g. (table.col LIKE %a%' OR ))
					if like.Or {
						whereStr = whereStr + " OR"
						// otherwise, append AND (e.g. (table2.col LIKE %a%' AND ))
					} else {
						whereStr = whereStr + " AND"
					}
				}
				// add field (e.g. table.col), then append LIKE followed by match str surrounded be wildcards (e.g. '%matchStr%')
				whereStr = whereStr + " " + field + " LIKE " + "\"%" + like.Match + "%\""
			}
			// add closing paranthesis after all fields
			whereStr = whereStr + ")"
		}
	}
	// if groupby is set, generate that part of query
	if groupBy != nil {
		groupStr = " GROUP BY " + *groupBy
	}
	// if sort is set, generate that part of query
	if sort != nil {
		sortStr = " ORDER BY " + *sort
	}
	// append statements strings to query
	q = q + whereStr + groupStr + sortStr
	// log and return query
	log.Printf("QueryBuilder: %v", q)
	return q
}

// LabelProcessor
// takes label cols and convert them into slice of objects
// ideally this is done as service layer but doing in db layer prevents relooping through results
func LabelProcessor(ids *string, names *string, colors *string, createdTimes *string, updatedTimes *string) ([]models.Label, error) {
	var labels []models.Label
	// if any are null, return empty labels slice
	if ids == nil || names == nil || colors == nil || createdTimes == nil || updatedTimes == nil {
		return labels, nil
	}
	// parse comma delimited values into slices by comma
	parsedIds := strings.Split(*ids, ",")
	parsedNames := strings.Split(*names, ",")
	parsedColors := strings.Split(*colors, ",")
	parsedCreatedTimes := strings.Split(*createdTimes, ",")
	parsedUpdatedTimes := strings.Split(*updatedTimes, ",")
	// throw error if slices are not the same length
	parsedIdsLen := len(parsedIds)
	if parsedIdsLen != len(parsedNames) || parsedIdsLen != len(parsedColors) || parsedIdsLen != len(parsedCreatedTimes) || parsedIdsLen != len(parsedUpdatedTimes) {
		return labels, errors.New("Label columns from db do not have same length, returning empty labels, please try again")
	}
	// for loop to minlength and generate a label object for each
	for i := 0; i < len(parsedIds); i++ {
		// convert id string into int64
		id, err := strconv.ParseInt(parsedIds[i], 10, 64)
		// if error skip iteration
		if err != nil {
			log.Printf("Error converting label id string to int64, skipping: %v", err)
			continue
		}
		// parse create and update times
		layout := "2006-01-02 15:04:05"
		createdTime, err := time.Parse(layout, parsedCreatedTimes[i])
		if err != nil {
			log.Printf("Error converting created time string to time value, skipping: %v", err)
			continue
		}
		updatedTime, err := time.Parse(layout, parsedUpdatedTimes[i])
		if err != nil {
			log.Printf("Error converting created time string to time value, skipping: %v", err)
			continue
		}
		// create Label
		label := &models.Label{
			ID:         id,
			Name:       parsedNames[i],
			Color:      &parsedColors[i],
			Created_at: createdTime,
			Updated_at: updatedTime,
		}
		// append to labels
		labels = append(labels, *label)
	}
	return labels, nil
}
