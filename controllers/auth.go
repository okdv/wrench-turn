package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/okdv/wrench-turn/models"
	"github.com/okdv/wrench-turn/services"
)

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

var jwtCookieName = "wrenchturn-jwt"

// Auth
// takes credentials, validates them, gets and returns JWT
func (ac *AuthController) Auth(w http.ResponseWriter, r *http.Request) {
	var creds *models.Credentials
	// get credentials from request body
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// retrieve user auth info
	userId, username, isAdmin, _, isValid, err, statusCode := services.RetrieveAuthInfo(creds)
	if err != nil || !isValid {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, err.Error())
		return
	}
	// generate new jwt
	jwtCookie, err := services.CreateJWT(*userId, *username, *isAdmin, jwtCookieName)
	if err != nil || jwtCookie == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	// set jwt as cookie
	http.SetCookie(w, jwtCookie)
	// return 200ok auth jwt and session
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwtCookie)
}

// Logout()
// Sets JWT and Session cookies to expire immediately
func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    jwtCookieName,
		Expires: time.Now(),
	})
}
