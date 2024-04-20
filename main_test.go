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
var createdTask *models.Task
var createdVehicle *models.Vehicle
var createdAlert *models.Alert
var createdLabel *models.Label
var req *http.Request
var w *httptest.ResponseRecorder
var testUsername string
var testPassword string
var jwtCookie *http.Cookie

func TestMain(m *testing.M) {
	// check if db file exists
	dbFilename := "test.db"
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) {
		// if not, create it
		log.Printf("database file %v does not exist, running db setup... ", dbFilename)
		file, err := os.Create(dbFilename)
		if err != nil {
			log.Fatal("failed to create database file")
		}
		file.Close()
		// read sql from schema
		sqlBytes, err := os.ReadFile("schema.sql")
		if err != nil {
			log.Fatalf("failed to read schema file: %v", err)
		}
		err = db.CreateDatabase(dbFilename, string(sqlBytes))
		if err != nil {
			log.Fatalf("failed to create database: %v", err)
		}
		log.Print("Will attempt to connect to db now...")
	}
	// test db connection
	DB, err := db.ConnectDatabase(dbFilename)
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
	taskController := controllers.NewTaskController()
	vehicleController := controllers.NewVehicleController()
	alertController := controllers.NewAlertController()
	labelController := controllers.NewLabelController()

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
	r.Post("/jobs/{jobId:[0-9]+}/assignLabel/{labelId:[0-9]+}", authController.Verify(jobController.AssignJobLabel))
	r.Post("/jobs/create", authController.Verify(jobController.CreateJob))
	r.Post("/jobs/edit", authController.Verify(jobController.EditJob))
	r.Delete("/jobs/{id:[0-9]+}", authController.Verify(jobController.DeleteJob))
	// task routes
	r.Get("/jobs/{jobId:[0-9]+}/tasks", taskController.ListTasks)
	r.Get("/jobs/{jobId:[0-9]+}/tasks/{taskId:[0-9]+}", taskController.GetTask)
	r.Patch("/jobs/{jobId:[0-9]+}/tasks/{taskId:[0-9]+}/complete", authController.Verify(taskController.MarkComplete))
	r.Post("/jobs/{jobId:[0-9]+}/tasks/create", authController.Verify(taskController.CreateTask))
	r.Post("/jobs/{jobId:[0-9]+}/tasks/edit", authController.Verify(taskController.EditTask))
	r.Delete("/jobs/{jobId:[0-9]+}/tasks/{taskId:[0-9]+}", authController.Verify(taskController.DeleteTask))
	r.Delete("/jobs/{jobId:[0-9]+}/tasks", authController.Verify(taskController.DeleteTask))
	// vehicle routes
	r.Get("/vehicles", vehicleController.ListVehicles)
	r.Get("/vehicles/{id:[0-9]+}", vehicleController.GetVehicle)
	r.Post("/vehicles/create", authController.Verify(vehicleController.CreateVehicle))
	r.Post("/vehicles/edit", authController.Verify(vehicleController.EditVehicle))
	r.Delete("/vehicles/{id:[0-9]+}", authController.Verify(vehicleController.DeleteVehicle))
	// alert routes
	r.Get("/alerts", authController.Verify(alertController.ListAlerts))
	r.Get("/alerts/{id:[0-9]+}", authController.Verify(alertController.GetAlert))
	r.Patch("/alerts/{id:[0-9]+}/read", authController.Verify(alertController.MarkRead))
	r.Post("/alerts/create", authController.Verify(alertController.CreateAlert))
	r.Post("/alerts/edit", authController.Verify(alertController.EditAlert))
	r.Delete("/alerts/{id:[0-9]+}", authController.Verify(alertController.DeleteAlert))
	// label routes
	r.Get("/labels", labelController.ListLabels)
	r.Get("/labels/{id:[0-9]+}", labelController.GetLabel)
	r.Post("/labels/create", authController.Verify(labelController.CreateLabel))
	r.Post("/labels/edit", authController.Verify(labelController.EditLabel))
	r.Delete("/labels/{id:[0-9]+}", authController.Verify(labelController.DeleteLabel))
	// run tests
	exitCode := m.Run()
	// Close the database connection explicitly
	if DB != nil {
		if err := DB.Close(); err != nil {
			panic(err)
		}
	}
	// Cleanup after all tests
	if err := os.Remove(dbFilename); err != nil {
		panic(err)
	}
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
// Tests getting and editing user created by TestCreateUser
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
	// create a second job, this one to be auto deleted on vehicle deletion
	req = httptest.NewRequest("POST", "/jobs/create", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	log.Print("Successfully created job")
}

// TestCreateLabel
// Tests createing a label with user created by TestCreateUser
func TestCreateLabel(t *testing.T) {
	defaultColor := "ff0000"
	// setuop new test label
	newLabel := &models.NewLabel{
		Name:  "wrench-turn go test label",
		Color: &defaultColor,
		User:  &createdUser.ID,
	}
	// convert to json
	jsonData, err := json.Marshal(newLabel)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// create via api
	req = httptest.NewRequest("POST", "/labels/create", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	if err := json.NewDecoder(w.Body).Decode(&createdLabel); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	log.Print("Successfully created label")
}

// TestAssignAndUnassignJobLabel
// Tests assigning and unassigning a label created by TestCreateLabel to job created by TestCreateJob
func TestAssignAndUnassignJobLabel(t *testing.T) {
	// assign via api
	req = httptest.NewRequest("POST", "/jobs/"+strconv.FormatInt(createdJob.ID, 10)+"/assignLabel/"+strconv.FormatInt(createdLabel.ID, 10), nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	log.Print("Successfully assigned label")
	// will be unassigned on job delete
}

// TestGetJob
// Tests getting job
func TestGetJob(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/jobs/"+strconv.FormatInt(createdJob.ID, 10), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var job *models.Job
	if err := json.NewDecoder(w.Body).Decode(&job); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if no returned users
	if job == nil {
		t.Errorf("No job retreived, at least one (test jobs from TestCreateJob) should exist")
	}
	log.Print("Successfully retrieved jobs")
}

// TestGetLabel
// Tests getting label
func TestGetLabel(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/labels/"+strconv.FormatInt(createdLabel.ID, 10), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var label *models.Label
	if err := json.NewDecoder(w.Body).Decode(&label); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if no returned users
	if label == nil {
		t.Errorf("No label retreived, at least one (test labels from TestCreateLabel) should exist")
	}
	log.Print("Successfully retrieved labels")
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

// TestGetAndEditLabel
// Tests getting and editing label created by TestCreateLabel
func TestGetAndEditLabel(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/labels/"+strconv.FormatInt(createdLabel.ID, 10), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var fetchedLabel *models.Label
	if err := json.NewDecoder(w.Body).Decode(&fetchedLabel); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// error if returned label ID is not the same as created label ID
	if fetchedLabel.ID != createdLabel.ID {
		t.Errorf("Label ID %d fetched a different label, ID %d, than expected", createdLabel.ID, fetchedLabel.ID)
	}
	log.Print("Successfully retrieved test label")
	// change label description
	testColor := "0400ff"
	fetchedLabel.Color = &testColor
	// convert to json
	jsonData, err := json.Marshal(fetchedLabel)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// edit via post req
	req = httptest.NewRequest("POST", "/labels/edit", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// decode body from json or error
	if err := json.NewDecoder(w.Body).Decode(&fetchedLabel); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// compare dscriptions to confirm successful edit
	if fetchedLabel.Color != &testColor {
		t.Error("Color was not updated")
	}
	log.Print("Successfully edited label")
}

// TestCreateTask
// Test creating task on job created by TestCreateJob
func TestCreateTask(t *testing.T) {
	// setuop new test job
	newTask := &models.NewTask{
		Name: "wrench-turn go test task",
	}
	// convert to json
	jsonData, err := json.Marshal(newTask)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	log.Println(createdJob.User)
	log.Println(createdUser.ID)
	// create via api
	req = httptest.NewRequest("POST", "/jobs/"+strconv.FormatInt(createdJob.ID, 10)+"/tasks/create", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusCreated, w.Code)
	}
	// error if unable to decode response
	if err := json.NewDecoder(w.Body).Decode(&createdTask); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// create second task to test deletion workflow when job is deleted
	// create via api
	req = httptest.NewRequest("POST", "/jobs/"+strconv.FormatInt(createdJob.ID, 10)+"/tasks/create", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusCreated, w.Code)
	}
	log.Print("Successfully created task")
}

// TestGetTask
// Tests getting all tasks for created job
func TestGetTask(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/jobs/"+strconv.FormatInt(createdJob.ID, 10)+"/tasks/"+strconv.FormatInt(createdTask.ID, 10), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var task *models.Job
	if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if no returned users
	if task == nil {
		t.Errorf("No task retreived, task (test task from TestCreateTask) should exist")
	}
	log.Print("Successfully retrieved tasks")
}

// TestGetAndMarkCompleteAndEditTask
// Tests getting and editing task created by TestCreateTask
func TestGetAndMarkCompleteAndEditTask(t *testing.T) {
	// mark complete in api
	req = httptest.NewRequest("PATCH", "/jobs/"+strconv.FormatInt(createdJob.ID, 10)+"/tasks/"+strconv.FormatInt(createdTask.ID, 10)+"/complete", nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// get from api
	req = httptest.NewRequest("GET", "/jobs/"+strconv.FormatInt(createdJob.ID, 10)+"/tasks/"+strconv.FormatInt(createdTask.ID, 10), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var fetchedTask *models.Task
	if err := json.NewDecoder(w.Body).Decode(&fetchedTask); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// error if returned taskID is not the same as created task ID
	if fetchedTask.ID != createdTask.ID {
		t.Errorf("Task ID %d fetched a different task, ID %d, than expected", createdTask.ID, fetchedTask.ID)
	}

	// error if task is not completed
	if fetchedTask.Is_complete != 1 {
		t.Error("Task should be completed, is still incomplete", http.StatusOK)
	}
	log.Print("Successfully retrieved test task")
	// change job description
	testDescription := "Test Description"
	fetchedTask.Description = &testDescription
	// convert to json
	jsonData, err := json.Marshal(fetchedTask)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// edit via post req
	req = httptest.NewRequest("POST", "/jobs/"+strconv.FormatInt(createdJob.ID, 10)+"/tasks/edit", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// decode body from json or error
	if err := json.NewDecoder(w.Body).Decode(&fetchedTask); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// compare dscriptions to confirm successful edit
	if fetchedTask.Description != &testDescription {
		t.Error("Description was not updated")
	}
	log.Print("Successfully edited task")
}

// TestCreateAlert
// Tests creating a alert with user created by TestCreateUser
func TestCreateAlert(t *testing.T) {
	alertName := "wrench-turn go test alert"
	// setuop new test alert
	newAlert := &models.NewAlert{
		Name:    &alertName,
		Type:    "notification",
		Vehicle: &createdVehicle.ID,
		Job:     &createdJob.ID,
		Task:    &createdTask.ID,
	}
	// convert to json
	jsonData, err := json.Marshal(newAlert)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// create via api
	req = httptest.NewRequest("POST", "/alerts/create", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	if err := json.NewDecoder(w.Body).Decode(&createdAlert); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// create another alert to be deleted at job deletion
	req = httptest.NewRequest("POST", "/alerts/create", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// create another alert to be deleted at vehicle deletion
	// set job and task nil to ensure this alert is not deleted by other tests
	newAlert.Job = nil
	newAlert.Task = nil
	req = httptest.NewRequest("POST", "/alerts/create", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusCreated {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	log.Print("Successfully created alerts")
}

// TestGetAlert
// Tests getting alert
func TestGetAlert(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/alerts/"+strconv.FormatInt(createdAlert.ID, 10), nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var alert *models.Alert
	if err := json.NewDecoder(w.Body).Decode(&alert); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if no returned alerts
	if alert == nil {
		t.Errorf("No alert retreived, at least one (test alerts from TestCreateAlert) should exist")
	}
	log.Print("Successfully retrieved alerts")
}

// TestGetAndMarkReadAndEditAlert
// Tests getting and editing alert created by TestCreateAlert
func TestGetAndMarkReadAndEditAlert(t *testing.T) {
	// mark read in api
	req = httptest.NewRequest("PATCH", "/alerts/"+strconv.FormatInt(createdAlert.ID, 10)+"/read", nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// get from api
	req = httptest.NewRequest("GET", "/alerts/"+strconv.FormatInt(createdAlert.ID, 10), nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var fetchedAlert *models.Alert
	if err := json.NewDecoder(w.Body).Decode(&fetchedAlert); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// error if returned taskID is not the same as created alert ID
	if fetchedAlert.ID != createdAlert.ID {
		t.Errorf("Alert ID %d fetched a different alert, ID %d, than expected", createdAlert.ID, fetchedAlert.ID)
	}

	// error if alert is not read
	if *fetchedAlert.Is_read != 1 {
		t.Error("Alert should be read, is still unread", http.StatusOK)
	}
	log.Print("Successfully retrieved test alert")
	// change job description
	testDescription := "Test Description"
	fetchedAlert.Description = &testDescription
	// convert to json
	jsonData, err := json.Marshal(fetchedAlert)
	if err != nil {
		t.Errorf("Error encoding request body: %v", err)
	}
	// edit via post req
	req = httptest.NewRequest("POST", "/alerts/edit", bytes.NewReader(jsonData))
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// decode body from json or error
	if err := json.NewDecoder(w.Body).Decode(&fetchedAlert); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// compare dscriptions to confirm successful edit
	if fetchedAlert.Description != &testDescription {
		t.Error("Description was not updated")
	}
	log.Print("Successfully edited alert")
}

// TestDeleteAlert
// Tests deleting the alert created by TestCreateTask
func TestDeleteAlert(t *testing.T) {
	// delete via api
	req = httptest.NewRequest("DELETE", "/alerts/"+strconv.FormatInt(createdAlert.ID, 10), nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected http status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
		// if issue deleting alert, log that id as it may need to be manually deleted from db
		log.Printf("Test alert ID %d may still exist, delete manually if so", createdAlert.ID)
	}
	log.Print("Successfully deleted alert")
}

// TestDeleteTask
// Tests deleting the task created by TestCreateTask
func TestDeleteTask(t *testing.T) {
	// delete via api
	req = httptest.NewRequest("DELETE", "/jobs/"+strconv.FormatInt(createdJob.ID, 10)+"/tasks/"+strconv.FormatInt(createdTask.ID, 10), nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected http status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
		// if issue deleting task, log that id as it may need to be manually deleted from db
		log.Printf("Test task ID %d may still exist, delete manually if so", createdJob.ID)
	}
	log.Print("Successfully deleted task")
}

// TestDeleteJob
// Tests deleting the job created by TestCreateJob
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

// TestListLabels
// Tests getting all labels
func TestListLabels(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/labels", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var labels *[]models.Label
	if err := json.NewDecoder(w.Body).Decode(&labels); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if no returned users
	if labels == nil || len(*labels) == 0 {
		t.Errorf("No labels retreived, at least one (test labels from TestCreateLabel) should exist")
	}
	log.Print("Successfully retrieved labels")
	// get from api ofr job only, should be nil
	req = httptest.NewRequest("GET", "/labels?job="+strconv.FormatInt(createdLabel.ID, 10), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	if err := json.NewDecoder(w.Body).Decode(&labels); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if no returned users
	if labels != nil && len(*labels) > 0 {
		t.Errorf("No labels should have been retreived, it was disassociated from job on TestDeleteJob")
	}
	log.Print("Successfully confirmed label was disassociated from job")
}

// TestDeleteLabel
// Tests deleting the alert created by TestCreateTask
func TestDeleteLabel(t *testing.T) {
	// delete via api
	req = httptest.NewRequest("DELETE", "/labels/"+strconv.FormatInt(createdLabel.ID, 10), nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected http status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
		// if issue deleting label, log that id as it may need to be manually deleted from db
		log.Printf("Test label ID %d may still exist, delete manually if so", createdLabel.ID)
	}
	log.Print("Successfully deleted label")
}

// TestListTasks
// Tests getting all tasks for created job
func TestListTasks(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/jobs/"+strconv.FormatInt(createdJob.ID, 10)+"/tasks", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var tasks *[]models.Task
	if err := json.NewDecoder(w.Body).Decode(&tasks); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if any tasks are returned, they should al be deleted by TestDeleteJob
	if tasks != nil && len(*tasks) > 0 {
		t.Errorf("No tasks should be retreived, should have been deleted by TestDeleteJob")
	}
	log.Print("Successfully confirmed tasks were deleted on job deletion")
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

// TestListJobs
// Tests getting all jobs
func TestListJobs(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/jobs?vehicle="+strconv.FormatInt(createdVehicle.ID, 10), nil)
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
	// error if jobs are returned
	if jobs != nil && len(*jobs) > 0 {
		t.Errorf("No jobs should be retreived")
	}
	log.Print("Successfully confirmed jobs were deleted on TestDeleteVehicle")
}

// TestListAlerts
// Tests getting all alerts
func TestListAlerts(t *testing.T) {
	// get from api
	req = httptest.NewRequest("GET", "/alerts", nil)
	req.Header.Add("Authorization", "Bearer "+jwtCookie.Value)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// error if unexpected HTTP status
	if w.Code != http.StatusOK {
		t.Errorf("Expted status code %d, got %d", http.StatusOK, w.Code)
	}
	// error if unable to decode response
	var alerts *[]models.Alert
	if err := json.NewDecoder(w.Body).Decode(&alerts); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
	// error if any alerts are returned, they should al be deleted by TestDeleteJob
	if alerts != nil && len(*alerts) > 0 {
		t.Errorf("No alerts should be retreived, should have been deleted by TestDeleteJob/Vehicle")
	}
	log.Print("Successfully confirmed alerts were deleted on job deletion")
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
