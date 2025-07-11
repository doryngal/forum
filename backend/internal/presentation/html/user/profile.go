package profile

import (
	"forum/internal/domain"
	"forum/internal/service/post"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"strings"
)

type ProfileHandler struct {
	tmpl           *template.Template
	userService    user.Service
	postService    post.Service
	sessionService session.Service
}

func NewProfileHandler(
	tmpl *template.Template,
	userService user.Service,
	postService post.Service,
	sessionService session.Service,
) *ProfileHandler {
	return &ProfileHandler{
		tmpl:           tmpl,
		userService:    userService,
		postService:    postService,
		sessionService: sessionService,
	}
}

func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch r.Method {
	case http.MethodGet:
		// /profile — текущий пользователь из сессии
		if path == "/profile" || path == "/profile/" {
			h.handleOwnProfile(w, r)
			return
		}

		// /profile/username — чужой профиль
		if strings.HasPrefix(path, "/profile/") {
			username := strings.TrimPrefix(path, "/profile/")
			h.handleGetProfile(w, r, username)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProfileHandler) handleOwnProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	session, err := h.sessionService.GetByToken(cookie.Value)
	if err != nil || session.UserID == uuid.Nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := h.userService.GetUserByID(session.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	h.renderProfilePage(w, user)
}

func (h *ProfileHandler) renderProfilePage(w http.ResponseWriter, user *domain.User) {
	stats, err := h.userService.GetUserStats(user.ID)
	if err != nil {
		http.Error(w, "Failed to get user stats", http.StatusInternalServerError)
		return
	}

	posts, err := h.postService.GetPostsByUser(user.ID)
	if err != nil {
		http.Error(w, "Failed to get user posts", http.StatusInternalServerError)
		return
	}

	data := domain.ProfileData{
		User:  user,
		Stats: stats,
		Posts: posts,
	}

	err = h.tmpl.ExecuteTemplate(w, "profile.html", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h *ProfileHandler) handleGetProfile(w http.ResponseWriter, r *http.Request, username string) {
	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	h.renderProfilePage(w, user)
}
