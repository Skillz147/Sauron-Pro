// configdb/links.go
package configdb

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

const slugLength = 8

var (
	tableOnce sync.Once
)

// ensureLinksTable creates the user_links table exactly once.
func ensureLinksTable(db *sql.DB) error {
	var err error
	tableOnce.Do(func() {
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS user_links (
				user_id   TEXT NOT NULL UNIQUE,
				slug      TEXT NOT NULL UNIQUE,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);
			CREATE INDEX IF NOT EXISTS idx_user_links_slug ON user_links(slug);
		`)
	})
	return err
}

// GetSlugForUser returns the existing slug or "".
func GetSlugForUser(userID string) (string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	if err := ensureLinksTable(db); err != nil {
		return "", err
	}

	var slug string
	err = db.QueryRow(`SELECT slug FROM user_links WHERE user_id = ?`, userID).Scan(&slug)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return slug, err
}

// CreateSlugForUser fetches or generates a unique slug for user_id.
func CreateSlugForUser(userID string) (string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	if err := ensureLinksTable(db); err != nil {
		return "", err
	}

	// 1) Already exists?
	var existing string
	if err := db.QueryRow(`SELECT slug FROM user_links WHERE user_id = ?`, userID).Scan(&existing); err == nil {
		return existing, nil
	}

	// 2) Generate (or reuse userID in dev mode)
	devMode := os.Getenv("DEV_MODE") == "true"
	slug := userID
	if !devMode {
		slug = generateRandomSlug(slugLength)
	}

	// 3) Insert with retry on slug collision (max 3 tries)
	for tries := 0; tries < 3; tries++ {
		_, err = db.Exec(
			`INSERT INTO user_links (user_id, slug, created_at) VALUES (?, ?, ?)`,
			userID, slug, time.Now(),
		)
		if err == nil {
			log.Printf("🔗 New link generated for %s: %s", userID, slug)
			return slug, nil
		}
		// UNIQUE constraint → slug collision, generate a new one
		if !devMode && strings.Contains(err.Error(), "UNIQUE") {
			slug = generateRandomSlug(slugLength)
			continue
		}
		return "", err
	}
	return "", errors.New("unable to generate unique slug after retries")
}

// IsValidSlug returns true if the slug exists and gives back its user_id.
func IsValidSlug(slug string) (string, bool) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", false
	}
	defer db.Close()

	if err := ensureLinksTable(db); err != nil {
		return "", false
	}

	var uid string
	err = db.QueryRow(`SELECT user_id FROM user_links WHERE slug = ?`, slug).Scan(&uid)
	return uid, err == nil
}

// GetUserIDForSlug returns the user_id for a given slug.
func GetUserIDForSlug(slug string) (string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var userID string
	err = db.QueryRow(`SELECT user_id FROM user_links WHERE slug = ?`, slug).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}
