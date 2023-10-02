package db

import (
	"errors"
	"log"

	"github.com/okdv/wrench-turn/models"
)

// Auth Queries

// GetAuthInfoByUsername
// Take username, retrieve auth info: user ID, isAdmin and the hashed password
func GetAuthInfoByUsername(username string) (*int64, *string, *int, *[]byte, error) {
	var userId int64
	var isAdmin int
	var hashed []byte
	err := DB.QueryRow("SELECT id, is_admin, hashed_pw FROM user WHERE username = ?", username).Scan(&userId, &isAdmin, &hashed)
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return nil, nil, nil, nil, err
	}
	return &userId, &username, &isAdmin, &hashed, nil
}

// User Queries

// GetUserById
// Take user ID, return entire user
func GetUserById(userId int64) (*models.User, error) {
	var user models.User
	// query db, return any errors
	err := DB.QueryRow("SELECT * FROM user WHERE id=?", userId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Description,
		&user.Hashed_pw,
		&user.Is_admin,
		&user.Created_at,
		&user.Updated_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername
// Take username, return entire user
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	// query db, return any errors
	err := DB.QueryRow("SELECT * FROM user WHERE username=?", username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Description,
		&user.Hashed_pw,
		&user.Is_admin,
		&user.Created_at,
		&user.Updated_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return &user, nil
}

// GetHashedPasswordByUsername
// Take username as arg, return Hashed_pw from user table where username matches
func GetHashedPasswordByUsername(username string) (*[]byte, error) {
	var hashed *[]byte
	err := DB.QueryRow("SELECT hashed_pw FROM user WHERE username=?", username).Scan(&hashed)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return hashed, nil
}

// ListUsers
// Take filters as args, return User list
func ListUsers(jobId *string, vehicleId *string, isAdmin *string, searchStr *string, sort *string) ([]*models.User, error) {
	var joins []string
	var wheres []string
	var likes []Like
	// establish default sort if not provided
	var orderBy = "u.updated_at DESC"
	// establish basic query
	q := "SELECT * FROM user AS u"
	// if isAdmin provided, add where to query
	if isAdmin != nil && len(*isAdmin) > 0 {
		wheres = append(wheres, "u.is_admin="+*isAdmin)
	}
	// if job ID provided join by userID where jobID is present
	if jobId != nil && len(*jobId) > 0 {
		joins = append(joins, "JOIN user_job AS uj ON u.id = uj.user")
		wheres = append(wheres, "uj.job="+*jobId)
	}
	// if vehicle ID provided join by userID where vehicleID is present
	if vehicleId != nil && len(*vehicleId) > 0 {
		joins = append(joins, "JOIN user_vehicle AS uv ON u.id = uv.user")
		wheres = append(wheres, "uv.vehicle="+*vehicleId)
	}
	// if search string provided, construct likes to query username, description cols
	if searchStr != nil && len(*searchStr) > 0 {
		var fields []string
		fields = append(fields, "u.username")
		fields = append(fields, "u.description")
		likes = append(likes, Like{
			Fields: fields,
			Match:  *searchStr,
			Or:     true,
		})
	}
	// if sort provided, append appropriate sort based on query param
	if sort != nil && len(*sort) > 0 {
		switch *sort {
		case "az":
			orderBy = "u.username ASC"
		case "za":
			orderBy = "u.username DESC"
		case "oldest":
			orderBy = "u.created_at ASC"
		case "newest":
			orderBy = "u.created_at DESC"
		case "last_updated":
			orderBy = "u.updated_at DESC"
		default:
			orderBy = "u.updated_at DESC"
		}
	}
	// generate query with QueryBuilder
	query := QueryBuilder(q, &joins, &wheres, &likes, orderBy)
	// retrieve all matching rows
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return nil, err
	}
	defer rows.Close()
	// create list of User
	users := make([]*models.User, 0)
	// loop through returned rows
	for rows.Next() {
		// attribute to User
		user := models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Description,
			&user.Hashed_pw,
			&user.Is_admin,
			&user.Created_at,
			&user.Updated_at,
		)
		if err != nil {
			log.Printf("Error scanning rows retrieved from DB: %s", err)
			return nil, err
		}
		// append User to list of User
		users = append(users, &user)
	}
	return users, nil
}

// CreateUser
// Take NewUser and hashed PW as arguments, insert them into db
func CreateUser(newUser models.NewUser, hashedPw []byte) (*int64, error) {
	// insert into db, return any errors
	res, err := DB.Exec("INSERT INTO user(Username, Email, Hashed_pw, Is_admin) VALUES (?,?,?,?)",
		newUser.Username,
		newUser.Email,
		hashedPw,
		newUser.Is_admin,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	// get inserted users id
	userId, err := res.LastInsertId()
	return &userId, err
}

// EditUser
// Take User as arg, update it in db
func EditUser(editedUser models.User) error {
	_, err := DB.Exec("UPDATE user SET username=?, email=?, description=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		editedUser.Username,
		editedUser.Email,
		editedUser.Description,
		editedUser.ID,
	)
	// throw SQL errors
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	return err
}

// DeleteUser
// Take username as arg, delete User from user table where username present
func DeleteUser(username string) error {
	res, err := DB.Exec("DELETE FROM user WHERE username = ?", username)
	// throw SQL errors
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	// retrieve rows affected count
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	// throw error if no rows affected
	if rows == 0 {
		log.Printf("No rows deleted")
		return errors.New("No rows deleted")
	}

	return nil
}

// UpdatePassword
// Take username and hashed pw as args, update it in db
func UpdatePassword(username string, password *[]byte) error {
	res, err := DB.Exec("UPDATE user SET hashed_pw=? WHERE username=?", password, username)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return err
	}
	// retrieve rows affected count, error if 0
	rowCount, err := res.RowsAffected()
	if rowCount == 0 || err != nil {
		log.Printf("No rows updated: %v", err)
		return errors.New("No rows updated")
	}
	return nil
}
