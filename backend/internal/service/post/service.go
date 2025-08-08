package post

import (
	"errors"
	"fmt"
	"forum/internal/domain"
	category_repo "forum/internal/repository/category"
	post_repo "forum/internal/repository/post"
	user_repo "forum/internal/repository/user"
	"forum/internal/service/post/validator"
	"forum/internal/service/user"
	"github.com/google/uuid"
	"time"
)

type service struct {
	repo         post_repo.Repository
	userRepo     user_repo.Repository
	categoryRepo category_repo.Repository
	validator    validator.PostValidator
}

func New(repo post_repo.Repository, userRepo user_repo.Repository,
	categoryRepo category_repo.Repository, validator validator.PostValidator) Service {
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

	if err := s.ensureUserExists(post.UserID); err != nil {
		return err
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

func (s *service) GetPostsByUserID(userID, sessionID uuid.UUID) ([]*domain.Post, error) {
	if err := s.ensureUserExists(userID); err != nil {
		return nil, err
	}

	return s.repo.GetByUserID(userID, sessionID)
}

func (s *service) LikePost(postID, userID uuid.UUID) error {
	if err := s.ensurePostExists(postID); err != nil {
		return err
	}
	if err := s.ensureUserExists(userID); err != nil {
		return err
	}

	_, err := s.repo.GetReaction(postID, userID)
	if err != nil && !errors.Is(err, post_repo.ErrReactionNotFound) {
		return err
	}
	return s.repo.Like(postID, userID)
}

func (s *service) DislikePost(postID, userID uuid.UUID) error {
	if err := s.ensurePostExists(postID); err != nil {
		return err
	}
	if err := s.ensureUserExists(userID); err != nil {
		return err
	}

	_, err := s.repo.GetReaction(postID, userID)
	if err != nil && !errors.Is(err, post_repo.ErrReactionNotFound) {
		return err
	}

	return s.repo.Dislike(postID, userID)
}

func (s *service) UpdatePost(post *domain.Post, userID uuid.UUID) error {
	_, err := s.authorizePostAccess(post.ID, userID)
	if err != nil {
		return err
	}
	return s.repo.Update(post)
}

func (s *service) DeletePost(postID, userID uuid.UUID) error {
	_, err := s.authorizePostAccess(postID, userID)
	if err != nil {
		return err
	}
	return s.repo.Delete(postID, userID)
}

func (s *service) GetLikedPosts(userID uuid.UUID) ([]*domain.Post, error) {
	if err := s.ensureUserExists(userID); err != nil {
		return nil, err
	}
	return s.repo.GetLikedPostsByUserID(userID)
}

func (s *service) GetDislikedPosts(userID uuid.UUID) ([]*domain.Post, error) {
	if err := s.ensureUserExists(userID); err != nil {
		return nil, err
	}

	return s.repo.GetDislikedPostsByUserID(userID)
}

func (s *service) ensureUserExists(userID uuid.UUID) error {
	if userID == uuid.Nil {
		return user.ErrUserNotFound
	}
	exists, err := s.userRepo.ExistsByID(userID)
	if err != nil {
		return fmt.Errorf("userRepo.ExistsByID: %w", err)
	}
	if !exists {
		return user.ErrUserNotFound
	}
	return nil
}

func (s *service) ensurePostExists(postID uuid.UUID) error {
	if postID == uuid.Nil {
		return ErrInvalidPost
	}
	exists, err := s.repo.ExistsByID(postID)
	if err != nil {
		return fmt.Errorf("postRepo.ExistsByID: %w", err)
	}
	if !exists {
		return ErrPostNotFound
	}
	return nil
}

func (s *service) authorizePostAccess(postID, userID uuid.UUID) (*domain.Post, error) {
	if err := s.ensurePostExists(postID); err != nil {
		return nil, err
	}
	if err := s.ensureUserExists(userID); err != nil {
		return nil, err
	}

	existingPost, err := s.repo.GetByID(postID)
	if err != nil {
		return nil, err
	}
	if existingPost.UserID != userID {
		return nil, ErrUnauthorized
	}

	return existingPost, nil
}
