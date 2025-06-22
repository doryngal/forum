package validator

import (
	"errors"
	"forum/internal/domain"
	"strings"
)

var (
	ErrEmptyTitle   = errors.New("validator: title cannot be empty")
	ErrEmptyContent = errors.New("validator: content cannot be empty")
	ErrTooShort     = errors.New("validator: post is too short")
	ErrTooLong      = errors.New("validator: post is too long")
)

type postValidator struct{}

func NewPostValidator() PostValidator {
	return &postValidator{}
}

func (v *postValidator) ValidatePost(post *domain.Post) error {
	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)

	if post.Title == "" {
		return ErrEmptyTitle
	}
	if post.Content == "" {
		return ErrEmptyContent
	}
	if len(post.Content) < 10 {
		return ErrTooShort
	}
	if len(post.Content) > 10000 {
		return ErrTooLong
	}
	return nil
}
