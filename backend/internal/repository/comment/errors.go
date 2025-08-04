package comment

import "errors"

var (
	ErrInsertCommentFailed  = errors.New("comment: failed to insert comment")
	ErrQueryFailed          = errors.New("comment: failed to query comments")
	ErrScanFailed           = errors.New("comment: failed to scan comment row")
	ErrUUIDParseFailed      = errors.New("comment: failed to parse UUID")
	ErrReactionUpdateFailed = errors.New("comment: failed to update reaction")
	ErrReactionNotFound     = errors.New("comment: reaction not found")
	ErrGetReactionFailed    = errors.New("comment: failed to get reaction")
)
