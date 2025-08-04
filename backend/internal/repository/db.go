package repository

import (
	"database/sql"
	category_repo "forum/internal/repository/category"
	category_sqlite "forum/internal/repository/category/sqlite"
	comment_repo "forum/internal/repository/comment"
	comment_sqlite "forum/internal/repository/comment/sqlite"
	post_repo "forum/internal/repository/post"
	post_sqlite "forum/internal/repository/post/sqlite"
	session_repo "forum/internal/repository/session"
	session_sqlite "forum/internal/repository/session/sqlite"
	user_repo "forum/internal/repository/user"
	user_sqlite "forum/internal/repository/user/sqlite"
)

type Repositories struct {
	Category category_repo.Repository
	Comment  comment_repo.Repository
	Session  session_repo.Repository
	Post     post_repo.Repository
	User     user_repo.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Category: category_sqlite.New(db),
		Comment:  comment_sqlite.New(db),
		Session:  session_sqlite.New(db),
		Post:     post_sqlite.New(db),
		User:     user_sqlite.New(db),
	}
}
