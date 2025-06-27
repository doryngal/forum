package auth

import (
	"forum/internal/domain"
	"forum/internal/service/user"
	"html/template"
	"net/http"
)

type LoginHandler struct {
	tmpl        *template.Template
	userService user.Service
}

func NewLoginHandler(tmpl *template.Template, userService user.Service) *LoginHandler {
	return &LoginHandler{
		tmpl:        tmpl,
		userService: userService,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.renderLogin(w, nil)
	case http.MethodPost:
		h.handleLogin(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *LoginHandler) renderLogin(w http.ResponseWriter, data interface{}) {
	err := h.tmpl.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h *LoginHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLogin(w, TemplateData{Error: "Invalid form data"})
		return
	}

	credentials := domain.LoginCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := h.userService.Login(credentials.Email, credentials.Password)
	if err != nil {
		h.renderLogin(w, TemplateData{Error: "Invalid email or password"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    user.ID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // HTTPS
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type TemplateData struct {
	Error string
}
