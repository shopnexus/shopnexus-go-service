package model

import "errors"

var (
	ErrForbidden    = errors.New("forbidden")
	ErrTokenInvalid = errors.New("token invalid")
)
