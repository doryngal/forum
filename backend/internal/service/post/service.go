package post

import (
	"forum/internal/domain"
	post_repo "forum/internal/repository/sqlite/post"
	"github.com/google/uuid"
)

type service struct {
	repo post_repo.Repository
}

func New(repo post_repo.Repository) Service {
	return &service{repo: repo}
}

func (s service) CreatePost(post *domain.Post) error {
	//TODO implement me
	panic("implement me")
}

func (s service) GetPostByID(id uuid.UUID) (*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetAllPosts() ([]*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetPostsByCategory(categoryID uuid.UUID) ([]*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetPostsByUser(userID uuid.UUID) ([]*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetLikedPostsByUser(userID uuid.UUID) ([]*domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) LikePost(postID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s service) DislikePost(postID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
