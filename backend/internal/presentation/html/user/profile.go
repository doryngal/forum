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
	"sync"
)

const (
	profilePath     = "/profile"
	loginPath       = "/login"
	sessionCookie   = "session_id"
	profileTemplate = "profile.html"
	defaultRedirect = "/"
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
	switch r.Method {
	case http.MethodGet:
		h.handleGetRequest(w, r)
	case http.MethodPost:
		h.handlePostRequest(w, r)
	default:
		h.errorHandler.HandleError(w, "Method not allowed", nil, http.StatusMethodNotAllowed)
	}
}

func (h *ProfileHandler) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == profilePath || path == profilePath+"/":
		h.handleOwnProfile(w, r)
	case strings.HasPrefix(path, profilePath+"/"):
		username := strings.TrimPrefix(path, profilePath+"/")
		h.handleViewProfile(w, r, username)
	default:
		http.Redirect(w, r, defaultRedirect, http.StatusFound)
	}
}

func (h *ProfileHandler) handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.errorHandler.HandleError(w, "Invalid form data", err, http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return
	}

	if err := h.processProfileAction(w, r, userID); err != nil {
		h.errorHandler.HandleError(w, "Failed to process action", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, profilePath, http.StatusSeeOther)
}

func (h *ProfileHandler) handleOwnProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		h.errorHandler.HandleError(w, "User not found", err, http.StatusNotFound)
		return
	}

	h.renderProfilePage(w, user, userID)
}

func (h *ProfileHandler) handleViewProfile(w http.ResponseWriter, r *http.Request, username string) {
	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		h.errorHandler.HandleError(w, "User not found", err, http.StatusNotFound)
		return
	}

	sessionID, _ := h.getUserIDFromSession(r)
	h.renderProfilePage(w, user, sessionID)
}

func (h *ProfileHandler) processProfileAction(w http.ResponseWriter, r *http.Request, userID uuid.UUID) error {
	action := r.FormValue("action")
	postIDStr := r.FormValue("post_id")

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return err
	}

	switch action {
	case "like":
		return h.postService.LikePost(postID, userID)
	case "dislike":
		return h.postService.DislikePost(postID, userID)
	default:
		return domain.ErrInvalidAction
	}
}

func (h *ProfileHandler) getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return uuid.Nil, err
	}

	sess, err := h.sessionService.GetByToken(cookie.Value)
	if err != nil {
		return uuid.Nil, err
	}

	if sess.UserID == uuid.Nil {
		return uuid.Nil, domain.ErrInvalidSession
	}

	return sess.UserID, nil
}

func (h *ProfileHandler) renderProfilePage(w http.ResponseWriter, user *domain.User, sessionID uuid.UUID) {
	data, err := h.prepareProfileData(user, sessionID)
	if err != nil {
		h.errorHandler.HandleError(w, "Failed to prepare profile data", err, http.StatusInternalServerError)
		return
	}

	if err := h.tmpl.ExecuteTemplate(w, profileTemplate, data); err != nil {
		h.errorHandler.HandleError(w, "Failed to render profile page", err, http.StatusInternalServerError)
	}
}

func (h *ProfileHandler) prepareProfileData(user *domain.User, sessionID uuid.UUID) (map[string]interface{}, error) {
	var (
		createdPosts  []*domain.Post
		likedPosts    []*domain.Post
		dislikedPosts []*domain.Post
		comments      []*domain.CommentWithPostTitle

		errCreatedPosts  error
		errLikedPosts    error
		errDislikedPosts error
		errComments      error
	)

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		createdPosts, errCreatedPosts = h.postService.GetPostsByUserID(user.ID, sessionID)
	}()

	go func() {
		defer wg.Done()
		likedPosts, errLikedPosts = h.postService.GetLikedPosts(user.ID)
	}()

	go func() {
		defer wg.Done()
		dislikedPosts, errDislikedPosts = h.postService.GetDislikedPosts(user.ID)
	}()

	go func() {
		defer wg.Done()
		comments, errComments = h.commentService.GetCommentsByUserID(user.ID)
	}()

	wg.Wait()

	if errCreatedPosts != nil {
		return nil, errCreatedPosts
	}
	if errLikedPosts != nil {
		return nil, errLikedPosts
	}
	if errDislikedPosts != nil {
		return nil, errDislikedPosts
	}
	if errComments != nil {
		return nil, errComments
	}

	return map[string]interface{}{
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
	}, nil
}
