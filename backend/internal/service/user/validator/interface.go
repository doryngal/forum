package validator

import "forum/internal/domain"

type UserValidator interface {
	ValidateUser(user *domain.User) error
	ValidatePassword(password string) error
}
