package post

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Service interface {
	CreatePost(post *domain.Post) error
	GetPostByID(id uuid.UUID) (*domain.Post, error)
	GetAllPosts() ([]*domain.Post, error)
	GetPostsByCategory(categoryID uuid.UUID) ([]*domain.Post, error)
	GetPostsByUser(userID uuid.UUID) ([]*domain.Post, error)
	GetLikedPostsByUser(userID uuid.UUID) ([]*domain.Post, error)
	LikePost(postID, userID uuid.UUID) error
	DislikePost(postID, userID uuid.UUID) error
	UpdatePost(post *domain.Post) error
	DeletePost(postID, userID uuid.UUID) error
}
