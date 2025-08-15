#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Sauron Installation Verification${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo

ERRORS=0
WARNINGS=0

# Function to report success
success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Function to report error
error() {
    echo -e "${RED}❌ $1${NC}"
    ((ERRORS++))
}

# Function to report warning
warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    ((WARNINGS++))
}

# Function to report info
info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Check if binary exists and is executable
echo -e "${BLUE}📦 BINARY INSTALLATION${NC}"
if [ -x "/usr/local/bin/sauron" ]; then
    success "Sauron binary installed at /usr/local/bin/sauron"
    
    # Check binary size (should be around 35MB)
    BINARY_SIZE=$(du -h /usr/local/bin/sauron | cut -f1)
    success "Binary size: $BINARY_SIZE"
    
    # Check version
    VERSION=$(/usr/local/bin/sauron --version 2>/dev/null || echo "unknown")
    success "Version: $VERSION"
else
    error "Sauron binary not found or not executable at /usr/local/bin/sauron"
fi
echo

# Check systemd service
echo -e "${BLUE}⚙️  SYSTEMD SERVICE${NC}"
if systemctl list-unit-files | grep -q "sauron.service"; then
    success "Sauron service file exists"
    
    if systemctl is-enabled sauron.service >/dev/null 2>&1; then
        success "Sauron service enabled (auto-start on boot)"
    else
        warning "Sauron service not enabled for auto-start"
    fi
    
    if systemctl is-active sauron.service >/dev/null 2>&1; then
        success "Sauron service running"
        
        # Check how long it's been running
        UPTIME=$(systemctl show sauron.service --property=ActiveEnterTimestamp --value)
        info "Service started: $UPTIME"
    else
        error "Sauron service not running"
        info "Try: sudo systemctl start sauron"
    fi
else
    error "Sauron service file not found"
fi
echo

# Check configuration
echo -e "${BLUE}🔧 CONFIGURATION${NC}"
if [ -f ".env" ]; then
    success "Environment file exists (.env)"
    
    # Source the .env file safely
    if source .env 2>/dev/null; then
        # Check required variables
        if [ -n "$SAURON_DOMAIN" ]; then
            success "SAURON_DOMAIN configured: $SAURON_DOMAIN"
        else
            error "SAURON_DOMAIN not set"
        fi
        
        if [ -n "$CLOUDFLARE_API_TOKEN" ]; then
            success "CLOUDFLARE_API_TOKEN configured"
        else
            warning "CLOUDFLARE_API_TOKEN not set (SSL automation disabled)"
        fi
        
        if [ -n "$TURNSTILE_SECRET" ]; then
            success "TURNSTILE_SECRET configured (bot protection enabled)"
        else
            info "TURNSTILE_SECRET not set (bot protection disabled)"
        fi
    else
        error "Failed to parse .env file"
    fi
else
    error ".env file not found"
    info "Create one with: cp .env.example .env"
fi
echo

# Check network and ports
echo -e "${BLUE}🌐 NETWORK & PORTS${NC}"
if command -v ss >/dev/null 2>&1; then
    if ss -tlnp | grep -q ":443"; then
        success "Port 443 (HTTPS) is open"
        HTTPS_PROCESS=$(ss -tlnp | grep ":443" | awk '{print $6}' | head -1)
        info "Process: $HTTPS_PROCESS"
    else
        error "Port 443 (HTTPS) not open"
    fi
    
    if ss -tlnp | grep -q ":80"; then
        success "Port 80 (HTTP) is open"
    else
        warning "Port 80 (HTTP) not open (redirects may not work)"
    fi
else
    warning "ss command not available, skipping port check"
fi
echo

# Check SSL certificates
echo -e "${BLUE}🔒 SSL CERTIFICATES${NC}"
if [ -f "tls/cert.pem" ]; then
    success "SSL certificate file exists (tls/cert.pem)"
    
    # Check certificate expiry
    if command -v openssl >/dev/null 2>&1; then
        CERT_EXPIRY=$(openssl x509 -in tls/cert.pem -noout -enddate 2>/dev/null | cut -d= -f2)
        if [ -n "$CERT_EXPIRY" ]; then
            success "Certificate expires: $CERT_EXPIRY"
            
            # Check if certificate is still valid
            if openssl x509 -in tls/cert.pem -noout -checkend 604800 >/dev/null 2>&1; then
                success "Certificate valid for at least 7 more days"
            else
                warning "Certificate expires within 7 days - renewal needed"
            fi
        fi
    fi
else
    warning "SSL certificate not found (tls/cert.pem)"
    info "Certificates will be auto-generated on first run"
fi

if [ -f "tls/key.pem" ]; then
    success "SSL private key exists (tls/key.pem)"
else
    warning "SSL private key not found (tls/key.pem)"
fi
echo

# Check acme.sh installation
echo -e "${BLUE}🏆 ACME.SH (SSL AUTOMATION)${NC}"
if command -v acme.sh >/dev/null 2>&1; then
    success "acme.sh installed"
    
    ACME_VERSION=$(acme.sh --version 2>/dev/null | head -1)
    info "Version: $ACME_VERSION"
    
    # List certificates
    CERT_COUNT=$(acme.sh --list 2>/dev/null | grep -c "Main_Domain" || echo "0")
    if [ "$CERT_COUNT" -gt 0 ]; then
        success "$CERT_COUNT SSL certificate(s) managed by acme.sh"
    else
        info "No certificates managed by acme.sh yet"
    fi
else
    warning "acme.sh not found (manual SSL management required)"
fi
echo

# Check Docker (if applicable)
echo -e "${BLUE}🐳 DOCKER${NC}"
if command -v docker >/dev/null 2>&1; then
    success "Docker installed"
    
    DOCKER_VERSION=$(docker --version 2>/dev/null)
    info "$DOCKER_VERSION"
    
    if docker ps >/dev/null 2>&1; then
        success "Docker daemon running"
        
        # Check for Sauron containers
        SAURON_CONTAINERS=$(docker ps --filter "name=sauron" --format "table {{.Names}}" | grep -v NAMES | wc -l)
        if [ "$SAURON_CONTAINERS" -gt 0 ]; then
            success "$SAURON_CONTAINERS Sauron container(s) running"
        else
            info "No Sauron containers running (binary deployment)"
        fi
    else
        warning "Docker daemon not running"
    fi
else
    info "Docker not installed (not required for binary deployment)"
fi
echo

# Check Redis (if applicable)
echo -e "${BLUE}🔴 REDIS${NC}"
if command -v redis-cli >/dev/null 2>&1; then
    if redis-cli ping >/dev/null 2>&1; then
        success "Redis server running"
    else
        warning "Redis server not responding"
    fi
elif systemctl is-active redis-server >/dev/null 2>&1; then
    success "Redis service running"
else
    info "Redis not detected (Telegram queue may not work)"
fi
echo

# Check domain connectivity (if domain is configured)
if [ -n "$SAURON_DOMAIN" ]; then
    echo -e "${BLUE}🌍 DOMAIN CONNECTIVITY${NC}"
    
    # DNS resolution
    if command -v dig >/dev/null 2>&1; then
        if dig "$SAURON_DOMAIN" +short | grep -q .; then
            success "DNS resolution working for $SAURON_DOMAIN"
            IP=$(dig "$SAURON_DOMAIN" +short | head -1)
            info "Resolves to: $IP"
        else
            error "DNS resolution failed for $SAURON_DOMAIN"
        fi
    fi
    
    # HTTP connectivity
    if command -v curl >/dev/null 2>&1; then
        if curl -s -I "https://$SAURON_DOMAIN" --connect-timeout 10 | grep -q "HTTP"; then
            success "HTTPS connectivity working"
            STATUS=$(curl -s -I "https://$SAURON_DOMAIN" --connect-timeout 10 | head -1)
            info "Response: $STATUS"
        else
            warning "HTTPS connectivity failed"
            info "This is normal during initial setup"
        fi
    fi
    echo
fi

# Check log files
echo -e "${BLUE}📊 LOG FILES${NC}"
if [ -d "logs" ]; then
    success "Logs directory exists"
    
    LOG_FILES=$(find logs -name "*.log" -type f 2>/dev/null | wc -l)
    if [ "$LOG_FILES" -gt 0 ]; then
        success "$LOG_FILES log file(s) found"
        
        # Check recent activity
        RECENT_LOGS=$(find logs -name "*.log" -type f -mtime -1 2>/dev/null | wc -l)
        if [ "$RECENT_LOGS" -gt 0 ]; then
            success "$RECENT_LOGS log file(s) updated in last 24 hours"
        else
            info "No recent log activity"
        fi
    else
        info "No log files found yet"
    fi
else
    info "Logs directory not found (will be created on first run)"
fi
echo

# Check working directory structure
echo -e "${BLUE}📁 DIRECTORY STRUCTURE${NC}"
REQUIRED_DIRS=("tls" "logs" "data" "geo")
for dir in "${REQUIRED_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        success "Directory exists: $dir"
    else
        info "Directory missing: $dir (will be created automatically)"
    fi
done
echo

# Summary
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📊 VERIFICATION SUMMARY${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    success "Perfect! All checks passed."
    echo -e "${GREEN}🎯 Sauron is properly installed and configured!${NC}"
elif [ $ERRORS -eq 0 ]; then
    success "Installation looks good with $WARNINGS minor warning(s)."
    echo -e "${YELLOW}⚠️  Address warnings above for optimal performance.${NC}"
else
    error "$ERRORS critical issue(s) found, $WARNINGS warning(s)."
    echo -e "${RED}🚨 Fix critical issues before using Sauron.${NC}"
fi

echo
echo -e "${BLUE}🔧 Quick Actions:${NC}"
if [ $ERRORS -gt 0 ]; then
    echo -e "${YELLOW}./configure-env.sh${NC}                   Fix configuration issues"
    echo -e "${YELLOW}sudo systemctl start sauron${NC}         Start the service"
fi
echo -e "${YELLOW}sudo systemctl status sauron${NC}         Check service status"
echo -e "${YELLOW}sudo journalctl -u sauron -f${NC}         View live logs"
echo -e "${YELLOW}./help.sh${NC}                           Show all commands"

echo
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Exit with appropriate code
if [ $ERRORS -gt 0 ]; then
    exit 1
else
    exit 0
fi
