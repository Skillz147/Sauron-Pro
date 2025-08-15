#!/bin/bash

# Admin Cleanup Script for Sauron MITM Framework
# Usage: ./admin_cleanup.sh [operation] [retention_days] [--dry-run]

set -e

# Configuration
ADMIN_KEY="${ADMIN_KEY:-$(grep 'ADMIN_KEY=' .env 2>/dev/null | cut -d'=' -f2)}"
SERVER_URL="${SERVER_URL:-https://localhost}"
RETENTION_DAYS="${2:-30}"
DRY_RUN=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if dry run is requested
if [[ "$3" == "--dry-run" ]] || [[ "$2" == "--dry-run" ]]; then
    DRY_RUN="true"
    echo -e "${YELLOW}🧪 DRY RUN MODE - No actual changes will be made${NC}"
fi

# Validate admin key
if [[ -z "$ADMIN_KEY" ]]; then
    echo -e "${RED}❌ Error: ADMIN_KEY not found${NC}"
    echo "Set ADMIN_KEY environment variable or add to .env file"
    exit 1
fi

# Function to call cleanup API
cleanup_api() {
    local operation="$1"
    local retention="$2"
    local dry_run="$3"
    
    echo -e "${BLUE}🧹 Initiating cleanup operation: $operation${NC}"
    
    curl -s -X POST "$SERVER_URL/admin/cleanup" \
        -H "Content-Type: application/json" \
        -d "{
            \"admin_key\": \"$ADMIN_KEY\",
            \"operations\": [\"$operation\"],
            \"retention_days\": $retention,
            \"dry_run\": $([ "$dry_run" = "true" ] && echo "true" || echo "false")
        }" | jq '.'
}

# Function to get cleanup status
cleanup_status() {
    echo -e "${BLUE}📊 Getting cleanup status...${NC}"
    
    curl -s -X GET "$SERVER_URL/admin/cleanup/status" \
        -H "X-Admin-Key: $ADMIN_KEY" | jq '.'
}

# Function to show usage
show_usage() {
    echo -e "${GREEN}Admin Cleanup Script for Sauron MITM Framework${NC}"
    echo ""
    echo "Usage: $0 [operation] [retention_days] [--dry-run]"
    echo ""
    echo "Operations:"
    echo "  logs        - Clean up old log files"
    echo "  database    - Clean up old database records"
    echo "  credentials - Clear captured credentials from memory"
    echo "  firestore   - Clean up old Firestore documents"
    echo "  all         - Run all cleanup operations"
    echo "  status      - Show current system status"
    echo ""
    echo "Options:"
    echo "  retention_days  - Keep data newer than N days (default: 30, 0 = delete all)"
    echo "  --dry-run      - Preview what would be deleted without making changes"
    echo ""
    echo "Examples:"
    echo "  $0 logs 7                  # Remove log files older than 7 days"
    echo "  $0 database 30 --dry-run   # Preview database cleanup (30 days retention)"
    echo "  $0 all 0                   # Delete all cleanable data (DANGER!)"
    echo "  $0 status                  # Show current system status"
}

# Main logic
case "${1:-help}" in
    "logs"|"database"|"credentials"|"firestore"|"all")
        if [[ "$1" == "all" ]] && [[ "$RETENTION_DAYS" == "0" ]] && [[ "$DRY_RUN" != "true" ]]; then
            echo -e "${RED}⚠️  WARNING: This will delete ALL data!${NC}"
            echo -e "${YELLOW}Are you sure? Type 'DELETE_ALL' to confirm:${NC}"
            read -r confirmation
            if [[ "$confirmation" != "DELETE_ALL" ]]; then
                echo -e "${GREEN}❌ Operation cancelled${NC}"
                exit 0
            fi
        fi
        cleanup_api "$1" "$RETENTION_DAYS" "$DRY_RUN"
        ;;
    "status")
        cleanup_status
        ;;
    "help"|"-h"|"--help"|*)
        show_usage
        ;;
esac
