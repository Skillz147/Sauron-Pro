package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := filepath.Join(".", "config.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT ip, country, timestamp FROM banned_ips ORDER BY timestamp DESC
	`)
	if err != nil {
		log.Fatalf("❌ Failed to query banned_ips: %v", err)
	}
	defer rows.Close()

	fmt.Println("🚫 Banned IPs:")
	fmt.Println("--------------------------------------------------")
	for rows.Next() {
		var ip, country string
		var ts string
		if err := rows.Scan(&ip, &country, &ts); err != nil {
			log.Printf("⚠️ Error scanning row: %v", err)
			continue
		}
		parsedTime, _ := time.Parse(time.RFC3339, ts)
		fmt.Printf("• %s [%s] — %s\n", ip, country, parsedTime.Format("2006-01-02 15:04:05"))
	}
}
