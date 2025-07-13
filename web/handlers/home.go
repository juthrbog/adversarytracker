package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
)

// Home handles the root path and renders the welcome page
func Home(w http.ResponseWriter, r *http.Request) {
	// Parse templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "home.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
