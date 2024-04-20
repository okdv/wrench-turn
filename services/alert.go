package services

import (
	"errors"
	"time"

	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/models"
)

// GetAlert
// Takes id as arg, passes to db query, returns Alert
func GetAlert(alertId int64) (*models.Alert, error) {
	alert, err := db.GetAlert(alertId)
	return alert, err
}

// CreateAlert
// Takes newAlert as arg, passes to db query, calls GetAlert, returns Alert
func CreateAlert(newAlert models.NewAlert) (*models.Alert, error) {
	// pass to db query, return new Alerts id
	alertId, err := db.CreateAlert(newAlert)
	if err != nil || alertId == nil {
		err = errors.Join(err, errors.New("No ID of new Alert found"))
		return nil, err
	}
	// pass to GetAlert, return Alert
	alert, err := GetAlert(*alertId)
	return alert, err
}

// EditAlert
// Takes User as arg, passes to EditAlert query, returns updated User
func EditAlert(editedAlert models.Alert) (*models.Alert, error) {
	err := db.EditAlert(editedAlert)
	if err != nil {
		return nil, err
	}
	alert, err := GetAlert(editedAlert.ID)
	return alert, err
}

// ListAlerts
// Takes URL query params as args, passes to ListAlerts query, returns Alert list
func ListAlerts(userId *string, vehicleId *string, jobId *string, taskId *string, typeStr *string, isRead *string, isAlerted *string, searchStr *string, sort *string) ([]*models.Alert, error) {
	var alertDate *string
	if isAlerted != nil {
		if *isAlerted == "true" {
			currentDatetime := time.Now()
			currentDatetimeStr := currentDatetime.Format("2006-01-02T15:04:05Z")
			alertDate = &currentDatetimeStr
		}
	}
	users, err := db.ListAlerts(userId, vehicleId, jobId, taskId, typeStr, isRead, alertDate, searchStr, sort)
	return users, err
}

// DeleteAlert
// Takes alert id as arg, passes to DeleteAlert query
func DeleteAlert(alertId int64, userId *int64) error {
	err := db.DeleteAlert(alertId, userId)
	return err
}

// MarkRead
// Takes job id, task id, complete status as args, passes to MarkRead query
func MarkRead(alertId int64, userId int64, status int) error {
	err := db.UpdatedAlertStatus(alertId, userId, status)
	return err
}
