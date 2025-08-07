package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/domain"
	session_repo "forum/internal/repository/session"
	"github.com/google/uuid"
	"time"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) session_repo.Repository {
	return &repository{db: db}
}

func (r *repository) Create(session *domain.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, session.ID, session.UserID.String(), session.ExpiresAt)
	if err != nil {
		return errors.Join(session_repo.ErrCreateSession, err)
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
			return nil, session_repo.ErrSessionNotFound
		}
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, session_repo.ErrSessionExpired
	}

	return &session, nil
}

func (r *repository) GetByUserID(userID uuid.UUID) (*domain.Session, error) {
	query := `SELECT id, user_id, expires_at FROM sessions WHERE user_id = $1`
	row := r.db.QueryRow(query, userID.String())

	var session domain.Session
	err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, session_repo.ErrSessionNotFound
		}
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, session_repo.ErrSessionExpired
	}

	return &session, nil
}

func (r *repository) Delete(token string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	res, err := r.db.Exec(query, token)
	if err != nil {
		return errors.Join(session_repo.ErrDeleteSession, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return session_repo.ErrSessionNotFound
	}

	return nil
}

func (r *repository) DeleteByUserID(userID uuid.UUID) error {
	const query = `DELETE FROM sessions WHERE user_id = ?` // SQLite

	res, err := r.db.Exec(query, userID.String())
	if err != nil {
		return fmt.Errorf("%w: %w", session_repo.ErrDeleteSession, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %w", session_repo.ErrAffectedRowsCheckFail, err)
	}

	return nil
}

func (r *repository) DeleteExpiredBefore(t time.Time) error {
	query := `DELETE FROM sessions WHERE expires_at < $1`
	_, err := r.db.Exec(query, t)
	return err
}
