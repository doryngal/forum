package comment

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	Create(comment *domain.Comment) error
	GetByPostID(postID uuid.UUID) ([]*domain.Comment, error)
	Like(commentID, userID uuid.UUID) error
	Dislike(commentID, userID uuid.UUID) error
}
