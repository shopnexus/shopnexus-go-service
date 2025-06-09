package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID int64       `json:"user_id"`
	Role   AccountType `json:"role"`
	jwt.RegisteredClaims
}
