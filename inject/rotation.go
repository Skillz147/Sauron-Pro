//inject/rotation.go

package inject

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const rotationKey = "last_rotation"
const dbFile = "config.db"

// StartAutoRotation runs a goroutine that rebuilds the obfuscated script every 24 hours
func StartAutoRotation() {
	go func() {
		log.Println("⏳ Obfuscation rotation started...")

		// Ensure DB exists and ready
		ensureMetaTable()

		// First check immediately
		if needsRotation() {
			log.Println("🔁 First-time rotation triggered...")
			rotate()
		}

		for {
			time.Sleep(30 * time.Minute)
			if needsRotation() {
				log.Println("🔁 Scheduled rotation triggered...")
				rotate()
			}
		}
	}()
}

func rotate() {
	if err := BuildObfuscatedScript(); err != nil {
		log.Printf("❌ Obfuscation failed: %v", err)
		return
	}
	if err := updateLastRotation(); err != nil {
		log.Printf("⚠️ Could not update last_rotation: %v", err)
	} else {
		log.Println("✅ Obfuscation script rotated successfully")
	}
}

func ensureMetaTable() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("❌ SQLite open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT)`)
	if err != nil {
		log.Fatalf("❌ Failed to create meta table: %v", err)
	}
}

func needsRotation() bool {
	db, err := openDB()
	if err != nil {
		log.Printf("❌ SQLite open failed: %v", err)
		return false
	}
	defer db.Close()

	var last string
	err = db.QueryRow(`SELECT value FROM meta WHERE key = ?`, rotationKey).Scan(&last)
	if err != nil {
		// Not found = needs first-time rotation
		return true
	}

	parsed, err := time.Parse(time.RFC3339, last)
	if err != nil {
		log.Printf("⚠️ Invalid timestamp in meta: %v", err)
		return true
	}

	return time.Since(parsed) > 24*time.Hour
}

func updateLastRotation() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	_, err = db.Exec(`
		INSERT INTO meta (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, rotationKey, now)
	return err
}

func openDB() (*sql.DB, error) {
	// Try relative to CWD
	if _, err := os.Stat(dbFile); err == nil {
		return sql.Open("sqlite3", dbFile)
	}

	// Fallback: try two levels up
	alt := filepath.Join("..", "..", dbFile)
	if _, err := os.Stat(alt); err == nil {
		return sql.Open("sqlite3", alt)
	}

	// Else just return default (may fail)
	return sql.Open("sqlite3", dbFile)
}
