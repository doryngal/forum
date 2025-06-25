package router

import (
	"forum/internal/config"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/pkg/database"
	"net/http"
)

func NewApp(cfg config.Config) (http.Handler, error) {
	sqlite, err := database.NewSQLite(cfg.Database)
	if err != nil {
		return nil, err
	}
	defer sqlite.Close()

	reps := repository.NewRepositories(sqlite.GetDB())
	services := service.NewServices(reps)
	panic(services)
}
