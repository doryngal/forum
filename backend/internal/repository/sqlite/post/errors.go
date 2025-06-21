package post

import "errors"

var (
	ErrInsertPostFailed     = errors.New("post: failed to insert post")
	ErrQueryFailed          = errors.New("post: failed to query posts")
	ErrScanFailed           = errors.New("post: failed to scan row")
	ErrUUIDParseFailed      = errors.New("post: failed to parse UUID")
	ErrLoadCategoriesFailed = errors.New("post: failed to load categories")
	ErrReactionUpdateFailed = errors.New("post: failed to update reaction")
)
