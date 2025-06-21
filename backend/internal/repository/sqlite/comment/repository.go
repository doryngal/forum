package comment

import (
	"database/sql"
	"forum/internal/domain"
	"github.com/google/uuid"
	"time"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r repository) Create(comment *domain.Comment) error {
	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()
	_, err := r.db.Exec(
		"INSERT INTO comments (id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
		comment.ID.String(), comment.PostID.String(), comment.UserID.String(), comment.Content, comment.CreatedAt,
	)
	return err
}

func (r repository) GetByPostID(postID uuid.UUID) ([]*domain.Comment, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, 
		       u.username,
		       COALESCE(SUM(CASE WHEN cr.reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN cr.reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id
		LEFT JOIN comment_reactions cr ON c.id = cr.comment_id
		WHERE c.post_id = ?
		GROUP BY c.id
		ORDER BY c.created_at DESC`, postID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		var idStr, postIDStr, userIDStr string
		if err := rows.Scan(&idStr, &postIDStr, &userIDStr, &c.Content,
			&c.CreatedAt, &c.AuthorUsername, &c.Likes, &c.Dislikes); err != nil {
			return nil, err
		}
		c.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		c.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, err
		}
		c.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (r repository) Like(commentID, userID uuid.UUID) error {
	return r.setReaction(commentID, userID, 1)
}

func (r repository) Dislike(commentID, userID uuid.UUID) error {
	return r.setReaction(commentID, userID, -1)
}

func (r repository) setReaction(commentID, userID uuid.UUID, reaction int) error {
	_, err := r.db.Exec(`
		INSERT INTO comment_reactions (user_id, comment_id, reaction) 
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, comment_id) DO UPDATE SET reaction = excluded.reaction`,
		userID.String(), commentID.String(), reaction)
	return err
}
