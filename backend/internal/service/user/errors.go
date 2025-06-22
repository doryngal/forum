package user

import "errors"

var (
	ErrInvalidUser  = errors.New("user: invalid user ID or data")
	ErrUserNotFound = errors.New("user: not found")
	ErrEmailTaken   = errors.New("user: email already in use")
)
