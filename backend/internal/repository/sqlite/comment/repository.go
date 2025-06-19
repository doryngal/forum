package comment

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

func (r repository) Create(comment *domain.Comment) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) GetByPostID(postID uuid.UUID) ([]*domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Like(commentID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Dislike(commentID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
