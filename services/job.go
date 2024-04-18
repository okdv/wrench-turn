package services

import (
	"errors"
	"log"
	"strconv"

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
func ListJobs(userId *string, vehicleId *string, isTemplate *string, isComplete *string, labelId *string, searchStr *string, sort *string) ([]*models.Job, error) {
	users, err := db.ListJobs(userId, vehicleId, isTemplate, isComplete, labelId, searchStr, sort)
	return users, err
}

// DeleteJob
// Takes job id as arg, passes to DeleteJob query
func DeleteJob(jobId int64, userId *int64) error {
	// init values
	// TODO make some more consistency on data types at different app layers, i was pretty sloppy with this...
	jobIdStr := strconv.FormatInt(jobId, 10)
	// get jobs tasks
	tasks, err := ListTasks(jobId, nil, nil, nil)
	if err != nil {
		log.Printf("Could not get jobs tasks: %v", err)
	}
	// if there are tasks, delete them
	if len(tasks) > 0 {
		// delete jobs tasks
		err = DeleteTask(jobId, nil)
		if err != nil {
			log.Printf("Could not delete jobs tasks: %v", err)
		}
	}
	// get jobs alerts
	alerts, err := ListAlerts(nil, nil, &jobIdStr, nil, nil, nil, nil, nil, nil)
	if err != nil {
		log.Printf("Could not get jobs alerts: %v", err)
	}
	// if there are alerts, delete them
	if len(alerts) > 0 {
		// loop through alerts
		for i := 0; i < len(alerts); i++ {
			err = DeleteAlert(alerts[i].ID, userId)
			if err != nil {
				log.Printf("Could not delete job alert ID %d: %v", alerts[i].ID, err)
			}
		}
	}
	// get jobs labels
	labels, err := ListLabels(nil, &jobIdStr, nil, nil)
	if err != nil {
		log.Printf("Could not get jobs labels: %v", err)
	}
	// if there are labels, delete them
	if len(labels) > 0 {
		// loop through labels
		for i := 0; i < len(labels); i++ {
			_, err = AssignJobLabel(jobId, labels[i].ID, 0)
			if err != nil {
				log.Printf("Could not remove job label ID %d: %v", labels[i].ID, err)
			}
		}
	}
	// delete job
	err = db.DeleteJob(jobId, userId)
	if err != nil {
		return err
	}
	return nil
}

// AssignJobLabel
// Takes job id, label id, creates an entry into job_label if assigning, otherwise deletes any existing entry
func AssignJobLabel(jobId int64, taskId int64, assign int) (*int64, error) {
	// if assigning, call that query and return
	if assign == 1 {
		relationshipId, err := db.AssignJobLabel(jobId, taskId)
		return relationshipId, err
	}
	// otherwise call unassign query
	err := db.UnassignJobLabel(jobId, taskId)
	return nil, err
}
