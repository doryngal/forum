package post

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	Create(post *domain.Post) error
	GetByID(id uuid.UUID) (*domain.Post, error)
	GetAll() ([]*domain.Post, error)
	GetByCategory(categoryID uuid.UUID) ([]*domain.Post, error)
	GetByUserID(userID, sessionID uuid.UUID) ([]*domain.Post, error)
	GetLikedByUser(userID uuid.UUID) ([]*domain.Post, error)
	Like(postID, userID uuid.UUID) error
	Dislike(postID, userID uuid.UUID) error
	ExistsByID(postID uuid.UUID) (bool, error)
	GetReaction(postID, userID uuid.UUID) (int, error)

	Update(post *domain.Post) error
	Delete(postID uuid.UUID, userID uuid.UUID) error

	GetLikedPostsByUserID(userID uuid.UUID) ([]*domain.Post, error)
	GetDislikedPostsByUserID(userID uuid.UUID) ([]*domain.Post, error)
}
