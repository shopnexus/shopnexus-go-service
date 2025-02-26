package model

import "errors"

var (
	ErrForbidden       = errors.New("forbidden")
	ErrMalformedParams = errors.New("malformed request parameters")
	ErrTokenInvalid    = errors.New("token invalid")
	ErrInvalidCreds    = errors.New("invalid credentials")

	//Account
	ErrWrongPassword = errors.New("wrong password")

	// Refund
	// ErrRefund
)
