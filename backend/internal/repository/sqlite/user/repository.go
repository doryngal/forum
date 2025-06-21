package user

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

func (r repository) Create(user *domain.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	_, err := r.db.Exec(
		"INSERT INTO users (id, email, username, password_hash, created_at) VALUES (?, ?, ?, ?, ?)",
		user.ID.String(), user.Email, user.Username, user.PasswordHash, user.CreatedAt,
	)
	return err
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
		return nil, err
	}
	u.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	u.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r repository) IsEmailTaken(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		email,
	).Scan(&exists)
	return exists, err
}
