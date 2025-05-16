package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dkr290/go-ai-code-agent/internal/handlers"
	"github.com/dkr290/go-ai-code-agent/internal/templates"
	"github.com/dkr290/go-ai-code-agent/internal/utils"
)

func main() {
	tmpl, err := templates.LoadTemplates("web/template/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	appHandler := handlers.NewAppHandler(tmpl)
	mux := http.NewServeMux()
	mux.HandleFunc("/", utils.MakeHandlers(appHandler.IndexHandler)) // For index.html
	mux.HandleFunc(
		"/agent",
		utils.MakeHandlers(appHandler.CallAgentHandler)) // call the actual ai stuff

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
