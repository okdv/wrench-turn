package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/okdv/wrench-turn/controllers"
	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/models"
)

var r *chi.Mux
var createdUser *models.User
var createdJob *models.Job
var createdVehicle *models.Vehicle
var req *http.Request
var w *httptest.ResponseRecorder
var testUsername string
var testPassword string
var jwtCookie *http.Cookie

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
	jobController := controllers.NewJobController()
	vehicleController := controllers.NewVehicleController()
	// create routes
	// auth routes
	r.Post("/auth", authController.Auth)
	r.Get("/verify", authController.Verify(authController.TestVerify))
	r.Get("/refresh", authController.Verify(authController.Refresh))
	// user routes
	r.Get("/users", userController.ListUsers)
	r.Get("/users/{username}", userController.GetUserByUsername)
	r.Post("/users/create", userController.CreateUser)
	r.Delete("/users/{username}", authController.Verify(userController.DeleteUser))
	r.Post("/users/edit", authController.Verify(userController.EditUser))
	// job routes
	r.Get("/jobs", jobController.ListJobs)
	r.Get("/jobs/{id:[0-9]+}", jobController.GetJob)
	r.Post("/jobs/create", authController.Verify(jobController.CreateJob))
	r.Post("/jobs/edit", authController.Verify(jobController.EditJob))
	r.Delete("/jobs/{id:[0-9]+}", authController.Verify(jobController.DeleteJob))
	// vehicle routes
	r.Get("/vehicles", vehicleController.ListVehicles)
	r.Get("/vehicles/{id:[0-9]+}", vehicleController.GetVehicle)
	r.Post("/vehicles/create", authController.Verify(vehicleController.CreateVehicle))
	r.Post("/vehicles/edit", authController.Verify(vehicleController.EditVehicle))
	r.Delete("/vehicles/{id:[0-9]+}", authController.Verify(vehicleController.DeleteVehicle))
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

// TestVerify
// Tests verify endpoint
func TestVerify(t *testing.T) {
	// test via api
	req = httptest.NewRequest("GET", "/verify", nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	log.Print("Successfully verified with endpoint")
}

// TestVerify
// Tests verify endpoint
func TestRefresh(t *testing.T) {
	// refresh via api (should not work since jwt doesnt expire within 24 hours)
	req = httptest.NewRequest("GET", "/refresh", nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusNoContent {
		t.Errorf("Expted status code %d, got %d", http.StatusNoContent, w.Code)
	}
	log.Print("Successfully called refresh")

	// force refresh via api
	req = httptest.NewRequest("GET", "/refresh?force=true", nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusNoContent, w.Code)
	}
	// error if unable to decode response
	if err := json.NewDecoder(w.Body).Decode(&jwtCookie); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if jwt is still empty
	if jwtCookie == nil {
		t.Errorf("JWT not present")
	}
	log.Print("Successfully forced refresh")
}

// TestGetAndEditUser
// Tests getting and editing user created by TestCreateJob
func TestGetAndEditUser(t *testing.T) {
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
	// change user description
	testDescription := "Test Description"
	fetchedUser.Description = &testDescription
	// convert to json for editing
	jsonData, err := json.Marshal(fetchedUser)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// edit via post req
	req = httptest.NewRequest("POST", "/users/edit", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// decode body from json or error
	if err := json.NewDecoder(w.Body).Decode(&fetchedUser); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// compare dscriptions to confirm successful edit
	if fetchedUser.Description != &testDescription {
		t.Error("Description was not updated")
	}
	log.Print("Successfully edited user")
}

// TestCreateVehicle
// Tests createing a vehicle with user created by TestCreateUser
func TestCreateVehicle(t *testing.T) {
	// setup new test vehicle
	newVehicle := &models.NewVehicle{
		Name: "wrench-turn go test vehicle",
	}
	// convert to json
	jsonData, err := json.Marshal(newVehicle)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// create via api
	req = httptest.NewRequest("POST", "/vehicles/create", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	if err := json.NewDecoder(w.Body).Decode(&createdVehicle); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	log.Print("Successfully created vehicle")
}

// TestListVehicles
// Tests getting all vehicles
func TestListVehicles(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/vehicles", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var vehicles *[]models.Vehicle
	if err := json.NewDecoder(w.Body).Decode(&vehicles); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if no returned users
	if vehicles == nil || len(*vehicles) == 0 {
		t.Errorf("No vehicles retreived, at least one (test vehicles from TestCreateVehicle) should exist")
	}
	log.Print("Successfully retrieved vehicles")
}

// TestGetAndEditVehicle
// Tests getting and editing job created by TestCreateVehicle
func TestGetAndEditVehicle(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/vehicles/"+strconv.FormatInt(createdVehicle.ID, 10), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var fetchedVehicle *models.Vehicle
	if err := json.NewDecoder(w.Body).Decode(&fetchedVehicle); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// error if returned vehicle ID is not the same as created vehicle ID
	if fetchedVehicle.ID != createdVehicle.ID {
		t.Errorf("Vehicle ID %d fetched a different vehicle, ID %d, than expected", createdVehicle.ID, fetchedVehicle.ID)
	}
	log.Print("Successfully retrieved test vehicle")
	// change vehicle description
	testDescription := "Test Description"
	fetchedVehicle.Description = &testDescription
	// convert to json
	jsonData, err := json.Marshal(fetchedVehicle)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// edit via post req
	req = httptest.NewRequest("POST", "/vehicles/edit", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// decode body from json or error
	if err := json.NewDecoder(w.Body).Decode(&fetchedVehicle); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// compare dscriptions to confirm successful edit
	if fetchedVehicle.Description != &testDescription {
		t.Error("Description was not updated")
	}
	log.Print("Successfully edited vehicle")
}

// TestCreateJob
// Tests createing a job with user created by TestCreateUser
func TestCreateJob(t *testing.T) {
	// setuop new test job
	newJob := &models.NewJob{
		Name:    "wrench-turn go test job",
		Vehicle: &createdVehicle.ID,
	}
	// convert to json
	jsonData, err := json.Marshal(newJob)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// create via api
	req = httptest.NewRequest("POST", "/jobs/create", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	if err := json.NewDecoder(w.Body).Decode(&createdJob); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	log.Print("Successfully created job")
}

// TestListJobs
// Tests getting all jobs
func TestListJobs(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/jobs", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var jobs *[]models.Job
	if err := json.NewDecoder(w.Body).Decode(&jobs); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if no returned users
	if jobs == nil || len(*jobs) == 0 {
		t.Errorf("No jobs retreived, at least one (test jobs from TestCreateJob) should exist")
	}
	log.Print("Successfully retrieved jobs")
}

// TestGetAndEditJob
// Tests getting and editing job created by TestCreateJob
func TestGetAndEditJob(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/jobs/"+strconv.FormatInt(createdJob.ID, 10), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var fetchedJob *models.Job
	if err := json.NewDecoder(w.Body).Decode(&fetchedJob); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// error if returned job ID is not the same as created job ID
	if fetchedJob.ID != createdJob.ID {
		t.Errorf("Job ID %d fetched a different job, ID %d, than expected", createdJob.ID, fetchedJob.ID)
	}
	log.Print("Successfully retrieved test job")
	// change job description
	testDescription := "Test Description"
	fetchedJob.Description = &testDescription
	// convert to json
	jsonData, err := json.Marshal(fetchedJob)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// edit via post req
	req = httptest.NewRequest("POST", "/jobs/edit", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// decode body from json or error
	if err := json.NewDecoder(w.Body).Decode(&fetchedJob); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// compare dscriptions to confirm successful edit
	if fetchedJob.Description != &testDescription {
		t.Error("Description was not updated")
	}
	log.Print("Successfully edited job")
}

// TestDeleteJob
// Tests deleting the job created by TestCreateUser
func TestDeleteJob(t *testing.T) {
	// delete via api
	req = httptest.NewRequest("DELETE", "/jobs/"+strconv.FormatInt(createdJob.ID, 10), nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected http status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
		// if issue deleting job, log that id as it may need to be manually deleted from db
		log.Printf("Test job ID %d may still exist, delete manually if so", createdJob.ID)
	}
	log.Print("Successfully deleted job")
}

// TestDeleteVehicle
// Tests deleting the vehicle created by TestCreateUser
func TestDeleteVehicle(t *testing.T) {
	// delete via api
	req = httptest.NewRequest("DELETE", "/vehicles/"+strconv.FormatInt(createdVehicle.ID, 10), nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected http status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
		// if issue deleting vehicle, log that id as it may need to be manually deleted from db
		log.Printf("Test vehicle ID %d may still exist, delete manually if so", createdVehicle.ID)
	}
	log.Print("Successfully deleted vehicle")
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
