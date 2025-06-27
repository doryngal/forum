package router

import (
	"forum/internal/config"
	"forum/internal/presentation/html"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/pkg/database"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func NewApp(cfg config.Config) (http.Handler, error) {
	// init database
	sqlite, err := database.NewSQLite(cfg.Database)
	if err != nil {
		return nil, err
	}

	tmpl, err := parseTemplates()
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// init repositories, services, handlers
	repos := repository.NewRepositories(sqlite.GetDB())
	services := service.NewServices(repos)
	templateHandlers := html.NewTemplateHandlers(tmpl, services)

	mux := http.NewServeMux()

	// HTML routes
	registerHTMLRoutes(mux, templateHandlers)

	// registerAPIRoutes(mux, services)

	// static files
	fs := http.FileServer(http.Dir("./../frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux, nil
}

func registerHTMLRoutes(mux *http.ServeMux, h *html.TemplateHandlers) {
	mux.HandleFunc("/", h.Home.ServeHTTP)

	mux.HandleFunc("/login", h.Login.ServeHTTP)
	mux.HandleFunc("/logout", h.Login.Logout)
	mux.HandleFunc("/register", h.Register.ServeHTTP)

	mux.HandleFunc("/create-post", h.CreatePost.ServeHTTP)
	mux.HandleFunc("/post/", h.Post.ServeHTTP)

	mux.HandleFunc("/profile/", h.Profile.ServeHTTP)
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func parseTemplates() (*template.Template, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"truncate": truncate,
	})

	err := filepath.Walk("./../frontend/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err := tmpl.ParseFiles(path)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
