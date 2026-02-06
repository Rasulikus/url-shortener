package service

import (
	"errors"
)

var (
	ErrInvalidInput  = errors.New("service: invalid input")
	ErrNotFound      = errors.New("service: not found")
	ErrConflict      = errors.New("service: conflict")
	ErrInternalError = errors.New("service: internal error")
)
