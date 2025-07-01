package proxy

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"o365/utils"
)

// LogProxyTraffic logs incoming MITM request details
func LogProxyTraffic(r *http.Request, status int, size int) {
	ip := getIP(r)
	method := r.Method
	path := r.URL.Path
	userAgent := r.UserAgent()
	timestamp := time.Now().Format(time.RFC3339)

	line := fmt.Sprintf(`[MITM] %s | %s | %s | %d bytes | %d | %s | %s`,
		timestamp, ip, method+" "+path, size, status, userAgent, r.Host)

	utils.AccessLogger.Info().
		Str("ip", ip).
		Str("method", method).
		Str("path", path).
		Int("size", size).
		Int("status", status).
		Str("ua", userAgent).
		Str("host", r.Host).
		Msg(line)
}

// getIP returns the client IP address from headers or RemoteAddr
func getIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	ip := strings.Split(r.RemoteAddr, ":")
	if len(ip) > 0 {
		return ip[0]
	}
	return r.RemoteAddr
}
