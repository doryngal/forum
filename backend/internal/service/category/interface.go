package category

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Service interface {
	CreateCategory(name string) error
	GetAllCategories() ([]*domain.Category, error)
	AssignCategoriesToPost(postID uuid.UUID, categoryIDs []uuid.UUID) error
}
