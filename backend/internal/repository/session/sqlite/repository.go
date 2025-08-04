package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/domain"
	session2 "forum/internal/repository/session"
	"github.com/google/uuid"
	"time"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) session2.Repository {
	return &repository{db: db}
}

func (r *repository) Create(session *domain.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, session.ID, session.UserID, session.ExpiresAt)
	if err != nil {
		return errors.Join(session.ErrCreateSession, err)
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
			return nil, session.ErrSessionNotFound
		}
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, session.ErrSessionExpired
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
			return nil, session.ErrSessionNotFound
		}
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, session.ErrSessionExpired
	}

	return &session, nil
}

func (r *repository) Delete(token string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	res, err := r.db.Exec(query, token)
	if err != nil {
		return errors.Join(session2.ErrDeleteSession, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return session2.ErrSessionNotFound
	}

	return nil
}

func (r *repository) DeleteByUserID(userID uuid.UUID) error {
	const query = `DELETE FROM sessions WHERE user_id = ?` // SQLite

	res, err := r.db.Exec(query, userID.String())
	if err != nil {
		return fmt.Errorf("%w: %w", session2.ErrDeleteSession, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %w", session2.ErrAffectedRowsCheckFail, err)
	}

	return nil
}

func (r *repository) DeleteExpiredBefore(t time.Time) error {
	query := `DELETE FROM sessions WHERE expires_at < $1`
	_, err := r.db.Exec(query, t)
	return err
}
