package honeypot

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// HoneypotManager manages the honeypot system for invalid slug access attempts
type HoneypotManager struct {
	rateLimiter *RateLimiter
	fail2ban    *Fail2BanManager
	templates   *HoneypotTemplate
	uuidPattern *regexp.Regexp
}

// NewHoneypotManager creates a new honeypot manager with all components
func NewHoneypotManager() *HoneypotManager {
	return &HoneypotManager{
		rateLimiter: NewRateLimiter(10*time.Minute, 10), // 10 attempts per 10 minutes
		fail2ban:    NewFail2BanManager(),
		templates:   NewHoneypotTemplate(),
		uuidPattern: regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
	}
}

// Global honeypot instance
var GlobalHoneypot = NewHoneypotManager()

// ProcessInvalidSlug handles invalid slug attempts with professional enumeration detection
func (h *HoneypotManager) ProcessInvalidSlug(w http.ResponseWriter, r *http.Request, invalidSlug, clientIP, userAgent string) {
	// Secure sanitized logging - leak-proof and spoof-proof
	GlobalSecureLogger.LogEnumerationAttempt(r, invalidSlug, clientIP, "initial_check", h.rateLimiter.GetAttemptCount(clientIP))

	// Advanced enumeration pattern detection
	enumerationType := h.detectEnumerationType(invalidSlug)

	if enumerationType != "none" {
		// Systematic enumeration detected - immediate ban
		reason := fmt.Sprintf("%s enumeration pattern detected", enumerationType)
		h.executeImmediateBan(clientIP, userAgent, invalidSlug, reason)

		// Serve minimal response to automated tools
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Access Denied"))
		return
	}

	// Rate limiting check with progressive penalties
	if h.rateLimiter.CheckAndRecord(clientIP) {
		// Rate limit exceeded - ban with detailed reason
		reason := fmt.Sprintf("Rate limit exceeded: %d invalid attempts", h.rateLimiter.GetAttemptCount(clientIP))
		h.executeImmediateBan(clientIP, userAgent, invalidSlug, reason)

		w.WriteHeader(http.StatusTooManyRequests)
		w.Header().Set("Retry-After", "3600")
		w.Write([]byte("Too Many Requests"))
		return
	}

	// Serve realistic honeypot page for intelligence gathering
	h.serveIntelligenceHoneypot(w, r, invalidSlug, clientIP, enumerationType)
}

// detectEnumerationType analyzes slug patterns to identify enumeration attacks
func (h *HoneypotManager) detectEnumerationType(invalidSlug string) string {
	// UUID pattern detection (most critical for our use case)
	if h.uuidPattern.MatchString(invalidSlug) {
		return "uuid"
	}

	// Secondary UUID validation (loose format check)
	if len(invalidSlug) == 36 && strings.Count(invalidSlug, "-") == 4 {
		parts := strings.Split(invalidSlug, "-")
		if len(parts) == 5 &&
			len(parts[0]) == 8 &&
			len(parts[1]) == 4 &&
			len(parts[2]) == 4 &&
			len(parts[3]) == 4 &&
			len(parts[4]) == 12 {
			return "uuid-format"
		}
	}

	// Common enumeration wordlists
	enumerationPatterns := map[string][]string{
		"admin":  {"admin", "administrator", "root", "sa", "sysadmin"},
		"auth":   {"auth", "login", "signin", "sso", "oauth", "auth2", "authentication"},
		"api":    {"api", "rest", "graphql", "endpoint", "service", "webhook"},
		"test":   {"test", "testing", "qa", "dev", "demo", "sandbox", "trial"},
		"config": {"config", "configuration", "settings", "setup", "install", "init"},
		"portal": {"portal", "dashboard", "panel", "console", "interface", "ui"},
		"backup": {"backup", "bak", "old", "temp", "tmp", "cache", "archive"},
		"debug":  {"debug", "trace", "log", "monitor", "health", "status", "ping"},
	}

	slugLower := strings.ToLower(invalidSlug)
	for category, patterns := range enumerationPatterns {
		for _, pattern := range patterns {
			if slugLower == pattern || strings.Contains(slugLower, pattern) {
				return category
			}
		}
	}

	// Sequential enumeration detection (admin1, admin2, test001, etc.)
	if len(invalidSlug) > 1 {
		// Check for numeric suffixes
		re := regexp.MustCompile(`^([a-zA-Z]+)(\d+)$`)
		if matches := re.FindStringSubmatch(invalidSlug); len(matches) == 3 {
			base := strings.ToLower(matches[1])
			for category, patterns := range enumerationPatterns {
				for _, pattern := range patterns {
					if base == pattern {
						return fmt.Sprintf("sequential-%s", category)
					}
				}
			}
		}
	}

	// Random string detection (likely brute force)
	if len(invalidSlug) >= 6 && h.looksLikeRandomString(invalidSlug) {
		return "brute-force"
	}

	return "none"
}

// looksLikeRandomString detects if a string appears to be randomly generated
func (h *HoneypotManager) looksLikeRandomString(s string) bool {
	// Simple entropy check - mixed case, numbers, no dictionary words
	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	// Random strings typically have mixed case and numbers
	return hasUpper && hasLower && hasDigit
}

// sanitizeSlugForLogging safely sanitizes slugs for logging without revealing actual values
func (h *HoneypotManager) sanitizeSlugForLogging(slug string) string {
	enumerationType := h.detectEnumerationType(slug)

	// For UUID patterns, show masked version like da92fa51-****-4228-****-e873f43f8c8f
	if enumerationType == "uuid" || enumerationType == "uuid-format" {
		if len(slug) == 36 && strings.Count(slug, "-") == 4 {
			// UUID format: 8-4-4-4-12
			parts := strings.Split(slug, "-")
			if len(parts) == 5 {
				return fmt.Sprintf("%s-****-%s-****-%s", parts[0], parts[2], parts[4])
			}
		}
		return fmt.Sprintf("uuid_pattern_len_%d", len(slug))
	}

	// For known enumeration patterns, show category and first few chars
	if enumerationType != "none" && enumerationType != "brute-force" {
		if len(slug) > 3 {
			return fmt.Sprintf("%s_pattern_%s***", enumerationType, slug[:3])
		}
		return fmt.Sprintf("%s_pattern_%s", enumerationType, slug)
	}

	// For brute force or unknown, just show length and character types
	charTypes := ""
	if regexp.MustCompile(`[a-z]`).MatchString(slug) {
		charTypes += "lower"
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(slug) {
		charTypes += "upper"
	}
	if regexp.MustCompile(`[0-9]`).MatchString(slug) {
		charTypes += "digit"
	}
	if regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(slug) {
		charTypes += "special"
	}

	return fmt.Sprintf("unknown_len_%d_types_%s", len(slug), charTypes)
}

// executeImmediateBan performs immediate banning for detected enumeration attacks
func (h *HoneypotManager) executeImmediateBan(clientIP, userAgent, invalidSlug, reason string) {
	// Secure sanitized logging for immediate ban
	GlobalSecureLogger.LogImmediateBan(clientIP, reason, h.detectEnumerationType(invalidSlug))

	// Execute ban in secure database if available
	// For now, use fail2ban as primary ban mechanism - sanitize slug to prevent leaks
	sanitizedSlug := h.sanitizeSlugForLogging(invalidSlug)
	fullReason := fmt.Sprintf("Honeypot enumeration: %s | Pattern: %s", reason, sanitizedSlug)
	h.fail2ban.TriggerBan(clientIP, "slug-enumeration", fullReason)
}

// serveIntelligenceHoneypot serves realistic business websites for intelligence gathering
func (h *HoneypotManager) serveIntelligenceHoneypot(w http.ResponseWriter, r *http.Request, invalidSlug, clientIP, enumerationType string) {
	// Choose business website template based on user agent and enumeration patterns
	templateName := "techstartup" // Default to tech startup - appears most legitimate

	userAgent := strings.ToLower(r.UserAgent())
	if strings.Contains(userAgent, "security") || strings.Contains(userAgent, "scan") || strings.Contains(userAgent, "bot") {
		templateName = "cybersec" // Show cybersecurity site to security scanners - ironic honeypot
	} else if strings.Contains(userAgent, "marketing") || strings.Contains(userAgent, "social") || strings.Contains(userAgent, "seo") {
		templateName = "marketing" // Show marketing agency to marketing-related bots
	}

	// Secure sanitized logging for intelligence gathering
	GlobalSecureLogger.LogHoneypotServed(r, invalidSlug, clientIP, templateName)

	// Add realistic timing delay based on slug complexity
	delay := time.Duration(200+len(invalidSlug)*10) * time.Millisecond
	time.Sleep(delay)

	// Serve the realistic business website
	h.templates.ServeBusinessSite(w, r, templateName)
}

// Legacy compatibility methods for backward compatibility with existing code
func (h *HoneypotManager) isSlugEnumeration(invalidSlug string) bool {
	return h.detectEnumerationType(invalidSlug) != "none"
}

func (h *HoneypotManager) banIP(clientIP, userAgent, invalidSlug, reason string) {
	h.executeImmediateBan(clientIP, userAgent, invalidSlug, reason)
}

func (h *HoneypotManager) serveHoneypotPage(w http.ResponseWriter, r *http.Request, invalidSlug, clientIP string) {
	h.serveIntelligenceHoneypot(w, r, invalidSlug, clientIP, "legacy")
}
