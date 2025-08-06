package post

import (
	"errors"
	"fmt"
	"forum/internal/domain"
	"forum/internal/presentation/html/errorhandler"
	"forum/internal/service/category"
	"forum/internal/service/post"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type CreateHandler struct {
	tmpl            *template.Template
	userService     user.Service
	postService     post.Service
	sessionService  session.Service
	categoryService category.Service
	errorHandler    errorhandler.Handler
}

func NewCreateHandler(tmpl *template.Template, userService user.Service, postService post.Service, sessionService session.Service, categoryService category.Service, errorHandler errorhandler.Handler) *CreateHandler {
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
		h.renderCreateForm(w, r, nil)
	case http.MethodPost:
		h.handleCreatePost(w, r)
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

func (h *CreateHandler) renderCreateForm(w http.ResponseWriter, r *http.Request, data *CreatePostData) {
	if data == nil {
		data = &CreatePostData{}
	}

	// Получаем пользователя
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		sess, err := h.sessionService.GetByToken(cookie.Value)
		if err == nil {
			data.User, _ = h.userService.GetUserByID(sess.UserID)
		}
	}

	// Получаем категории
	categories, err := h.categoryService.GetAllCategories()
	if err == nil {
		data.Categories = categories
	}

	if err := h.tmpl.ExecuteTemplate(w, "create-post.html", data); err != nil {
		log.Printf("Template rendering error: %v", err)
	}
}

func (h *CreateHandler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderCreateForm(w, r, &CreatePostData{
			Error: "Invalid form data",
		})
		return
	}

	formData := map[string]string{
		"title":   strings.TrimSpace(r.FormValue("title")),
		"tags":    strings.TrimSpace(r.FormValue("tags")),
		"content": r.FormValue("message"),
	}
	selectedCategories := r.Form["categories"]

	if len(formData["title"]) == 0 {
		h.renderCreateForm(w, r, &CreatePostData{
			Error:              "Title is required",
			Form:               formData,
			SelectedCategories: selectedCategories,
		})
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	categoryIDs := r.Form["categories"] // slice of selected category IDs (as strings)
	fmt.Println(categoryIDs)
	var categoryUUIDs []uuid.UUID
	for _, idStr := range categoryIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("Invalid UUID: %s, error: %v", idStr, err)
			continue // или return, если важно прервать
		}
		categoryUUIDs = append(categoryUUIDs, id)
	}

	if len(categoryUUIDs) == 0 {
		h.renderCreateForm(w, r, &CreatePostData{
			Error: "At least one category must be selected",
			Form:  formData,
		})
		return
	}

	post := &domain.Post{
		ID:      uuid.New(),
		Title:   formData["title"],
		Content: formData["content"],
		UserID:  user.ID,
		Tags:    strings.Fields(formData["tags"]),
	}

	if err := h.postService.CreatePost(post); err != nil {
		log.Printf("Failed to create post: %v", err)
		h.renderCreateForm(w, r, &CreatePostData{
			Error: "Failed to create post: " + err.Error(),
			Form:  formData,
		})
		return
	}

	err = h.categoryService.AssignCategoriesToPost(post.ID, categoryUUIDs)
	if err != nil {
		log.Printf("Failed to assign categories: %v", err)
		h.renderCreateForm(w, r, &CreatePostData{
			Error: "Failed to assign categories to post: " + err.Error(),
			Form:  formData,
		})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *CreateHandler) getUserFromSession(r *http.Request) (*domain.User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	sess, err := h.sessionService.GetByToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.GetUserByID(sess.UserID)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	return user, nil
}
