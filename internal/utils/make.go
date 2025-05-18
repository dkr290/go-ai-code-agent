package utils

import (
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dkr290/go-ai-code-agent/internal/templates"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statuscode"`
}

func MakeHandlers(
	handlerFunc func(w http.ResponseWriter, r *http.Request) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handlerFunc(w, r); err != nil {
			slog.Error("Internal server error", "err", err, "path", r.URL.Path)

			var statusCode int
			var message string

			switch {
			case strings.Contains(err.Error(), "method not allowed"):
				statusCode = http.StatusMethodNotAllowed
				message = "Method Not Allowed"
			case strings.Contains(err.Error(), "error parsing form"):
				statusCode = http.StatusBadRequest
				message = "Error Parsing Web Form"
			case strings.Contains(err.Error(), "Error calling AI agent"):
				statusCode = http.StatusInternalServerError
				message = "error calling AI agent"

			default:
				statusCode = http.StatusInternalServerError
				message = "Internal Server Error" + err.Error()
			}

			RenderError(w, statusCode, message)
		}
	}
}

func RenderError(w http.ResponseWriter, status int, message string) {
	tmpl, err := templates.LoadTemplates("web/template/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	data := struct {
		Flash ErrorResponse
	}{
		Flash: ErrorResponse{
			Message:    message,
			StatusCode: status,
		},
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		slog.Error("Failed to render error template", "err", err)
		http.Error(w, message, status)
	}
}
