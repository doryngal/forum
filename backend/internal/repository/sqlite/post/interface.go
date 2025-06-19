package post

import "forum/internal/domain"

type Repository interface {
	Create(post *domain.Post) error
	GetByID(id int) (*domain.Post, error)
	GetAll() ([]*domain.Post, error)
	GetByCategory(categoryID int) ([]*domain.Post, error)
	GetByUserID(userID int) ([]*domain.Post, error)
	GetLikedByUser(userID int) ([]*domain.Post, error)
	Like(postID, userID int) error
	Dislike(postID, userID int) error
}
