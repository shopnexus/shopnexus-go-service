package model

import "errors"

var (
	ErrForbidden     = errors.New("forbidden")
	ErrTokenInvalid  = errors.New("token invalid")
	ErrInvalidCreds  = errors.New("invalid credentials")
	ErrWrongPassword = errors.New("wrong password")
)
