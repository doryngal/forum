package post

import (
	"errors"
	"forum/internal/domain"
	"forum/internal/presentation/html/errorhandler"
	"forum/internal/service/category"
	"forum/internal/service/comment"
	"forum/internal/service/post"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type PostHandler struct {
	tmpl            *template.Template
	userService     user.Service
	postService     post.Service
	commentService  comment.Service
	sessionService  session.Service
	categoryService category.Service
	errorHandler    errorhandler.Handler
}

func NewPostHandler(tmpl *template.Template, us user.Service, ps post.Service, cs comment.Service, ss session.Service, cts category.Service, errorHandler errorhandler.Handler) *PostHandler {
	return &PostHandler{
		tmpl:            tmpl,
		userService:     us,
		postService:     ps,
		commentService:  cs,
		sessionService:  ss,
		categoryService: cts,
		errorHandler:    errorHandler,
	}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		h.errorHandler.HandleError(w, "Invalid post ID", err, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetPost(w, r, postID)
	case http.MethodPost:
		h.handlePostAction(w, r, postID)
	default:
		h.errorHandler.HandleError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
	}
}

type PostData struct {
	Post       *domain.Post
	Comments   []*domain.Comment
	Error      string
	Success    string
	User       *domain.User
	Categories []*domain.Category
}

func (h *PostHandler) handleGetPost(w http.ResponseWriter, r *http.Request, postID uuid.UUID) {
	userID, err := h.getUserIDFromSession(r)

	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	comments, err := h.commentService.GetCommentsByPost(postID, userID)
	if err != nil {
		h.errorHandler.HandleError(w, "Failed to load comments", err, http.StatusInternalServerError)
		return
	}

	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		categories = nil
	}

	data := PostData{
		Post:       post,
		Comments:   comments,
		Categories: categories,
	}

	// Добавляем информацию о пользователе из сессии
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		sess, err := h.sessionService.GetByToken(cookie.Value)
		if err == nil {
			data.User, _ = h.userService.GetUserByID(sess.UserID)
		}
	}

	if err := h.tmpl.ExecuteTemplate(w, "post.html", data); err != nil {
		log.Printf("template render failed: %v", err)
	}
}

func (h *PostHandler) handlePostAction(w http.ResponseWriter, r *http.Request, postID uuid.UUID) {
	if err := r.ParseForm(); err != nil {
		h.errorHandler.HandleError(w, "Invalid form data", err, http.StatusBadRequest)
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
	case "like_comment":
		commentIDStr := r.FormValue("comment_id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			h.errorHandler.HandleError(w, "Invalid comment ID", err, http.StatusBadRequest)
			return
		}
		err = h.commentService.LikeComment(commentID, userID)

	case "dislike_comment":
		commentIDStr := r.FormValue("comment_id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			h.errorHandler.HandleError(w, "Invalid comment ID", err, http.StatusBadRequest)
			return
		}
		err = h.commentService.DislikeComment(commentID, userID)
	default:
		h.errorHandler.HandleError(w, "Unknown action", err, http.StatusBadRequest)
		return
	}

	if err != nil {
		h.errorHandler.HandleError(w, "Action failed", err, http.StatusInternalServerError)
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
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return uuid.Nil, err
	}

	sess, err := h.sessionService.GetByToken(cookie.Value)
	if err != nil {
		return uuid.Nil, err
	}

	return sess.UserID, nil
}
