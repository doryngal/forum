package user

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uuid.UUID) (*domain.User, error)
	IsEmailTaken(email string) (bool, error)
	ExistsByID(id uuid.UUID) (bool, error)
	FindByUsername(username string) (*domain.User, error)
	GetStats(userID uuid.UUID) (*domain.UserStats, error)
}
