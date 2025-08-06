package router

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/presentation/html"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/pkg/database"
	"forum/pkg/logger"
	"forum/pkg/templates"
	"net/http"
)

func NewApp(cfg config.Config) (http.Handler, error) {
	log := logger.New()

	db, err := database.NewSQLite(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	tmpl, err := templates.Parse(cfg.Server.TemplatePath)
	if err != nil {
		return nil, fmt.Errorf("init templates: %w", err)
	}

	repos := repository.NewRepositories(db.GetDB())
	services := service.NewServices(repos, log)
	handlers := html.NewTemplateHandlers(tmpl, services, log)

	router := InitRouter(handlers, cfg.Server.StaticPath, log)

	return router, nil
}
