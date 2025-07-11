package post

import (
	"errors"
	"forum/internal/domain"
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
	tmpl           *template.Template
	userService    user.Service
	postService    post.Service
	commentService comment.Service
	sessionService session.Service
}

func NewPostHandler(tmpl *template.Template, userService user.Service, postService post.Service, commentService comment.Service, sessionService session.Service) *PostHandler {
	return &PostHandler{
		tmpl:           tmpl,
		userService:    userService,
		postService:    postService,
		commentService: commentService,
		sessionService: sessionService,
	}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	User     *domain.User
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
	case "like_comment":
		commentIDStr := r.FormValue("comment_id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		err = h.commentService.LikeComment(commentID, userID)

	case "dislike_comment":
		commentIDStr := r.FormValue("comment_id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		err = h.commentService.DislikeComment(commentID, userID)
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
