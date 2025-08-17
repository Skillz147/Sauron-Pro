#!/bin/bash

#############################################################################
# Sauron-Pro Fleet Master Controller Deployment
# 
# This script deploys the master controller server that manages multiple
# VPS instances in a distributed MITM architecture.
#
# Master Controller Functions:
# - VPS registration and heartbeat management
# - Command dispatch to VPS fleet
# - Fleet monitoring and statistics
# - Centralized script management
# - Admin interface for fleet control
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
DOMAIN="${DOMAIN:-your-domain.com}"
FLEET_PORT="${FLEET_PORT:-8443}"

if [[ "$DOMAIN" == "your-domain.com" ]]; then
    echo -e "${RED}❌ Please set DOMAIN environment variable or edit this script${NC}"
    echo -e "${YELLOW}💡 Example: DOMAIN=admin.yourdomain.com ./deploy-fleet-master.sh${NC}"
    exit 1
fi

SERVER_NAME="sauron-fleet-master"
SERVICE_USER="sauron"
INSTALL_DIR="/opt/sauron-pro"
CONFIG_DIR="/etc/sauron-pro"
LOG_DIR="/var/log/sauron-pro"

echo -e "${BLUE}🚀 Sauron-Pro Fleet Master Controller Deployment${NC}"
echo -e "${BLUE}=================================================${NC}"
echo -e "${CYAN}📍 Domain: ${DOMAIN}${NC}"
echo -e "${CYAN}🔌 Fleet Port: ${FLEET_PORT}${NC}"
echo -e "${CYAN}📁 Install Directory: ${INSTALL_DIR}${NC}"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}❌ This script must be run as root (use sudo)${NC}"
   exit 1
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
echo -e "${YELLOW}🔨 Building Sauron-Pro Fleet Master...${NC}"
if [[ ! -f "go.mod" ]]; then
    echo -e "${RED}❌ go.mod not found. Run from project root directory.${NC}"
    exit 1
fi

# Build with fleet management tags
go build -ldflags "-X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.Version=fleet-master-$(date +%Y%m%d)" -o "$INSTALL_DIR/bin/sauron-fleet-master" .

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

# Create fleet master systemd service
echo -e "${YELLOW}🔧 Creating systemd service...${NC}"
cat > /etc/systemd/system/sauron-fleet-master.service << EOF
[Unit]
Description=Sauron-Pro Fleet Master Controller
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=5
User=$SERVICE_USER
Group=$SERVICE_USER
WorkingDirectory=$INSTALL_DIR

# Environment variables
Environment=SAURON_MODE=fleet-master
Environment=DOMAIN=$DOMAIN
Environment=FLEET_PORT=$FLEET_PORT
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
MemoryMax=1G

# Executable
ExecStart=$INSTALL_DIR/bin/sauron-fleet-master
ExecReload=/bin/kill -HUP \$MAINPID

# Logging
StandardOutput=append:$LOG_DIR/fleet-master.log
StandardError=append:$LOG_DIR/fleet-master-error.log

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and enable service
systemctl daemon-reload
systemctl enable sauron-fleet-master.service

# Create fleet master configuration
echo -e "${YELLOW}⚙️ Creating fleet master configuration...${NC}"
cat > "$CONFIG_DIR/fleet-config.json" << EOF
{
  "fleet_master": {
    "domain": "$DOMAIN",
    "port": $FLEET_PORT,
    "max_vps_instances": 100,
    "heartbeat_timeout": 600,
    "command_timeout": 30,
    "database": {
      "path": "$INSTALL_DIR/data/fleet.db",
      "backup_interval": 3600
    },
    "security": {
      "require_vps_auth": true,
      "max_command_rate": 10,
      "allowed_commands": ["status", "restart", "script", "config", "update"]
    },
    "monitoring": {
      "enable_metrics": true,
      "metrics_port": 9090,
      "alert_thresholds": {
        "vps_down_alert": 5,
        "response_time_alert": 5000
      }
    }
  }
}
EOF

chown "$SERVICE_USER:$SERVICE_USER" "$CONFIG_DIR/fleet-config.json"
chmod 640 "$CONFIG_DIR/fleet-config.json"

# Create fleet management scripts
echo -e "${YELLOW}🛠️ Creating fleet management utilities...${NC}"

# Fleet status script
cat > "$INSTALL_DIR/bin/fleet-status" << EOF
#!/bin/bash
# Fleet Status Utility
curl -s "https://\${DOMAIN}:\${FLEET_PORT}/fleet/instances" | jq '.'
EOF

# Fleet command script
cat > "$INSTALL_DIR/bin/fleet-command" << EOF
#!/bin/bash
# Fleet Command Utility
if [[ \$# -lt 2 ]]; then
    echo "Usage: \$0 <vps-id> <command> [payload]"
    echo "Commands: status, restart, script, config, update"
    exit 1
fi

VPS_ID="\$1"
COMMAND="\$2"
PAYLOAD="\${3:-{}}"

curl -X POST "https://\${DOMAIN}:\${FLEET_PORT}/fleet/command" \\
  -H "Content-Type: application/json" \\
  -d "{\\"vps_id\\":\\"\$VPS_ID\\", \\"command\\":\\"\$COMMAND\\", \\"payload\\":\$PAYLOAD}" | jq '.'
EOF

chmod +x "$INSTALL_DIR/bin/fleet-status" "$INSTALL_DIR/bin/fleet-command"

# Create firewall rules
echo -e "${YELLOW}🔥 Configuring firewall...${NC}"
if command -v ufw >/dev/null 2>&1; then
    ufw allow "$FLEET_PORT/tcp" comment "Sauron Fleet Master"
    ufw allow 22/tcp comment "SSH"
    ufw allow 80/tcp comment "HTTP"
    ufw allow 443/tcp comment "HTTPS"
    echo -e "${GREEN}✅ UFW rules configured${NC}"
elif command -v firewall-cmd >/dev/null 2>&1; then
    firewall-cmd --permanent --add-port="$FLEET_PORT/tcp"
    firewall-cmd --reload
    echo -e "${GREEN}✅ Firewalld rules configured${NC}"
fi

# Create log rotation
echo -e "${YELLOW}📋 Setting up log rotation...${NC}"
cat > /etc/logrotate.d/sauron-fleet-master << EOF
$LOG_DIR/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    sharedscripts
    postrotate
        systemctl reload sauron-fleet-master || true
    endscript
}
EOF

# Set executable permissions
chown "$SERVICE_USER:$SERVICE_USER" "$INSTALL_DIR/bin/sauron-fleet-master"
chmod 755 "$INSTALL_DIR/bin/sauron-fleet-master"

# Start the service
echo -e "${YELLOW}🚀 Starting Sauron Fleet Master...${NC}"
systemctl start sauron-fleet-master.service

# Wait a moment for startup
sleep 3

# Check service status
if systemctl is-active --quiet sauron-fleet-master.service; then
    echo -e "${GREEN}✅ Sauron Fleet Master started successfully${NC}"
    echo ""
    echo -e "${PURPLE}📊 Fleet Master Controller Information:${NC}"
    echo -e "${CYAN}🌐 Domain: ${DOMAIN}:${FLEET_PORT}${NC}"
    echo -e "${CYAN}📁 Install Directory: ${INSTALL_DIR}${NC}"
    echo -e "${CYAN}📋 Config Directory: ${CONFIG_DIR}${NC}"
    echo -e "${CYAN}📝 Log Directory: ${LOG_DIR}${NC}"
    echo ""
    echo -e "${PURPLE}🛠️ Fleet Management Commands:${NC}"
    echo -e "${CYAN}   ${INSTALL_DIR}/bin/fleet-status          - View all VPS instances${NC}"
    echo -e "${CYAN}   ${INSTALL_DIR}/bin/fleet-command <id> <cmd> - Send command to VPS${NC}"
    echo -e "${CYAN}   systemctl status sauron-fleet-master    - Check service status${NC}"
    echo -e "${CYAN}   journalctl -u sauron-fleet-master -f    - View live logs${NC}"
    echo ""
    echo -e "${PURPLE}🔗 Fleet Management Endpoints:${NC}"
    echo -e "${CYAN}   POST ${DOMAIN}:${FLEET_PORT}/fleet/register   - VPS registration${NC}"
    echo -e "${CYAN}   GET  ${DOMAIN}:${FLEET_PORT}/fleet/instances  - List VPS instances${NC}"
    echo -e "${CYAN}   POST ${DOMAIN}:${FLEET_PORT}/fleet/command    - Send VPS commands${NC}"
    echo -e "${CYAN}   GET  ${DOMAIN}:${FLEET_PORT}/admin/scripts    - Script management${NC}"
    echo ""
    echo -e "${GREEN}🎯 Fleet Master Controller deployment complete!${NC}"
    echo -e "${YELLOW}💡 Next: Deploy VPS agents using deploy-vps-agent.sh${NC}"
else
    echo -e "${RED}❌ Failed to start Sauron Fleet Master${NC}"
    echo -e "${YELLOW}📋 Check logs: journalctl -u sauron-fleet-master.service${NC}"
    exit 1
fi
