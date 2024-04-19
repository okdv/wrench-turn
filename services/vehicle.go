package services

import (
	"errors"
	"log"
	"strconv"

	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/models"
)

// GetVehicle
// Takes id as arg, passes to db query, returns Job
func GetVehicle(vehicleId int64) (*models.Vehicle, error) {
	vehicle, err := db.GetVehicle(vehicleId)
	return vehicle, err
}

// ListVehicles
// Takes URL query params as args, passes to ListVehicles query, returns Vehicle list
func ListVehicles(userId *string, jobId *string, searchStr *string, sort *string) ([]*models.Vehicle, error) {
	vehicles, err := db.ListVehicles(userId, jobId, searchStr, sort)
	return vehicles, err
}

// CreateVehicle
// Takes newVehicle as arg, passes to db query, calls GetVehicle, returns Vehicle
func CreateVehicle(newVehicle models.NewVehicle) (*models.Vehicle, error) {
	// set default values
	defaultBool := 0
	if newVehicle.Is_metric == nil {
		newVehicle.Is_metric = &defaultBool
	}
	// pass to db query, return new Vehicles id
	vehicleId, err := db.CreateVehicle(newVehicle)
	if err != nil || vehicleId == nil {
		err = errors.Join(err, errors.New("No ID of new Vehicle found"))
		return nil, err
	}
	// pass to GetVehicle, return Vehicle
	vehicle, err := GetVehicle(*vehicleId)
	return vehicle, err
}

// EditVehicle
// Takes User as arg, passes to EditVehicle query, returns updated User
func EditVehicle(editedVehicle models.Vehicle) (*models.Vehicle, error) {
	err := db.EditVehicle(editedVehicle)
	if err != nil {
		return nil, err
	}
	vehicle, err := GetVehicle(editedVehicle.ID)
	return vehicle, err
}

// DeleteVehicle
// Takes vehicle id as arg, passes to DeleteVehicle query
func DeleteVehicle(vehicleId int64, userId *int64) error {
	vehicleIdStr := strconv.FormatInt(vehicleId, 10)
	// get vehicles jobs
	jobs, err := ListJobs(nil, &vehicleIdStr, nil, nil, nil, nil, nil)
	if err != nil {
		log.Printf("Could not get vehicles jobs: %v", err)
	}
	// if there are jobs, delete them
	if len(jobs) > 0 {
		// loop through jobs
		for i := 0; i < len(jobs); i++ {
			err = DeleteJob(jobs[i].ID, nil)
			if err != nil {
				log.Printf("Could not delete job ID %d: %v", jobs[i].ID, err)
			}
		}
	}
	err = db.DeleteVehicle(vehicleId, userId)
	return err
}
