package errorhandler

import (
	"bytes"
	"forum/pkg/logger"
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl   *template.Template
	logger logger.Logger
}

func NewErrorHandler(tmpl *template.Template, logger logger.Logger) *Handler {
	return &Handler{
		tmpl:   tmpl,
		logger: logger,
	}
}

func (h *Handler) HandleError(w http.ResponseWriter, errorMsg string, errLog error, statusCode int) {
	var buf bytes.Buffer

	data := struct {
		ErrorCode int
		ErrorMsg  string
	}{
		ErrorCode: statusCode,
		ErrorMsg:  errorMsg,
	}

	if err := h.tmpl.ExecuteTemplate(&buf, "error.html", data); err != nil {
		h.logger.Errorf("failed to render error page: %v", err)
		http.Error(w, "Error rendering error page", http.StatusInternalServerError)
		return
	}

	if errLog != nil {
		h.logger.Error(errLog.Error())
	}

	w.WriteHeader(statusCode)
	if _, err := buf.WriteTo(w); err != nil {
		h.logger.Errorf("failed to write response: %v", err)
	}
}
