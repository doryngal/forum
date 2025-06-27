package html

import (
	"forum/internal/presentation/html/auth"
	"forum/internal/presentation/html/home"
	"forum/internal/presentation/html/post"
	"forum/internal/service"
	"html/template"
)

type TemplateHandlers struct {
	Home       *home.HomeHandler
	Login      *auth.LoginHandler
	Register   *auth.RegisterHandler
	Post       *post.PostHandler
	CreatePost *post.CreateHandler
}

func NewTemplateHandlers(tmpl *template.Template, services *service.Service) *TemplateHandlers {
	return &TemplateHandlers{
		Home:       home.NewHomeHandler(tmpl, services.Post, services.User, services.Category),
		Login:      auth.NewLoginHandler(tmpl, services.User),
		Register:   auth.NewRegisterHandler(tmpl, services.User),
		Post:       post.NewPostHandler(tmpl, services.User, services.Post, services.Comment),
		CreatePost: post.NewCreateHandler(tmpl, services.User, services.Post),
	}
}
