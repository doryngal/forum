package auth

import (
	"forum/internal/domain"
	"forum/internal/service/user"
	"html/template"
	"net/http"
	"strings"
)

type RegisterHandler struct {
	tmpl        *template.Template
	userService user.Service
}

func NewRegisterHandler(tmpl *template.Template, us user.Service) *RegisterHandler {
	return &RegisterHandler{
		tmpl:        tmpl,
		userService: us,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.renderRegisterForm(w, nil)
	case http.MethodPost:
		h.handleRegister(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type RegisterFormData struct {
	Error   string
	Success string
	Form    map[string]string
}

func (h *RegisterHandler) renderRegisterForm(w http.ResponseWriter, data *RegisterFormData) {
	if data == nil {
		data = &RegisterFormData{}
	}
	err := h.tmpl.ExecuteTemplate(w, "register.html", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h *RegisterHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderRegisterForm(w, &RegisterFormData{
			Error: "Invalid form data",
		})
		return
	}

	formData := map[string]string{
		"username":        strings.TrimSpace(r.FormValue("username")),
		"email":           strings.TrimSpace(r.FormValue("email")),
		"password":        r.FormValue("password"),
		"repeat-password": r.FormValue("repeat-password"),
	}

	// Validate form data
	if formData["password"] != formData["repeat-password"] {
		h.renderRegisterForm(w, &RegisterFormData{
			Error: "Passwords do not match",
			Form:  formData,
		})
		return
	}

	user := &domain.User{
		Username:     formData["username"],
		Email:        formData["email"],
		PasswordHash: formData["password"],
	}

	err := h.userService.RegisterUser(user)
	if err != nil {
		h.renderRegisterForm(w, &RegisterFormData{
			Error: "Registration failed: " + err.Error(),
			Form:  formData,
		})
		return
	}

	http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
}
