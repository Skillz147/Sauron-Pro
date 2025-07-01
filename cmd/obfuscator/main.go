//cmd/obfuscator/main.go

package main

import (
	"log"
	"o365/inject"
)

func main() {
	if err := inject.BuildObfuscatedScript(); err != nil {
		log.Fatal("❌ Failed to build obfuscated script:", err)
	}
	log.Println("✅ Obfuscated script saved to inject/obfuscated.go")
}
