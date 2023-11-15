package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/okdv/wrench-turn/controllers"
	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/models"
)

var r *chi.Mux
var createdUser *models.User
var req *http.Request
var w *httptest.ResponseRecorder
var testUsername string
var testPassword string
var jwtCookie *http.Cookie

/* Order of events
==============================================================================
1 TestMain starts, defers to after rest of tests to close and exit
2 TestCreateUser creates a new user
3 TestListUsers lists all users, should be at least one (from TestCreateUser)
4 TestGetUser retrieves user from TestCreateUser
5 TestAuth authorized with credentials of user from TestCreateUser
6 TestDeleteUser deletes user from TestCreateUser with JWT from TestAuth
7 TestMain finishes
*/

func TestMain(m *testing.M) {
	// test db connection
	DB, err := db.ConnectDatabase()
	if err != nil {
		panic("Unable to open database: " + err.Error())
	}
	log.Print("Successfully connected to database")
	// declare router
	r = chi.NewRouter()
	// create controllers
	userController := controllers.NewUserController()
	authController := controllers.NewAuthController()
	// create routes
	// auth routes
	r.Post("/auth", authController.Auth)
	// user routes
	r.Get("/users", userController.ListUsers)
	r.Get("/users/{username}", userController.GetUserByUsername)
	r.Post("/users/create", userController.CreateUser)
	r.Delete("/users/{username}", authController.Verify(userController.DeleteUser))
	// after all tests, close db
	defer DB.Close()
	// run tests
	exitCode := m.Run()
	// cleanup
	os.Exit(exitCode)

}

// TestCreateUser
// Tests creating a new user using default credentials and user controller
func TestCreateUser(t *testing.T) {
	// setup user
	testUsername = "wrench-turn_go_test_user"
	testPassword = "Password123"
	newUser := &models.NewUser{
		Username: testUsername,
		Password: &testPassword,
	}
	// convert to json
	jsonData, err := json.Marshal(newUser)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// post user json to api
	req = httptest.NewRequest("POST", "/users/create", bytes.NewReader(jsonData))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	if err := json.NewDecoder(w.Body).Decode(&createdUser); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	log.Print("Successfully created user")
}

// TestListUsers
// Tests getting all users
func TestListUsers(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/users", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var users *[]models.User
	if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if no returned users
	if users == nil || len(*users) == 0 {
		t.Errorf("No users retreived, at least one (test user from TestCreateUser) should exist")
	}
	log.Print("Successfully retrieved users")
}

// TestGetUser
// Tests getting all user created by TestCreateUser
func TestGetUser(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/users/"+createdUser.Username, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var fetchedUser *models.User
	if err := json.NewDecoder(w.Body).Decode(&fetchedUser); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// error if returned user ID is not the same as created user ID
	if fetchedUser.ID != createdUser.ID {
		t.Errorf("Username %v fetched a different user ID, %d, than expected, %d", createdUser.Username, fetchedUser.ID, createdUser.ID)
	}
	log.Print("Successfully retrieved test user")
}

// TestAuth
// Tests auth with user created by TestCreateUser
func TestAuth(t *testing.T) {
	creds := &models.Credentials{
		Username: testUsername,
		Password: testPassword,
	}
	jsonData, err := json.Marshal(creds)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// post auth credentials to api
	req = httptest.NewRequest("POST", "/auth", bytes.NewReader(jsonData))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unable to decode response
	if err := json.NewDecoder(w.Body).Decode(&jwtCookie); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if jwt is still empty
	if jwtCookie == nil {
		t.Errorf("JWT not present")
	}
	log.Print("Successfully logged in")
}

// TestDeleteUser
// Tests deleting the user created by TestCreateUser
func TestDeleteUser(t *testing.T) {
	// delete from api
	req = httptest.NewRequest("DELETE", "/users/"+createdUser.Username, nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected http status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
		// if issue deleting user, log that user may need to be manually deleted from db (can conflict future tests)
		log.Print("Test user wrench-turn_go_test_user may still exist, delete manually if so")
	}
	log.Print("Successfully deleted user")
}
