package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
		fmt.Fprintf(w, "Unable to parse request body: %v", err)
		return
	}
	// retrieve user auth info
	userId, username, isAdmin, _, isValid, err, statusCode := services.RetrieveAuthInfo(creds)
	if err != nil || !isValid {
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "Unable to retrieve user auth info: %v", err)
		return
	}
	// generate new jwt
	jwtCookie, err := services.CreateJWT(*userId, *username, *isAdmin, jwtCookieName)
	if err != nil || jwtCookie == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to generate JWT: %v", err)
		return
	}
	// set jwt as cookie
	http.SetCookie(w, jwtCookie)
	// return 200ok auth jwt and session
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwtCookie)
}

// Verify
// Takes another controller as arg, verifies active JWT Bearer as Auth header
func (ac *AuthController) Verify(endpointHandler func(w http.ResponseWriter, r *http.Request, c *models.Claims)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve auth header
		jwt := r.Header.Get("Authorization")
		// if no auth header, return err
		if len(jwt) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "No JWT provided as Bearer token in Authorization header")
			return
		}
		// trim "Bearer " prefix from auth header if present
		jwt = strings.TrimPrefix(jwt, "Bearer ")
		// retrieve JWT claims from VerifyJWT service
		claims, err := services.VerifyJWT(jwt)
		if err != nil || claims == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unable to verify JWT: %v", err)
			return
		}
		// return controller provided as arg with JWT claims attached
		endpointHandler(w, r, claims)
	})
}

// TestVerify
// Returns a simple success if JWT was verified by verify controller
func (ac *AuthController) TestVerify(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "JWT successfully verified")
	return
}

// Logout
// Sets JWT and Session cookies to expire immediately
func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    jwtCookieName,
		Expires: time.Now(),
	})
}
