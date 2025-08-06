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

// EditHandler handles editing posts.
type EditHandler struct {
	tmpl            *template.Template
	postService     post.Service
	userService     user.Service
	sessionService  session.Service
	categoryService category.Service
	errorHandler    errorhandler.Handler
}

func NewEditHandler(t *template.Template, ps post.Service, us user.Service, ss session.Service, cs category.Service, errorHandler errorhandler.Handler) *EditHandler {
	return &EditHandler{t, ps, us, ss, cs, errorHandler}
}

func (h *EditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/edit-post/")
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

	savedPost, err := h.postService.GetPostByID(postID)
	if err != nil || savedPost.UserID != user.ID {
		h.errorHandler.HandleError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.renderEditForm(w, r, savedPost, nil)
	case http.MethodPost:
		h.handleEditPost(w, r, savedPost)
	default:
		h.errorHandler.HandleError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
	}
}

func (h *EditHandler) renderEditForm(w http.ResponseWriter, r *http.Request, post *domain.Post, errMsg *string) {
	categories, _ := h.categoryService.GetAllCategories()

	user, err := h.getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	data := struct {
		Error              string
		Post               *domain.Post
		User               *domain.User
		Categories         []*domain.Category
		SelectedCategories []*domain.Category
	}{
		Post:               post,
		User:               user,
		Categories:         categories,
		SelectedCategories: post.Categories,
	}
	if errMsg != nil {
		data.Error = *errMsg
	}
	h.tmpl.ExecuteTemplate(w, "edit-post.html", data)
}

func (h *EditHandler) handleEditPost(w http.ResponseWriter, r *http.Request, post *domain.Post) {
	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("message")
	catIDs := r.Form["categories"]

	if title == "" || content == "" {
		err := "Title and content are required."
		h.renderEditForm(w, r, post, &err)
		return
	}

	post.Title = title
	post.Content = content
	_ = h.postService.UpdatePost(post)

	var categoryUUIDs []uuid.UUID
	for _, idStr := range catIDs {
		id, err := uuid.Parse(idStr)
		if err == nil {
			categoryUUIDs = append(categoryUUIDs, id)
		}
	}
	h.categoryService.AssignCategoriesToPost(post.ID, categoryUUIDs)

	http.Redirect(w, r, "/post/"+post.ID.String(), http.StatusSeeOther)
}

func (h *EditHandler) getUserFromSession(r *http.Request) (*domain.User, error) {
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
