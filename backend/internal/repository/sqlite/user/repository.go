package user

import (
	"database/sql"
	"fmt"
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

func (r repository) Create(user *domain.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()

	_, err := r.db.Exec(
		"INSERT INTO users (id, email, username, password_hash, created_at) VALUES (?, ?, ?, ?, ?)",
		user.ID.String(), user.Email, user.Username, user.PasswordHash, user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInsertUserFailed, err)
	}
	return nil
}

func (r repository) FindByEmailORUsername(emailOrUsername string) (*domain.User, error) {
	var u domain.User
	var idStr string

	err := r.db.QueryRow(
		`SELECT id, email, username, password_hash, created_at 
		 FROM users 
		 WHERE email = ? OR username = ?`,
		emailOrUsername, emailOrUsername,
	).Scan(&idStr, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", ErrQueryFailed, err)
	}

	u.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUUIDParseFailed, err)
	}

	return &u, nil
}

func (r repository) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	var idStr string

	err := r.db.QueryRow(
		"SELECT id, email, username, password_hash, created_at FROM users WHERE email = ?",
		email,
	).Scan(&idStr, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", ErrQueryFailed, err)
	}

	u.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUUIDParseFailed, err)
	}

	return &u, nil
}

func (r repository) FindByID(id uuid.UUID) (*domain.User, error) {
	var u domain.User
	var idStr string

	err := r.db.QueryRow(
		"SELECT id, email, username, password_hash, created_at FROM users WHERE id = ?",
		id.String(),
	).Scan(&idStr, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", ErrQueryFailed, err)
	}

	u.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUUIDParseFailed, err)
	}

	return &u, nil
}

func (r repository) IsEmailTaken(email string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		email,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrCheckExistsFailed, err)
	}

	return exists, nil
}

func (r *repository) ExistsByID(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)",
		id.String(),
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrQueryFailed, err)
	}
	return exists, nil
}

func (r *repository) FindByUsername(username string) (*domain.User, error) {
	var u domain.User
	var idStr string

	err := r.db.QueryRow(
		"SELECT id, email, username, password_hash, created_at FROM users WHERE username = ?",
		username,
	).Scan(&idStr, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrQueryFailed, err)
	}

	u.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUUIDParseFailed, err)
	}

	return &u, nil
}

func (r *repository) GetStats(userID uuid.UUID) (*domain.UserStats, error) {
	var stats domain.UserStats

	err := r.db.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM posts WHERE user_id = ?),
			(SELECT COUNT(*) FROM comments WHERE user_id = ?),
			(SELECT COUNT(*) FROM post_reactions pr 
			 JOIN posts p ON pr.post_id = p.id 
			 WHERE p.user_id = ? AND pr.reaction = 1),
			(SELECT COUNT(*) FROM post_reactions pr 
			 JOIN posts p ON pr.post_id = p.id 
			 WHERE p.user_id = ? AND pr.reaction = -1)
	`, userID.String(), userID.String(), userID.String(), userID.String(),
	).Scan(&stats.PostCount, &stats.CommentCount, &stats.LikeCount, &stats.DislikeCount)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQueryFailed, err)
	}

	return &stats, nil
}
