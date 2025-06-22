package comment

import "errors"

var (
	ErrInvalidComment  = errors.New("comment: invalid comment data")
	ErrCommentNotFound = errors.New("comment: not found")
	ErrUnauthorized    = errors.New("comment: unauthorized action")
	ErrAlreadyReacted  = errors.New("comment: already reacted")
)
