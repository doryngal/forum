package post

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
	"strings"
)

const (
	editPathPrefix   = "/edit-post/"
	editPostTemplate = "edit-post.html"
	titleField       = "title"
	messageField     = "message"
	categoriesField  = "categories"
)

// EditHandler handles editing posts.
type EditHandler struct {
	tmpl            *template.Template
	postService     post.Service
	userService     user.Service
	sessionService  session.Service
	categoryService category.Service
	errorHandler    errorhandler.Handler
}

func NewEditHandler(
	tmpl *template.Template,
	postService post.Service,
	userService user.Service,
	sessionService session.Service,
	categoryService category.Service,
	errorHandler errorhandler.Handler,
) *EditHandler {
	return &EditHandler{
		tmpl:            tmpl,
		postService:     postService,
		userService:     userService,
		sessionService:  sessionService,
		categoryService: categoryService,
		errorHandler:    errorHandler,
	}
}

func (h *EditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postID, err := h.extractPostID(r)
	if err != nil {
		h.errorHandler.HandleError(w, "Invalid post ID", nil, http.StatusNotFound)
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return
	}

	savedPost, err := h.postService.GetPostByID(postID)
	if err != nil {
		h.errorHandler.HandleError(w, "Post not found", err, http.StatusNotFound)
		return
	}

	if savedPost.UserID != user.ID {
		h.errorHandler.HandleError(w, "Forbidden", nil, http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetEditForm(w, r, savedPost)
	case http.MethodPost:
		h.handlePostEditForm(w, r, savedPost, user.ID)
	default:
		h.errorHandler.HandleError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
	}
}

type EditFormData struct {
	Error              string
	Post               *domain.Post
	User               *domain.User
	Categories         []*domain.Category
	SelectedCategories []*domain.Category
}

func (h *EditHandler) handleGetEditForm(w http.ResponseWriter, r *http.Request, post *domain.Post) {
	data, err := h.prepareEditFormData(r, post, "")
	if err != nil {
		h.errorHandler.HandleError(w, "Failed to prepare form data", err, http.StatusInternalServerError)
		return
	}

	if err := h.tmpl.ExecuteTemplate(w, editPostTemplate, data); err != nil {
		h.errorHandler.HandleError(w, "Failed to render edit form", err, http.StatusInternalServerError)
	}
}

func (h *EditHandler) handlePostEditForm(w http.ResponseWriter, r *http.Request, post *domain.Post, userID uuid.UUID) {
	if err := r.ParseForm(); err != nil {
		h.errorHandler.HandleError(w, "Invalid form data", err, http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue(titleField))
	content := strings.TrimSpace(r.FormValue(messageField))
	categoryIDs := r.Form[categoriesField]

	if title == "" || content == "" {
		data, err := h.prepareEditFormData(r, post, "Title and content are required")
		if err != nil {
			h.errorHandler.HandleError(w, "Failed to prepare form data", err, http.StatusInternalServerError)
			return
		}

		if err := h.tmpl.ExecuteTemplate(w, editPostTemplate, data); err != nil {
			h.errorHandler.HandleError(w, "Failed to render edit form", err, http.StatusInternalServerError)
		}
		return
	}

	// Update post
	post.Title = title
	post.Content = content
	if err := h.postService.UpdatePost(post, userID); err != nil {
		h.errorHandler.HandleError(w, "Failed to update post", err, http.StatusInternalServerError)
		return
	}

	// Update categories
	if err := h.updatePostCategories(post.ID, categoryIDs); err != nil {
		h.errorHandler.HandleError(w, "Failed to update categories", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+post.ID.String(), http.StatusSeeOther)
}

func (h *EditHandler) extractPostID(r *http.Request) (uuid.UUID, error) {
	idStr := strings.TrimPrefix(r.URL.Path, editPathPrefix)
	return uuid.Parse(idStr)
}

func (h *EditHandler) prepareEditFormData(r *http.Request, post *domain.Post, errorMsg string) (*EditFormData, error) {
	user, err := h.getUserFromSession(r)
	if err != nil {
		return nil, err
	}

	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		categories = []*domain.Category{}
	}

	return &EditFormData{
		Error:              errorMsg,
		Post:               post,
		User:               user,
		Categories:         categories,
		SelectedCategories: post.Categories,
	}, nil
}

func (h *EditHandler) updatePostCategories(postID uuid.UUID, categoryIDs []string) error {
	var validCategoryIDs []uuid.UUID

	for _, idStr := range categoryIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			validCategoryIDs = append(validCategoryIDs, id)
		}
	}

	return h.categoryService.AssignCategoriesToPost(postID, validCategoryIDs)
}

func (h *EditHandler) getUserFromSession(r *http.Request) (*domain.User, error) {
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
