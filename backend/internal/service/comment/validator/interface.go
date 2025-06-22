package validator

import "forum/internal/domain"

type CommentValidator interface {
	ValidateComment(comment *domain.Comment) error
}
