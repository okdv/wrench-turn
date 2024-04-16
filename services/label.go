package services

import (
	"errors"

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
	// set default values
	if newLabel.Color == nil {
		defaultColor := "blue"
		newLabel.Color = &defaultColor
	}
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
func ListLabels(userId *string, searchStr *string, sort *string) ([]*models.Label, error) {
	users, err := db.ListLabels(userId, searchStr, sort)
	return users, err
}

// DeleteLabel
// Takes label id as arg, passes to DeleteLabel query
func DeleteLabel(labelId int64, userId *int64) error {
	err := db.DeleteLabel(labelId, userId)
	return err
}
