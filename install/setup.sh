#!/bin/bash

set -e

echo "🛠️ Starting Sauron setup..."

# Update & install base tools
apt update
apt install -y curl unzip git redis wget software-properties-common

# Install Go
GO_VERSION=1.22.3
if ! command -v go >/dev/null; then
  echo "📦 Installing Go $GO_VERSION..."
  wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz
  echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
  export PATH=$PATH:/usr/local/go/bin
fi
go version

# Ensure we're in the project root
if [ ! -f "main.go" ]; then
  echo "❌ Please run this script from the root of the Sauron project."
  exit 1
fi

# Build and install the binary
echo "🔨 Building Sauron binary..."
go mod download
go build -o /usr/local/bin/sauron

# Copy systemd service
echo "⚙️ Installing systemd service..."
cp install/sauron.service /etc/systemd/system/sauron.service

# Reload and start the service
systemctl daemon-reexec
systemctl daemon-reload
systemctl enable sauron.service
systemctl restart sauron.service

echo "✅ Sauron setup complete. Status:"
systemctl status sauron.service
