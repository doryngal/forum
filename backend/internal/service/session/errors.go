package session

import "errors"

var (
	ErrSessionCreationFailed = errors.New("session: failed to create session")
	ErrInvalidSession        = errors.New("session: invalid session")
	ErrSessionExpired        = errors.New("session: session expired")
	ErrSessionNotFound       = errors.New("session: session not found")
)
