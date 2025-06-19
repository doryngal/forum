package user

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

func (r repository) Create(user *domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindByEmail(email string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindByID(id uuid.UUID) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) IsEmailTaken(email string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
