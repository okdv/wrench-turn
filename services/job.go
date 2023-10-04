package services

import (
	"errors"

	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/models"
)

// GetJob
// Takes id as arg, passes to db query, returns Job
func GetJob(jobId int64) (*models.Job, error) {
	job, err := db.GetJob(jobId)
	return job, err
}

// CreateJob
// Takes newJob as arg, passes to db query, calls GetJob, returns Job
func CreateJob(newJob models.NewJob) (*models.Job, error) {
	// set default values
	defaultBool := 0
	if newJob.Is_template == nil {
		newJob.Is_template = &defaultBool
	}
	if newJob.Repeats == nil {
		newJob.Repeats = &defaultBool
	}
	// pass to db query, return new Jobs id
	jobId, err := db.CreateJob(newJob)
	if err != nil || jobId == nil {
		err = errors.Join(err, errors.New("No ID of new Job found"))
		return nil, err
	}
	// pass to GetJob, return Job
	job, err := GetJob(*jobId)
	return job, err
}

// EditJob
// Takes User as arg, passes to EditJob query, returns updated User
func EditJob(editedJob models.Job) (*models.Job, error) {
	err := db.EditJob(editedJob)
	if err != nil {
		return nil, err
	}
	job, err := GetJob(editedJob.ID)
	return job, err
}

// ListJobs
// Takes URL query params as args, passes to ListJobs query, returns Job list
func ListJobs(userId *string, vehicleId *string, isTemplate *string, searchStr *string, sort *string) ([]*models.Job, error) {
	users, err := db.ListJobs(userId, vehicleId, isTemplate, searchStr, sort)
	return users, err
}

// DeleteJob
// Takes job id as arg, passes to DeleteJob query
func DeleteJob(jobId int64, userId *int64) error {
	err := db.DeleteJob(jobId, userId)
	return err
}
