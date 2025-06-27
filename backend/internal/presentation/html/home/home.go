package home

import (
	"forum/internal/domain"
	"forum/internal/service/category"
	"forum/internal/service/post"
	"forum/internal/service/user"
	"html/template"
	"net/http"
)

type HomeHandler struct {
	tmpl            *template.Template
	postService     post.Service
	userService     user.Service
	categoryService category.Service
}

func NewHomeHandler(tmpl *template.Template, ps post.Service, us user.Service, cs category.Service) *HomeHandler {
	return &HomeHandler{
		tmpl:            tmpl,
		postService:     ps,
		userService:     us,
		categoryService: cs,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	h.handleHome(w, r)
}

func (h *HomeHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		categories = nil
	}

	data := struct {
		Posts      []*domain.Post
		Categories []*domain.Category
	}{
		Posts:      posts,
		Categories: categories,
	}

	if err := h.tmpl.ExecuteTemplate(w, "home/index.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
