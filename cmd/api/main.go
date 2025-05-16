package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dkr290/go-ai-code-agent/internal/handlers"
	"github.com/dkr290/go-ai-code-agent/internal/templates"
)

func main() {
	tmpl, err := templates.LoadTemplates("web/template/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	// Create HTTP handlers
	appHandler := handlers.NewAppHandler(tmpl)
	// Create a new ServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/", appHandler.IndexHandler)          // For index.html
	mux.HandleFunc("/agent", appHandler.CallAgentHandler) // Endpoint to call the AI agent

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Create the HTTP server with custom configuration
	s := &http.Server{
		Addr:         ":8080",
		Handler:      mux, // Use the ServeMux
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the HTTP server
	fmt.Printf("Server listening on port %s...\n", s.Addr)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
