package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message"`
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
				message = "Internal Server Error"
			}

			w.Header().Set("Content-Type", "application/json") // Set content type to JSON
			w.WriteHeader(statusCode)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: message}); err != nil {
				slog.Error("Error encoding message", "err", err)
			}
		}
	}
}
