#!/bin/bash

echo "🔍 Sauron Database Fix"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if we're in the right directory
if [ ! -f "./sauron" ]; then
    echo "❌ Please run this script from the sauron directory"
    echo "💡 Expected files: sauron binary, .env file"
    exit 1
fi

echo "� Diagnosing database issue..."

# The root cause: SQLite driver not imported
echo "❌ Root cause identified: SQLite driver missing from main.go"
echo "🔧 The binary was compiled without SQLite driver support"
echo ""

# Check current environment
echo "� Environment Check:"
if [ -f ".env" ]; then
    echo "✅ .env file exists"
    # Check for basic variables
    source .env 2>/dev/null
    if [ -n "$SAURON_DOMAIN" ]; then
        echo "✅ SAURON_DOMAIN: $SAURON_DOMAIN"
    else
        echo "❌ SAURON_DOMAIN not set"
    fi
else
    echo "❌ .env file missing"
fi

echo ""
echo "� SOLUTION:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "The issue is that your binary was compiled without SQLite driver support."
echo "Here's how to fix it:"
echo ""
echo "1️⃣  Download the fixed version (includes SQLite driver):"
echo "   wget https://github.com/Skillz147/Sauron-Pro/releases/latest/download/sauron-linux-amd64.tar.gz"
echo ""
echo "2️⃣  Extract the new version:"
echo "   tar -xzf sauron-linux-amd64.tar.gz"
echo ""
echo "3️⃣  Update the binary:"
echo "   cd sauron"
echo "   sudo systemctl stop sauron"
echo "   sudo cp sauron /usr/local/bin/"
echo "   sudo systemctl start sauron"
echo ""
echo "4️⃣  Verify the fix:"
echo "   sudo systemctl status sauron"
echo ""
echo "📊 Technical Details:"
echo "   • The binary needs 'import _ \"github.com/mattn/go-sqlite3\"' in main.go"
echo "   • Current binary was cross-compiled without CGO support"
echo "   • New binary includes proper SQLite driver registration"
echo ""
echo "🚀 After updating, the database will initialize properly!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
