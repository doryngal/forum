package post

import (
	"errors"
	"forum/internal/domain"
	category_repo "forum/internal/repository/sqlite/category"
	post_repo "forum/internal/repository/sqlite/post"
	user_repo "forum/internal/repository/sqlite/user"
	"forum/internal/service/category"
	"forum/internal/service/user"
	"github.com/google/uuid"
	"time"
)

type service struct {
	repo         post_repo.Repository
	userRepo     user_repo.Repository
	categoryRepo category_repo.Repository
	validator    Validator
}

type Validator interface {
	ValidatePost(post *domain.Post) error
}

func New(repo post_repo.Repository, userRepo user_repo.Repository,
	categoryRepo category_repo.Repository, validator Validator) Service {
	return &service{
		repo:         repo,
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
		validator:    validator,
	}
}

func (s *service) CreatePost(post *domain.Post) error {
	if err := s.validator.ValidatePost(post); err != nil {
		return err
	}

	exists, err := s.userRepo.ExistsByID(post.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return user.ErrUserNotFound
	}

	post.ID = uuid.New()
	post.CreatedAt = time.Now().UTC()

	return s.repo.Create(post)
}

func (s *service) GetPostByID(id uuid.UUID) (*domain.Post, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidPost
	}

	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, post_repo.ErrScanFailed) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return post, nil
}

func (s *service) GetAllPosts() ([]*domain.Post, error) {
	posts, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *service) GetPostsByCategory(categoryID uuid.UUID) ([]*domain.Post, error) {
	if categoryID == uuid.Nil {
		return nil, category.ErrInvalidCategoryID
	}

	exists, err := s.categoryRepo.ExistsByID(categoryID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, category.ErrInvalidCategoryID
	}

	return s.repo.GetByCategory(categoryID)
}

func (s *service) GetPostsByUser(userID uuid.UUID) ([]*domain.Post, error) {
	if userID == uuid.Nil {
		return nil, user.ErrUserNotFound
	}

	exists, err := s.userRepo.ExistsByID(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, user.ErrUserNotFound
	}

	return s.repo.GetByUserID(userID)
}

func (s *service) GetLikedPostsByUser(userID uuid.UUID) ([]*domain.Post, error) {
	if userID == uuid.Nil {
		return nil, user.ErrUserNotFound
	}

	exists, err := s.userRepo.ExistsByID(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, user.ErrUserNotFound
	}

	return s.repo.GetLikedByUser(userID)
}

func (s *service) LikePost(postID, userID uuid.UUID) error {
	if postID == uuid.Nil || userID == uuid.Nil {
		return ErrInvalidPost
	}

	if exists, err := s.repo.ExistsByID(postID); err != nil || !exists {
		return ErrPostNotFound
	}
	if exists, err := s.userRepo.ExistsByID(userID); err != nil || !exists {
		return user.ErrUserNotFound
	}

	reaction, err := s.repo.GetReaction(postID, userID)
	if err != nil && !errors.Is(err, post_repo.ErrReactionNotFound) {
		return err
	}
	if reaction == 1 {
		return ErrAlreadyReacted
	}

	return s.repo.Like(postID, userID)
}

func (s *service) DislikePost(postID, userID uuid.UUID) error {
	if postID == uuid.Nil || userID == uuid.Nil {
		return ErrInvalidPost
	}

	if exists, err := s.repo.ExistsByID(postID); err != nil || !exists {
		return ErrPostNotFound
	}
	if exists, err := s.userRepo.ExistsByID(userID); err != nil || !exists {
		return user.ErrUserNotFound
	}

	reaction, err := s.repo.GetReaction(postID, userID)
	if err != nil && !errors.Is(err, post_repo.ErrReactionNotFound) {
		return err
	}
	if reaction == -1 {
		return ErrAlreadyReacted
	}

	return s.repo.Dislike(postID, userID)
}
