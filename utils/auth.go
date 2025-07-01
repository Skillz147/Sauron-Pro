package utils

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

// Redis client (must be initialized at startup)
var RedisClient *redis.Client

// InitRedis sets up the global Redis connection
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",        // Adjust if needed
		Password: os.Getenv("REDIS_PASS"), // Optional
		DB:       0,
	})
}

// RequireAuth checks if user's token is present in Redis
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow Turnstile and verify endpoints
		if strings.HasPrefix(r.URL.Path, "/turnstile") || strings.HasPrefix(r.URL.Path, "/verify") {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			redirectToTurnstile(w, r)
			return
		}

		ok, err := RedisClient.Exists(ctx, "turnstile:"+cookie.Value).Result()
		if err != nil {
			log.Error().Err(err).Str("token", cookie.Value).Msg("Redis lookup failed")
			redirectToTurnstile(w, r)
			return
		}

		if ok == 0 {
			redirectToTurnstile(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func redirectToTurnstile(w http.ResponseWriter, r *http.Request) {
	target := "/turnstile?return=" + url.QueryEscape(r.URL.RequestURI())
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func CheckRedisToken(token string) bool {
	if token == "" {
		return false
	}
	ok, err := RedisClient.Exists(ctx, "turnstile:"+token).Result()
	return err == nil && ok > 0
}
