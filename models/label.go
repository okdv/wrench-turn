package models

import "time"

// used for new label forms
type NewLabel struct {
	// meta data
	Name  string  `json:"name"`
	Color *string `json:"color"`
	// ownership
	User *int64 `json:"user"`
}

// used for existing label data
type Label struct {
	// meta data
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Color *string `json:"color"`
	// ownership
	User *int64 `json:"user"`
	// times
	Created_at time.Time `json:"createdAt"`
	Updated_at time.Time `json:"updatedAt"`
}
