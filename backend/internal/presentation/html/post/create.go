package post

import (
	"errors"
	"forum/internal/domain"
	"forum/internal/service/post"
	"forum/internal/service/user"
	"html/template"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type CreateHandler struct {
	tmpl        *template.Template
	userService user.Service
	postService post.Service
}

func NewCreateHandler(tmpl *template.Template, userService user.Service, postService post.Service) *CreateHandler {
	return &CreateHandler{
		tmpl:        tmpl,
		userService: userService,
		postService: postService,
	}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.renderCreateForm(w, nil)
	case http.MethodPost:
		h.handleCreatePost(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type CreatePostData struct {
	Error string
	Form  map[string]string
}

func (h *CreateHandler) renderCreateForm(w http.ResponseWriter, data *CreatePostData) {
	if data == nil {
		data = &CreatePostData{}
	}
	if err := h.tmpl.ExecuteTemplate(w, "create-post.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h *CreateHandler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderCreateForm(w, &CreatePostData{
			Error: "Invalid form data",
		})
		return
	}

	formData := map[string]string{
		"title":   strings.TrimSpace(r.FormValue("title")),
		"tags":    strings.TrimSpace(r.FormValue("tags")),
		"content": r.FormValue("message"),
	}

	if len(formData["title"]) == 0 {
		h.renderCreateForm(w, &CreatePostData{
			Error: "Title is required",
			Form:  formData,
		})
		return
	}

	user, err := h.getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	post := &domain.Post{
		Title:   formData["title"],
		Content: formData["content"],
		UserID:  user.ID,
		Tags:    strings.Fields(formData["tags"]),
	}

	if err := h.postService.CreatePost(post); err != nil {
		h.renderCreateForm(w, &CreatePostData{
			Error: "Failed to create post: " + err.Error(),
			Form:  formData,
		})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *CreateHandler) getUserFromSession(r *http.Request) (*domain.User, error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(cookie.Value)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	return user, nil
}
