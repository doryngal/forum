package session

import (
	"database/sql"
	"errors"
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

func (r *repository) Create(session *domain.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, session.ID, session.UserID, session.ExpiresAt)
	if err != nil {
		return errors.Join(ErrCreateSession, err)
	}
	return nil
}

func (r *repository) GetByToken(token string) (*domain.Session, error) {
	query := `SELECT id, user_id, expires_at FROM sessions WHERE id = $1`
	row := r.db.QueryRow(query, token)

	var session domain.Session
	err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, ErrSessionExpired
	}

	return &session, nil
}

func (r *repository) GetByUserID(userID uuid.UUID) (*domain.Session, error) {
	query := `SELECT id, user_id, expires_at FROM sessions WHERE user_id = $1`
	row := r.db.QueryRow(query, userID)

	var session domain.Session
	err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, ErrSessionExpired
	}

	return &session, nil
}

func (r *repository) Delete(token string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	res, err := r.db.Exec(query, token)
	if err != nil {
		return errors.Join(ErrDeleteSession, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSessionNotFound
	}

	return nil
}

func (r *repository) DeleteByUserID(userID uuid.UUID) error {
	const query = `DELETE FROM sessions WHERE user_id = ?` // SQLite

	res, err := r.db.Exec(query, userID.String())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDeleteSession, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAffectedRowsCheckFail, err)
	}

	return nil
}

func (r *repository) DeleteExpiredBefore(t time.Time) error {
	query := `DELETE FROM sessions WHERE expires_at < $1`
	_, err := r.db.Exec(query, t)
	return err
}
