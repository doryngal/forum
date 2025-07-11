package post

import (
	"errors"
	"forum/internal/domain"
	"forum/internal/service/post"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"html/template"
	"net/http"
	"strings"
)

type CreateHandler struct {
	tmpl           *template.Template
	userService    user.Service
	postService    post.Service
	sessionService session.Service
}

func NewCreateHandler(tmpl *template.Template, userService user.Service, postService post.Service, sessionService session.Service) *CreateHandler {
	return &CreateHandler{
		tmpl:           tmpl,
		userService:    userService,
		postService:    postService,
		sessionService: sessionService,
	}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.renderCreateForm(w, r, nil)
	case http.MethodPost:
		h.handleCreatePost(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type CreatePostData struct {
	Error string
	Form  map[string]string
	User  *domain.User
}

func (h *CreateHandler) renderCreateForm(w http.ResponseWriter, r *http.Request, data *CreatePostData) {
	if data == nil {
		data = &CreatePostData{}
	}

	// Добавляем информацию о пользователе из сессии
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		sess, err := h.sessionService.GetByToken(cookie.Value)
		if err == nil {
			data.User, _ = h.userService.GetUserByID(sess.UserID)
		}
	}

	if err := h.tmpl.ExecuteTemplate(w, "create-post.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
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

	if len(formData["title"]) == 0 {
		h.renderCreateForm(w, r, &CreatePostData{
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
		h.renderCreateForm(w, r, &CreatePostData{
			Error: "Failed to create post: " + err.Error(),
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
