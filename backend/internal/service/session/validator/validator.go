package validator

import (
	"errors"
	"forum/internal/domain"
	"github.com/google/uuid"
	"time"
	"unicode/utf8"
)

var (
	ErrInvalidToken  = errors.New("invalid session token")
	ErrInvalidUserID = errors.New("invalid user ID")
	ErrInvalidExpiry = errors.New("invalid expiry time")
	ErrTokenTooShort = errors.New("token is too short")
	ErrTokenTooLong  = errors.New("token is too long")
	ErrExpiryInPast  = errors.New("expiry time is in the past")
)

type sessionValidator struct{}

func NewSessionValidator() SessionValidator {
	return &sessionValidator{}
}

func (v *sessionValidator) ValidateCreate(session *domain.Session) error {
	if session == nil {
		return ErrInvalidToken
	}

	if err := v.ValidateToken(session.ID); err != nil {
		return err
	}

	if uuid.Nil == session.UserID {
		return ErrInvalidUserID
	}

	if session.ExpiresAt.IsZero() || session.ExpiresAt.Before(time.Now()) {
		return ErrExpiryInPast
	}

	return nil
}

func (v *sessionValidator) ValidateToken(token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	length := utf8.RuneCountInString(token)
	if length < 32 {
		return ErrTokenTooShort
	}
	if length > 512 {
		return ErrTokenTooLong
	}

	return nil
}
