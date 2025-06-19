package comment

import "forum/internal/domain"

type Repository interface {
	Create(comment *domain.Comment) error
	GetByPostID(postID int) ([]*domain.Comment, error)
	Like(commentID, userID int) error
	Dislike(commentID, userID int) error
}
