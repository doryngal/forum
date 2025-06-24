package service

import (
	"forum/internal/repository"
	category_service "forum/internal/service/category"
	comment_service "forum/internal/service/comment"
	post_service "forum/internal/service/post"
	user_service "forum/internal/service/user"

	category_validator "forum/internal/service/category/validator"
	comment_validator "forum/internal/service/comment/validator"
	post_validator "forum/internal/service/post/validator"
	user_validator "forum/internal/service/user/validator"
)

type Service struct {
	Category category_service.Service
	Comment  comment_service.Service
	Post     post_service.Service
	User     user_service.Service
}

func NewServices(reps *repository.Repositories) *Service {
	return &Service{
		Category: category_service.New(reps.Category, reps.Post, category_validator.NewCategoryValidator()),
		Comment:  comment_service.New(reps.Comment, reps.Post, reps.User, comment_validator.NewCommentValidator()),
		Post:     post_service.New(reps.Post, reps.User, reps.Category, post_validator.NewPostValidator()),
		User:     user_service.New(reps.User, user_validator.NewUserValidator()),
	}
}
