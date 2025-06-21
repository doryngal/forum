package comment

import (
	"forum/internal/domain"
	comment_repo "forum/internal/repository/sqlite/comment"
	"github.com/google/uuid"
)

type service struct {
	repo comment_repo.Repository
}

func New(repo comment_repo.Repository) Service {
	return &service{repo: repo}
}

func (s service) CreateComment(comment *domain.Comment) error {
	//TODO implement me
	panic("implement me")
}

func (s service) GetCommentsByPost(postID uuid.UUID) ([]*domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) LikeComment(commentID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s service) DislikeComment(commentID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
