package honeypot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// SecureHoneypotLogger provides sanitized, leak-proof and spoof-proof logging for honeypot events
type SecureHoneypotLogger struct {
	sessionKey   []byte
	ipHashSalt   []byte
	maxLogLength int
}

// NewSecureHoneypotLogger creates a new secure logger instance
func NewSecureHoneypotLogger() *SecureHoneypotLogger {
	// Generate deterministic but secure keys for this session
	sessionKey := sha256.Sum256([]byte(fmt.Sprintf("honeypot-session-%d", time.Now().Unix()/3600))) // Changes hourly
	ipHashSalt := sha256.Sum256([]byte("honeypot-ip-salt-2024"))

	return &SecureHoneypotLogger{
		sessionKey:   sessionKey[:],
		ipHashSalt:   ipHashSalt[:],
		maxLogLength: 200, // Maximum length for any logged field
	}
}

// SanitizedAttackEvent represents a sanitized honeypot attack event
type SanitizedAttackEvent struct {
	EventID          string    `json:"event_id"`
	Timestamp        time.Time `json:"timestamp"`
	IPHash           string    `json:"ip_hash"`
	UserAgentPattern string    `json:"user_agent_pattern"`
	SlugPattern      string    `json:"slug_pattern"`
	EnumerationType  string    `json:"enumeration_type"`
	AttackSeverity   string    `json:"attack_severity"`
	GeoCountry       string    `json:"geo_country,omitempty"`
	TemplateServed   string    `json:"template_served"`
	AttemptCount     int       `json:"attempt_count"`
	ThreatLevel      int       `json:"threat_level"`
}

// LogEnumerationAttempt logs a sanitized enumeration attempt
func (s *SecureHoneypotLogger) LogEnumerationAttempt(r *http.Request, invalidSlug, clientIP, enumerationType string, attemptCount int) {
	event := SanitizedAttackEvent{
		EventID:          s.generateEventID(clientIP, time.Now()),
		Timestamp:        time.Now().UTC(),
		IPHash:           s.SanitizeIP(clientIP),
		UserAgentPattern: s.sanitizeUserAgent(r.UserAgent()),
		SlugPattern:      s.sanitizeSlug(invalidSlug),
		EnumerationType:  s.sanitizeEnumType(enumerationType),
		AttackSeverity:   s.calculateSeverity(enumerationType, attemptCount),
		GeoCountry:       s.getGeoCountry(clientIP),
		AttemptCount:     attemptCount,
		ThreatLevel:      s.calculateThreatLevel(enumerationType, attemptCount),
	}

	log.Warn().
		Str("event_id", event.EventID).
		Time("timestamp", event.Timestamp).
		Str("ip_hash", event.IPHash).
		Str("user_agent_pattern", event.UserAgentPattern).
		Str("slug_pattern", event.SlugPattern).
		Str("enumeration_type", event.EnumerationType).
		Str("attack_severity", event.AttackSeverity).
		Str("geo_country", event.GeoCountry).
		Int("attempt_count", event.AttemptCount).
		Int("threat_level", event.ThreatLevel).
		Msg("🍯 SANITIZED: Enumeration attempt detected")
}

// LogHoneypotServed logs when a honeypot template is served
func (s *SecureHoneypotLogger) LogHoneypotServed(r *http.Request, invalidSlug, clientIP, templateName string) {
	event := SanitizedAttackEvent{
		EventID:          s.generateEventID(clientIP, time.Now()),
		Timestamp:        time.Now().UTC(),
		IPHash:           s.SanitizeIP(clientIP),
		UserAgentPattern: s.sanitizeUserAgent(r.UserAgent()),
		SlugPattern:      s.sanitizeSlug(invalidSlug),
		TemplateServed:   templateName,
		GeoCountry:       s.getGeoCountry(clientIP),
		ThreatLevel:      1, // Base threat level for honeypot serving
	}

	log.Info().
		Str("event_id", event.EventID).
		Time("timestamp", event.Timestamp).
		Str("ip_hash", event.IPHash).
		Str("user_agent_pattern", event.UserAgentPattern).
		Str("slug_pattern", event.SlugPattern).
		Str("template_served", event.TemplateServed).
		Str("geo_country", event.GeoCountry).
		Int("threat_level", event.ThreatLevel).
		Msg("🍯 SANITIZED: Honeypot template served")
}

// LogImmediateBan logs a sanitized immediate ban event
func (s *SecureHoneypotLogger) LogImmediateBan(clientIP, reason, enumerationType string) {
	event := SanitizedAttackEvent{
		EventID:         s.generateEventID(clientIP, time.Now()),
		Timestamp:       time.Now().UTC(),
		IPHash:          s.SanitizeIP(clientIP),
		EnumerationType: s.sanitizeEnumType(enumerationType),
		AttackSeverity:  "CRITICAL",
		GeoCountry:      s.getGeoCountry(clientIP),
		ThreatLevel:     10, // Maximum threat level
	}

	log.Error().
		Str("event_id", event.EventID).
		Time("timestamp", event.Timestamp).
		Str("ip_hash", event.IPHash).
		Str("enumeration_type", event.EnumerationType).
		Str("attack_severity", event.AttackSeverity).
		Str("geo_country", event.GeoCountry).
		Str("sanitized_reason", s.sanitizeString(reason)).
		Int("threat_level", event.ThreatLevel).
		Msg("🔨 SANITIZED: Immediate ban executed")
}

// SanitizeIP creates a consistent hash of the IP address for tracking without exposing real IPs
func (s *SecureHoneypotLogger) SanitizeIP(ip string) string {
	// Parse IP to handle both IPv4 and IPv6
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "invalid_ip"
	}

	// For IPv4, mask last octet; for IPv6, mask last 64 bits
	var maskedIP string
	if parsedIP.To4() != nil {
		// IPv4: mask last octet (e.g., 192.168.1.xxx -> 192.168.1.0)
		ipv4 := parsedIP.To4()
		maskedIP = fmt.Sprintf("%d.%d.%d.0", ipv4[0], ipv4[1], ipv4[2])
	} else {
		// IPv6: mask last 64 bits
		ipv6 := parsedIP.To16()
		for i := 8; i < 16; i++ {
			ipv6[i] = 0
		}
		maskedIP = net.IP(ipv6).String()
	}

	// Create HMAC hash for consistent tracking without revealing real IPs
	h := hmac.New(sha256.New, s.ipHashSalt)
	h.Write([]byte(maskedIP))
	return hex.EncodeToString(h.Sum(nil))[:16] // Use first 16 chars for readability
}

// sanitizeUserAgent extracts safe patterns from user agent without revealing full strings
func (s *SecureHoneypotLogger) sanitizeUserAgent(userAgent string) string {
	if len(userAgent) == 0 {
		return "empty"
	}

	// Truncate if too long
	if len(userAgent) > s.maxLogLength {
		userAgent = userAgent[:s.maxLogLength]
	}

	// Extract safe patterns
	patterns := []string{}

	// Browser patterns
	browsers := []string{"Chrome", "Firefox", "Safari", "Edge", "Opera"}
	for _, browser := range browsers {
		if strings.Contains(userAgent, browser) {
			patterns = append(patterns, browser)
		}
	}

	// OS patterns
	oses := []string{"Windows", "macOS", "Linux", "Android", "iOS"}
	for _, os := range oses {
		if strings.Contains(userAgent, os) {
			patterns = append(patterns, os)
		}
	}

	// Bot/crawler patterns
	bots := []string{"bot", "crawler", "spider", "scan", "curl", "wget", "python", "go-http", "apache"}
	for _, bot := range bots {
		if strings.Contains(strings.ToLower(userAgent), bot) {
			patterns = append(patterns, "bot_like")
			break
		}
	}

	// Suspicious patterns
	suspicious := []string{"script", "hack", "exploit", "enum", "fuzz", "test"}
	for _, susp := range suspicious {
		if strings.Contains(strings.ToLower(userAgent), susp) {
			patterns = append(patterns, "suspicious")
			break
		}
	}

	if len(patterns) == 0 {
		return "unknown_pattern"
	}

	return strings.Join(patterns, ",")
}

// sanitizeSlug creates a safe pattern from the invalid slug without revealing exact attempts
func (s *SecureHoneypotLogger) sanitizeSlug(slug string) string {
	if len(slug) == 0 {
		return "empty"
	}

	// Truncate if too long
	if len(slug) > s.maxLogLength {
		slug = slug[:s.maxLogLength]
	}

	// URL decode first to handle encoded attacks
	decoded, err := url.QueryUnescape(slug)
	if err == nil {
		slug = decoded
	}

	patterns := []string{}

	// Check for UUID pattern
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if uuidPattern.MatchString(slug) {
		patterns = append(patterns, "uuid_format")
	}

	// Check for sequential patterns
	if matched, _ := regexp.MatchString(`\d+`, slug); matched {
		patterns = append(patterns, "contains_numbers")
	}

	// Check for common attack patterns
	attackPatterns := map[string]string{
		`\.\.\/`:        "path_traversal",
		`<script`:       "xss_attempt",
		`union.*select`: "sql_injection",
		`\${.*}`:        "template_injection",
		`eval\(`:        "code_injection",
		`cmd\.exe`:      "command_injection",
	}

	for pattern, name := range attackPatterns {
		if matched, _ := regexp.MatchString(`(?i)`+pattern, slug); matched {
			patterns = append(patterns, name)
		}
	}

	// Check length-based patterns
	if len(slug) > 50 {
		patterns = append(patterns, "long_string")
	}
	if len(slug) < 5 {
		patterns = append(patterns, "short_string")
	}

	// Check character patterns
	if matched, _ := regexp.MatchString(`[^a-zA-Z0-9\-_]`, slug); matched {
		patterns = append(patterns, "special_chars")
	}

	if len(patterns) == 0 {
		return fmt.Sprintf("pattern_len_%d", len(slug))
	}

	return strings.Join(patterns, ",")
}

// sanitizeEnumType ensures enumeration type values are safe
func (s *SecureHoneypotLogger) sanitizeEnumType(enumType string) string {
	safeTypes := map[string]string{
		"uuid":       "uuid",
		"sequential": "sequential",
		"dictionary": "dictionary",
		"fuzzing":    "fuzzing",
		"none":       "none",
	}

	if safe, exists := safeTypes[enumType]; exists {
		return safe
	}
	return "unknown_enum"
}

// sanitizeString provides general string sanitization
func (s *SecureHoneypotLogger) sanitizeString(input string) string {
	if len(input) == 0 {
		return "empty"
	}

	// Truncate if too long
	if len(input) > s.maxLogLength {
		input = input[:s.maxLogLength]
	}

	// Remove potential log injection characters
	sanitized := regexp.MustCompile(`[\r\n\t]`).ReplaceAllString(input, "_")

	// Replace other control characters
	sanitized = regexp.MustCompile(`[^\x20-\x7E]`).ReplaceAllString(sanitized, "?")

	return sanitized
}

// calculateSeverity determines attack severity based on patterns
func (s *SecureHoneypotLogger) calculateSeverity(enumerationType string, attemptCount int) string {
	if enumerationType == "uuid" || attemptCount > 20 {
		return "CRITICAL"
	}
	if enumerationType == "sequential" || attemptCount > 10 {
		return "HIGH"
	}
	if attemptCount > 5 {
		return "MEDIUM"
	}
	return "LOW"
}

// calculateThreatLevel assigns numeric threat level
func (s *SecureHoneypotLogger) calculateThreatLevel(enumerationType string, attemptCount int) int {
	base := 1

	switch enumerationType {
	case "uuid":
		base = 8
	case "sequential":
		base = 5
	case "dictionary":
		base = 4
	case "fuzzing":
		base = 6
	}

	// Add attempt count factor
	factor := attemptCount / 5
	if factor > 5 {
		factor = 5
	}

	level := base + factor
	if level > 10 {
		level = 10
	}

	return level
}

// getGeoCountry safely extracts country info without revealing exact location
func (s *SecureHoneypotLogger) getGeoCountry(ip string) string {
	// This would integrate with your existing geo lookup
	// For now, return safe placeholder
	return "unknown"
}

// generateEventID creates a unique event ID for correlation
func (s *SecureHoneypotLogger) generateEventID(ip string, timestamp time.Time) string {
	h := hmac.New(sha256.New, s.sessionKey)
	h.Write([]byte(fmt.Sprintf("%s-%d", ip, timestamp.Unix())))
	return hex.EncodeToString(h.Sum(nil))[:12]
}

// Global secure logger instance
var GlobalSecureLogger = NewSecureHoneypotLogger()
