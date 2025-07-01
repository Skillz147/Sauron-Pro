package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

type SessionData struct {
	UserID  string `json:"user_id"`
	KeyHash string `json:"key_hash"`
}

func GetSessionFromCookie(r *http.Request) (SessionData, error) {
	cookie, err := r.Cookie("sauron_session")
	if err != nil || cookie.Value == "" {
		return SessionData{}, errors.New("missing session")
	}

	var parsed SessionData
	if err := json.Unmarshal([]byte(cookie.Value), &parsed); err != nil {
		return SessionData{}, errors.New("invalid session format")
	}

	return parsed, nil
}
