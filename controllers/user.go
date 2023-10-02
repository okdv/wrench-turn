package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/okdv/wrench-turn/models"
	"github.com/okdv/wrench-turn/services"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

// ListUsers
// Retrieves any URL query params, calls ListUsers service, returns User list
func (uc *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []*models.User
	// get URL query params
	jobId := r.URL.Query().Get("job")
	vehicleId := r.URL.Query().Get("vehicle")
	isAdmin := r.URL.Query().Get("admin")
	searchStr := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	// call ListUsers service
	users, err := services.ListUsers(&jobId, &vehicleId, &isAdmin, &searchStr, &sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to retrieve any users")
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to convert vehicle to JSON response")
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// CreateUser
// Takes NewUser as request body, validates it, calls CreateUser service - adding admin status to user if no admin exists, returns User
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser *models.NewUser
	isAdmin := 0
	// get user data from request body
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// confirm username and password are not null
	// get list of admin users, upgrade newUser to be admin if none exist
	adminStr := "1"
	adminUsers, err := services.ListUsers(nil, nil, &adminStr, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable check if existing admin users: %v", err)
		return
	}
	if len(adminUsers) == 0 {
		isAdmin = 1
	}
	newUser.Is_admin = &isAdmin
	// insert into db and return created user via corresponding service
	user, err := services.CreateUser(*newUser)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to create user: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to convert user to JSON response")
		return
	}
	// respond with json
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// EditUser
// Takes User as request body, validates calls EditUser service, returns User
func (uc *UserController) EditUser(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var user models.User
	// get user data from request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request body")
		return
	}
	// if requesting users id doesnt match id in request body, and they are not an admin, throw error
	if (c.ID != user.ID) && (c.Is_admin != true) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Users can only be edited by admins and themselves")
		return
	}
	// call EditUser service, return updated User
	updatedUser, err := services.EditUser(user)
	if err != nil || updatedUser == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to edit user: %v", err)
		return
	}
	// convert to JSON response
	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to convert user to JSON response")
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// GetUserByUsername
// Retrieves username param, calls GetUserByUsername service, returns User
func (uc *UserController) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	// get username from url params
	username := chi.URLParam(r, "username")
	// pass to service for user retrieval
	user, err := services.GetUserByUsername(username)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "User not found")
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to convert user to JSON response")
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// DeleteUser
// Retrieves username param, validates request, calls DeleteUser service
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get username from url params
	username := chi.URLParam(r, "username")
	// if requesting users username doesnt match username param, and they are not an admin, throw error
	if (c.Username != username) && (c.Is_admin != true) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Users can only be deleted by admins and themselves")
		return
	}
	// call DeleteUser service
	err := services.DeleteUser(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %v has been deleted", username)
}

// UpdatePassword
// Takes Passwords as arg, validates request, calls UpdatePassword service
func (uc *UserController) UpdatePassword(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var passwords *models.Passwords
	// get Passwords data from request body
	err := json.NewDecoder(r.Body).Decode(&passwords)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request body")
		return
	}
	// if request is not admin, run additional validations
	if c.Is_admin == false {
		// err if no current password
		if passwords.CurrentPassword == nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Unless admin, current password must be provided")
			return
		}
		// err if requester user is not requested user
		if passwords.Username != c.Username {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "User passwords can only be updated by admins and the user themselves")
			return
		}
		// retrieve auth info (including bool for if passwords match)
		_, _, _, _, valid, err, _ := services.RetrieveAuthInfo(&models.Credentials{
			Username: passwords.Username,
			Password: *passwords.CurrentPassword,
		})
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Incorrect password")
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unable to validate current password: %v", err)
			return
		}
	}
	// call UpdatePassword service
	err = services.UpdatePassword(passwords.Username, passwords.NewPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to update password: %v", err)
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Password updated")
	return
}
