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

// AdversaryRoutes returns a router with all adversary routes
func AdversaryRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ListAdversaries)
	r.Get("/new", NewAdversaryForm)
	r.Post("/", CreateAdversary)
	
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", ViewAdversary)
		r.Get("/edit", EditAdversaryForm)
		r.Post("/", UpdateAdversary)
		r.Delete("/", DeleteAdversary)
		// HTMX specific route for deletion with POST
		r.Post("/delete", DeleteAdversary)
	})

	return r
}

// ListAdversaries displays all adversaries
func ListAdversaries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get all adversaries from the database
	adversaries, err := db.GetAllAdversaries(ctx, app.DB)
	if err != nil {
		slog.Error("Failed to get adversaries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "adversaries", "list.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	data := map[string]interface{}{
		"Adversaries": adversaries,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ViewAdversary displays a single adversary
func ViewAdversary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get adversary ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid adversary ID", http.StatusBadRequest)
		return
	}

	// Get adversary from database
	adversary, err := db.GetAdversaryByID(ctx, app.DB, id)
	if err != nil {
		slog.Error("Failed to get adversary", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if adversary == nil {
		http.Error(w, "Adversary not found", http.StatusNotFound)
		return
	}

	// Parse templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "adversaries", "view.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	data := map[string]interface{}{
		"Adversary": adversary,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewAdversaryForm displays the form to create a new adversary
func NewAdversaryForm(w http.ResponseWriter, r *http.Request) {
	// Parse templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "adversaries", "form.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template with empty adversary for the form
	data := map[string]interface{}{
		"Adversary": &db.Adversary{},
		"IsNew":     true,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateAdversary handles the form submission to create a new adversary
func CreateAdversary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Create adversary from form data
	adv := &db.Adversary{
		Name:            r.FormValue("name"),
		Type:            r.FormValue("type"),
		ChallengeRating: r.FormValue("challenge_rating"),
		Size:            r.FormValue("size"),
		Speed:           r.FormValue("speed"),
		Abilities:       r.FormValue("abilities"),
		Actions:         r.FormValue("actions"),
		Reactions:       r.FormValue("reactions"),
		Description:     r.FormValue("description"),
	}

	// Parse numeric values
	armorClass, _ := strconv.Atoi(r.FormValue("armor_class"))
	adv.ArmorClass = armorClass

	hitPoints, _ := strconv.Atoi(r.FormValue("hit_points"))
	adv.HitPoints = hitPoints

	strength, _ := strconv.Atoi(r.FormValue("strength"))
	adv.Strength = strength

	dexterity, _ := strconv.Atoi(r.FormValue("dexterity"))
	adv.Dexterity = dexterity

	constitution, _ := strconv.Atoi(r.FormValue("constitution"))
	adv.Constitution = constitution

	intelligence, _ := strconv.Atoi(r.FormValue("intelligence"))
	adv.Intelligence = intelligence

	wisdom, _ := strconv.Atoi(r.FormValue("wisdom"))
	adv.Wisdom = wisdom

	charisma, _ := strconv.Atoi(r.FormValue("charisma"))
	adv.Charisma = charisma

	// Save to database
	id, err := db.CreateAdversary(ctx, app.DB, adv)
	if err != nil {
		slog.Error("Failed to create adversary", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX, redirect via response headers
		w.Header().Set("HX-Redirect", "/adversaries/"+strconv.FormatInt(id, 10))
		return
	}

	// Regular form submission, redirect to the new adversary
	http.Redirect(w, r, "/adversaries/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

// EditAdversaryForm displays the form to edit an existing adversary
func EditAdversaryForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get adversary ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid adversary ID", http.StatusBadRequest)
		return
	}

	// Get adversary from database
	adversary, err := db.GetAdversaryByID(ctx, app.DB, id)
	if err != nil {
		slog.Error("Failed to get adversary", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if adversary == nil {
		http.Error(w, "Adversary not found", http.StatusNotFound)
		return
	}

	// Parse templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "adversaries", "form.html"),
	)
	if err != nil {
		slog.Error("Failed to parse template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render template
	data := map[string]interface{}{
		"Adversary": adversary,
		"IsNew":     false,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateAdversary handles the form submission to update an existing adversary
func UpdateAdversary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get adversary ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid adversary ID", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Create adversary from form data
	adv := &db.Adversary{
		ID:              id,
		Name:            r.FormValue("name"),
		Type:            r.FormValue("type"),
		ChallengeRating: r.FormValue("challenge_rating"),
		Size:            r.FormValue("size"),
		Speed:           r.FormValue("speed"),
		Abilities:       r.FormValue("abilities"),
		Actions:         r.FormValue("actions"),
		Reactions:       r.FormValue("reactions"),
		Description:     r.FormValue("description"),
	}

	// Parse numeric values
	armorClass, _ := strconv.Atoi(r.FormValue("armor_class"))
	adv.ArmorClass = armorClass

	hitPoints, _ := strconv.Atoi(r.FormValue("hit_points"))
	adv.HitPoints = hitPoints

	strength, _ := strconv.Atoi(r.FormValue("strength"))
	adv.Strength = strength

	dexterity, _ := strconv.Atoi(r.FormValue("dexterity"))
	adv.Dexterity = dexterity

	constitution, _ := strconv.Atoi(r.FormValue("constitution"))
	adv.Constitution = constitution

	intelligence, _ := strconv.Atoi(r.FormValue("intelligence"))
	adv.Intelligence = intelligence

	wisdom, _ := strconv.Atoi(r.FormValue("wisdom"))
	adv.Wisdom = wisdom

	charisma, _ := strconv.Atoi(r.FormValue("charisma"))
	adv.Charisma = charisma

	// Update in database
	err = db.UpdateAdversary(ctx, app.DB, adv)
	if err != nil {
		slog.Error("Failed to update adversary", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX, redirect via response headers
		w.Header().Set("HX-Redirect", "/adversaries/"+idStr)
		return
	}

	// Regular form submission, redirect to the adversary
	http.Redirect(w, r, "/adversaries/"+idStr, http.StatusSeeOther)
}

// DeleteAdversary handles the deletion of an adversary
func DeleteAdversary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get adversary ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid adversary ID", http.StatusBadRequest)
		return
	}

	// Delete from database
	err = db.DeleteAdversary(ctx, app.DB, id)
	if err != nil {
		slog.Error("Failed to delete adversary", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX, redirect via response headers
		w.Header().Set("HX-Redirect", "/adversaries")
		return
	}

	// Regular form submission, redirect to the adversary list
	http.Redirect(w, r, "/adversaries", http.StatusSeeOther)
}
