package comment

import (
	"errors"
	"forum/internal/domain"
	comment_repo "forum/internal/repository/sqlite/comment"
	post_repo "forum/internal/repository/sqlite/post"
	user_repo "forum/internal/repository/sqlite/user"
	"forum/internal/service/comment/validator"
	"forum/internal/service/post"
	"forum/internal/service/user"
	"github.com/google/uuid"
	"time"
)

type service struct {
	repo      comment_repo.Repository
	postRepo  post_repo.Repository
	userRepo  user_repo.Repository
	validator validator.CommentValidator
}

func New(repo comment_repo.Repository, postRepo post_repo.Repository,
	userRepo user_repo.Repository, validator validator.CommentValidator) Service {
	return &service{
		repo:      repo,
		postRepo:  postRepo,
		userRepo:  userRepo,
		validator: validator,
	}
}

func (s *service) CreateComment(comment *domain.Comment) error {
	if err := s.validator.ValidateComment(comment); err != nil {
		return err
	}

	if exists, err := s.postRepo.ExistsByID(comment.PostID); err != nil || !exists {
		return post.ErrPostNotFound
	}
	if exists, err := s.userRepo.ExistsByID(comment.UserID); err != nil || !exists {
		return user.ErrUserNotFound
	}

	comment.ID = uuid.New()
	comment.CreatedAt = time.Now().UTC()

	return s.repo.Create(comment)
}

func (s *service) GetCommentsByPost(postID uuid.UUID) ([]*domain.Comment, error) {
	if postID == uuid.Nil {
		return nil, post.ErrInvalidPost
	}

	if exists, err := s.postRepo.ExistsByID(postID); err != nil || !exists {
		return nil, post.ErrPostNotFound
	}

	return s.repo.GetByPostID(postID)
}

func (s *service) LikeComment(commentID, userID uuid.UUID) error {
	if commentID == uuid.Nil || userID == uuid.Nil {
		return ErrInvalidComment
	}

	if exists, err := s.repo.ExistsByID(commentID); err != nil || !exists {
		return ErrCommentNotFound
	}
	if exists, err := s.userRepo.ExistsByID(userID); err != nil || !exists {
		return user.ErrUserNotFound
	}

	return s.repo.Like(commentID, userID)
}

func (s *service) DislikeComment(commentID, userID uuid.UUID) error {
	if commentID == uuid.Nil || userID == uuid.Nil {
		return ErrInvalidComment
	}

	if exists, err := s.repo.ExistsByID(commentID); err != nil || !exists {
		return ErrCommentNotFound
	}
	if exists, err := s.userRepo.ExistsByID(userID); err != nil || !exists {
		return user.ErrUserNotFound
	}

	reaction, err := s.repo.GetReaction(commentID, userID)
	if err != nil && !errors.Is(err, comment_repo.ErrReactionNotFound) {
		return err
	}
	if reaction == -1 {
		return ErrAlreadyReacted
	}

	return s.repo.Dislike(commentID, userID)
}
