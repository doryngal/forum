package home

import (
	"forum/internal/domain"
	"forum/internal/presentation/html/errorhandler"
	"forum/internal/service/category"
	"forum/internal/service/post"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"github.com/google/uuid"
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl            *template.Template
	postService     post.Service
	userService     user.Service
	categoryService category.Service
	sessionService  session.Service
	errorHandler    errorhandler.Handler
}

func NewHomeHandler(tmpl *template.Template, ps post.Service, us user.Service, cs category.Service, ss session.Service, errorHandler errorhandler.Handler) *Handler {
	return &Handler{
		tmpl:            tmpl,
		postService:     ps,
		userService:     us,
		categoryService: cs,
		sessionService:  ss,
		errorHandler:    errorHandler,
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
		h.errorHandler.HandleError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
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
		h.errorHandler.HandleError(w, "Failed to fetch posts", err, http.StatusInternalServerError)
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
		h.errorHandler.HandleError(w, "Template rendering failed", err, http.StatusInternalServerError)
	}
}

func (h *Handler) handleAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.errorHandler.HandleError(w, "Invalid form data", err, http.StatusBadRequest)
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
		h.errorHandler.HandleError(w, "Invalid post ID", err, http.StatusBadRequest)
		return
	}

	switch action {
	case "like":
		err = h.postService.LikePost(postID, userID)
	case "dislike":
		err = h.postService.DislikePost(postID, userID)
	default:
		h.errorHandler.HandleError(w, "Unknown action", nil, http.StatusBadRequest)
		return
	}

	if err != nil {
		h.errorHandler.HandleError(w, "Action failed", err, http.StatusInternalServerError)
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
