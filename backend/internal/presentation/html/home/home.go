package home

import (
	"forum/internal/domain"
	"forum/internal/service/category"
	"forum/internal/service/post"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"html/template"
	"net/http"
)

type HomeHandler struct {
	tmpl            *template.Template
	postService     post.Service
	userService     user.Service
	categoryService category.Service
	sessionService  session.Service
}

func NewHomeHandler(tmpl *template.Template, ps post.Service, us user.Service, cs category.Service, ss session.Service) *HomeHandler {
	return &HomeHandler{
		tmpl:            tmpl,
		postService:     ps,
		userService:     us,
		categoryService: cs,
		sessionService:  ss,
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
	var currentUser *domain.User

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		sess, err := h.sessionService.GetByToken(cookie.Value)
		if err == nil {
			currentUser, _ = h.userService.GetUserByID(sess.UserID)
		}
	}

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
		User       *domain.User
		Posts      []*domain.Post
		Categories []*domain.Category
	}{
		User:       currentUser, // ← добавим в шаблон
		Posts:      posts,
		Categories: categories,
	}

	if err := h.tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
