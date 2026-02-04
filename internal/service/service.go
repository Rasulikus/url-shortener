package service

import (
	"errors"
)

var (
	ErrInvalidInput        = errors.New("service: invalid input")
	ErrUserAlreadyExists   = errors.New("service: user already exists")
	ErrInvalidCredentials  = errors.New("service: invalid credentials")
	ErrRefreshTokenInvalid = errors.New("service: refresh token invalid")
)

type URLService interface {
}
