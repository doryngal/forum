package user

import (
	"forum/internal/domain"
	user_repo "forum/internal/repository/sqlite/user"
	"github.com/google/uuid"
)

type service struct {
	repo user_repo.Repository
}

func New(repo user_repo.Repository) Service {
	return &service{repo: repo}
}

func (s service) RegisterUser(user *domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (s service) GetUserByEmail(email string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetUserByID(id uuid.UUID) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) IsEmailTaken(email string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
