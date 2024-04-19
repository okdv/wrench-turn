package services

import (
	"errors"
	"log"
	"strconv"

	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/models"
)

// GetLabel
// Takes id as arg, passes to db query, returns Label
func GetLabel(labelId int64) (*models.Label, error) {
	label, err := db.GetLabel(labelId)
	return label, err
}

// CreateLabel
// Takes newLabel as arg, passes to db query, calls GetLabel, returns Label
func CreateLabel(newLabel models.NewLabel) (*models.Label, error) {
	// pass to db query, return new Labels id
	labelId, err := db.CreateLabel(newLabel)
	if err != nil || labelId == nil {
		err = errors.Join(err, errors.New("No ID of new Label found"))
		return nil, err
	}
	// pass to GetLabel, return Label
	label, err := GetLabel(*labelId)
	return label, err
}

// EditLabel
// Takes User as arg, passes to EditLabel query, returns updated User
func EditLabel(editedLabel models.Label) (*models.Label, error) {
	err := db.EditLabel(editedLabel)
	if err != nil {
		return nil, err
	}
	label, err := GetLabel(editedLabel.ID)
	return label, err
}

// ListLabels
// Takes URL query params as args, passes to ListLabels query, returns Label list
func ListLabels(userId *string, jobId *string, searchStr *string, sort *string) ([]*models.Label, error) {
	users, err := db.ListLabels(userId, jobId, searchStr, sort)
	return users, err
}

// DeleteLabel
// Takes label id as arg, passes to DeleteLabel query
func DeleteLabel(labelId int64, userId *int64) error {
	labelIdStr := strconv.FormatInt(labelId, 10)
	// get labels jobs
	jobs, err := ListJobs(nil, nil, nil, nil, &labelIdStr, nil, nil)
	if err != nil {
		log.Printf("Could not get jobs labels: %v", err)
	}
	// if there are labels, delete them
	if len(jobs) > 0 {
		// loop through labels
		for i := 0; i < len(jobs); i++ {
			_, err = AssignJobLabel(jobs[i].ID, labelId, 0)
			if err != nil {
				log.Printf("Could not remove label from job ID %d: %v", jobs[i].ID, err)
			}
		}
	}
	err = db.DeleteLabel(labelId, userId)
	return err
}
