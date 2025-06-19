package category

import "forum/internal/domain"

type Repository interface {
	Create(name string) error
	GetAll() ([]*domain.Category, error)
	AssignToPost(postID int, categoryIDs []int) error
}
