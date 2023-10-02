package models

import "time"

type Passwords struct {
	Username           string  `json:"username"`
	CurrentPassword    *string `json:"currentPassword"`
	NewPassword        string  `json:"newPassword"`
	ConfirmNewPassword string  `json:"confirmNewPassword"`
}

type NewUser struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Is_admin *int    `json:"isAdmin"`
	Email    *string `json:"email"`
}

type User struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       *string    `json:"email"`
	Description *string    `json:"description"`
	Hashed_pw   *[]byte    `json:"hashedPw"`
	Is_admin    *int       `json:"isAdmin"`
	Created_at  *time.Time `json:"createdAt"`
	Updated_at  *time.Time `json:"updatedAt"`
}
