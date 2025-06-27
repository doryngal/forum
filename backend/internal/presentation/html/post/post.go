package post

import (
	"errors"
	"forum/internal/domain"
	"forum/internal/service/comment"
	"forum/internal/service/post"
	"forum/internal/service/user"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type PostHandler struct {
	tmpl           *template.Template
	userService    user.Service
	postService    post.Service
	commentService comment.Service
}

func NewPostHandler(tmpl *template.Template, userService user.Service, postService post.Service, commentService comment.Service) *PostHandler {
	return &PostHandler{
		tmpl:           tmpl,
		userService:    userService,
		postService:    postService,
		commentService: commentService,
	}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ожидается маршрут вида /post/{uuid}
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetPost(w, r, postID)
	case http.MethodPost:
		h.handlePostAction(w, r, postID)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type PostData struct {
	Post     *domain.Post
	Comments []*domain.Comment
	Error    string
	Success  string
}

func (h *PostHandler) handleGetPost(w http.ResponseWriter, r *http.Request, postID uuid.UUID) {
	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	comments, err := h.commentService.GetCommentsByPost(postID)
	if err != nil {
		http.Error(w, "Failed to load comments", http.StatusInternalServerError)
		return
	}

	data := PostData{
		Post:     post,
		Comments: comments,
	}

	if err := h.tmpl.ExecuteTemplate(w, "post.html", data); err != nil {
		log.Printf("template render failed: %v", err)
	}
}

func (h *PostHandler) handlePostAction(w http.ResponseWriter, r *http.Request, postID uuid.UUID) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	action := r.FormValue("action")
	switch action {
	case "like":
		err = h.postService.LikePost(postID, userID)
	case "dislike":
		err = h.postService.DislikePost(postID, userID)
	case "comment":
		err = h.handleCreateComment(r, postID, userID)
	default:
		http.Error(w, "Unknown action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Action failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func (h *PostHandler) handleCreateComment(r *http.Request, postID, userID uuid.UUID) error {
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		return errors.New("comment cannot be empty")
	}
	comment := &domain.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}
	return h.commentService.CreateComment(comment)
}

func (h *PostHandler) getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(cookie.Value)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
