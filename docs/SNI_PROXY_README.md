# SNI Reverse Proxy Setup

This setup allows both Sauron and Shield to use port 443 for clean URLs without conflicts.

## Architecture

```
Internet (Port 443)
       ↓
   SNI Proxy (Port 443)
       ↓
   ┌─────────┬─────────┐
   ↓         ↓         ↓
Sauron   Shield   Fallback
:8080    :8444    to Sauron
```

## How It Works

1. **SNI Proxy** listens on port 443 with SNI support
2. **Sauron** runs on internal port 8080
3. **Shield** runs on internal port 8444 (dev) or 8443 (prod)
4. **SNI Proxy** routes traffic based on domain:
   - `*.sauron-domain.com` → Sauron (port 8080)
   - `*.shield-domain.com` → Shield (port 8444/8443)
   - Unknown domains → Sauron (fallback)

## Setup

### 1. Environment Variables

Ensure your `.env` file has:
```bash
SAURON_DOMAIN=your-sauron-domain.com
SHIELD_DOMAIN=your-shield-domain.com
DEV_MODE=true
# ... other existing variables
```

### 2. Start Services

```bash
# Start all services
./start-sni-proxy.sh

# Or start individually:
# Terminal 1: Sauron
go run main.go

# Terminal 2: Shield
cd shield-domain && go run main.go

# Terminal 3: SNI Proxy
cd sni-proxy && go run main.go
```

### 3. Test URLs

- **Sauron**: `https://your-sauron-domain.com`
- **Shield**: `https://your-shield-domain.com`
- **Subdomains**: `https://login.your-sauron-domain.com`

## Benefits

✅ **Clean URLs**: No port numbers in campaign URLs
✅ **No Conflicts**: Both services can use port 443
✅ **SNI Support**: Correct certificates served based on domain
✅ **Fallback**: Unknown domains default to Sauron
✅ **Internal Services**: Sauron and Shield run on internal ports

## Troubleshooting

### Check Service Status
```bash
# Check if services are running
netstat -tlnp | grep -E ":(443|8080|8444|8443)"

# Check logs
tail -f logs/system.log
```

### Common Issues

1. **Port 443 Permission**: Run with `sudo` if needed
2. **Certificate Issues**: Ensure certificates are loaded properly
3. **Domain Resolution**: Check `/etc/hosts` in dev mode

## Development vs Production

- **Dev Mode**: Shield uses port 8444, Sauron uses port 8080
- **Production**: Shield uses port 8443, Sauron uses port 8080
- **SNI Proxy**: Always uses port 443 (external)
