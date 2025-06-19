package post

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

func (r repository) Create(post *domain.Post) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetByID(id uuid.UUID) (*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetAll() ([]*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetByCategory(categoryID uuid.UUID) ([]*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetByUserID(userID uuid.UUID) ([]*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetLikedByUser(userID uuid.UUID) ([]*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Like(postID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Dislike(postID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
