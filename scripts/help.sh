#!/bin/bash

# Colors for better readability
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Header
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${WHITE}🎯 SAURON MITM PROXY - COMMAND REFERENCE${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo

# Check if we're in a Sauron directory
if [ ! -f "sauron" ] && [ ! -f "/usr/local/bin/sauron" ]; then
    echo -e "${RED}⚠️  Warning: Not in Sauron directory or Sauron not installed${NC}"
    echo
fi

# Installation & Setup
echo -e "${CYAN}📦 INSTALLATION & SETUP${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}sudo ./install-production.sh${NC}          Install Sauron (first time setup)"
echo -e "${GREEN}./verify-installation.sh${NC}             Verify installation is working"
echo -e "${GREEN}./configure-env.sh${NC}                   Interactive configuration wizard"
echo

# Configuration Management
echo -e "${CYAN}🔧 CONFIGURATION MANAGEMENT${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}./configure-env.sh${NC}                   Run configuration wizard"
echo -e "${GREEN}./configure-env.sh --check${NC}           Validate current configuration"
echo -e "${GREEN}./configure-env.sh --status${NC}          Show configuration status"
echo -e "${GREEN}./configure-env.sh --test-domain${NC}     Test domain connectivity"
echo -e "${GREEN}./configure-env.sh --check-ssl${NC}       Verify SSL certificates"
echo -e "${GREEN}nano .env${NC}                            Edit configuration manually"
echo

# System Management
echo -e "${CYAN}⚙️  SYSTEM MANAGEMENT${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}sudo systemctl status sauron${NC}         Check service status"
echo -e "${GREEN}sudo systemctl start sauron${NC}          Start Sauron service"
echo -e "${GREEN}sudo systemctl stop sauron${NC}           Stop Sauron service"
echo -e "${GREEN}sudo systemctl restart sauron${NC}        Restart Sauron service"
echo -e "${GREEN}sudo systemctl enable sauron${NC}         Enable auto-start on boot"
echo -e "${GREEN}sudo systemctl disable sauron${NC}        Disable auto-start on boot"
echo

# Monitoring & Logs
echo -e "${CYAN}📊 MONITORING & LOGS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}sudo journalctl -u sauron -f${NC}         Follow live logs (Ctrl+C to exit)"
echo -e "${GREEN}sudo journalctl -u sauron --no-pager${NC} View all logs"
echo -e "${GREEN}sudo journalctl -u sauron -n 50${NC}      View last 50 log entries"
echo -e "${GREEN}sudo journalctl -u sauron --since today${NC} View today's logs"
echo -e "${GREEN}tail -f logs/system.log${NC}              Follow application logs"
echo

# Updates
echo -e "${CYAN}🔄 UPDATES${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}sudo ./update-sauron.sh --force${NC}              Auto-update to latest version"
echo -e "${GREEN}sudo ./update-sauron.sh --check${NC}      Check for updates (no install)"
echo -e "${GREEN}/usr/local/bin/sauron --version${NC}      Show current version"
echo

# SSL & Certificates
echo -e "${CYAN}🔒 SSL & CERTIFICATES${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}sudo acme.sh --list${NC}                  List SSL certificates"
echo -e "${GREEN}sudo acme.sh --renew -d yourdomain.com${NC} Renew SSL certificate"
echo -e "${GREEN}openssl x509 -in tls/cert.pem -text${NC}  View certificate details"
echo -e "${GREEN}./configure-env.sh --check-ssl${NC}       Verify SSL configuration"
echo

# Network & Domain Testing
echo -e "${CYAN}🌐 NETWORK & DOMAIN TESTING${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}dig yourdomain.com${NC}                   Check DNS resolution"
echo -e "${GREEN}curl -I https://yourdomain.com${NC}       Test domain connectivity"
echo -e "${GREEN}nslookup yourdomain.com${NC}              DNS lookup"
echo -e "${GREEN}ping yourdomain.com${NC}                  Test basic connectivity"
echo -e "${GREEN}ss -tlnp | grep :443${NC}                 Check if port 443 is open"
echo

# Access Points
echo -e "${CYAN}🎯 ACCESS POINTS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Admin Panel:${NC}      ${GREEN}https://yourdomain.com/admin${NC}"
echo -e "${YELLOW}Statistics:${NC}       ${GREEN}https://yourdomain.com/stats${NC}"
echo -e "${YELLOW}WebSocket:${NC}        ${GREEN}wss://yourdomain.com/ws${NC}"
echo -e "${YELLOW}Phishing Link:${NC}    ${GREEN}https://login.yourdomain.com/your-slug${NC}"
echo

# Troubleshooting
echo -e "${CYAN}🆘 TROUBLESHOOTING${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}./verify-installation.sh${NC}             Full system check"
echo -e "${GREEN}sudo systemctl status sauron${NC}         Service status"
echo -e "${GREEN}sudo journalctl -u sauron --no-pager -l${NC} Detailed error logs"
echo -e "${GREEN}./configure-env.sh --check${NC}           Configuration validation"
echo -e "${GREEN}ps aux | grep sauron${NC}                 Check running processes"
echo -e "${GREEN}netstat -tlnp | grep :443${NC}            Check port usage"
echo

# File Operations
echo -e "${CYAN}📁 FILE OPERATIONS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}ls -la${NC}                               List directory contents"
echo -e "${GREEN}cat .env${NC}                             View configuration file"
echo -e "${GREEN}ls -la /usr/local/bin/sauron*${NC}        List installed binaries"
echo -e "${GREEN}ls -la /etc/systemd/system/sauron*${NC}   List service files"
echo -e "${GREEN}du -h .${NC}                              Check directory size"
echo

# Docker (if applicable)
echo -e "${CYAN}🐳 DOCKER OPERATIONS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}docker ps${NC}                            List running containers"
echo -e "${GREEN}docker ps -a${NC}                         List all containers"
echo -e "${GREEN}docker logs sauron${NC}                   View container logs"
echo -e "${GREEN}docker restart sauron${NC}                Restart container"
echo -e "${GREEN}docker system prune${NC}                  Clean up Docker"
echo

# Emergency Commands
echo -e "${CYAN}🚨 EMERGENCY COMMANDS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${RED}sudo systemctl stop sauron${NC}           Stop service immediately"
echo -e "${RED}sudo pkill -f sauron${NC}                 Kill all Sauron processes"
echo -e "${RED}sudo rm /usr/local/bin/sauron${NC}        Remove binary (emergency)"
echo -e "${RED}sudo systemctl disable sauron${NC}        Disable service"
echo -e "${RED}cp .env.example .env${NC}                 Reset configuration"
echo

# Quick Actions
echo -e "${CYAN}⚡ QUICK ACTIONS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}./help.sh${NC}                           Show this help (you are here)"
echo -e "${GREEN}sudo ./update-sauron.sh && sudo systemctl status sauron${NC}"
echo -e "                                      Update and check status"
echo -e "${GREEN}sudo systemctl restart sauron && sudo journalctl -u sauron -f${NC}"
echo -e "                                      Restart and follow logs"
echo

# Current Status Check
echo -e "${CYAN}📊 CURRENT STATUS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Check if Sauron is installed
if [ -f "/usr/local/bin/sauron" ]; then
    VERSION=$(/usr/local/bin/sauron --version 2>/dev/null || echo "unknown")
    echo -e "${GREEN}✅ Sauron installed:${NC} $VERSION"
else
    echo -e "${RED}❌ Sauron not installed${NC}"
fi

# Check service status
if systemctl is-active sauron.service >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Service running${NC}"
else
    echo -e "${RED}❌ Service not running${NC}"
fi

# Check configuration
if [ -f ".env" ]; then
    echo -e "${GREEN}✅ Configuration file exists${NC}"
else
    echo -e "${YELLOW}⚠️  No .env file found${NC}"
fi

# Check domain configuration
if [ -f ".env" ]; then
    DOMAIN=$(grep "SAURON_DOMAIN" .env 2>/dev/null | cut -d'=' -f2 | tr -d '"' | tr -d "'" | xargs)
    if [ -n "$DOMAIN" ]; then
        echo -e "${GREEN}✅ Domain configured:${NC} $DOMAIN"
    else
        echo -e "${YELLOW}⚠️  No domain configured${NC}"
    fi
fi

echo
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${WHITE}💡 TIP: Use 'sudo ./update-sauron.sh' to get the latest features and fixes${NC}"
echo -e "${WHITE}📖 Full documentation: https://github.com/Skillz147/Sauron-Pro${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
