package domain

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")

	InternalError = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
)
