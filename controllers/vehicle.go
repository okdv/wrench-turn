package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/okdv/wrench-turn/models"
	"github.com/okdv/wrench-turn/services"
)

type VehicleController struct {
}

func NewVehicleController() *VehicleController {
	return &VehicleController{}
}

// GetVehicle
// Retrieves id param, calls GetVehicle services, returns Vehicle
func (vc *VehicleController) GetVehicle(w http.ResponseWriter, r *http.Request) {
	// get vehicle id from url params, parse into int
	vehicleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// call GetVehicle service, return Vehicle
	vehicle, err := services.GetVehicle(vehicleId)
	if err != nil || vehicle == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Vehicle not found: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(vehicle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert vehicle to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// ListVehicles
// Retrieves any URL query params, calls ListVehicles service, returns Vehicle list
func (vc *VehicleController) ListVehicles(w http.ResponseWriter, r *http.Request) {
	var vehicles []*models.Vehicle
	// get URL query params
	userId := r.URL.Query().Get("user")
	jobId := r.URL.Query().Get("job")
	searchStr := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	// call ListVehicles service
	vehicles, err := services.ListVehicles(&userId, &jobId, &searchStr, &sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to retrieve any vehicles: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(vehicles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to convert vehicles to JSON response")
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// CreateVehicle
// Takes NewVehicle as request body, validates it, calls CreateVehicle service, return Vehicle
func (vc *VehicleController) CreateVehicle(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var newVehicle *models.NewVehicle
	// get vehicle data from request body
	err := json.NewDecoder(r.Body).Decode(&newVehicle)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// set newVehicle.user is nil, set to current user
	if newVehicle.User == nil {
		newVehicle.User = &c.ID
	}
	// if newVehicle user is not requesting user, check if admin
	if newVehicle.User != &c.ID && c.Is_admin == false {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to create vehicles for other users")
		return
	}
	// send to NewVehicle service, return Vehicle
	vehicle, err := services.CreateVehicle(*newVehicle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to create vehicle: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(vehicle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert vehicle to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// EditVehicle
// Takes Vehicle as request body, calls EditVehicle service, return Vehicle
func (vc *VehicleController) EditVehicle(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var vehicle models.Vehicle
	// get user data from request body
	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// if requesting users id doesnt match user id in request body, and they are not an admin, throw error
	if (c.ID != vehicle.User) && (c.Is_admin != true) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to edit vehicles of other users")
		return
	}
	// call EditVehicle service, return updated Vehicle
	updatedVehicle, err := services.EditVehicle(vehicle)
	if err != nil || updatedVehicle == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to edit vehicle: %v", err)
		return
	}
	// convert to JSON response
	jsonData, err := json.Marshal(updatedVehicle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert vehicle to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// DeleteVehicle
// Retrieves username param, validates request, calls DeleteVehicle service
func (vc *VehicleController) DeleteVehicle(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get vehicle id from url params, parse into int
	vehicleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// if admin, call DeleteVehicle service, otherwise call DeleteVehicle to only allow vehicle deletion for requesting users vehicles
	var userId *int64 = nil
	if c.Is_admin != true {
		userId = &c.ID
	}
	err = services.DeleteVehicle(vehicleId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to delete vehicle: %v", err)
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Vehicle ID %v has been deleted", vehicleId)
}
