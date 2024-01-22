package services

import (
	"errors"

	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/models"
)

// GetTask
// Takes ids as args, passes to db query, returns Task
func GetTask(jobId int64, taskId int64) (*models.Task, error) {
	task, err := db.GetTask(jobId, taskId)
	return task, err
}

// CreateTask
// Takes newTask as arg, passes to db query, calls GetTask, returns Task
func CreateTask(newTask models.NewTask, jobId int64) (*models.Task, error) {
	// pass to db query, return new Tasks id
	taskId, err := db.CreateTask(newTask, jobId)
	if err != nil || taskId == nil {
		err = errors.Join(err, errors.New("No ID of new Task found"))
		return nil, err
	}
	// pass to GetTask, return Task
	task, err := GetTask(jobId, *taskId)
	return task, err
}

// EditTask
// Takes edited task, jobid as args, passes to EditTask query, returns updated Task
func EditTask(editedTask models.Task, jobId int64) (*models.Task, error) {
	err := db.EditTask(editedTask, jobId)
	if err != nil {
		return nil, err
	}
	task, err := GetTask(jobId, editedTask.ID)
	return task, err
}

// MarkComplete
// Takes job id, task id, complete status as args, passes to MarkComplete query
func MarkComplete(jobId int64, taskId int64, status int) error {
	err := db.UpdateTaskStatus(jobId, taskId, status)
	return err
}

// ListTasks
// Takes URL query params as args, passes to ListTasks query, returns Task list
func ListTasks(jobId int64, isComplete *string, searchStr *string, sort *string) ([]*models.Task, error) {
	users, err := db.ListTasks(jobId, isComplete, searchStr, sort)
	return users, err
}

// DeleteTask
// Takes job id, task id as args, passes to DeleteTask query
func DeleteTask(jobId int64, taskId int64) error {
	err := db.DeleteTask(jobId, taskId)
	return err
}
