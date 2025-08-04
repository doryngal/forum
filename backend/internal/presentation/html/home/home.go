package home

import (
	"forum/internal/domain"
	"forum/internal/service/category"
	"forum/internal/service/post"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	tmpl            *template.Template
	postService     post.Service
	userService     user.Service
	categoryService category.Service
	sessionService  session.Service
}

func NewHomeHandler(tmpl *template.Template, ps post.Service, us user.Service, cs category.Service, ss session.Service) *Handler {
	return &Handler{
		tmpl:            tmpl,
		postService:     ps,
		userService:     us,
		categoryService: cs,
		sessionService:  ss,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleHome(w, r)
	case http.MethodPost:
		h.handleAction(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleHome(w http.ResponseWriter, r *http.Request) {
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
		User:       currentUser,
		Posts:      posts,
		Categories: categories,
	}

	if err := h.tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
	}
}

func (h *Handler) handleAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	action := r.FormValue("action")
	postIDStr := r.FormValue("post_id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	switch action {
	case "like":
		err = h.postService.LikePost(postID, userID)
	case "dislike":
		err = h.postService.DislikePost(postID, userID)
	default:
		http.Error(w, "Unknown action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Action failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return uuid.Nil, err
	}

	sess, err := h.sessionService.GetByToken(cookie.Value)
	if err != nil {
		return uuid.Nil, err
	}

	return sess.UserID, nil
}
