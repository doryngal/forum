package html

import (
	"forum/internal/presentation/html/auth"
	"forum/internal/presentation/html/errorhandler"
	"forum/internal/presentation/html/home"
	"forum/internal/presentation/html/post"
	profile "forum/internal/presentation/html/user"
	"forum/internal/service"
	"forum/pkg/logger"
	"html/template"
)

type TemplateHandlers struct {
	Home     *home.Handler
	Login    *auth.LoginHandler
	Register *auth.RegisterHandler

	Post       *post.PostHandler
	CreatePost *post.CreateHandler
	EditPost   *post.EditHandler
	DeletePost *post.DeleteHandler

	Profile *profile.ProfileHandler
}

func NewTemplateHandlers(tmpl *template.Template, services *service.Service, logger logger.Logger) *TemplateHandlers {
	errorHandler := *errorhandler.NewErrorHandler(tmpl, logger)
	return &TemplateHandlers{
		Home:       home.NewHomeHandler(tmpl, services.Post, services.User, services.Category, services.Session, errorHandler),
		Login:      auth.NewLoginHandler(tmpl, services.User, services.Session, errorHandler),
		Register:   auth.NewRegisterHandler(tmpl, services.User, errorHandler),
		Post:       post.NewPostHandler(tmpl, services.User, services.Post, services.Comment, services.Session, services.Category, errorHandler),
		CreatePost: post.NewCreateHandler(tmpl, services.User, services.Post, services.Session, services.Category, errorHandler),
		EditPost:   post.NewEditHandler(tmpl, services.Post, services.User, services.Session, services.Category, errorHandler),
		DeletePost: post.NewDeleteHandler(services.Post, services.Session, services.User, errorHandler),
		Profile:    profile.NewProfileHandler(tmpl, services.User, services.Post, services.Comment, services.Session, errorHandler),
	}
}
