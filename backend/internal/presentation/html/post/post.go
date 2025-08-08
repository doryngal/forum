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
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	postPathPrefix = "/post/"
	loginPath      = "/login"
	sessionCookie  = "session_id"
	postTemplate   = "post.html"
	commentAction  = "comment"
	likeAction     = "like"
	dislikeAction  = "dislike"
	likeComment    = "like_comment"
	dislikeComment = "dislike_comment"
	commentIDField = "comment_id"
	contentField   = "content"
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

func NewPostHandler(
	tmpl *template.Template,
	userService user.Service,
	postService post.Service,
	commentService comment.Service,
	sessionService session.Service,
	categoryService category.Service,
	errorHandler errorhandler.Handler,
) *PostHandler {
	return &PostHandler{
		tmpl:            tmpl,
		userService:     userService,
		postService:     postService,
		commentService:  commentService,
		sessionService:  sessionService,
		categoryService: categoryService,
		errorHandler:    errorHandler,
	}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postID, err := h.extractPostID(r)
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
	userID, _ := h.getUserIDFromSession(r)

	data, err := h.preparePostData(postID, userID, r)
	if err != nil {
		h.handlePostDataError(w, err)
		return
	}

	if err := h.tmpl.ExecuteTemplate(w, postTemplate, data); err != nil {
		h.errorHandler.HandleError(w, "Failed to render post page", err, http.StatusInternalServerError)
	}
}

func (h *PostHandler) handlePostAction(w http.ResponseWriter, r *http.Request, postID uuid.UUID) {
	if err := r.ParseForm(); err != nil {
		h.errorHandler.HandleError(w, "Invalid form data", err, http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromSession(r)
	if err != nil {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return
	}

	action := r.FormValue("action")
	if err := h.processAction(r, action, postID, userID); err != nil {
		h.errorHandler.HandleError(w, "Action failed", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func (h *PostHandler) extractPostID(r *http.Request) (uuid.UUID, error) {
	postIDStr := strings.TrimPrefix(r.URL.Path, postPathPrefix)
	return uuid.Parse(postIDStr)
}

func (h *PostHandler) preparePostData(postID, userID uuid.UUID, r *http.Request) (*PostData, error) {
	var (
		p   *domain.Post
		c   []*domain.Comment
		u   *domain.User
		err error
	)

	// Получаем основные данные поста
	p, err = h.postService.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	// Получаем комментарии
	c, err = h.commentService.GetCommentsByPost(postID, userID)
	if err != nil {
		return nil, err
	}

	// Получаем данные пользователя из сессии
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		if sess, err := h.sessionService.GetByToken(cookie.Value); err == nil {
			u, _ = h.userService.GetUserByID(sess.UserID)
		}
	}

	return &PostData{
		Post:       p,
		Comments:   c,
		Categories: p.Categories,
		User:       u,
	}, nil
}

func (h *PostHandler) processAction(r *http.Request, action string, postID, userID uuid.UUID) error {
	switch action {
	case likeAction:
		return h.postService.LikePost(postID, userID)
	case dislikeAction:
		return h.postService.DislikePost(postID, userID)
	case commentAction:
		return h.handleCreateComment(r, postID, userID)
	case likeComment, dislikeComment:
		return h.handleCommentReaction(r, action, userID)
	default:
		return domain.ErrInvalidAction
	}
}

func (h *PostHandler) handleCommentReaction(r *http.Request, action string, userID uuid.UUID) error {
	commentIDStr := r.FormValue(commentIDField)
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		return err
	}

	switch action {
	case likeComment:
		return h.commentService.LikeComment(commentID, userID)
	case dislikeComment:
		return h.commentService.DislikeComment(commentID, userID)
	default:
		return domain.ErrInvalidAction
	}
}

func (h *PostHandler) handleCreateComment(r *http.Request, postID, userID uuid.UUID) error {
	content := strings.TrimSpace(r.FormValue(contentField))
	if content == "" {
		return comment.ErrInvalidComment
	}

	comment := &domain.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}
	return h.commentService.CreateComment(comment)
}

func (h *PostHandler) getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil {
		return uuid.Nil, err
	}

	sess, err := h.sessionService.GetByToken(cookie.Value)
	if err != nil {
		return uuid.Nil, err
	}

	if sess.UserID == uuid.Nil {
		return uuid.Nil, domain.ErrInvalidSession
	}

	return sess.UserID, nil
}

func (h *PostHandler) handlePostDataError(w http.ResponseWriter, err error) {
	if errors.Is(err, post.ErrPostNotFound) {
		h.errorHandler.HandleError(w, "Post not found", err, http.StatusNotFound)
	} else {
		h.errorHandler.HandleError(w, "Failed to load post data", err, http.StatusInternalServerError)
	}
}
