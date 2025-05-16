package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Params struct {
	openaiKey   string
	deepSeekKey string
	geminiKey   string
	useLLM      string
	outputDir   string
	basePackage string
	language    string
	addTemplate string
	userPrompt  string
}

// AppHandler struct holds dependencies for HTTP handlers
type AppHandler struct {
	templates *template.Template
}

// NewAppHandler creates a new AppHandler instance
func NewAppHandler(tmpl *template.Template) *AppHandler {
	return &AppHandler{
		templates: tmpl,
	}
}

// IndexHandler handles the root path and renders the index.html template
func (h *AppHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := h.templates.ExecuteTemplate(w, "index.html", nil) // You might pass data here
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *AppHandler) CallAgentHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Extract all the parameters from the form
		params := Params{
			openaiKey:   r.FormValue("openai-key"),
			deepSeekKey: r.FormValue("deepseek-key"),
			geminiKey:   r.FormValue("gemini-key"),
			useLLM:      r.FormValue("use-llm"),
			outputDir:   r.FormValue("output-dir"),
			basePackage: r.FormValue("base-package"),
			language:    r.FormValue("use-language"),
			addTemplate: r.FormValue("use-template"),
			userPrompt:  r.FormValue("user-prompt"),
		}

		// Call the AI agent with all the parameters
		if err := run(ctx, &params); err != nil {
			http.Error(w, "Error calling AI agent", http.StatusInternalServerError)
			return
		}

		// For this example, we'll just send the AI response back as plain text
		// HTMX will then update the #ai_response div with this content.
		_, _ = fmt.Fprintf(w, "%s", "OK")
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
