package user

import "errors"

var (
	ErrInvalidUser        = errors.New("user: invalid user ID or data")
	ErrInvalidCredentials = errors.New("user: invalid email or password")
	ErrUserNotFound       = errors.New("user: not found")
	ErrEmailTaken         = errors.New("user: email already in use")
)
