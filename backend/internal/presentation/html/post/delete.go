package post

import (
	"forum/internal/domain"
	"forum/internal/presentation/html/errorhandler"
	"forum/internal/service/post"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

const (
	deletePathPrefix = "/delete-post/"
	profilePath      = "/profile/"
)

// DeleteHandler handles post deletion.
type DeleteHandler struct {
	postService    post.Service
	sessionService session.Service
	userService    user.Service
	errorHandler   errorhandler.Handler
}

func NewDeleteHandler(
	postService post.Service,
	sessionService session.Service,
	userService user.Service,
	errorHandler errorhandler.Handler,
) *DeleteHandler {
	return &DeleteHandler{
		postService:    postService,
		sessionService: sessionService,
		userService:    userService,
		errorHandler:   errorHandler,
	}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorHandler.HandleError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
		return
	}

	postID, err := h.extractPostID(r)
	if err != nil {
		h.errorHandler.HandleError(w, "Invalid post ID", err, http.StatusBadRequest)
		return
	}

	user, err := h.authenticateUser(r)
	if err != nil {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return
	}

	if err := h.validatePostOwnership(postID, user.ID); err != nil {
		status := httpStatusFromError(err)
		h.errorHandler.HandleError(w, "Forbidden", err, status)
		return
	}

	if err := h.postService.DeletePost(postID, user.ID); err != nil {
		h.errorHandler.HandleError(w, "Failed to delete post", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, profilePath, http.StatusSeeOther)
}

func (h *DeleteHandler) extractPostID(r *http.Request) (uuid.UUID, error) {
	idStr := strings.TrimPrefix(r.URL.Path, deletePathPrefix)
	return uuid.Parse(idStr)
}

func (h *DeleteHandler) authenticateUser(r *http.Request) (*domain.User, error) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return nil, err
	}

	sess, err := h.sessionService.GetByToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	if sess.UserID == uuid.Nil {
		return nil, domain.ErrInvalidSession
	}

	return h.userService.GetUserByID(sess.UserID)
}

func (h *DeleteHandler) validatePostOwnership(postID, userID uuid.UUID) error {
	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return domain.ErrForbidden
	}

	return nil
}
