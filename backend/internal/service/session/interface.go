package session

import (
	"forum/internal/domain"
	"github.com/google/uuid"
)

type Service interface {
	Create(userID uuid.UUID) (*domain.Session, error)
	GetByToken(token string) (*domain.Session, error)
	GetByUserID(userID uuid.UUID) (*domain.Session, error)
	Delete(token string) error
	CleanupExpired() error
}
