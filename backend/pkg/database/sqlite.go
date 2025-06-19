package database

import "database/sql"

type SQLite struct {
	db *sql.DB
}

func NewSQLite() (*SQLite, error) {
	db, err := sql.Open("sqlite3", "./store.db")
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
