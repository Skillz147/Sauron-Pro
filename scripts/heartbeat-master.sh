#!/bin/bash

# 💓 Master Heartbeat Script for Dead Man's Switch
# Keeps VPS instances alive by sending regular heartbeats
# Author: Sauron Pro Security Team

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="/var/log/sauron/heartbeat.log"
PID_FILE="/var/run/sauron/heartbeat.pid"
ADMIN_KEY="${SAURON_ADMIN_KEY:-}"
HEARTBEAT_INTERVAL="${HEARTBEAT_INTERVAL:-30}" # seconds
VPS_LIST_FILE="/tmp/sauron_vps_list.json"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Logging function
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# Error function
error() {
    echo -e "${RED}❌ ERROR: $1${NC}" >&2
    log "ERROR: $1"
    exit 1
}

# Success function
success() {
    echo -e "${GREEN}✅ $1${NC}"
    log "SUCCESS: $1"
}

# Warning function
warn() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
    log "WARNING: $1"
}

# Check if already running
check_running() {
    if [[ -f "$PID_FILE" ]]; then
        local pid
        pid=$(cat "$PID_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            error "Heartbeat service already running with PID $pid"
        else
            rm -f "$PID_FILE"
        fi
    fi
}

# Create necessary directories
setup_directories() {
    mkdir -p "$(dirname "$LOG_FILE")"
    mkdir -p "$(dirname "$PID_FILE")"
}

# Validate prerequisites
validate_prerequisites() {
    log "Validating prerequisites..."
    
    # Check required tools
    local required_tools=("curl" "jq")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            error "Required tool '$tool' not found"
        fi
    done
    
    # Check admin key
    if [[ -z "$ADMIN_KEY" ]]; then
        echo -e "${YELLOW}Admin key not set in environment.${NC}"
        read -rsp "Enter admin key: " ADMIN_KEY
        echo
        if [[ -z "$ADMIN_KEY" ]]; then
            error "Admin key is required for heartbeat service"
        fi
    fi
    
    success "Prerequisites validated"
}

# Get VPS fleet list from master
get_vps_fleet() {
    log "Fetching VPS fleet list..."
    
    local master_url="${SAURON_MASTER_URL:-}"
    local response
    
    response=$(curl -s -k -X GET \
        -H "Authorization: Bearer $ADMIN_KEY" \
        "$master_url/fleet/instances" 2>/dev/null || echo "ERROR")
    
    if [[ "$response" == "ERROR" ]]; then
        warn "Failed to fetch VPS list from master"
        return 1
    fi
    
    echo "$response" > "$VPS_LIST_FILE"
    
    local count
    count=$(echo "$response" | jq -r '.instances | length' 2>/dev/null || echo "0")
    success "Retrieved $count VPS instances"
    
    return 0
}

# Send heartbeat to single VPS
send_heartbeat_to_vps() {
    local vps_domain=$1
    local vps_id=$2
    
    local response
    response=$(curl -s -k -X POST \
        -H "Authorization: Bearer $ADMIN_KEY" \
        -H "Content-Type: application/json" \
        -m 10 \
        "https://$vps_domain/vps/heartbeat" 2>/dev/null || echo "ERROR")
    
    if [[ "$response" == "ERROR" ]]; then
        warn "Failed to send heartbeat to VPS $vps_id ($vps_domain)"
        return 1
    fi
    
    local success_status
    success_status=$(echo "$response" | jq -r '.success // false' 2>/dev/null || echo "false")
    
    if [[ "$success_status" == "true" ]]; then
        log "💓 Heartbeat sent to VPS $vps_id ($vps_domain)"
        return 0
    else
        warn "Heartbeat failed for VPS $vps_id ($vps_domain)"
        return 1
    fi
}

# Send heartbeats to all VPS instances
send_fleet_heartbeats() {
    if [[ ! -f "$VPS_LIST_FILE" ]]; then
        warn "No VPS list available, skipping heartbeats"
        return 1
    fi
    
    local success_count=0
    local total_count=0
    
    # Parse VPS list and send heartbeats
    while IFS= read -r line; do
        local vps_id domain status
        vps_id=$(echo "$line" | jq -r '.id // empty')
        domain=$(echo "$line" | jq -r '.domain // empty')
        status=$(echo "$line" | jq -r '.status // empty')
        
        if [[ -n "$vps_id" && -n "$domain" ]]; then
            ((total_count++))
            
            if [[ "$status" == "active" ]]; then
                if send_heartbeat_to_vps "$domain" "$vps_id"; then
                    ((success_count++))
                fi
            else
                log "Skipping inactive VPS $vps_id"
            fi
        fi
    done < <(jq -c '.instances[]?' "$VPS_LIST_FILE" 2>/dev/null || echo "")
    
    log "💓 Heartbeat summary: $success_count/$total_count VPS instances"
    
    if [[ $success_count -eq 0 && $total_count -gt 0 ]]; then
        error "All heartbeats failed - potential network issues"
    fi
}

# Main heartbeat loop
heartbeat_loop() {
    log "Starting heartbeat service with ${HEARTBEAT_INTERVAL}s interval"
    
    while true; do
        # Refresh VPS list periodically (every 10 heartbeats)
        if [[ $(($(date +%s) % (HEARTBEAT_INTERVAL * 10))) -lt HEARTBEAT_INTERVAL ]]; then
            get_vps_fleet || warn "Failed to refresh VPS fleet list"
        fi
        
        # Send heartbeats
        send_fleet_heartbeats
        
        # Sleep until next heartbeat
        sleep "$HEARTBEAT_INTERVAL"
    done
}

# Signal handlers
cleanup() {
    log "Stopping heartbeat service..."
    rm -f "$PID_FILE"
    exit 0
}

# Start service
start_service() {
    check_running
    setup_directories
    validate_prerequisites
    
    # Write PID file
    echo $$ > "$PID_FILE"
    
    # Set up signal handlers
    trap cleanup SIGTERM SIGINT EXIT
    
    # Get initial VPS list
    get_vps_fleet || warn "Failed to get initial VPS list"
    
    success "Heartbeat service started with PID $$"
    
    # Start main loop
    heartbeat_loop
}

# Stop service
stop_service() {
    if [[ -f "$PID_FILE" ]]; then
        local pid
        pid=$(cat "$PID_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid"
            sleep 2
            if kill -0 "$pid" 2>/dev/null; then
                kill -9 "$pid"
            fi
            rm -f "$PID_FILE"
            success "Heartbeat service stopped"
        else
            warn "Heartbeat service not running (stale PID file)"
            rm -f "$PID_FILE"
        fi
    else
        warn "Heartbeat service not running"
    fi
}

# Status check
check_status() {
    if [[ -f "$PID_FILE" ]]; then
        local pid
        pid=$(cat "$PID_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            success "Heartbeat service running with PID $pid"
            
            # Show recent heartbeat activity
            if [[ -f "$LOG_FILE" ]]; then
                echo -e "${YELLOW}Recent activity:${NC}"
                tail -5 "$LOG_FILE"
            fi
        else
            error "Heartbeat service not running (stale PID file)"
        fi
    else
        warn "Heartbeat service not running"
    fi
}

# Usage information
usage() {
    cat << EOF
💓 Sauron Pro Master Heartbeat Service

USAGE:
    $0 {start|stop|restart|status}

COMMANDS:
    start     - Start the heartbeat service
    stop      - Stop the heartbeat service  
    restart   - Restart the heartbeat service
    status    - Check service status

ENVIRONMENT VARIABLES:
    SAURON_ADMIN_KEY      - Admin authentication key (required)
    SAURON_MASTER_URL     - Master server URL (default: https://master.sauron.pro)
    HEARTBEAT_INTERVAL    - Heartbeat interval in seconds (default: 30)

EXAMPLES:
    # Start service with custom interval
    HEARTBEAT_INTERVAL=60 $0 start
    
    # Check status
    $0 status
    
    # Restart service
    $0 restart

FILES:
    $LOG_FILE      - Service log file
    $PID_FILE           - Service PID file
    $VPS_LIST_FILE     - Cached VPS list
EOF
}

# Main execution
case "${1:-}" in
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        stop_service
        sleep 1
        start_service
        ;;
    status)
        check_status
        ;;
    *)
        usage
        exit 1
        ;;
esac
