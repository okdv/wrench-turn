package db

import (
	"log"
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
