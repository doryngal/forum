package router

import (
	"forum/internal/presentation/html"
	"forum/pkg/logger"
	"net/http"
)

func InitRouter(h *html.TemplateHandlers, staticPath string, log logger.Logger) http.Handler {
	mux := http.NewServeMux()

	// HTML routes
	mux.HandleFunc("/", h.Home.ServeHTTP)
	mux.HandleFunc("/login", h.Login.ServeHTTP)
	mux.HandleFunc("/logout", h.Login.Logout)
	mux.HandleFunc("/register", h.Register.ServeHTTP)
	mux.HandleFunc("/create-post", h.CreatePost.ServeHTTP)
	mux.HandleFunc("/post/", h.Post.ServeHTTP)
	mux.HandleFunc("/profile", h.Profile.ServeHTTP)
	mux.HandleFunc("/profile/", h.Profile.ServeHTTP)

	mux.HandleFunc("/edit-post/", h.EditPost.ServeHTTP)
	mux.HandleFunc("/delete-post/", h.DeletePost.ServeHTTP)

	// Static files
	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
