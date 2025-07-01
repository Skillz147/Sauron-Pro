package handlers

import (
	"encoding/json"
	"net/http"
	"o365/configdb"
	"o365/utils"
	"strings"

	"github.com/rs/zerolog/log"
)

type BotReport struct {
	Reason    string `json:"reason"`
	UserAgent string `json:"ua"`
	Timestamp string `json:"ts"`
	Slug      string `json:"slug"` // new field (optional)
}

func HandleJSCheck(w http.ResponseWriter, r *http.Request) {
	ip := utils.GetRealIP(r)
	remote := r.RemoteAddr
	ip = strings.TrimPrefix(ip, "::ffff:")

	log.Info().Str("real_ip", ip).Str("remote", remote).Msg("JSCheck incoming")

	// Bypass for localhost
	localAddrs := map[string]bool{
		"127.0.0.1": true,
		"::1":       true,
		"localhost": true,
	}
	if localAddrs[ip] {
		log.Info().Str("ip", ip).Msg("🛡️ Localhost JS check bypassed")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	// Decode the BotReport from the request body
	var report BotReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		log.Warn().Err(err).Msg("⚠️ Invalid bot report JSON")
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Extract the slug from the request context
	slug, ok := r.Context().Value("slug").(string)
	if !ok {
		http.Error(w, "Slug not found in request context", http.StatusBadRequest)
		return
	}

	// Log the bot report, now using the extracted slug
	userID := ""
	if validSlug, ok := configdb.IsValidSlug(slug); ok {
		// Optional: Get the user ID if the slug is valid
		userID = validSlug
	}

	log.Warn().
		Str("ip", ip).
		Str("ua", report.UserAgent).
		Str("reason", report.Reason).
		Str("slug", slug). // Use the slug from the request context
		Str("user", userID).
		Msg("🤖 Headless/browser automation detected")

	// Mark the IP as banned for bot activity
	configdb.SaveBannedIP(ip, "bot-headless")

	// Redirect the user to a neutral page
	http.Redirect(w, r, "https://www.google.com", http.StatusFound)
}
