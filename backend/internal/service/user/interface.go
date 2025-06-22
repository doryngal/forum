package user

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Service interface {
	RegisterUser(user *domain.User) error
	Login(email, password string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
	IsEmailTaken(email string) (bool, error)
}
