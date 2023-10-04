package models

import "time"

type NewJob struct {
	// meta data
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Instructions *string `json:"instructions"`
	Is_template  *int    `json:"isTemplate"`
	// ownership
	Vehicle    *int64 `json:"vehicle"`
	User       *int64 `json:"user"`
	Origin_job *int64 `json:"originJob"`
	// repeats
	Repeats            *int    `json:"repeats"`
	Odo_interval       *int64  `json:"odoInterval"`
	Time_interval      *int64  `json:"timeInterval"`
	Time_interval_unit *string `json:"timeIntervalUnit"`
	// times
	Due_date *time.Time `json:"dueDate"`
}

type Job struct {
	// meta data
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Instructions *string `json:"instructions"`
	Is_template  int     `json:"isTemplate"`
	Is_complete  int     `json:"isComplete"`
	// ownership
	Vehicle    *int64 `json:"vehicle"`
	User       int64  `json:"user"`
	Origin_job *int64 `json:"originJob"`
	// repeats
	Repeats            int     `json:"repeats"`
	Odo_interval       *int64  `json:"odoInterval"`
	Time_interval      *int64  `json:"timeInterval"`
	Time_interval_unit *string `json:"timeIntervalUnit"`
	// times
	Due_date     *time.Time `json:"dueDate"`
	Completed_at *time.Time `json:"completedAt"`
	Created_at   time.Time  `json:"createdAt"`
	Updated_at   time.Time  `json:"updatedAt"`
}
