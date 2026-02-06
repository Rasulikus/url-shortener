package validate

import (
	"errors"
	"net/url"
)

var ErrInvalidURL = errors.New("invalid url")

func URL(u string) error {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return ErrInvalidURL
	}

	return nil
}
