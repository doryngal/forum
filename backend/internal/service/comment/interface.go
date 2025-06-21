package comment

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Service interface {
	CreateComment(comment *domain.Comment) error
	GetCommentsByPost(postID uuid.UUID) ([]*domain.Comment, error)
	LikeComment(commentID, userID uuid.UUID) error
	DislikeComment(commentID, userID uuid.UUID) error
}
