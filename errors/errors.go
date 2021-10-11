package errors

import (
	"errors"
)

func New(s string) error {
	return errors.New(s)
}

type Error string

func (err Error) Error() string {
	return string(err)
}

const (
	ErrNotFound   Error = "not found"
	ErrUserExists Error = "user already exists "
)
