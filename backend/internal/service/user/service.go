package user

import (
	"errors"
	"forum/internal/domain"
	user_repo "forum/internal/repository/sqlite/user"
	"forum/internal/service/user/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo      user_repo.Repository
	validator validator.UserValidator
}

func New(repo user_repo.Repository, validator validator.UserValidator) Service {
	return &service{
		repo:      repo,
		validator: validator,
	}
}

func (s *service) RegisterUser(user *domain.User) error {
	if err := s.validator.ValidateUser(user); err != nil {
		return err
	}

	taken, err := s.repo.IsEmailTaken(user.Email)
	if err != nil {
		return err
	}
	if taken {
		return ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *service) Login(email, password string) (*domain.User, error) {
	if err := s.validator.ValidateEmail(email); err != nil {
		return nil, ErrInvalidCredentials
	}

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, user_repo.ErrQueryFailed) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *service) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, user_repo.ErrQueryFailed) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *service) GetUserByID(id uuid.UUID) (*domain.User, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidUser
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, user_repo.ErrQueryFailed) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *service) IsEmailTaken(email string) (bool, error) {
	return s.repo.IsEmailTaken(email)
}

func (s *service) GetUserByUsername(username string) (*domain.User, error) {
	if username == "" {
		return nil, ErrInvalidUser
	}

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, user_repo.ErrQueryFailed) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *service) GetUserStats(userID uuid.UUID) (*domain.UserStats, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUser
	}

	stats, err := s.repo.GetStats(userID)
	if err != nil {
		if errors.Is(err, user_repo.ErrQueryFailed) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return stats, nil
}
