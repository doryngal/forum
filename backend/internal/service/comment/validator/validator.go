package validator

import (
	"errors"
	"forum/internal/domain"
	"strings"
)

var (
	ErrEmptyContent = errors.New("validator: comment content cannot be empty")
	ErrTooShort     = errors.New("validator: comment content is too short")
	ErrTooLong      = errors.New("validator: comment content is too long")
)

type commentValidator struct{}

func NewCommentValidator() CommentValidator {
	return &commentValidator{}
}

func (v *commentValidator) ValidateComment(comment *domain.Comment) error {
	content := strings.TrimSpace(comment.Content)

	if content == "" {
		return ErrEmptyContent
	}
	if len(content) < 2 {
		return ErrTooShort
	}
	if len(content) > 1000 {
		return ErrTooLong
	}
	return nil
}
