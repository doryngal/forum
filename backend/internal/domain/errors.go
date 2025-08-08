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

var (
	ErrInsertUserFailed  = errors.New("user: failed to insert user")
	ErrQueryFailed       = errors.New("user: failed to query user")
	ErrUUIDParseFailed   = errors.New("user: failed to parse UUID")
	ErrCheckExistsFailed = errors.New("user: failed to check email existence")
)

var (
	ErrInsertPostFailed     = errors.New("post: failed to insert post")
	ErrQueryFailedP         = errors.New("post: failed to query posts")
	ErrScanFailed           = errors.New("post: failed to scan row")
	ErrUUIDParseFailedP     = errors.New("post: failed to parse UUID")
	ErrLoadCategoriesFailed = errors.New("post: failed to load categories")
	ErrReactionUpdateFailed = errors.New("post: failed to update reaction")
	ErrReactionNotFound     = errors.New("post: reaction not found")
	ErrGetReactionFailed    = errors.New("post: failed to get reaction")
)

var (
	ErrInvalidUser    = errors.New("user: invalid user ID or data")
	ErrUserNotFound   = errors.New("user: not found")
	ErrInvalidComment = errors.New("comment: invalid comment data")
	ErrPostNotFound   = errors.New("post: not found")
	ErrEmailTaken     = errors.New("user: email already in use")
	ErrUsernameTaken  = errors.New("user: username already in use")
)
