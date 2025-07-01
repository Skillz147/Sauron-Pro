package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"o365/utils"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	TurnstileVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	TurnstileSecret    = "0x4AAAAAABh4ix44s4LMgjPejfwrj6rKHwk" // ✅ hardcoded
)

type turnstileResponse struct {
	Success bool `json:"success"`
}

func HandleTurnstileVerify(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	token := r.FormValue("cf-turnstile-response")
	ret := r.FormValue("return")
	if ret == "" {
		ret = "/"
	}

	resp, err := http.PostForm(TurnstileVerifyURL, url.Values{
		"secret":   {TurnstileSecret},
		"response": {token},
		"remoteip": {r.RemoteAddr},
	})
	log.Info().Str("key", TurnstileSecret).Msg("🔑 Using Turnstile secret")

	if err != nil {
		log.Warn().Err(err).Msg("Turnstile request failed")
		http.Error(w, "verification failed", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	var result turnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Success {
		log.Warn().Err(err).Msg("Turnstile decode or failure")
		http.Error(w, "invalid token", http.StatusForbidden)
		return
	}

	// ✅ Generate and store session token
	sessionToken := uuid.NewString()
	err = utils.RedisClient.Set(context.Background(), "turnstile:"+sessionToken, "ok", 5*time.Minute).Err()
	if err != nil {
		log.Error().Err(err).Msg("Redis set failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// ✅ Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   300, // 5 minutes
	})

	http.Redirect(w, r, ret, http.StatusSeeOther)
}
