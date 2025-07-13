package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/juthrbog/adversarytracker/db"
	"github.com/juthrbog/adversarytracker/internal/app"
)

// EncounterRoutes returns a router with all encounter routes
func EncounterRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ListEncounters)
	r.Get("/new", NewEncounterForm)
	r.Post("/", CreateEncounter)
	
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", ViewEncounter)
		r.Get("/edit", EditEncounterForm)
		r.Post("/", UpdateEncounter)
		r.Delete("/", DeleteEncounter)
		r.Post("/delete", DeleteEncounter) // For form submissions
		
		// Adversary management within encounter
		r.Post("/adversaries", AddAdversaryToEncounter)
		r.Delete("/adversaries/{adversaryId}", RemoveAdversaryFromEncounter)
		r.Post("/adversaries/{adversaryId}/delete", RemoveAdversaryFromEncounter) // For form submissions
	})

	// HTMX specific routes
	r.Get("/add-adversary/{adversaryId}", AddAdversaryModal)

	return r
}

// ListEncounters displays all encounters
func ListEncounters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all encounters from the database
	encounters, err := db.GetAllEncounters(ctx, app.DB)
	if err != nil {
		slog.Error("Failed to get encounters", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "encounters", "list.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	data := map[string]interface{}{
		"Encounters": encounters,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ViewEncounter displays a single encounter
func ViewEncounter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get encounter ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	// Get encounter from database
	encounter, err := db.GetEncounterByID(ctx, app.DB, id)
	if err != nil {
		slog.Error("Failed to get encounter", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if encounter == nil {
		http.Error(w, "Encounter not found", http.StatusNotFound)
		return
	}

	// Parse templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "encounters", "view.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	data := map[string]interface{}{
		"Encounter": encounter,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewEncounterForm displays the form to create a new encounter
func NewEncounterForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all adversaries for selection
	adversaries, err := db.GetAllAdversaries(ctx, app.DB)
	if err != nil {
		slog.Error("Failed to get adversaries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "encounters", "form.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template with empty encounter for the form
	data := map[string]interface{}{
		"Encounter":   &db.Encounter{},
		"Adversaries": adversaries,
		"IsNew":       true,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateEncounter handles the form submission to create a new encounter
func CreateEncounter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Create encounter from form data
	enc := &db.Encounter{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	// Save to database
	id, err := db.CreateEncounter(ctx, app.DB, enc)
	if err != nil {
		slog.Error("Failed to create encounter", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX, redirect via response headers
		w.Header().Set("HX-Redirect", "/encounters/"+strconv.FormatInt(id, 10))
		return
	}

	// Regular form submission, redirect to the new encounter
	http.Redirect(w, r, "/encounters/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

// EditEncounterForm displays the form to edit an existing encounter
func EditEncounterForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get encounter ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	// Get encounter from database
	encounter, err := db.GetEncounterByID(ctx, app.DB, id)
	if err != nil {
		slog.Error("Failed to get encounter", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if encounter == nil {
		http.Error(w, "Encounter not found", http.StatusNotFound)
		return
	}

	// Get all adversaries for selection
	adversaries, err := db.GetAllAdversaries(ctx, app.DB)
	if err != nil {
		slog.Error("Failed to get adversaries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "encounters", "form.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	data := map[string]interface{}{
		"Encounter":   encounter,
		"Adversaries": adversaries,
		"IsNew":       false,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateEncounter handles the form submission to update an existing encounter
func UpdateEncounter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get encounter ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Create encounter from form data
	enc := &db.Encounter{
		ID:          id,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	// Update in database
	err = db.UpdateEncounter(ctx, app.DB, enc)
	if err != nil {
		slog.Error("Failed to update encounter", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX, redirect via response headers
		w.Header().Set("HX-Redirect", "/encounters/"+idStr)
		return
	}

	// Regular form submission, redirect to the encounter
	http.Redirect(w, r, "/encounters/"+idStr, http.StatusSeeOther)
}

// DeleteEncounter handles the deletion of an encounter
func DeleteEncounter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get encounter ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	// Delete from database
	err = db.DeleteEncounter(ctx, app.DB, id)
	if err != nil {
		slog.Error("Failed to delete encounter", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX, redirect via response headers
		w.Header().Set("HX-Redirect", "/encounters")
		return
	}

	// Regular form submission, redirect to the encounter list
	http.Redirect(w, r, "/encounters", http.StatusSeeOther)
}

// AddAdversaryModal displays a modal for adding an adversary to an encounter
func AddAdversaryModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get adversary ID from URL
	adversaryIdStr := chi.URLParam(r, "adversaryId")
	adversaryId, err := strconv.ParseInt(adversaryIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid adversary ID", http.StatusBadRequest)
		return
	}

	// Get adversary from database
	adversary, err := db.GetAdversaryByID(ctx, app.DB, adversaryId)
	if err != nil {
		slog.Error("Failed to get adversary", "error", err, "id", adversaryId)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if adversary == nil {
		http.Error(w, "Adversary not found", http.StatusNotFound)
		return
	}

	// Get all encounters for selection
	encounters, err := db.GetAllEncounters(ctx, app.DB)
	if err != nil {
		slog.Error("Failed to get encounters", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse template
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "encounters", "add_adversary_modal.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	data := map[string]interface{}{
		"Adversary":  adversary,
		"Encounters": encounters,
	}

	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// AddAdversaryToEncounter handles adding an adversary to an encounter
func AddAdversaryToEncounter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get encounter ID from URL
	encounterIdStr := chi.URLParam(r, "id")
	encounterId, err := strconv.ParseInt(encounterIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get adversary ID and count from form
	adversaryIdStr := r.FormValue("adversary_id")
	adversaryId, err := strconv.ParseInt(adversaryIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid adversary ID", http.StatusBadRequest)
		return
	}

	countStr := r.FormValue("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 1 // Default to 1 if invalid
	}

	// Begin transaction
	tx, err := app.DB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Add adversary to encounter
	ea := &db.EncounterAdversary{
		EncounterID: encounterId,
		AdversaryID: adversaryId,
		Count:       count,
	}

	_, err = db.AddAdversaryToEncounter(ctx, tx, ea)
	if err != nil {
		slog.Error("Failed to add adversary to encounter", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX, redirect via response headers
		w.Header().Set("HX-Redirect", "/encounters/"+encounterIdStr)
		return
	}

	// Regular form submission, redirect to the encounter
	http.Redirect(w, r, "/encounters/"+encounterIdStr, http.StatusSeeOther)
}

// RemoveAdversaryFromEncounter handles removing an adversary from an encounter
func RemoveAdversaryFromEncounter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get encounter ID from URL
	encounterIdStr := chi.URLParam(r, "id")
	_, err := strconv.ParseInt(encounterIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid encounter ID", http.StatusBadRequest)
		return
	}

	// Get adversary ID from URL (this is actually the encounter_adversaries.id)
	adversaryIdStr := chi.URLParam(r, "adversaryId")
	encounterAdversaryID, err := strconv.ParseInt(adversaryIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid adversary ID", http.StatusBadRequest)
		return
	}

	// Remove adversary from encounter
	err = db.RemoveAdversaryFromEncounter(ctx, app.DB, encounterAdversaryID)
	if err != nil {
		slog.Error("Failed to remove adversary from encounter", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX, redirect via response headers
		w.Header().Set("HX-Redirect", "/encounters/"+encounterIdStr)
		return
	}

	// Regular form submission, redirect to the encounter
	http.Redirect(w, r, "/encounters/"+encounterIdStr, http.StatusSeeOther)
}
