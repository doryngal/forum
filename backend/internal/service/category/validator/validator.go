package validator

import (
	"errors"
	"strings"
)

var (
	ErrEmptyCategoryName = errors.New("validator: category name cannot be empty")
	ErrTooShortName      = errors.New("validator: category name is too short")
	ErrTooLongName       = errors.New("validator: category name is too long")
)

type categoryValidator struct{}

func NewCategoryValidator() CategoryValidator {
	return &categoryValidator{}
}

func (v *categoryValidator) ValidateCategoryName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrEmptyCategoryName
	}
	if len(name) < 2 {
		return ErrTooShortName
	}
	if len(name) > 50 {
		return ErrTooLongName
	}
	return nil
}
