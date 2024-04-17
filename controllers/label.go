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

type LabelController struct {
}

func NewLabelController() *LabelController {
	return &LabelController{}
}

// GetLabel
// Retrieves id param, calls GetLabel services, returns Label
func (jc *LabelController) GetLabel(w http.ResponseWriter, r *http.Request) {
	// get label id from url params, parse into int
	labelId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// call GetLabel service, return Label
	label, err := services.GetLabel(labelId)
	if err != nil || label == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Label not found: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(label)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert label to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// ListLabels
// Retrieves any URL query params, calls ListLabels service, returns Label list
func (jc *LabelController) ListLabels(w http.ResponseWriter, r *http.Request) {
	var labels []*models.Label
	// get URL query params
	userId := r.URL.Query().Get("user")
	jobId := r.URL.Query().Get("job")
	searchStr := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	// call ListLabels service
	labels, err := services.ListLabels(&userId, &jobId, &searchStr, &sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to retrieve any labels: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(labels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to convert labels to JSON response")
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// CreateLabel
// Takes NewLabel as request body, validates it, calls CreateLabel service, return Label
func (jc *LabelController) CreateLabel(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var newLabel *models.NewLabel
	// get label data from request body
	err := json.NewDecoder(r.Body).Decode(&newLabel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// set newLabel.user is nil, set to current user
	if newLabel.User == nil {
		newLabel.User = &c.ID
	}
	// if newLabel user is not requesting user, check if admin
	if newLabel.User != &c.ID && c.Is_admin == false {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to create labels for other users")
		return
	}
	// send to newLabel service, return Label
	label, err := services.CreateLabel(*newLabel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to create label: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(label)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert label to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// EditLabel
// Takes Label as request body, calls EditLabel service, return Label
func (jc *LabelController) EditLabel(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var label models.Label
	// get user data from request body
	err := json.NewDecoder(r.Body).Decode(&label)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// if requesting users id doesnt match user id in request body, and they are not an admin, throw error
	if (c.ID != label.User) && (c.Is_admin != true) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to edit labels of other users")
		return
	}
	// call EditLabel service, return updated Label
	updatedLabel, err := services.EditLabel(label)
	if err != nil || updatedLabel == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to edit label: %v", err)
		return
	}
	// convert to JSON response
	jsonData, err := json.Marshal(updatedLabel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert label to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// DeleteLabel
// Retrieves username param, validates request, calls DeleteLabel service
func (jc *LabelController) DeleteLabel(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get label id from url params, parse into int
	labelId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// if admin not admin, only allow deleting users own labels
	var userId *int64 = nil
	if c.Is_admin != true {
		userId = &c.ID
	}
	// call delete label
	err = services.DeleteLabel(labelId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to delete label: %v", err)
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Label ID %v has been deleted", labelId)
}
