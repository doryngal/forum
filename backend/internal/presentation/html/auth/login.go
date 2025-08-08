package auth

import (
	"forum/internal/domain"
	"forum/internal/presentation/html/errorhandler"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"html/template"
	"net/http"
)

const (
	loginTemplate   = "login.html"
	homePath        = "/"
	sessionCookie   = "session_id"
	emailOrUsername = "emailOrUsername"
)

type LoginHandler struct {
	tmpl           *template.Template
	userService    user.Service
	sessionService session.Service
	errorHandler   errorhandler.Handler
}

func NewLoginHandler(
	tmpl *template.Template,
	userService user.Service,
	sessionService session.Service,
	errorHandler errorhandler.Handler,
) *LoginHandler {
	return &LoginHandler{
		tmpl:           tmpl,
		userService:    userService,
		sessionService: sessionService,
		errorHandler:   errorHandler,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetLogin(w, r)
	case http.MethodPost:
		h.handlePostLogin(w, r)
	default:
		h.errorHandler.HandleError(w, "Method not allowed", nil, http.StatusMethodNotAllowed)
	}
}

type LoginTemplateData struct {
	Error   string
	Success string
}

func (h *LoginHandler) handleGetLogin(w http.ResponseWriter, r *http.Request) {
	data := &LoginTemplateData{}
	if r.URL.Query().Get("registered") == "true" {
		data.Success = "Registration successful! Please log in."
	}
	h.renderLogin(w, data)
}

func (h *LoginHandler) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLogin(w, &LoginTemplateData{Error: "Invalid form submission"})
		return
	}

	credentials := r.FormValue(emailOrUsername)
	password := r.FormValue(passwordField)

	user, err := h.authenticateUser(credentials, password)
	if err != nil {
		h.renderLogin(w, &LoginTemplateData{Error: "Invalid credentials"})
		return
	}

	if err := h.createUserSession(w, user); err != nil {
		h.errorHandler.HandleError(w, "Failed to create session", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, homePath, http.StatusSeeOther)
}

func (h *LoginHandler) authenticateUser(credentials, password string) (*domain.User, error) {
	user, err := h.userService.Login(credentials, password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	return user, nil
}

func (h *LoginHandler) createUserSession(w http.ResponseWriter, user *domain.User) error {
	// Clear existing sessions
	if err := h.sessionService.DeleteByUserID(user.ID); err != nil {
		return err
	}

	// Create new session
	session, err := h.sessionService.Create(user.ID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Should be true in production with HTTPS
		Expires:  session.ExpiresAt,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func (h *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookie)
	if err == nil && cookie.Value != "" {
		_ = h.sessionService.Delete(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, homePath, http.StatusFound)
}

func (h *LoginHandler) renderLogin(w http.ResponseWriter, data *LoginTemplateData) {
	if err := h.tmpl.ExecuteTemplate(w, loginTemplate, data); err != nil {
		h.errorHandler.HandleError(w, "Failed to render login page", err, http.StatusInternalServerError)
	}
}
