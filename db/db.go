package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// type Like
// Used by QueryBuilder
type Like struct {
	Fields []string `json:"fields"`
	Match  string   `json:"match"`
	Or     bool     `json:"or"`
}

// ConnectDatabase
// Use sqlite pkg to establish a connection
func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		return err
	}
	DB = db
	return nil
}

// QueryBuilder
// Take basic parts of SQL query, construct into usable query
// May need ORM in the future if it gets too complex, but should work fine for basic queries used thus far
func QueryBuilder(query string, joins *[]string, wheres *[]string, likes *[]Like, sort string) string {
	// start with main query (e.g. SELECT ... FROM ... AS ...)
	var q string = query
	// start where str to be appended to query after joins, wheres, likes
	var whereStr string = ""
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
	// append whereStr to query, add ORDER BY to end to fully construct query
	q = q + whereStr + " ORDER BY " + sort
	// log and return query
	log.Printf("QueryBuilder: %v", q)
	return q
}
