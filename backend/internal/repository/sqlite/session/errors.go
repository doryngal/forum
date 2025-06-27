package session

import "errors"

var (
	ErrSessionNotFound       = errors.New("session: not found")
	ErrSessionExpired        = errors.New("session: expired")
	ErrCreateSession         = errors.New("session: failed to create session")
	ErrDeleteSession         = errors.New("session: failed to delete session")
	ErrAffectedRowsCheckFail = errors.New("session: failed to get affected rows")
)
