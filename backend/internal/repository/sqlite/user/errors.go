package user

import "errors"

var (
	ErrInsertUserFailed  = errors.New("user: failed to insert user")
	ErrQueryFailed       = errors.New("user: failed to query user")
	ErrUUIDParseFailed   = errors.New("user: failed to parse UUID")
	ErrCheckExistsFailed = errors.New("user: failed to check email existence")
)
