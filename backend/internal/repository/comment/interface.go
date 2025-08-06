package comment

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	Create(comment *domain.Comment) error
	GetByPostID(postID, userID uuid.UUID) ([]*domain.Comment, error)
	Like(commentID, userID uuid.UUID) error
	Dislike(commentID, userID uuid.UUID) error
	ExistsByID(id uuid.UUID) (bool, error)
	GetReaction(commentID, userID uuid.UUID) (int, error)
	GetCommentsByUserID(userID uuid.UUID) ([]*domain.CommentWithPostTitle, error)
}
