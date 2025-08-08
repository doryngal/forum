package user

import (
	"errors"
	"forum/internal/domain"
	user_repo "forum/internal/repository/user"
	"forum/internal/service/user/validator"
	"forum/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo      user_repo.Repository
	validator validator.UserValidator
	log       logger.Logger
}

func New(repo user_repo.Repository, validator validator.UserValidator, log logger.Logger) Service {
	return &service{
		repo:      repo,
		validator: validator,
		log:       log,
	}
}

func (s *service) RegisterUser(user *domain.User) error {
	s.log.Info("RegisterUser started", logger.F("email", user.Email))

	if err := s.validator.ValidateUser(user); err != nil {
		s.log.Error("User validation failed", logger.F("error", err))
		return err
	}

	taken, err := s.repo.IsEmailTaken(user.Email)
	if err != nil {
		s.log.Error("Failed to check if email is taken", logger.F("error", err))
		return err
	}
	if taken {
		s.log.Info("Email is already taken", logger.F("email", user.Email))
		return ErrEmailTaken
	}

	taken, err = s.repo.IsUsernameTaken(user.Username)
	if err != nil {
		s.log.Error("Failed to check if username is taken", logger.F("error", err))
		return err
	}
	if taken {
		s.log.Info("Username is already taken", logger.F("username", user.Username))
		return ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("Password hashing failed", logger.F("error", err))
		return err
	}
	user.PasswordHash = string(hashedPassword)

	if err := s.repo.Create(user); err != nil {
		s.log.Error("User creation failed", logger.F("error", err))
		return err
	}

	s.log.Info("User registered successfully", logger.F("email", user.Email))
	return nil
}

func (s *service) Login(emailOrUsername, password string) (*domain.User, error) {
	s.log.Info("Login attempt", logger.F("email", emailOrUsername))

	user, err := s.repo.FindByEmailORUsername(emailOrUsername)
	if err != nil {
		if errors.Is(err, user_repo.ErrQueryFailed) {
			s.log.Info("User not found", logger.F("email/username", emailOrUsername))
			return nil, ErrInvalidCredentials
		}
		s.log.Error("Find user by email failed", logger.F("error", err))
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.log.Info("Invalid password", logger.F("email/username", emailOrUsername))
		return nil, ErrInvalidCredentials
	}

	s.log.Info("Login successful", logger.F("email/username", emailOrUsername))
	return user, nil
}

func (s *service) GetUserByEmail(email string) (*domain.User, error) {
	s.log.Debug("GetUserByEmail", logger.F("email", email))

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, user_repo.ErrQueryFailed) {
			s.log.Info("User not found by email", logger.F("email", email))
			return nil, ErrUserNotFound
		}
		s.log.Error("GetUserByEmail failed", logger.F("error", err))
		return nil, err
	}
	return user, nil
}

func (s *service) GetUserByID(id uuid.UUID) (*domain.User, error) {
	s.log.Debug("GetUserByID", logger.F("user_id", id))

	if id == uuid.Nil {
		s.log.Error("Invalid user ID", logger.F("user_id", id))
		return nil, ErrInvalidUser
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, user_repo.ErrQueryFailed) {
			s.log.Info("User not found by ID", logger.F("user_id", id))
			return nil, ErrUserNotFound
		}
		s.log.Error("GetUserByID failed", logger.F("error", err))
		return nil, err
	}
	return user, nil
}

func (s *service) IsEmailTaken(email string) (bool, error) {
	s.log.Debug("IsEmailTaken", logger.F("email", email))

	taken, err := s.repo.IsEmailTaken(email)
	if err != nil {
		s.log.Error("IsEmailTaken failed", logger.F("error", err))
	}
	return taken, err
}

func (s *service) GetUserByUsername(username string) (*domain.User, error) {
	s.log.Debug("GetUserByUsername", logger.F("username", username))

	if username == "" {
		s.log.Error("Username is empty")
		return nil, ErrInvalidUser
	}

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		s.log.Info("User not found by username", logger.F("username", username))
		return nil, ErrUserNotFound
	}
	return user, nil
}
