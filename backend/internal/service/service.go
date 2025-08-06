package service

import (
	"forum/internal/repository"
	category_service "forum/internal/service/category"
	category_validator "forum/internal/service/category/validator"
	comment_service "forum/internal/service/comment"
	comment_validator "forum/internal/service/comment/validator"
	post_service "forum/internal/service/post"
	post_validator "forum/internal/service/post/validator"
	session_service "forum/internal/service/session"
	session_validator "forum/internal/service/session/validator"
	user_service "forum/internal/service/user"
	user_validator "forum/internal/service/user/validator"
	"forum/pkg/logger"
)

type Service struct {
	Category category_service.Service
	Comment  comment_service.Service
	Session  session_service.Service
	Post     post_service.Service
	User     user_service.Service
}

func NewServices(reps *repository.Repositories, logger logger.Logger) *Service {
	return &Service{
		Category: category_service.New(reps.Category, reps.Post, category_validator.NewCategoryValidator()),
		Comment:  comment_service.New(reps.Comment, reps.Post, reps.User, comment_validator.NewCommentValidator()),
		Session:  session_service.New(reps.Session, session_validator.NewSessionValidator()),
		Post:     post_service.New(reps.Post, reps.User, reps.Category, post_validator.NewPostValidator()),
		User:     user_service.New(reps.User, user_validator.NewUserValidator(), logger),
	}
}
