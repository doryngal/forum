package category

import (
	"forum/internal/domain"
	category_repo "forum/internal/repository/category"
	post_repo "forum/internal/repository/post"
	"forum/internal/service/category/validator"
	"forum/internal/service/post"
	"github.com/google/uuid"
)

type service struct {
	repo      category_repo.Repository
	postRepo  post_repo.Repository
	validator validator.CategoryValidator
}

func New(repo category_repo.Repository, postRepo post_repo.Repository, validator validator.CategoryValidator) Service {
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

	if exists, err := s.postRepo.ExistsByID(postID); err != nil || !exists {
		return post.ErrPostNotFound
	}

	if len(categoryIDs) == 0 {
		return ErrNoCategories
	}

	categoryIDs = removeDuplicateUUIDs(categoryIDs)

	for _, catID := range categoryIDs {
		if catID == uuid.Nil {
			return ErrInvalidCategoryID
		}
		exists, err := s.repo.ExistsByID(catID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrNotFound
		}
	}

	return s.repo.AssignToPost(postID, categoryIDs)
}

func removeDuplicateUUIDs(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{})
	result := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}
	return result
}
