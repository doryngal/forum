package html

import (
	"forum/internal/presentation/html/auth"
	"forum/internal/presentation/html/home"
	"forum/internal/presentation/html/post"
	profile "forum/internal/presentation/html/user"
	"forum/internal/service"
	"html/template"
)

type TemplateHandlers struct {
	Home     *home.HomeHandler
	Login    *auth.LoginHandler
	Register *auth.RegisterHandler

	Post       *post.PostHandler
	CreatePost *post.CreateHandler

	Profile *profile.ProfileHandler
}

func NewTemplateHandlers(tmpl *template.Template, services *service.Service) *TemplateHandlers {
	return &TemplateHandlers{
		Home:     home.NewHomeHandler(tmpl, services.Post, services.User, services.Category, services.Session),
		Login:    auth.NewLoginHandler(tmpl, services.User, services.Session),
		Register: auth.NewRegisterHandler(tmpl, services.User),

		Post:       post.NewPostHandler(tmpl, services.User, services.Post, services.Comment),
		CreatePost: post.NewCreateHandler(tmpl, services.User, services.Post),

		Profile: profile.NewProfileHandler(tmpl, services.User, services.Post),
	}
}
