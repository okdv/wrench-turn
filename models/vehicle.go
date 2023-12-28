package models

import "time"

type NewVehicle struct {
	// meta data
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	Is_metric   *int    `json:"isMetric"`
	// vehicle data
	Vin   *string `json:"vin"`
	Year  *int64  `json:"year"`
	Make  *string `json:"make"`
	Model *string `json:"model"`
	Trim  *string `json:"trim"`
	// life data
	Odometer *int `json:"odometer"`
	// ownership
	User *int64 `json:"user"`
}

type Vehicle struct {
	// meta data
	ID          int64   `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	Is_metric   *int    `json:"isMetric"`
	// vehicle data
	Vin   *string `json:"vin"`
	Year  *int64  `json:"year"`
	Make  *string `json:"make"`
	Model *string `json:"model"`
	Trim  *string `json:"trim"`
	// life data
	Odometer *int64 `json:"odometer"`
	// ownership
	User int64 `json:"user"`
	// times
	Created_at time.Time `json:"createdAt"`
	Updated_at time.Time `json:"updatedAt"`
}
