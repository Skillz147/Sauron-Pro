#!/bin/bash
cd "$(dirname "$0")/.." || exit 1

echo "📥 Checking for updates..."

# Backup current .env file if it exists
if [ -f ".env" ]; then
  echo "💾 Backing up current .env configuration..."
  cp .env .env.backup.$(date +%Y%m%d_%H%M%S)
fi

# Backup current configuration database if it exists
if [ -f "config.db" ]; then
  echo "💾 Backing up current config.db..."
  cp config.db config.db.backup.$(date +%Y%m%d_%H%M%S)
fi

git fetch origin main
git reset --hard origin/main

# Restore .env file if backup exists
if [ -f ".env.backup."* ]; then
  latest_backup=$(ls -t .env.backup.* | head -n1)
  echo "🔄 Restoring .env configuration from $latest_backup"
  cp "$latest_backup" .env
fi

# Restore config.db if backup exists
if [ -f "config.db.backup."* ]; then
  latest_backup=$(ls -t config.db.backup.* | head -n1)
  echo "🔄 Restoring config.db from $latest_backup"
  cp "$latest_backup" config.db
fi

echo "✅ Updated to latest version."

# Validate configuration if possible
if [ -f ".env" ] && [ -f "scripts/configure-env.sh" ]; then
  echo "🔍 Validating environment configuration..."
  if ./scripts/configure-env.sh validate >/dev/null 2>&1; then
    echo "✅ Environment configuration is valid"
  else
    echo "⚠️  Environment configuration needs attention:"
    ./scripts/configure-env.sh validate
    echo ""
    echo "Run './scripts/configure-env.sh setup' to fix configuration"
  fi
fi

# Optional: rebuild and restart
go build -o sauron .
sudo systemctl restart sauron