package configdb

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AppConfig struct {
	Status             string
	DomainRoot         string
	VPSServer          string
	CertEmail          string
	CloudflareAPIToken string
}

type TelegramConfig struct {
	UserID   string
	BotToken string
	ChatID   string
}

type LicenseConfig struct {
	UserID string
	Key    string
	Expiry time.Time
}

var dbPath = filepath.Join(".", "config.db")
var Global AppConfig

// IsIPBanned checks if an IP is in the banned list
func IsIPBanned(ip string) bool {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return false
	}
	defer db.Close()

	var count int
	err = db.QueryRow(`SELECT COUNT(1) FROM banned_ips WHERE ip = ?`, ip).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// GetLicenseKey fetches license key data from DB for a user
func GetLicenseKey(userID string) (LicenseConfig, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return LicenseConfig{}, err
	}
	defer db.Close()

	var key string
	var expiryStr string

	err = db.QueryRow(`SELECT key, expiry FROM licenses WHERE user_id = ?`, userID).Scan(&key, &expiryStr)
	if err != nil {
		return LicenseConfig{}, err
	}

	expiry, err := time.Parse(time.RFC3339, expiryStr)
	if err != nil {
		return LicenseConfig{}, err
	}

	return LicenseConfig{
		UserID: userID,
		Key:    key,
		Expiry: expiry,
	}, nil
}

// ValidateUserKey checks if the key matches and is not expired
func ValidateUserKey(userID, key string) (bool, error) {
	record, err := GetLicenseKey(userID)
	if err != nil {
		return false, err
	}
	if record.Key != key {
		return false, nil
	}
	if record.Expiry.Before(time.Now()) {
		return false, nil
	}
	return true, nil
}

// GetDomain returns domain_root from ENV or Global fallback
func GetDomain() string {
	val := strings.TrimSpace(os.Getenv("SAURON_DOMAIN"))
	if val != "" {
		return val
	}
	return Global.DomainRoot
}

// GetCloudflareToken returns token from ENV or Global fallback
func GetCloudflareToken() string {
	val := strings.TrimSpace(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if val != "" {
		return val
	}
	return Global.CloudflareAPIToken
}
