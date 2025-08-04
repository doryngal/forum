package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/domain"
	comment_repo "forum/internal/repository/comment"
	"github.com/google/uuid"
	"time"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) comment_repo.Repository {
	return &repository{db: db}
}

func (r repository) Create(comment *domain.Comment) error {
	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()

	_, err := r.db.Exec(
		"INSERT INTO comments (id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
		comment.ID.String(), comment.PostID.String(), comment.UserID.String(), comment.Content, comment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", comment_repo.ErrInsertCommentFailed, err)
	}
	return nil
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
		return nil, fmt.Errorf("%w: %v", comment_repo.ErrQueryFailed, err)
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		var idStr, postIDStr, userIDStr string
		if err := rows.Scan(&idStr, &postIDStr, &userIDStr, &c.Content,
			&c.CreatedAt, &c.AuthorUsername, &c.Likes, &c.Dislikes); err != nil {
			return nil, fmt.Errorf("%w: %v", comment_repo.ErrScanFailed, err)
		}
		c.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", comment_repo.ErrUUIDParseFailed, err)
		}
		c.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", comment_repo.ErrUUIDParseFailed, err)
		}
		c.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", comment_repo.ErrUUIDParseFailed, err)
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
	var existingReaction int
	err := r.db.QueryRow(`
		SELECT reaction FROM comment_reactions 
		WHERE user_id = ? AND comment_id = ?`,
		userID.String(), commentID.String(),
	).Scan(&existingReaction)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("%w: %v", comment_repo.ErrReactionUpdateFailed, err)
	}

	if err == nil {
		if existingReaction == reaction {
			_, err := r.db.Exec(`
				DELETE FROM comment_reactions 
				WHERE user_id = ? AND comment_id = ?`,
				userID.String(), commentID.String())
			if err != nil {
				return fmt.Errorf("%w: %v", comment_repo.ErrReactionUpdateFailed, err)
			}
			return nil
		}
		_, err := r.db.Exec(`
			UPDATE comment_reactions 
			SET reaction = ? 
			WHERE user_id = ? AND comment_id = ?`,
			reaction, userID.String(), commentID.String())
		if err != nil {
			return fmt.Errorf("%w: %v", comment_repo.ErrReactionUpdateFailed, err)
		}
		return nil
	}

	_, err = r.db.Exec(`
		INSERT INTO comment_reactions (user_id, comment_id, reaction) 
		VALUES (?, ?, ?)`,
		userID.String(), commentID.String(), reaction)
	if err != nil {
		return fmt.Errorf("%w: %v", comment_repo.ErrReactionUpdateFailed, err)
	}
	return nil
}

func (r *repository) ExistsByID(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM comments WHERE id = ?)",
		id.String()).Scan(&exists)
	return exists, err
}

func (r *repository) GetReaction(commentID, userID uuid.UUID) (int, error) {
	var reaction int
	err := r.db.QueryRow(
		"SELECT reaction FROM comment_reactions WHERE comment_id = ? AND user_id = ?",
		commentID.String(), userID.String(),
	).Scan(&reaction)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, comment_repo.ErrReactionNotFound
		}
		return 0, fmt.Errorf("%w: %v", comment_repo.ErrGetReactionFailed, err)
	}

	return reaction, nil
}
