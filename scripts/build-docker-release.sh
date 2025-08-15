#!/bin/bash
set -e

echo "🐳 Building Sauron Docker Production Release..."

VERSION=${1:-"v1.0.0"}
IMAGE_NAME="sauron-proxy"
RELEASE_DIR="docker-release-$VERSION"

# Clean previous builds
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

# ───────────── Build Docker Image ─────────────
echo "🔨 Building Docker image..."
docker build -f Dockerfile.production -t "$IMAGE_NAME:$VERSION" .
docker tag "$IMAGE_NAME:$VERSION" "$IMAGE_NAME:latest"

# ───────────── Save Docker Image ─────────────
echo "💾 Saving Docker image to file..."
docker save "$IMAGE_NAME:$VERSION" | gzip > "$RELEASE_DIR/sauron-$VERSION.tar.gz"

# ───────────── Copy Deployment Files ─────────────
echo "📁 Copying deployment files..."
cp docker-compose.production.yml "$RELEASE_DIR/docker-compose.yml"
cp .env.example "$RELEASE_DIR/"

# ───────────── Create Docker Deployment Script ─────────────
cat > "$RELEASE_DIR/deploy-docker.sh" << 'DEPLOY_EOF'
#!/bin/bash
set -e

echo "🐳 Deploying Sauron with Docker..."

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "❌ This script must be run as root (use sudo)"
   exit 1
fi

# Check for Docker
if ! command -v docker &> /dev/null; then
    echo "📦 Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    systemctl enable docker
    systemctl start docker
    rm get-docker.sh
fi

# Check for Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "📦 Installing Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi

# Check for .env file
if [ ! -f ".env" ]; then
    echo "❌ Missing .env file. Please create it with your configuration."
    echo "📝 Use: cp .env.example .env && nano .env"
    exit 1
fi

# Load Docker image
if [ -f "sauron-*.tar.gz" ]; then
    echo "📥 Loading Docker image..."
    docker load < sauron-*.tar.gz
else
    echo "❌ Docker image file not found"
    exit 1
fi

# Stop existing containers
echo "⏹️ Stopping existing containers..."
docker-compose down || true

# Start services
echo "🚀 Starting Sauron services..."
docker-compose up -d

echo "✅ Sauron deployed successfully!"
echo "🔍 Check status: docker-compose ps"
echo "📋 View logs: docker-compose logs -f sauron"

DEPLOY_EOF

chmod +x "$RELEASE_DIR/deploy-docker.sh"

# ───────────── Create Management Scripts ─────────────
cat > "$RELEASE_DIR/manage-sauron.sh" << 'MANAGE_EOF'
#!/bin/bash

case "$1" in
    start)
        echo "🚀 Starting Sauron..."
        docker-compose up -d
        ;;
    stop)
        echo "⏹️ Stopping Sauron..."
        docker-compose down
        ;;
    restart)
        echo "🔄 Restarting Sauron..."
        docker-compose restart
        ;;
    logs)
        echo "📋 Showing Sauron logs..."
        docker-compose logs -f sauron
        ;;
    status)
        echo "🔍 Checking Sauron status..."
        docker-compose ps
        ;;
    update)
        echo "🔄 Updating Sauron..."
        docker-compose pull
        docker-compose up -d
        ;;
    backup)
        echo "💾 Creating backup..."
        mkdir -p backups
        docker-compose exec sauron tar -czf - /app/config.db /app/logs | cat > "backups/sauron-backup-$(date +%Y%m%d_%H%M%S).tar.gz"
        echo "✅ Backup created in backups/ directory"
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|logs|status|update|backup}"
        echo ""
        echo "Commands:"
        echo "  start   - Start Sauron services"
        echo "  stop    - Stop Sauron services"
        echo "  restart - Restart Sauron services"
        echo "  logs    - Show live logs"
        echo "  status  - Show container status"
        echo "  update  - Pull and restart with latest image"
        echo "  backup  - Create backup of data"
        exit 1
        ;;
esac
MANAGE_EOF

chmod +x "$RELEASE_DIR/manage-sauron.sh"

# ───────────── Create Documentation ─────────────
cat > "$RELEASE_DIR/README.md" << 'README_EOF'
# Sauron Docker Production Release

Production-ready Docker deployment of Sauron MITM proxy system.

## Quick Deployment

1. **Upload this folder to your VPS**
2. **Configure environment:**
   ```bash
   cp .env.example .env
   nano .env  # Fill in your configuration
   ```
3. **Deploy:**
   ```bash
   sudo ./deploy-docker.sh
   ```

## Prerequisites

- Ubuntu 20.04+ VPS with root access
- Domain name pointed to your server IP
- Cloudflare account with API token
- Wildcard DNS: `*.yourdomain.com → your_server_ip`

## Management

Use the management script for common operations:

```bash
# Start services
./manage-sauron.sh start

# View logs
./manage-sauron.sh logs

# Check status
./manage-sauron.sh status

# Restart services
./manage-sauron.sh restart

# Create backup
./manage-sauron.sh backup
```

## Manual Docker Commands

```bash
# View all containers
docker-compose ps

# View logs
docker-compose logs -f sauron

# Execute commands inside container
docker-compose exec sauron /bin/sh

# Update services
docker-compose pull && docker-compose up -d
```

## Configuration

Edit `.env` file with your settings:

- `SAURON_DOMAIN`: Your phishing domain
- `CLOUDFLARE_API_TOKEN`: For automatic SSL certificates
- `TURNSTILE_SECRET`: Cloudflare Turnstile secret
- `ADMIN_KEY`: Strong password for admin access

## Data Persistence

The following data is persisted in Docker volumes:
- TLS certificates: `./tls/`
- Application logs: `./logs/`
- Database: `./config.db`
- Redis data: Docker volume
- ACME data: Docker volume

## Security Features

- Non-root container execution
- Security options enabled
- Network isolation
- Health checks
- Resource limits

## Troubleshooting

```bash
# Check container health
docker-compose ps

# View detailed logs
docker-compose logs sauron

# Restart specific service
docker-compose restart sauron

# Rebuild and restart
docker-compose up -d --build
```

## Backup and Recovery

```bash
# Create backup
./manage-sauron.sh backup

# Restore from backup (manual)
tar -xzf backups/sauron-backup-YYYYMMDD_HHMMSS.tar.gz
```
README_EOF

echo ""
echo "🎉 Docker release package created successfully!"
echo "📁 Location: $RELEASE_DIR/"
echo "🐳 Image: $IMAGE_NAME:$VERSION"
echo "💾 Archive size: $(du -h "$RELEASE_DIR/sauron-$VERSION.tar.gz" | cut -f1)"
echo ""
echo "🚀 Deployment Instructions:"
echo "1. Upload $RELEASE_DIR/ to your VPS"
echo "2. Configure: cp .env.example .env && nano .env"
echo "3. Deploy: sudo ./deploy-docker.sh"
echo ""
