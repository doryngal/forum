package category

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	Create(name string) error
	GetAll() ([]*domain.Category, error)
	AssignToPost(postID uuid.UUID, categoryIDs []uuid.UUID) error
	ExistsByName(name string) (bool, error)
	ExistsByID(id uuid.UUID) (bool, error)
}
