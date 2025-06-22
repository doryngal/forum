package validator

import "forum/internal/domain"

type PostValidator interface {
	ValidatePost(post *domain.Post) error
}
