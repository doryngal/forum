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

// DeleteHandler handles post deletion.
type DeleteHandler struct {
	postService    post.Service
	sessionService session.Service
	userService    user.Service
	errorHandler   errorhandler.Handler
}

func NewDeleteHandler(ps post.Service, ss session.Service, us user.Service, errorHandler errorhandler.Handler) *DeleteHandler {
	return &DeleteHandler{ps, ss, us, errorHandler}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.errorHandler.HandleError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/delete-post/")
	postID, err := uuid.Parse(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	p, err := h.postService.GetPostByID(postID)
	if err != nil || p.UserID != user.ID {
		h.errorHandler.HandleError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	h.postService.DeletePost(postID, user.ID)
	http.Redirect(w, r, "/profile/", http.StatusSeeOther)
}

func (h *DeleteHandler) getUserFromSession(r *http.Request) (*domain.User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}
	sess, err := h.sessionService.GetByToken(cookie.Value)
	if err != nil {
		return nil, err
	}
	return h.userService.GetUserByID(sess.UserID)
}
