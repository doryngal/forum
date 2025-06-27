package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"forum/internal/domain"
	session_repo "forum/internal/repository/sqlite/session"
	"forum/internal/service/session/validator"
	"github.com/google/uuid"
	"time"
)

type service struct {
	repo      session_repo.Repository
	validator validator.SessionValidator
	duration  time.Duration
}

func New(repo session_repo.Repository, validator validator.SessionValidator) Service {
	return &service{
		repo:      repo,
		validator: validator,
		duration:  time.Hour * 24 * 7,
	}
}

func (s *service) Create(userID uuid.UUID) (*domain.Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, errors.Join(ErrSessionCreationFailed, err)
	}

	session := &domain.Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.duration),
	}

	if err := s.validator.ValidateCreate(session); err != nil {
		return nil, errors.Join(ErrInvalidSession, err)
	}

	if err := s.repo.Create(session); err != nil {
		return nil, errors.Join(ErrSessionCreationFailed, err)
	}

	return session, nil
}

func (s *service) GetByToken(token string) (*domain.Session, error) {
	if err := s.validator.ValidateToken(token); err != nil {
		return nil, errors.Join(ErrInvalidSession, err)
	}

	session, err := s.repo.GetByToken(token)
	if err != nil {
		if errors.Is(err, session_repo.ErrSessionExpired) {
			return nil, errors.Join(ErrSessionExpired, err)
		}
		return nil, errors.Join(ErrSessionNotFound, err)
	}

	return session, nil
}

func (s *service) GetByUserID(userID uuid.UUID) (*domain.Session, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidSession
	}

	session, err := s.repo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, session_repo.ErrSessionExpired) {
			return nil, errors.Join(ErrSessionExpired, err)
		}
		return nil, errors.Join(ErrSessionNotFound, err)
	}

	return session, nil
}

func (s *service) Delete(token string) error {
	if err := s.validator.ValidateToken(token); err != nil {
		return errors.Join(ErrInvalidSession, err)
	}

	if err := s.repo.Delete(token); err != nil {
		if errors.Is(err, session_repo.ErrSessionNotFound) {
			return errors.Join(ErrSessionNotFound, err)
		}
		return err
	}

	return nil
}

func (s *service) DeleteByUserID(userID uuid.UUID) error {
	if err := s.repo.DeleteByUserID(userID); err != nil {
		if errors.Is(err, session_repo.ErrSessionNotFound) {
			return errors.Join(ErrSessionNotFound, err)
		}
		return err
	}

	return nil
}

func (s *service) CleanupExpired() error {
	return s.repo.DeleteExpiredBefore(time.Now())
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
