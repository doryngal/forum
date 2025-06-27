package validator

import "forum/internal/domain"

type SessionValidator interface {
	ValidateCreate(session *domain.Session) error
	ValidateToken(token string) error
}
