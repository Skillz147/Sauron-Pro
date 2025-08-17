#!/bin/bash

# 🔴 Kill Switch Emergency Activation Script
# Use with EXTREME CAUTION - This will DESTROY the target VPS
# Author: Sauron Pro Security Team
# Usage: ./kill-switch.sh [VPS_ID] [LEVEL] [DELAY] [REASON]

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="/tmp/killswitch-$(date +%s).log"
MASTER_ENDPOINT="${SAURON_MASTER_URL:-}"
ADMIN_KEY="${SAURON_ADMIN_KEY:-}"
CONFIRMATION_CODE="OMEGA-DESTROY"

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

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

# Warning function
warn() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
    log "WARNING: $1"
}

# Success function
success() {
    echo -e "${GREEN}✅ $1${NC}"
    log "SUCCESS: $1"
}

# Show banner
show_banner() {
    echo -e "${RED}"
    cat << 'EOF'
╔══════════════════════════════════════════════════════════════╗
║                    🔴 KILL SWITCH ACTIVATED 🔴                ║
║                                                              ║
║  ⚠️  THIS WILL PERMANENTLY DESTROY THE TARGET VPS  ⚠️        ║
║                                                              ║
║  • All data will be irrecoverably wiped                     ║
║  • System will be rendered inoperable                       ║
║  • No forensic traces will remain                           ║
║                                                              ║
║             USE ONLY IN EMERGENCY SITUATIONS                ║
╚══════════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
}

# Validate prerequisites
validate_prerequisites() {
    log "Validating prerequisites..."
    
    # Check if running as root/sudo (for system operations)
    if [[ $EUID -eq 0 ]]; then
        warn "Running as root - maximum destruction capability enabled"
    fi
    
    # Check required tools
    local required_tools=("curl" "jq" "date")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            error "Required tool '$tool' not found. Please install it."
        fi
    done
    
    # Check admin key
    if [[ -z "$ADMIN_KEY" ]]; then
        echo -e "${YELLOW}Admin key not set in environment.${NC}"
        read -rsp "Enter admin key: " ADMIN_KEY
        echo
        if [[ -z "$ADMIN_KEY" ]]; then
            error "Admin key is required for kill switch activation"
        fi
    fi
    
    success "Prerequisites validated"
}

# Get VPS list
get_vps_list() {
    log "Fetching VPS fleet list..."
    
    local response
    response=$(curl -s -k -X GET \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ADMIN_KEY" \
        "$MASTER_ENDPOINT/fleet/instances" || echo "ERROR")
    
    if [[ "$response" == "ERROR" ]]; then
        error "Failed to fetch VPS list from master"
    fi
    
    echo "$response" | jq -r '.instances[]? | "\(.id) - \(.domain) (\(.status))"' 2>/dev/null || {
        warn "Could not parse VPS list, proceeding with manual input"
        return 1
    }
}

# Interactive VPS selection
select_vps() {
    echo -e "${YELLOW}Available VPS instances:${NC}"
    
    if get_vps_list; then
        echo
        echo "Enter VPS ID (or 'ALL' for fleet-wide destruction):"
    else
        echo "Enter VPS ID manually (or 'ALL' for fleet-wide destruction):"
    fi
    
    read -r vps_id
    
    if [[ -z "$vps_id" ]]; then
        error "VPS ID cannot be empty"
    fi
    
    echo "$vps_id"
}

# Interactive destruction level selection
select_destruction_level() {
    echo -e "${YELLOW}Select destruction level:${NC}"
    echo "1 - Memory Purge Only (RAM cleanup)"
    echo "2 - Data Obliteration (application data + logs)"
    echo "3 - System Corruption (system files + config)"
    echo "4 - Hardware Destruction (full disk wipe)"
    echo "5 - Stealth Exit (complete annihilation + panic)"
    echo
    read -rp "Destruction level (1-5): " level
    
    if ! [[ "$level" =~ ^[1-5]$ ]]; then
        error "Invalid destruction level. Must be 1-5."
    fi
    
    echo "$level"
}

# Confirmation prompts
confirm_destruction() {
    local vps_id=$1
    local level=$2
    
    echo -e "${RED}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    FINAL CONFIRMATION                         ║"
    echo "╠══════════════════════════════════════════════════════════════╣"
    echo "║ Target VPS: $vps_id"
    echo "║ Destruction Level: $level"
    echo "║ This action is IRREVERSIBLE and DESTRUCTIVE"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    # Triple confirmation
    read -rp "Type 'CONFIRM' to proceed: " confirm1
    if [[ "$confirm1" != "CONFIRM" ]]; then
        error "Operation cancelled by user"
    fi
    
    read -rp "Type 'DESTROY' to continue: " confirm2
    if [[ "$confirm2" != "DESTROY" ]]; then
        error "Operation cancelled by user"
    fi
    
    read -rp "Type 'OMEGA' for final confirmation: " confirm3
    if [[ "$confirm3" != "OMEGA" ]]; then
        error "Operation cancelled by user"
    fi
    
    success "Triple confirmation received - proceeding with destruction"
}

# Execute kill switch
execute_kill_switch() {
    local vps_id=$1
    local level=$2
    local delay=${3:-0}
    local reason=${4:-"Emergency kill switch activation"}
    
    log "Executing kill switch - VPS: $vps_id, Level: $level, Delay: ${delay}s"
    
    # Prepare payload
    local payload
    payload=$(jq -n \
        --arg admin_key "$ADMIN_KEY" \
        --arg vps_id "$vps_id" \
        --argjson level "$level" \
        --argjson delay "$delay" \
        --arg reason "$reason" \
        --arg confirmation_code "$CONFIRMATION_CODE" \
        '{
            admin_key: $admin_key,
            vps_id: ($vps_id | if . == "ALL" then "" else . end),
            destruction_level: $level,
            delay_seconds: $delay,
            reason: $reason,
            confirmation_code: $confirmation_code
        }'
    )
    
    log "Sending kill switch command to master..."
    
    # Send kill switch command
    local response
    response=$(curl -s -k -X POST \
        -H "Content-Type: application/json" \
        -d "$payload" \
        "$MASTER_ENDPOINT/admin/killswitch" || echo "ERROR")
    
    if [[ "$response" == "ERROR" ]]; then
        error "Failed to send kill switch command"
    fi
    
    # Parse response
    local success_status
    success_status=$(echo "$response" | jq -r '.success // false' 2>/dev/null || echo "false")
    
    if [[ "$success_status" == "true" ]]; then
        success "Kill switch command sent successfully"
        
        local message
        message=$(echo "$response" | jq -r '.message // "Unknown"' 2>/dev/null)
        log "Response: $message"
        
        if [[ "$delay" -gt 0 ]]; then
            warn "Destruction will begin in $delay seconds"
            countdown "$delay"
        fi
        
        success "🔴 KILL SWITCH ACTIVATED - DESTRUCTION IN PROGRESS"
    else
        local error_msg
        error_msg=$(echo "$response" | jq -r '.error // "Unknown error"' 2>/dev/null)
        error "Kill switch activation failed: $error_msg"
    fi
}

# Countdown function
countdown() {
    local seconds=$1
    while [[ $seconds -gt 0 ]]; do
        echo -ne "${RED}💀 Destruction in: $seconds seconds\r${NC}"
        sleep 1
        ((seconds--))
    done
    echo -e "${RED}💀 DESTRUCTION INITIATED${NC}"
}

# Main execution flow
main() {
    local vps_id="${1:-}"
    local level="${2:-}"
    local delay="${3:-0}"
    local reason="${4:-Emergency manual activation}"
    
    # Show banner
    show_banner
    
    # Validate prerequisites
    validate_prerequisites
    
    # Interactive mode if parameters not provided
    if [[ -z "$vps_id" ]]; then
        vps_id=$(select_vps)
    fi
    
    if [[ -z "$level" ]]; then
        level=$(select_destruction_level)
    fi
    
    # Confirm destruction
    confirm_destruction "$vps_id" "$level"
    
    # Execute
    execute_kill_switch "$vps_id" "$level" "$delay" "$reason"
    
    # Final log
    log "Kill switch operation completed"
    success "Log saved to: $LOG_FILE"
}

# Handle script interruption
trap 'error "Kill switch operation interrupted"' INT TERM

# Ensure we have proper signal handling
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
