package html

import (
	"forum/internal/presentation/html/home"
	"forum/internal/service"
)

type TemplateHandlers struct {
	Home *home.HomeHandler
}

func NewTemplateHandlers(services *service.Service) *TemplateHandlers {
	return &TemplateHandlers{
		Home: home.NewHomeHandler(services.Post, services.User, services.Category),
	}
}
