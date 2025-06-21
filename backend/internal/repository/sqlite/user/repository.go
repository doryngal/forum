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
