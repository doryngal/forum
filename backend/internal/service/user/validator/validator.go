package validator

import (
	"errors"
	"forum/internal/domain"
	"net/mail"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail     = errors.New("validator: invalid email format")
	ErrEmptyUsername    = errors.New("validator: username cannot be empty")
	ErrEmptyPassword    = errors.New("validator: password cannot be empty")
	ErrPasswordTooShort = errors.New("validator: password is too short")
	ErrPasswordTooWeak  = errors.New("validator: password must contain at least one number or symbol")
)

type userValidator struct{}

func NewUserValidator() UserValidator {
	return &userValidator{}
}

func (v *userValidator) ValidateUser(user *domain.User) error {
	if _, err := mail.ParseAddress(strings.TrimSpace(user.Email)); err != nil {
		return ErrInvalidEmail
	}
	if strings.TrimSpace(user.Username) == "" {
		return ErrEmptyUsername
	}
	return v.ValidatePassword(user.PasswordHash)
}

func (v *userValidator) ValidatePassword(password string) error {
	password = strings.TrimSpace(password)

	if password == "" {
		return ErrEmptyPassword
	}
	if len(password) < 6 {
		return ErrPasswordTooShort
	}

	hasSymbolOrDigit := false
	for _, c := range password {
		if c >= '0' && c <= '9' || strings.ContainsRune("!@#$%^&*()-_=+[]{}<>.,", c) {
			hasSymbolOrDigit = true
			break
		}
	}
	if !hasSymbolOrDigit {
		return ErrPasswordTooWeak
	}
	return nil
}

func (v *userValidator) ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrInvalidEmail
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}
