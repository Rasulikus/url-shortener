package service

import (
	"errors"
)

var (
	ErrInvalidInput   = errors.New("service: invalid input")
	ErrNotFound       = errors.New("service: not found")
	ErrAliasCollision = errors.New("service: failed to generate unique alias")
	ErrInternalError  = errors.New("service: internal error")
)
