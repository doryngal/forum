package category

import (
	"errors"
	"forum/internal/domain"
	category_repo "forum/internal/repository/sqlite/category"
	post_repo "forum/internal/repository/sqlite/post"
	"github.com/google/uuid"
)

type service struct {
	repo      category_repo.Repository
	postRepo  post_repo.Repository
	validator Validator
}

type Validator interface {
	ValidateCategoryName(name string) error
}

func New(repo category_repo.Repository, postRepo post_repo.Repository, validator Validator) Service {
	return &service{
		repo:      repo,
		postRepo:  postRepo,
		validator: validator,
	}
}

func (s *service) CreateCategory(name string) error {
	if err := s.validator.ValidateCategoryName(name); err != nil {
		return err
	}

	exists, err := s.repo.ExistsByName(name)
	if err != nil {
		return err
	}
	if exists {
		return ErrCategoryExists
	}

	return s.repo.Create(name)
}

func (s *service) GetAllCategories() ([]*domain.Category, error) {
	categories, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, ErrNoCategories
	}
	return categories, nil
}

func (s *service) AssignCategoriesToPost(postID uuid.UUID, categoryIDs []uuid.UUID) error {
	if postID == uuid.Nil {
		return ErrInvalidPostID
	}

	if len(categoryIDs) == 0 {
		return ErrNoCategories
	}

	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrInvalidPostID
		}
		return err
	}

	for _, catID := range categoryIDs {
		if catID == uuid.Nil {
			return ErrInvalidCategoryID
		}
		exists, err := s.repo.ExistsByID(catID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrInvalidCategoryID
		}
	}

	return s.repo.AssignToPost(postID, categoryIDs)
}
