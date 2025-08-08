package auth

import (
	"errors"
	"forum/internal/domain"
	"forum/internal/presentation/html/errorhandler"
	"forum/internal/service/user"
	"forum/internal/service/user/validator"
	"html/template"
	"net/http"
	"strings"
)

const (
	registerTemplate    = "register.html"
	loginPath           = "/login"
	usernameField       = "username"
	emailField          = "email"
	passwordField       = "password"
	repeatPasswordField = "repeat-password"
	registeredParam     = "registered"
)

type RegisterHandler struct {
	tmpl         *template.Template
	userService  user.Service
	errorHandler errorhandler.Handler
}

func NewRegisterHandler(
	tmpl *template.Template,
	userService user.Service,
	errorHandler errorhandler.Handler,
) *RegisterHandler {
	return &RegisterHandler{
		tmpl:         tmpl,
		userService:  userService,
		errorHandler: errorHandler,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetRegister(w, r)
	case http.MethodPost:
		h.handlePostRegister(w, r)
	default:
		h.errorHandler.HandleError(w, "Method not allowed", nil, http.StatusMethodNotAllowed)
	}
}

type RegisterFormData struct {
	Error   string
	Success string
	Form    map[string]string
}

func (h *RegisterHandler) handleGetRegister(w http.ResponseWriter, r *http.Request) {
	data := &RegisterFormData{}
	if r.URL.Query().Get(registeredParam) == "true" {
		data.Success = "Registration successful! Please log in."
	}
	h.renderRegisterForm(w, data)
}

func (h *RegisterHandler) handlePostRegister(w http.ResponseWriter, r *http.Request) {
	formData, err := h.parseRegisterForm(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.renderRegisterForm(w, &RegisterFormData{
			Error: "Invalid form data",
		})
		return
	}

	if err := h.validateRegisterForm(formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.renderRegisterForm(w, &RegisterFormData{
			Error: err.Error(),
			Form:  formData,
		})
		return
	}

	user := h.createUserFromForm(formData)
	if err := h.registerUser(user); err != nil {
		w.WriteHeader(httpStatusFromError(err))
		h.renderRegisterForm(w, &RegisterFormData{
			Error: "Registration failed: " + err.Error(),
			Form:  formData,
		})
		return
	}

	http.Redirect(w, r, loginPath+"?"+registeredParam+"=true", http.StatusSeeOther)
}

func (h *RegisterHandler) parseRegisterForm(r *http.Request) (map[string]string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return map[string]string{
		usernameField:       strings.TrimSpace(r.FormValue(usernameField)),
		emailField:          strings.TrimSpace(r.FormValue(emailField)),
		passwordField:       r.FormValue(passwordField),
		repeatPasswordField: r.FormValue(repeatPasswordField),
	}, nil
}

func (h *RegisterHandler) validateRegisterForm(formData map[string]string) error {
	if formData[passwordField] != formData[repeatPasswordField] {
		return domain.ErrPasswordsNotMatch
	}

	if formData[usernameField] == "" {
		return domain.ErrUsernameRequired
	}

	if formData[emailField] == "" {
		return domain.ErrEmailRequired
	}

	if len(formData[passwordField]) < 8 {
		return domain.ErrPasswordTooShort
	}

	return nil
}

func (h *RegisterHandler) createUserFromForm(formData map[string]string) *domain.User {
	return &domain.User{
		Username:     formData[usernameField],
		Email:        formData[emailField],
		PasswordHash: formData[passwordField],
	}
}

func (h *RegisterHandler) registerUser(user *domain.User) error {
	return h.userService.RegisterUser(user)
}

func (h *RegisterHandler) renderRegisterForm(w http.ResponseWriter, data *RegisterFormData) {
	if err := h.tmpl.ExecuteTemplate(w, registerTemplate, data); err != nil {
		h.errorHandler.HandleError(w, "Failed to render registration form", err, http.StatusInternalServerError)
	}
}

func httpStatusFromError(err error) int {
	switch {
	case errors.Is(err, validator.ErrInvalidEmail),
		errors.Is(err, validator.ErrEmptyUsername),
		errors.Is(err, validator.ErrEmptyPassword),
		errors.Is(err, validator.ErrPasswordTooShort),
		errors.Is(err, validator.ErrPasswordTooWeak):
		return http.StatusBadRequest

	case errors.Is(err, user.ErrEmailTaken),
		errors.Is(err, user.ErrUsernameTaken):
		return http.StatusConflict

	case errors.Is(err, domain.ErrInsertUserFailed),
		errors.Is(err, domain.ErrQueryFailed),
		errors.Is(err, domain.ErrUUIDParseFailed):
		return http.StatusInternalServerError

	default:
		return http.StatusInternalServerError
	}
}
