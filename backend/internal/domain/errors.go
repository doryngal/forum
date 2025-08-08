package domain

import "errors"

var (
	ErrInvalidAction  = errors.New("invalid action")
	ErrInvalidSession = errors.New("invalid session")
	ErrForbidden      = errors.New("forbidden")
	ErrUnauthorized   = errors.New("unauthorized")

	ErrTitleRequired     = errors.New("title is required")
	ErrContentRequired   = errors.New("content is required")
	ErrCategoryRequired  = errors.New("category is required")
	ErrInvalidCategoryID = errors.New("invalid category ID")

	ErrInvalidPostID      = errors.New("invalid post id")
	ErrPasswordTooShort   = errors.New("password is too short")
	ErrEmailRequired      = errors.New("email is required")
	ErrUsernameRequired   = errors.New("username is required")
	ErrPasswordsNotMatch  = errors.New("passwords do not match")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
