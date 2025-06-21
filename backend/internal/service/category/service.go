package category

import (
	"forum/internal/domain"
	category_repo "forum/internal/repository/sqlite/category"
	"github.com/google/uuid"
)

type service struct {
	repo category_repo.Repository
}

func New(repo category_repo.Repository) Service {
	return &service{repo: repo}
}

func (s service) CreateCategory(name string) error {
	//TODO implement me
	panic("implement me")
}

func (s service) GetAllCategories() ([]*domain.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) AssignCategoriesToPost(postID uuid.UUID, categoryIDs []uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
