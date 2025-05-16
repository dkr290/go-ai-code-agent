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
	modelName   string
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
func (h *AppHandler) IndexHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		err := h.templates.ExecuteTemplate(w, "index.html", nil) // You might pass data here
		if err != nil {
			return fmt.Errorf("erro partsing templates %v", err)
		}
	}
	return nil
}

func (h *AppHandler) CallAgentHandler(w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return fmt.Errorf("error parsing form %v", err)
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
			modelName:   r.FormValue("model"),
		}

		// Call the AI agent with all the parameters
		if err := run(ctx, &params); err != nil {
			return fmt.Errorf("error calling AI agent %s", err)
		}

		// For this example, we'll just send the AI response back as plain text
		// HTMX will then update the #ai_response div with this content.
		_, _ = fmt.Fprintf(w, "%s", "OK")
		return nil
	}
	return fmt.Errorf("method not allowed")
}
