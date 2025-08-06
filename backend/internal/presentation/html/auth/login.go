package auth

import (
	"forum/internal/presentation/html/errorhandler"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"html/template"
	"net/http"
)

type LoginHandler struct {
	tmpl           *template.Template
	userService    user.Service
	sessionService session.Service
	errorHandler   errorhandler.Handler
}

func NewLoginHandler(tmpl *template.Template, userService user.Service, sessionService session.Service, errorHandler errorhandler.Handler) *LoginHandler {
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
		h.renderLogin(w, nil)
	case http.MethodPost:
		h.handleLogin(w, r)
	default:
		h.errorHandler.HandleError(w, "Method not allowed", nil, http.StatusMethodNotAllowed)
	}
}

func (h *LoginHandler) renderLogin(w http.ResponseWriter, data *TemplateData) {
	if data == nil {
		data = &TemplateData{}
	}

	if err := h.tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		h.errorHandler.HandleError(w, "Failed to render login page", err, http.StatusInternalServerError)
	}
}

func (h *LoginHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLogin(w, &TemplateData{Error: "Invalid form submission"})
		return
	}

	email := r.FormValue("emailOrUsername")
	password := r.FormValue("password")

	user, err := h.userService.Login(email, password)
	if err != nil {
		h.renderLogin(w, &TemplateData{Error: "Incorrect email or password"})
		return
	}

	h.sessionService.DeleteByUserID(user.ID)

	// Create session
	session, err := h.sessionService.Create(user.ID)
	if err != nil {
		h.errorHandler.HandleError(w, "Failed to create session", err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true for HTTPS
		Expires:  session.ExpiresAt,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		_ = h.sessionService.Delete(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

type TemplateData struct {
	Error string
}
