package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Try to read from .env file first (new method)
	domain := getDomainFromEnv()
	if domain != "" {
		fmt.Printf("Domain (from .env): %s\n", domain)
		return
	}

	// Fallback to config.db (legacy method)
	fmt.Println("No domain found in .env file")
	fmt.Println("Run './scripts/configure-env.sh setup' to configure domain")
}

// getDomainFromEnv reads SAURON_DOMAIN from .env file
func getDomainFromEnv() string {
	file, err := os.Open(".env")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "SAURON_DOMAIN=") {
			// Extract value after =
			value := strings.TrimPrefix(line, "SAURON_DOMAIN=")
			// Remove quotes if present
			value = strings.Trim(value, `"'`)
			return value
		}
	}
	return ""
}
