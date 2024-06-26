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

type JobController struct {
}

func NewJobController() *JobController {
	return &JobController{}
}

// GetJob
// Retrieves id param, calls GetJob services, returns Job
func (jc *JobController) GetJob(w http.ResponseWriter, r *http.Request) {
	// get job id from url params, parse into int
	jobId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// call GetJob service, return Job
	job, err := services.GetJob(jobId)
	if err != nil || job == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Job not found: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert job to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// ListJobs
// Retrieves any URL query params, calls ListJobs service, returns Job list
func (jc *JobController) ListJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []*models.Job
	// get URL query params
	userId := r.URL.Query().Get("user")
	vehicleId := r.URL.Query().Get("vehicle")
	isTemplate := r.URL.Query().Get("template")
	isComplete := r.URL.Query().Get("complete")
	labelId := r.URL.Query().Get("label")
	searchStr := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	// call ListJobs service
	jobs, err := services.ListJobs(&userId, &vehicleId, &isTemplate, &isComplete, &labelId, &searchStr, &sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to retrieve any jobs: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to convert jobs to JSON response")
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// CreateJob
// Takes NewJob as request body, validates it, calls CreateJob service, return Job
func (jc *JobController) CreateJob(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var newJob *models.NewJob
	// get job data from request body
	err := json.NewDecoder(r.Body).Decode(&newJob)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// set newJob.user is nil, set to current user
	if newJob.User == nil {
		newJob.User = &c.ID
	}
	// if newJob user is not requesting user, check if admin
	if newJob.User != &c.ID && !c.Is_admin {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to create jobs for other users")
		return
	}
	// send to newJob service, return Job
	job, err := services.CreateJob(*newJob)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to create job: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert job to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// EditJob
// Takes Job as request body, calls EditJob service, return Job
func (jc *JobController) EditJob(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var job models.Job
	// get user data from request body
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// if requesting users id doesnt match user id in request body, and they are not an admin, throw error
	if (c.ID != job.User) && !c.Is_admin {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to edit jobs of other users")
		return
	}
	// call EditJob service, return updated Job
	updatedJob, err := services.EditJob(job)
	if err != nil || updatedJob == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to edit job: %v", err)
		return
	}
	// convert to JSON response
	jsonData, err := json.Marshal(updatedJob)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert job to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// DeleteJob
// Retrieves username param, validates request, calls DeleteJob service
func (jc *JobController) DeleteJob(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get job id from url params, parse into int
	jobId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// if admin, call DeleteJob service, otherwise call DeleteUsersJob to only allow job deletion for requesting users jobs
	var userId *int64 = nil
	if !c.Is_admin {
		userId = &c.ID
	}
	err = services.DeleteJob(jobId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to delete job: %v", err)
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Job ID %v has been deleted", jobId)
}

// AssignJobLabel
// Assign label to a job
func (jc *JobController) AssignJobLabel(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get URL query params, convert to int
	unassign := r.URL.Query().Get("unassign")
	assign := 1
	if unassign == "true" {
		assign = 0
	}
	// get job from url
	jobId, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Job ID must be an integer: %v", err)
		return
	}
	// get label id from url params, parse into int
	labelId, err := strconv.ParseInt(chi.URLParam(r, "labelId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be an integer: %v", err)
		return
	}
	// get Job Data
	job, err := services.GetJob(jobId)
	if job == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Job ID %d not found: %v", jobId, err)
		return
	}
	// if requesting users id doesnt match user from job, and they are not an admin, throw error
	if (c.ID != job.User) && !c.Is_admin {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to assign labels to other users jobs")
		return
	}
	// get Label Data
	label, err := services.GetLabel(labelId)
	if label == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Label ID %d not found: %v", labelId, err)
		return
	}
	// if requesting users id doesnt match user from label, and its not an unowned label, and they are not an admin, throw error
	if (label.User != nil) && (c.ID != *label.User) && !c.Is_admin {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to assign other users labels to jobs")
		return
	}
	_, err = services.AssignJobLabel(jobId, labelId, assign)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to assign label to job: %v", err)
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Label ID %v has been assigned to Job ID %v", labelId, jobId)
}
