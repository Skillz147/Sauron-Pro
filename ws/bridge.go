package ws

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"sync"
	"time"

	"o365/configdb"

	"github.com/gorilla/websocket"
)

type UserWS struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

// ──────────────────────────────── Message Send ────────────────────────────────

func SendToAdmin(msgType string, payload interface{}) {
	if AdminConn == nil {
		log.Println("⚠️ No admin WebSocket connected")
		return
	}
	sendUnsafe(AdminConn, msgType, payload)
}

func SendToUser(userID string, msgType string, payload interface{}) {
	const maxRetries = 6
	const retryDelay = 500 * time.Millisecond

	for attempt := 0; attempt <= maxRetries; attempt++ {
		Mu.Lock()
		uws, ok := UserConns[userID]
		Mu.Unlock()

		if ok && uws != nil && uws.Conn != nil {
			if attempt > 0 {
				log.Printf("✅ [WS] WebSocket became ready after %d retries for user: %s\n", attempt, userID)
			}
			sendSafe(uws, msgType, payload)
			return
		}
		log.Printf("🔍 [WS] Attempting to send '%s' to user %s...\n", msgType, userID)

		if attempt < maxRetries {
			log.Printf("⚠️ [WS] Not ready for user: %s — retrying... (%d/%d)\n", userID, attempt+1, maxRetries)
			time.Sleep(retryDelay)
		}
	}

	log.Printf("❌ [WS] Failed to send %q to user %s — no WebSocket after %d retries\n", msgType, userID, maxRetries)
}

func sendSafe(uws *UserWS, msgType string, payload interface{}) {
	msg := SocketMessage{Type: msgType, Data: payload}
	raw, err := json.Marshal(msg)
	if err != nil {
		log.Printf("❌ Failed to encode message '%s': %v\n", msgType, err)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(raw)

	uws.Mutex.Lock()
	defer uws.Mutex.Unlock()

	err = uws.Conn.WriteMessage(websocket.TextMessage, []byte(encoded))
	if err != nil {
		log.Printf("❌ Failed to send '%s' to user WS: %v\n", msgType, err)
	} else {
		log.Printf("✅ Successfully sent '%s' to user WS\n", msgType)
	}
}

func sendUnsafe(conn *websocket.Conn, msgType string, payload interface{}) {
	msg := SocketMessage{
		Type: msgType,
		Data: payload,
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		log.Printf("❌ Failed to encode admin message: %v\n", err)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(raw)

	if err := conn.WriteMessage(websocket.TextMessage, []byte(encoded)); err != nil {
		log.Printf("⚠️ Failed to send message to admin WS: %v\n", err)
	}
}

// ──────────────────────────────── Slug Updates ────────────────────────────────

func HandleSlugStatUpdate(slug string, stats configdb.SlugStats) {
	if userID, ok := GetUserIDForSlug(slug); ok {
		SendToUser(userID, "slug_stats", configdb.SlugStatsPayload{
			Slug:  slug,
			Stats: stats,
		})
	}
}
