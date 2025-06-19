package category

import (
	"database/sql"
	"forum/internal/domain"
	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r repository) Create(name string) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetAll() ([]*domain.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) AssignToPost(postID uuid.UUID, categoryIDs []uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
