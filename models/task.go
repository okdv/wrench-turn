package models

import "time"

// used for new job forms
type NewTask struct {
	// meta data
	Name        string  `json:"name"`
	Description *string `json:"description"`
	// part
	Part_name *string `json:"partName"`
	Part_link *string `json:"partLink"`
	// times
	Due_date *time.Time `json:"dueDate"`
}

// used for existing job data
type Task struct {
	// meta data
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Is_complete int     `json:"isComplete"`
	// ownership
	Job *int64 `json:"job"`
	// part
	Part_name *string `json:"partName"`
	Part_link *string `json:"partLink"`
	// times
	Due_date     *time.Time `json:"dueDate"`
	Completed_at *time.Time `json:"completedAt"`
	Created_at   time.Time  `json:"createdAt"`
	Updated_at   time.Time  `json:"updatedAt"`
}
