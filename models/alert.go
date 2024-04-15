package models

import (
	"time"
)

type NewAlert struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Type        string     `json:"type"`
	User        *int64     `json:"user"`
	Vehicle     *int64     `json:"vehicle"`
	Job         *int64     `json:"job"`
	Task        *int64     `json:"task"`
	Alert_at    *time.Time `json:"alertAt"`
}

type Alert struct {
	ID          int64      `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Type        string     `json:"type"`
	User        int64      `json:"user"`
	Vehicle     *int64     `json:"vehicle"`
	Job         *int64     `json:"job"`
	Task        *int64     `json:"task"`
	Is_read     *int       `json:"isRead"`
	Read_at     *time.Time `json:"readAt"`
	Alert_at    *time.Time `json:"alertAt"`
	Created_at  time.Time  `json:"createdAt"`
	Updated_at  time.Time  `json:"updatedAt"`
}
