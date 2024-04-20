package models

import "github.com/golang-jwt/jwt/v5"

// used for auth purposes
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// used for storing data in JWT
type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Is_admin bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}
