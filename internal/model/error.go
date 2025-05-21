package model

import (
	"errors"
	"fmt"
)

var (
	ErrForbidden        = errors.New("forbidden")
	ErrMalformedParams  = errors.New("malformed request parameters")
	ErrTokenInvalid     = errors.New("token invalid")
	ErrInvalidCreds     = errors.New("invalid credentials")
	ErrPermissionDenied = errors.New("permission denied")

	//Account
	ErrWrongPassword = errors.New("wrong password")

	// Refund
	// ErrRefund
)

type ErrorWithCode struct {
	Code string
	Msg  string
	Err  error // optional wrapped error
}

func (e *ErrorWithCode) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e *ErrorWithCode) Unwrap() error {
	return e.Err
}

var (
	ErrWrongCredentials = &ErrorWithCode{
		Code: "wrong_credentials",
		Msg:  "Wrong credentials",
	}
)
