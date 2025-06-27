package session

import (
	"forum/internal/domain"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	Create(session *domain.Session) error
	GetByToken(token string) (*domain.Session, error)
	GetByUserID(userID uuid.UUID) (*domain.Session, error)
	Delete(token string) error
	DeleteExpiredBefore(t time.Time) error
}
