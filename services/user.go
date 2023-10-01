package services

import (
	"errors"

	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/models"
	"github.com/okdv/wrench-turn/utils"
)

// CreateUser
// Takes NewUser as arg, validates and prepares it for db, inserts into user table
func CreateUser(newUser models.NewUser) (*models.User, error) {
	// validate and generate hashed pw
	hashed, err := utils.ValidateAndHashPassword(newUser.Password)
	if err != nil {
		return nil, err
	}
	// insert into db, return userID
	userId, err := db.CreateUser(newUser, hashed)
	if err != nil || userId == nil {
		return nil, err
	}
	// retrieve User from db by userID, return User
	user, err := GetUserById(*userId)
	return user, err
}

// ListUsers
// Takes URL query params as args, passes to ListUsers query, returns User list
func ListUsers(jobId *string, vehicleId *string, searchStr *string, sort *string) ([]*models.User, error) {
	users, err := db.ListUsers(jobId, vehicleId, searchStr, sort)
	return users, err
}

// GetUserByUsername
// Takes username as arg, passes to GetUserByUsername query, returns User
func GetUserByUsername(username string) (*models.User, error) {
	user, err := db.GetUserByUsername(username)
	return user, err
}

// GetUserById
// Takes userID as arg, passes to GetUserById query, returns User
func GetUserById(userID int64) (*models.User, error) {
	user, err := db.GetUserById(userID)
	return user, err
}

// DeleteUser
// Takes username as arg, passes to DeleteUser query
func DeleteUser(username string) error {
	err := db.DeleteUser(username)
	return err
}

// EditUser
// Takes User as arg, passes to EditUser query, returns updated User
func EditUser(editedUser models.User) (*models.User, error) {
	err := db.EditUser(editedUser)
	if err != nil {
		return nil, err
	}
	user, err := GetUserById(editedUser.ID)
	return user, err
}

// UpdatePassword
// Take Passwords as arg, process it, pass to UpdatePassword query
func UpdatePassword(passwords *models.Passwords) error {
	// if current password provided, validate it
	if passwords.CurrentPassword != nil {
		// create credentials
		creds := &models.Credentials{
			Username: passwords.Username,
			Password: *passwords.CurrentPassword,
		}
		// retrieve auth info (including bool for if passwords match)
		_, _, _, _, valid, err, _ := RetrieveAuthInfo(creds)
		if err != nil {
			return err
		}
		if !valid {
			return errors.New("Incorrect password")
		}
	}
	// check if new passwords match
	if passwords.NewPassword != passwords.ConfirmNewPassword {
		return errors.New("New passwords do not match")
	}
	// if no new passwords, query as NULL, otherwise hash and query as hash
	var newPassword *[]byte
	if len(passwords.NewPassword) == 0 {
		newPassword = nil
	} else {
		// hash new password
		hashedPw, err := utils.ValidateAndHashPassword(passwords.NewPassword)
		if err != nil {
			return err
		}
		newPassword = &hashedPw
	}
	// call db query
	err := db.UpdatePassword(passwords.Username, newPassword)
	return err
}
