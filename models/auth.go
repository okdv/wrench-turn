package models

import "github.com/golang-jwt/jwt/v5"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Is_admin bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}
