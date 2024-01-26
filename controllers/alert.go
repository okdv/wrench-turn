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

type AlertController struct {
}

func NewAlertController() *AlertController {
	return &AlertController{}
}

// GetAlert
// Retrieves id param, calls GetAlert services, returns Alert
func (ac *AlertController) GetAlert(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get alert id from url params, parse into int
	alertId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// call GetAlert service, return Alert
	alert, err := services.GetAlert(alertId)
	if err != nil || alert == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Alert not found: %v", err)
		return
	}
	// if retrieved alert user is not requesting user, check if admin
	if alert.User != c.ID && c.Is_admin == false {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to get alerts for other users")
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(alert)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert alert to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// ListAlerts
// Retrieves any URL query params, calls ListAlerts service, returns Alert list
func (ac *AlertController) ListAlerts(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var alerts []*models.Alert
	// get URL query params
	userId := r.URL.Query().Get("user")
	vehicleId := r.URL.Query().Get("vehicle")
	jobId := r.URL.Query().Get("job")
	taskId := r.URL.Query().Get("task")
	typeStr := r.URL.Query().Get("type")
	isRead := r.URL.Query().Get("read")
	searchStr := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	// set newAlert.user is nil, set to current user
	if len(userId) == 0 {
		userId = strconv.FormatInt(c.ID, 10)
	}
	// if newAlert user is not requesting user, check if admin
	if userId != strconv.FormatInt(c.ID, 10) && c.Is_admin == false {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to list alerts for other users")
		return
	}
	// call ListAlerts service
	alerts, err := services.ListAlerts(&userId, &vehicleId, &jobId, &taskId, &typeStr, &isRead, &searchStr, &sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to retrieve any alerts: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(alerts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to convert alerts to JSON response")
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// CreateAlert
// Takes NewAlert as request body, validates it, calls CreateAlert service, return Alert
func (ac *AlertController) CreateAlert(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var newAlert *models.NewAlert
	// get alert data from request body
	err := json.NewDecoder(r.Body).Decode(&newAlert)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// set newAlert.user is nil, set to current user
	if newAlert.User == nil {
		newAlert.User = &c.ID
	}
	// if newAlert user is not requesting user, check if admin
	if newAlert.User != &c.ID && c.Is_admin == false {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to create alerts for other users")
		return
	}
	// if type is invalid throw error
	if newAlert.Type != "notification" && newAlert.Type != "reminder" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Type must be notification or reminder")
		return
	}
	// send to newAlert service, return Alert
	alert, err := services.CreateAlert(*newAlert)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to create alert: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(alert)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert alert to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// EditAlert
// Takes Alert as request body, calls EditAlert service, return Alert
func (ac *AlertController) EditAlert(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var alert models.Alert
	// get user data from request body
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// if requesting users id doesnt match user id in request body, and they are not an admin, throw error
	if (c.ID != alert.User) && (c.Is_admin != true) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to edit alerts of other users")
		return
	}
	// if type is invalid throw error
	if alert.Type != "notification" && alert.Type != "reminder" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Type must be notification or reminder")
		return
	}
	// call EditAlert service, return updated Alert
	updatedAlert, err := services.EditAlert(alert)
	if err != nil || updatedAlert == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to edit alert: %v", err)
		return
	}
	// convert to JSON response
	jsonData, err := json.Marshal(updatedAlert)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert alert to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// DeleteAlert
// Retrieves username param, validates request, calls DeleteAlert service
func (ac *AlertController) DeleteAlert(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get alert id from url params, parse into int
	alertId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// if admin, call DeleteAlert service, otherwise call DeleteUsersAlert to only allow alert deletion for requesting users alerts
	var userId *int64 = nil
	if c.Is_admin != true {
		userId = &c.ID
	}
	err = services.DeleteAlert(alertId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to delete alert: %v", err)
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Alert ID %v has been deleted", alertId)
}

// MarkRead
// Marks task read, or unread
func (ac *AlertController) MarkRead(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get URL query params, convert to int
	unread := r.URL.Query().Get("unread")
	status := 1
	if unread == "true" {
		status = 0
	}
	// get task id from url params, parse into int
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// get Alert Data
	alert, err := services.GetAlert(id)
	if alert == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Alert ID %d not found: %v", id, err)
		return
	}
	// if alert user is not requesting user, check if admin
	if alert.User != c.ID && c.Is_admin == false {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to create alerts for other users")
		return
	}
	err = services.MarkRead(id, alert.User, status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to mark task read: %v", err)
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task ID %v has been marked as read", id)
}
