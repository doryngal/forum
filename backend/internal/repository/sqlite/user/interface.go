package user

import "forum/internal/domain"

type Repository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id int) (*domain.User, error)
	IsEmailTaken(email string) (bool, error)
}
