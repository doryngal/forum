package post

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Service interface {
	CreatePost(post *domain.Post) error
	GetPostByID(id uuid.UUID) (*domain.Post, error)
	GetAllPosts() ([]*domain.Post, error)
	GetPostsByUserID(userID, sessionID uuid.UUID) ([]*domain.Post, error)
	LikePost(postID, userID uuid.UUID) error
	DislikePost(postID, userID uuid.UUID) error
	UpdatePost(post *domain.Post, userID uuid.UUID) error
	DeletePost(postID, userID uuid.UUID) error

	GetLikedPosts(userID uuid.UUID) ([]*domain.Post, error)
	GetDislikedPosts(userID uuid.UUID) ([]*domain.Post, error)
}
