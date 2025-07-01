// slug/slug.go
package slug

import (
	"fmt"
	"net/http"
	"o365/configdb"
	"regexp"
	"strings"
)

/* ── regexes ────────── */

var shortSlug = regexp.MustCompile(`^[A-Za-z0-9_-]{4,16}$`)
var uuidSlug = regexp.MustCompile(`^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$`)

func looksValid(s string) bool { return shortSlug.MatchString(s) || uuidSlug.MatchString(s) }

/* ── central helper ─── */

func getValidSlug(raw string) (string, bool) {
	if !looksValid(raw) {
		return "", false
	}
	if _, ok := configdb.IsValidSlug(raw); ok {
		return raw, true
	}
	return "", false
}

/* ── exported ───────── */

// GetSlugFromRequest now checks:  path ▶ query ▶ cookie
func GetSlugFromRequest(r *http.Request) (string, error) {
	// 1️⃣ first path seg
	if seg, ok := getValidSlug(strings.Split(strings.Trim(r.URL.Path, "/"), "/")[0]); ok {
		return seg, nil
	}
	// 2️⃣ ?slug=…
	if q, ok := getValidSlug(r.URL.Query().Get("slug")); ok {
		return q, nil
	}
	// 3️⃣ cookie set by our middleware
	if c, err := r.Cookie("o365_slug"); err == nil {
		if s, ok := getValidSlug(c.Value); ok {
			return s, nil
		}
	}
	return "", fmt.Errorf("slug not found or invalid")
}

// DB passthrough
func GetSlugForUser(userID string) (string, error) { return configdb.GetSlugForUser(userID) }
