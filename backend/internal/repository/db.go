package repository

import (
	"database/sql"
	category_repo "forum/internal/repository/sqlite/category"
	comment_repo "forum/internal/repository/sqlite/comment"
	post_repo "forum/internal/repository/sqlite/post"
	session_repo "forum/internal/repository/sqlite/session"
	user_repo "forum/internal/repository/sqlite/user"
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
		Category: category_repo.New(db),
		Comment:  comment_repo.New(db),
		Session:  session_repo.New(db),
		Post:     post_repo.New(db),
		User:     user_repo.New(db),
	}
}
