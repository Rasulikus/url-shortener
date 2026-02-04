package validate

import (
	"errors"
	"net/url"
)

var InvalidURL = errors.New("invalid url")

func URL(u string) error {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return InvalidURL
	}
	return nil
}
