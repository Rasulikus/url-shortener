package generator

import (
	"errors"
)

const (
	DefaultLength = 10
	alphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var (
	ErrInvalidLength = errors.New("alias generator: invalid length")
	ErrOverflow      = errors.New("alias generator: counter value overflow")
)
