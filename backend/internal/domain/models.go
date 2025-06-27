// domain/models.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Post struct {
	ID             uuid.UUID   `json:"id"`
	UserID         uuid.UUID   `json:"user_id"`
	ImageURL       string      `json:"image_url"`
	Title          string      `json:"title"`
	Content        string      `json:"content"`
	CreatedAt      time.Time   `json:"created_at"`
	AuthorUsername string      `json:"author_username"`
	Likes          int         `json:"likes"`
	Dislikes       int         `json:"dislikes"`
	CommentsCount  int         `json:"comments_count"`
	Categories     []*Category `json:"categories,omitempty"`
	Tags           []string    `json:"tags,omitempty"`
}

type Category struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type PostCategory struct {
	PostID     uuid.UUID `json:"post_id"`
	CategoryID uuid.UUID `json:"category_id"`
}

type Comment struct {
	ID             uuid.UUID `json:"id"`
	PostID         uuid.UUID `json:"post_id"`
	UserID         uuid.UUID `json:"user_id"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
	AuthorUsername string    `json:"author_username"`
	Likes          int       `json:"likes"`
	Dislikes       int       `json:"dislikes"`
}

type PostReaction struct {
	UserID   uuid.UUID `json:"user_id"`
	PostID   uuid.UUID `json:"post_id"`
	Reaction int       `json:"reaction"` // 1 = like, -1 = dislike
}

type CommentReaction struct {
	UserID    uuid.UUID `json:"user_id"`
	CommentID uuid.UUID `json:"comment_id"`
	Reaction  int       `json:"reaction"` // 1 = like, -1 = dislike
}
