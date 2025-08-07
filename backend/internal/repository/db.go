package repository

import (
	"database/sql"
	category_repo "forum/internal/repository/category"
	sqlite_category "forum/internal/repository/category/sqlite"
	comment_repo "forum/internal/repository/comment"
	sqlite_comment "forum/internal/repository/comment/sqlite"
	post_repo "forum/internal/repository/post"
	sqlite_post "forum/internal/repository/post/sqlite"
	session_repo "forum/internal/repository/session"
	sqlite_session "forum/internal/repository/session/sqlite"
	user_repo "forum/internal/repository/user"
	sqlite_user "forum/internal/repository/user/sqlite"
)

type Repositories struct {
	Category category_repo.Repository
	Comment  comment_repo.Repository
	Session  session_repo.Repository
	Post     post_repo.Repository
	User     user_repo.Repository
}

func NewRepositories(db *sql.DB, provider string) *Repositories {
	switch provider {
	case "sqlite":
		return &Repositories{
			Category: sqlite_category.New(db),
			Comment:  sqlite_comment.New(db),
			Session:  sqlite_session.New(db),
			Post:     sqlite_post.New(db),
			User:     sqlite_user.New(db),
		}
	default:
		panic("unknown provider")
	}
}
