package services

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/models"
	"github.com/okdv/wrench-turn/utils"
	"golang.org/x/crypto/bcrypt"
)

// CreateJWT
// Takes userID, username, admin rights and cookie name, creates and returns JWT containing auth info
func CreateJWT(id int64, username string, isAdmin bool, cookieName string) (*http.Cookie, error) {
	// create timestamp 120 hours (5 days) from now
	newExpirationTime := time.Now().Add(120 * time.Hour)
	// create claim with auth info
	claims := &models.Claims{
		ID:       id,
		Username: username,
		Is_admin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(newExpirationTime),
		},
	}
	// create new JWT from claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// get signed JWT as string
	tokenStr, err := token.SignedString(os.Getenv("JWT_KEY"))
	if err != nil {
		return nil, err
	}
	// create a JWT cookie
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    tokenStr,
		Path:     "/",
		Domain:   os.Getenv("API_DOMAIN"),
		Expires:  newExpirationTime,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   false,
	}
	// return cookie
	return cookie, nil
}

// RetrieveAuthInfo
// Take username, retrieve userID, username, admin status, encrypted pw from db for auth purposes
func RetrieveAuthInfo(creds *models.Credentials) (*int64, *string, *bool, *[]byte, bool, error, int) {
	// retrieve auth info from db
	userId, username, isAdminInt, hashed, err := db.GetAuthInfoByUsername(creds.Username)
	isAdmin := utils.IntToBool(*isAdminInt)
	// throw error if present
	if err != nil {
		return nil, nil, nil, nil, false, err, 500
		// throw error if no auth info retrieved
	} else if userId == nil || username == nil {
		return nil, nil, nil, nil, false, errors.New("No user found"), 404
		// throw exception auth if there is no password in db
	} else if hashed == nil || len(*hashed) == 0 {
		log.Println("No existing password for user, automatic authentication done based on username")
		return userId, username, &isAdmin, hashed, true, nil, 200
	}
	// if all auth info is present with no errors, check if passwords match
	err = bcrypt.CompareHashAndPassword(*hashed, []byte(creds.Password))
	if err != nil {
		log.Println(err)
		return nil, nil, nil, nil, false, err, 401
	}
	// return success
	return userId, username, &isAdmin, hashed, true, nil, 200
}
