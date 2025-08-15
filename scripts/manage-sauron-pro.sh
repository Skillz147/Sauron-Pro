#!/bin/bash

# Sauron-Pro Management Script
# Provides comprehensive management for Sauron-Pro deployments

COMMAND="$1"
DEPLOYMENT_TYPE=""

# Detect deployment type
if docker-compose -f docker-compose.pro.yml ps >/dev/null 2>&1 && [ "$(docker-compose -f docker-compose.pro.yml ps -q)" ]; then
    DEPLOYMENT_TYPE="docker"
elif systemctl is-active sauron >/dev/null 2>&1; then
    DEPLOYMENT_TYPE="binary"
else
    DEPLOYMENT_TYPE="none"
fi

case "$COMMAND" in
    "start")
        echo "🚀 Starting Sauron-Pro..."
        case "$DEPLOYMENT_TYPE" in
            "docker")
                docker-compose -f docker-compose.pro.yml up -d
                ;;
            "binary")
                sudo systemctl start sauron
                ;;
            "none")
                echo "❌ No deployment detected. Use ./deploy-production.sh first"
                exit 1
                ;;
        esac
        echo "✅ Sauron-Pro started"
        ;;
        
    "stop")
        echo "⏹️ Stopping Sauron-Pro..."
        case "$DEPLOYMENT_TYPE" in
            "docker")
                docker-compose -f docker-compose.pro.yml stop
                ;;
            "binary")
                sudo systemctl stop sauron
                ;;
        esac
        echo "✅ Sauron-Pro stopped"
        ;;
        
    "restart")
        echo "🔄 Restarting Sauron-Pro..."
        case "$DEPLOYMENT_TYPE" in
            "docker")
                docker-compose -f docker-compose.pro.yml restart sauron-pro
                ;;
            "binary")
                sudo systemctl restart sauron
                ;;
        esac
        echo "✅ Sauron-Pro restarted"
        ;;
        
    "status")
        echo "🔍 Sauron-Pro Status ($DEPLOYMENT_TYPE deployment):"
        echo "=============================================="
        
        case "$DEPLOYMENT_TYPE" in
            "docker")
                echo "🐳 Container Status:"
                docker-compose -f docker-compose.pro.yml ps
                echo ""
                echo "📊 Resource Usage:"
                docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}" $(docker-compose -f docker-compose.pro.yml ps -q) 2>/dev/null || echo "No containers running"
                ;;
            "binary")
                echo "📦 Service Status:"
                systemctl status sauron --no-pager
                echo ""
                echo "📊 Resource Usage:"
                ps aux | grep sauron | grep -v grep | awk '{print $11": CPU "$3"%, MEM "$4"%"}'
                ;;
            "none")
                echo "❌ No active deployment detected"
                ;;
        esac
        
        echo ""
        echo "🔌 Network Ports:"
        netstat -tlnp | grep -E ":(80|443|53)\s" || echo "No standard ports in use"
        ;;
        
    "logs")
        LINES="${2:-100}"
        echo "📋 Sauron-Pro Logs (last $LINES lines):"
        echo "======================================"
        
        case "$DEPLOYMENT_TYPE" in
            "docker")
                docker-compose -f docker-compose.pro.yml logs --tail="$LINES" -f sauron-pro
                ;;
            "binary")
                sudo journalctl -u sauron -n "$LINES" -f
                ;;
            "none")
                echo "❌ No deployment detected"
                exit 1
                ;;
        esac
        ;;
        
    "backup")
        BACKUP_DIR="backups/sauron-pro-$(date +%Y%m%d_%H%M%S)"
        echo "💾 Creating backup in $BACKUP_DIR..."
        
        mkdir -p "$BACKUP_DIR"
        
        case "$DEPLOYMENT_TYPE" in
            "docker")
                # Backup volumes and configs
                docker-compose -f docker-compose.pro.yml exec -T sauron-pro tar -czf - /app/config.db /app/logs 2>/dev/null | cat > "$BACKUP_DIR/app-data.tar.gz" || echo "Warning: Could not backup app data"
                docker run --rm -v docker-release-v2.0.0-pro_redis_data:/data -v $(pwd)/$BACKUP_DIR:/backup alpine tar -czf /backup/redis-data.tar.gz -C /data . 2>/dev/null || echo "Warning: Could not backup Redis data"
                cp .env "$BACKUP_DIR/" 2>/dev/null || echo "Warning: Could not backup .env"
                cp docker-compose.pro.yml "$BACKUP_DIR/" 2>/dev/null || echo "Warning: Could not backup docker-compose"
                ;;
            "binary")
                # Backup systemd and data
                sudo tar -czf "$BACKUP_DIR/system-data.tar.gz" /etc/systemd/system/sauron.service /usr/local/bin/sauron 2>/dev/null || echo "Warning: Could not backup system files"
                tar -czf "$BACKUP_DIR/app-data.tar.gz" config.db logs/ tls/ .env 2>/dev/null || echo "Warning: Could not backup app data"
                ;;
        esac
        
        echo "✅ Backup created: $BACKUP_DIR"
        echo "📊 Backup size: $(du -sh $BACKUP_DIR | cut -f1)"
        ;;
        
    "update")
        NEW_VERSION="${2:-v2.0.0-pro}"
        echo "🔄 Updating Sauron-Pro to $NEW_VERSION..."
        
        case "$DEPLOYMENT_TYPE" in
            "docker")
                echo "🐳 Pulling new Docker image..."
                ./scripts/build-docker-release.sh "$NEW_VERSION"
                docker-compose -f docker-compose.pro.yml up -d sauron-pro
                ;;
            "binary")
                echo "📦 Building new binary..."
                ./scripts/build-release.sh "$NEW_VERSION"
                echo "⏹️ Stopping service for update..."
                sudo systemctl stop sauron
                sudo cp "release-$NEW_VERSION/sauron/sauron" /usr/local/bin/sauron
                sudo chmod +x /usr/local/bin/sauron
                sudo systemctl start sauron
                ;;
        esac
        
        echo "✅ Update complete to $NEW_VERSION"
        ;;
        
    "config")
        echo "⚙️ Sauron-Pro Configuration:"
        echo "============================"
        
        if [ -f ".env" ]; then
            echo "📁 Environment Configuration:"
            grep -v -E '^#|^$|TOKEN|SECRET|PASSWORD' .env | sed 's/=.*/=***/' || echo "Could not read .env safely"
        else
            echo "❌ No .env file found"
        fi
        
        echo ""
        echo "🔧 System Configuration:"
        case "$DEPLOYMENT_TYPE" in
            "docker")
                echo "Deployment: Docker Compose"
                echo "Image: $(docker images sauron-proxy:* --format 'table {{.Repository}}:{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}')"
                ;;
            "binary")
                echo "Deployment: Systemd Service"
                echo "Binary: $(ls -lh /usr/local/bin/sauron 2>/dev/null || echo 'Not found')"
                ;;
        esac
        ;;
        
    "health")
        echo "🏥 Sauron-Pro Health Check:"
        echo "=========================="
        
        # Check if domain is set
        if [ -f ".env" ]; then
            DOMAIN=$(grep SAURON_DOMAIN .env | cut -d'=' -f2)
            if [ -n "$DOMAIN" ]; then
                echo "🌐 Testing domain: $DOMAIN"
                
                # Test HTTP
                if curl -s -f "http://$DOMAIN/health" >/dev/null 2>&1; then
                    echo "✅ HTTP health check passed"
                else
                    echo "❌ HTTP health check failed"
                fi
                
                # Test HTTPS
                if curl -s -f "https://$DOMAIN/health" >/dev/null 2>&1; then
                    echo "✅ HTTPS health check passed"
                else
                    echo "❌ HTTPS health check failed"
                fi
                
                # Test WebSocket
                if command -v wscat >/dev/null 2>&1; then
                    echo "🔌 WebSocket test..."
                    timeout 5 wscat -c "wss://$DOMAIN/ws" </dev/null >/dev/null 2>&1 && echo "✅ WebSocket accessible" || echo "❌ WebSocket connection failed"
                fi
            else
                echo "❌ SAURON_DOMAIN not configured"
            fi
        fi
        
        # Check Redis
        case "$DEPLOYMENT_TYPE" in
            "docker")
                if docker-compose -f docker-compose.pro.yml exec -T redis redis-cli ping >/dev/null 2>&1; then
                    echo "✅ Redis connection healthy"
                else
                    echo "❌ Redis connection failed"
                fi
                ;;
            "binary")
                if redis-cli ping >/dev/null 2>&1; then
                    echo "✅ Redis connection healthy"
                else
                    echo "❌ Redis connection failed"
                fi
                ;;
        esac
        ;;
        
    "stats")
        echo "📊 Sauron-Pro Statistics:"
        echo "========================"
        
        # Slug stats
        if [ -f "data/slug_stats.json" ]; then
            echo "🏷️ Slug Statistics:"
            jq -r '.slugs | to_entries[] | "\(.key): \(.value.total_requests) requests"' data/slug_stats.json 2>/dev/null | head -10 || echo "Could not parse slug stats"
        fi
        
        # Log statistics
        if [ -d "logs" ]; then
            echo ""
            echo "📋 Log Statistics:"
            echo "System logs: $(wc -l logs/system.log 2>/dev/null | cut -d' ' -f1 || echo 0) lines"
            echo "Bot logs: $(wc -l logs/bot.log 2>/dev/null | cut -d' ' -f1 || echo 0) lines"
            echo "Emit logs: $(wc -l logs/emits.log 2>/dev/null | cut -d' ' -f1 || echo 0) lines"
        fi
        ;;
        
    "clean")
        echo "🧹 Cleaning Sauron-Pro logs and temporary files..."
        
        # Clean logs
        if [ -d "logs" ]; then
            find logs/ -name "*.log" -mtime +7 -exec rm {} \; 2>/dev/null || true
            echo "✅ Old logs cleaned"
        fi
        
        # Clean Docker
        if [ "$DEPLOYMENT_TYPE" = "docker" ]; then
            docker system prune -f >/dev/null 2>&1
            echo "✅ Docker cleanup complete"
        fi
        
        # Clean builds
        rm -rf release-* docker-release-* 2>/dev/null || true
        echo "✅ Build artifacts cleaned"
        ;;
        
    *)
        echo "🏢 Sauron-Pro Management Console"
        echo "==============================="
        echo ""
        echo "Current deployment: $DEPLOYMENT_TYPE"
        echo ""
        echo "Usage: $0 <command> [options]"
        echo ""
        echo "Service Management:"
        echo "  start      - Start Sauron-Pro services"
        echo "  stop       - Stop Sauron-Pro services"  
        echo "  restart    - Restart Sauron-Pro services"
        echo "  status     - Show service status and resource usage"
        echo ""
        echo "Operations:"
        echo "  logs [n]   - Show last n lines of logs (default: 100)"
        echo "  backup     - Create full backup"
        echo "  update [v] - Update to version (default: v2.0.0-pro)"
        echo "  health     - Run health checks"
        echo "  stats      - Show statistics"
        echo ""
        echo "Maintenance:"
        echo "  config     - Show configuration"
        echo "  clean      - Clean logs and temporary files"
        echo ""
        echo "Examples:"
        echo "  $0 logs 50       # Show last 50 log lines"
        echo "  $0 update v2.1   # Update to v2.1"
        echo "  $0 backup        # Create backup"
        echo ""
        exit 1
        ;;
esac
