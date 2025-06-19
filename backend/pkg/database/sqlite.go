package database

import (
	"database/sql"
	"forum/internal/config"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(cfg config.DatabaseConfig) (*SQLite, error) {
	db, err := sql.Open(cfg.Driver, cfg.Path)
	if err != nil {
		return nil, err
	}
	return &SQLite{db: db}, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) GetDB() *sql.DB {
	return s.db
}
