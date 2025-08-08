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
	createPostTemplate = "create-post.html"
	homePath           = "/"
	tagsField          = "tags"
)

type CreateHandler struct {
	tmpl            *template.Template
	userService     user.Service
	postService     post.Service
	sessionService  session.Service
	categoryService category.Service
	errorHandler    errorhandler.Handler
}

func NewCreateHandler(
	tmpl *template.Template,
	userService user.Service,
	postService post.Service,
	sessionService session.Service,
	categoryService category.Service,
	errorHandler errorhandler.Handler,
) *CreateHandler {
	return &CreateHandler{
		tmpl:            tmpl,
		userService:     userService,
		postService:     postService,
		sessionService:  sessionService,
		categoryService: categoryService,
		errorHandler:    errorHandler,
	}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetCreateForm(w, r)
	case http.MethodPost:
		h.handlePostCreateForm(w, r)
	default:
		h.errorHandler.HandleError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
	}
}

type CreatePostData struct {
	Error              string
	Form               map[string]string
	SelectedCategories []string
	User               *domain.User
	Categories         []*domain.Category
}

func (h *CreateHandler) handleGetCreateForm(w http.ResponseWriter, r *http.Request) {
	data, err := h.prepareCreateFormData(r, "", nil, nil)
	if err != nil {
		h.errorHandler.HandleError(w, "Failed to prepare form data", err, http.StatusInternalServerError)
		return
	}

	if err := h.renderTemplate(w, data); err != nil {
		h.errorHandler.HandleError(w, "Failed to render template", err, http.StatusInternalServerError)
	}
}

func (h *CreateHandler) handlePostCreateForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.errorHandler.HandleError(w, "Invalid form data", err, http.StatusBadRequest)
		return
	}

	formData := h.extractFormData(r)
	if err := h.validateFormData(formData); err != nil {
		data, _ := h.prepareCreateFormData(r, err.Error(), formData, r.Form[categoriesField])
		h.renderTemplate(w, data)
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return
	}

	categoryUUIDs, err := h.parseCategoryIDs(r.Form[categoriesField])
	if err != nil {
		data, _ := h.prepareCreateFormData(r, "Invalid category selection", formData, r.Form[categoriesField])
		h.renderTemplate(w, data)
		return
	}

	post := h.createPostObject(formData, user.ID)
	if err := h.createPostWithCategories(post, categoryUUIDs); err != nil {
		data, _ := h.prepareCreateFormData(r, "Failed to create post: "+err.Error(), formData, r.Form[categoriesField])
		h.renderTemplate(w, data)
		return
	}

	http.Redirect(w, r, homePath, http.StatusSeeOther)
}

func (h *CreateHandler) prepareCreateFormData(r *http.Request, errorMsg string, formData map[string]string, selectedCategories []string) (*CreatePostData, error) {
	user, err := h.getUserFromSession(r)
	if err != nil && errorMsg == "" {
		return nil, err
	}

	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		categories = []*domain.Category{}
	}

	return &CreatePostData{
		Error:              errorMsg,
		Form:               formData,
		SelectedCategories: selectedCategories,
		User:               user,
		Categories:         categories,
	}, nil
}

func (h *CreateHandler) extractFormData(r *http.Request) map[string]string {
	return map[string]string{
		titleField:   strings.TrimSpace(r.FormValue(titleField)),
		tagsField:    strings.TrimSpace(r.FormValue(tagsField)),
		messageField: strings.TrimSpace(r.FormValue(messageField)),
	}
}

func (h *CreateHandler) validateFormData(formData map[string]string) error {
	if formData[titleField] == "" {
		return domain.ErrTitleRequired
	}
	if formData[messageField] == "" {
		return domain.ErrContentRequired
	}
	return nil
}

func (h *CreateHandler) parseCategoryIDs(categoryIDs []string) ([]uuid.UUID, error) {
	if len(categoryIDs) == 0 {
		return nil, domain.ErrCategoryRequired
	}

	var categoryUUIDs []uuid.UUID
	for _, idStr := range categoryIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, domain.ErrInvalidCategoryID
		}
		categoryUUIDs = append(categoryUUIDs, id)
	}
	return categoryUUIDs, nil
}

func (h *CreateHandler) createPostObject(formData map[string]string, userID uuid.UUID) *domain.Post {
	return &domain.Post{
		ID:      uuid.New(),
		Title:   formData[titleField],
		Content: formData[messageField],
		UserID:  userID,
		Tags:    strings.Fields(formData[tagsField]),
	}
}

func (h *CreateHandler) createPostWithCategories(post *domain.Post, categoryIDs []uuid.UUID) error {
	if err := h.postService.CreatePost(post); err != nil {
		return err
	}
	return h.categoryService.AssignCategoriesToPost(post.ID, categoryIDs)
}

func (h *CreateHandler) renderTemplate(w http.ResponseWriter, data *CreatePostData) error {
	return h.tmpl.ExecuteTemplate(w, createPostTemplate, data)
}

func (h *CreateHandler) getUserFromSession(r *http.Request) (*domain.User, error) {
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
