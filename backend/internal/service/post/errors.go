package post

import "errors"

var (
	ErrInvalidPost    = errors.New("post: invalid post ID or data")
	ErrPostNotFound   = errors.New("post: not found")
	ErrAlreadyReacted = errors.New("post: already reacted")
)
