package home

import (
	"errors"
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

const (
	homePath      = "/"
	loginPath     = "/login"
	sessionCookie = "session_id"
	homeTemplate  = "home.html"
	actionField   = "action"
	postIDField   = "post_id"
	likeAction    = "like"
	dislikeAction = "dislike"
)

type Handler struct {
	tmpl            *template.Template
	postService     post.Service
	userService     user.Service
	categoryService category.Service
	sessionService  session.Service
	errorHandler    errorhandler.Handler
}

func NewHomeHandler(
	tmpl *template.Template,
	postService post.Service,
	userService user.Service,
	categoryService category.Service,
	sessionService session.Service,
	errorHandler errorhandler.Handler,
) *Handler {
	return &Handler{
		tmpl:            tmpl,
		postService:     postService,
		userService:     userService,
		categoryService: categoryService,
		sessionService:  sessionService,
		errorHandler:    errorHandler,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != homePath {
		h.errorHandler.HandleError(w, "Not Found", nil, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleHomePage(w, r)
	case http.MethodPost:
		h.handlePostAction(w, r)
	default:
		h.errorHandler.HandleError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
	}
}

type HomePageData struct {
	User       *domain.User
	Posts      []*domain.Post
	Categories []*domain.Category
}

func (h *Handler) handleHomePage(w http.ResponseWriter, r *http.Request) {
	data, err := h.prepareHomeData(r)
	if err != nil {
		h.errorHandler.HandleError(w, "Failed to prepare home data", err, http.StatusInternalServerError)
		return
	}

	if err := h.renderHomePage(w, data); err != nil {
		h.errorHandler.HandleError(w, "Failed to render home page", err, http.StatusInternalServerError)
	}
}

func (h *Handler) handlePostAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.errorHandler.HandleError(w, "Invalid form data", err, http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return
	}

	if err := h.processPostAction(r, userID); err != nil {
		h.errorHandler.HandleError(w, "Action failed", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, homePath, http.StatusSeeOther)
}

func (h *Handler) prepareHomeData(r *http.Request) (*HomePageData, error) {
	user, err := h.getCurrentUser(r)
	if err != nil && !errors.Is(err, domain.ErrUnauthorized) {
		return nil, err
	}

	posts, err := h.postService.GetAllPosts()
	if err != nil {
		return nil, err
	}

	categories, _ := h.categoryService.GetAllCategories() // Ignore category error

	return &HomePageData{
		User:       user,
		Posts:      posts,
		Categories: categories,
	}, nil
}

func (h *Handler) renderHomePage(w http.ResponseWriter, data *HomePageData) error {
	return h.tmpl.ExecuteTemplate(w, homeTemplate, data)
}

func (h *Handler) getCurrentUser(r *http.Request) (*domain.User, error) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	sess, err := h.sessionService.GetByToken(cookie.Value)
	if err != nil {
		return nil, domain.ErrInvalidSession
	}

	if sess.UserID == uuid.Nil {
		return nil, domain.ErrInvalidSession
	}

	return h.userService.GetUserByID(sess.UserID)
}

func (h *Handler) processPostAction(r *http.Request, userID uuid.UUID) error {
	action := r.FormValue(actionField)
	postID, err := uuid.Parse(r.FormValue(postIDField))
	if err != nil {
		return domain.ErrInvalidPostID
	}

	switch action {
	case likeAction:
		return h.postService.LikePost(postID, userID)
	case dislikeAction:
		return h.postService.DislikePost(postID, userID)
	default:
		return domain.ErrInvalidAction
	}
}

func (h *Handler) getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	user, err := h.getCurrentUser(r)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}
