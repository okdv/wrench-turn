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

type TaskController struct {
}

func NewTaskController() *TaskController {
	return &TaskController{}
}

// GetTask
// Retrieves id param, calls GetTask services, returns Task
func (tc *TaskController) GetTask(w http.ResponseWriter, r *http.Request) {
	// get task id from url params, parse into int
	jobId, jobErr := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	taskId, taskErr := strconv.ParseInt(chi.URLParam(r, "taskId"), 10, 64)
	if jobErr != nil || taskErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Ids must be an integer: %v %v", jobErr, taskErr)
		return
	}
	// call GetTask service, return Task
	task, err := services.GetTask(jobId, taskId)
	if err != nil || task == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Task not found: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert task to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// ListTasks
// Retrieves any URL query params, calls ListTasks service, returns Task list
func (tc *TaskController) ListTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []*models.Task
	// get job from url
	jobId, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	// get URL query params
	isComplete := r.URL.Query().Get("template")
	searchStr := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	// call ListTasks service
	tasks, err = services.ListTasks(jobId, &isComplete, &searchStr, &sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to retrieve any tasks: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to convert tasks to JSON response")
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// CreateTask
// Takes NewTask as request body, validates it, calls CreateTask service, return Task
func (tc *TaskController) CreateTask(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var newTask *models.NewTask
	// get job from url
	jobId, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Job ID must be an integer: %v", err)
		return
	}
	// get Job Data
	job, err := services.GetJob(jobId)
	if job == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Job ID %d not found: %v", jobId, err)
		return
	}
	// get task data from request body
	err = json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// if newTask user is not requesting user, check if admin
	if job.User != c.ID && c.Is_admin == false {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to create tasks for other users jobs")
		return
	}
	// send to newTask service, return Task
	task, err := services.CreateTask(*newTask, jobId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to create task: %v", err)
		return
	}
	// covnert to JSON response
	jsonData, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert task to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// EditTask
// Takes Task as request body, calls EditTask service, return Task
func (tc *TaskController) EditTask(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	var task models.Task
	// get job from url
	jobId, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Job ID must be an integer: %v", err)
		return
	}
	// get Job Data
	job, err := services.GetJob(jobId)
	if job == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Job ID %d not found: %v", jobId, err)
		return
	}
	// get task data from request body
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}
	// get existing task Data
	currentTask, err := services.GetTask(jobId, task.ID)
	if currentTask == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Current task ID %d not found: %v", task.ID, err)
		return
	}
	// if requesting users id doesnt match user from job, and they are not an admin, throw error
	if (c.ID != job.User) && (c.Is_admin != true) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to edit tasks of other users")
		return
	}
	// call EditTask service, return updated Task
	updatedTask, err := services.EditTask(task, jobId)
	if err != nil || updatedTask == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to edit task: %v", err)
		return
	}
	// convert to JSON response
	jsonData, err := json.Marshal(updatedTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to convert task to JSON response: %v", err)
		return
	}
	// respond with json
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// MarkComplete
// Marks task complete, or incomplete
func (tc *TaskController) MarkComplete(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get URL query params, convert to int
	incomplete := r.URL.Query().Get("incomplete")
	status := 1
	if incomplete == "true" {
		status = 0
	}
	// get job from url
	jobId, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Job ID must be an integer: %v", err)
		return
	}
	// get task id from url params, parse into int
	taskId, err := strconv.ParseInt(chi.URLParam(r, "taskId"), 10, 64)
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
	if (c.ID != job.User) && (c.Is_admin != true) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to edit tasks of other users")
		return
	}
	err = services.MarkComplete(jobId, taskId, status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to mark task complete: %v", err)
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task ID %v has been marked as complete", taskId)
}

// DeleteTask
// Retrieves username param, validates request, calls DeleteTask service
func (tc *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request, c *models.Claims) {
	// get job from url
	jobId, err := strconv.ParseInt(chi.URLParam(r, "jobId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Job ID must be an integer: %v", err)
		return
	}
	// get task id from url params, parse into int
	taskId, err := strconv.ParseInt(chi.URLParam(r, "taskId"), 10, 64)
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
	if (c.ID != job.User) && (c.Is_admin != true) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Must be admin to edit tasks of other users")
		return
	}
	err = services.DeleteTask(jobId, taskId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to delete task: %v", err)
		return
	}
	// respond with text
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task ID %v has been deleted", taskId)
}
