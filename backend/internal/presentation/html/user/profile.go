package profile

import (
	"forum/internal/domain"
	"forum/internal/presentation/html/errorhandler"
	"forum/internal/service/comment"
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
	commentService comment.Service
	sessionService session.Service
	errorHandler   errorhandler.Handler
}

func NewProfileHandler(
	tmpl *template.Template,
	userService user.Service,
	postService post.Service,
	commentService comment.Service,
	sessionService session.Service,
	errorHandler errorhandler.Handler,
) *ProfileHandler {
	return &ProfileHandler{
		tmpl:           tmpl,
		userService:    userService,
		postService:    postService,
		commentService: commentService,
		sessionService: sessionService,
		errorHandler:   errorHandler,
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
	case http.MethodPost:
		h.handleAction(w, r)
	default:
		h.errorHandler.HandleError(w, "Method not allowed", nil, http.StatusMethodNotAllowed)
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
		h.errorHandler.HandleError(w, "User not found", err, http.StatusNotFound)
		return
	}

	h.renderProfilePage(w, user, user.ID)
}

func (h *ProfileHandler) handleAction(w http.ResponseWriter, r *http.Request) {
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
		h.errorHandler.HandleError(w, "Unknown action", err, http.StatusBadRequest)
		return
	}

	if err != nil {
		h.errorHandler.HandleError(w, "Action failed", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *ProfileHandler) getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
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

func (h *ProfileHandler) renderProfilePage(w http.ResponseWriter, user *domain.User, sessionID uuid.UUID) {
	createdPosts, _ := h.postService.GetPostsByUserID(user.ID, sessionID)
	likedPosts, _ := h.postService.GetLikedPosts(user.ID)
	dislikedPosts, _ := h.postService.GetDislikedPosts(user.ID)
	comments, _ := h.commentService.GetCommentsByUserID(user.ID) // Или commentService

	h.tmpl.ExecuteTemplate(w, "profile.html", map[string]interface{}{
		"User":          user,
		"Posts":         createdPosts,
		"LikedPosts":    likedPosts,
		"DislikedPosts": dislikedPosts,
		"Comments":      comments,
		"Stats": map[string]int{
			"PostCount":    len(createdPosts),
			"LikeCount":    len(likedPosts),
			"DislikeCount": len(dislikedPosts),
			"CommentCount": len(comments),
		},
	})
}

func (h *ProfileHandler) handleGetProfile(w http.ResponseWriter, r *http.Request, username string) {
	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		h.errorHandler.HandleError(w, "User not found", err, http.StatusNotFound)
		return
	}

	sessionID, _ := h.getUserIDFromSession(r)

	h.renderProfilePage(w, user, sessionID)
}
