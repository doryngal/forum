package errorhandler

import (
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl *template.Template
}

func NewErrorHandler(tmpl *template.Template) *Handler {
	return &Handler{
		tmpl: tmpl,
	}
}

func (h *Handler) HandleError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)

	data := struct {
		ErrorCode int
		ErrorMsg  string
	}{
		ErrorCode: statusCode,
		ErrorMsg:  errorMsg,
	}

	if err := h.tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, "Error rendering error page", http.StatusInternalServerError)
	}
}
