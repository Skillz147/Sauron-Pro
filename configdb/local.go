package configdb

import (
	"database/sql"
	"log"
	"time"
)

// AppConfig holds the core configuration

// InitLocalDB ensures all required tables exist before app runs
func InitLocalDB() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS config (
		id INTEGER PRIMARY KEY,
		status TEXT,
		domain_root TEXT,
		vps_server TEXT,
		cert_email TEXT,
		cloudflare_api_token TEXT
	);
	CREATE TABLE IF NOT EXISTS banned_ips (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL UNIQUE,
		country TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS telegram_settings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL UNIQUE,
		bot_token TEXT,
		chat_id TEXT
	);
	CREATE TABLE IF NOT EXISTS license_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL UNIQUE,
		key TEXT
	);
	CREATE TABLE IF NOT EXISTS user_links (
		user_id TEXT NOT NULL UNIQUE,
		slug TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS meta (
		key TEXT PRIMARY KEY,
		value TEXT
	);
	`)
	if err != nil {
		log.Printf("❌ Table creation error: %v", err)
	}
	return err
}

// SaveToLocalDB stores main config settings into the config table
func SaveToLocalDB(cfg AppConfig) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, _ = db.Exec(`DELETE FROM config`)
	_, err = db.Exec(`
		INSERT INTO config (status, domain_root, vps_server, cert_email, cloudflare_api_token)
		VALUES (?, ?, ?, ?, ?)`,
		cfg.Status, cfg.DomainRoot, cfg.VPSServer, cfg.CertEmail, cfg.CloudflareAPIToken,
	)
	if err == nil {
		log.Println("💾 App config saved to SQLite")
	}
	return err
}

// LoadFromLocalDB loads settings into Global
func LoadFromLocalDB() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	row := db.QueryRow(`SELECT status, domain_root, vps_server, cert_email, cloudflare_api_token FROM config LIMIT 1`)
	var cfg AppConfig
	if err := row.Scan(&cfg.Status, &cfg.DomainRoot, &cfg.VPSServer, &cfg.CertEmail, &cfg.CloudflareAPIToken); err != nil {
		return err
	}

	Global = cfg
	log.Printf("📦 App config loaded from local DB: %+v", cfg)
	return nil
}

// SaveBannedIP stores a blocked IP address
func SaveBannedIP(ip, country string) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("❌ DB open failed:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT OR IGNORE INTO banned_ips (ip, country)
		VALUES (?, ?)`, ip, country)

	if err != nil {
		log.Printf("❌ Failed to save banned IP %s (%s): %v", ip, country, err)
	} else {
		log.Printf("🛑 Banned IP saved: %s (%s)", ip, country)
	}
}

// SaveTelegramSettings persists Telegram bot settings for user
func SaveTelegramSettings(cfg TelegramConfig) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("❌ DB open failed:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT OR REPLACE INTO telegram_settings (user_id, bot_token, chat_id)
		VALUES (?, ?, ?)`,
		cfg.UserID, cfg.BotToken, cfg.ChatID,
	)

	if err != nil {
		log.Printf("❌ Failed to save telegram config for user %s: %v", cfg.UserID, err)
	} else {
		log.Printf("📬 Telegram config saved for user: %s", cfg.UserID)
	}
}

// SaveLicenseKey stores license key against user
func SaveLicenseKey(cfg LicenseConfig) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("❌ DB open failed:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT OR REPLACE INTO license_keys (user_id, key)
		VALUES (?, ?)`,
		cfg.UserID, cfg.Key,
	)

	if err != nil {
		log.Printf("❌ Failed to save license key for user %s: %v", cfg.UserID, err)
	} else {
		log.Printf("🔑 License key saved for user: %s", cfg.UserID)
	}
}

// generateRandomSlug returns a pseudo-random string
func generateRandomSlug(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
