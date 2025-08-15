# 🚨 Bad Customer Detection System

## Overview

The Bad Customer Detection System automatically monitors your framework customers for suspicious behavior patterns that could indicate:

- **Law enforcement operations** targeting your service
- **Bad actors** using your framework inappropriately  
- **Automated attacks** or abuse of your infrastructure
- **High-risk targets** that could compromise your service

## 🎯 Detection Methods

### 1. **Volume-Based Anomalies**

- High traffic volume in short periods (spray attacks)
- Massive credential capture attempts
- Unusually high success rates (suspicious)

### 2. **Target-Based Red Flags**

- Government domains (.gov, .mil)
- Law enforcement agencies
- Critical infrastructure
- Security companies

### 3. **Technical Signatures**

- Known law enforcement IP ranges
- Bot/automation tool signatures
- Rapid IP address cycling
- Honeypot interactions

### 4. **Behavioral Patterns**

- Geographic distribution anomalies
- Timing patterns suggesting coordination
- Framework fingerprinting attempts

## 📊 Risk Scoring

Customers are automatically scored based on detected patterns:

- **LOW (0-29)**: Normal customer behavior
- **MEDIUM (30-59)**: Some suspicious activity
- **HIGH (60-99)**: Significant risk indicators  
- **CRITICAL (100+)**: Automatic blocking triggered

## 🔧 API Endpoints

### Get Customer Risk Report

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     https://your-domain.com/admin/risk
```

### Get High-Risk Customers Only

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     https://your-domain.com/admin/risk?high_risk=true
```

### Get Specific Customer Risk

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     https://your-domain.com/admin/risk?slug=CUSTOMER_SLUG
```

### Take Action on Customer

```bash
curl -X POST -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"action":"block","slug":"CUSTOMER_SLUG","reason":"Gov targeting"}' \
     https://your-domain.com/admin/risk
```

Available actions: `block`, `investigate`, `whitelist`

### Get Metrics Dashboard

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     https://your-domain.com/admin/metrics
```

## 🖥️ Admin Dashboard

1. **Open** `admin_dashboard.html` in your browser
2. **Update** the `API_TOKEN` variable with your admin token
3. **Monitor** customer risks in real-time

The dashboard shows:

- Risk distribution statistics
- High-risk customer details  
- Security alerts
- Quick action buttons

## ⚙️ Configuration

### Customizing Risk Rules

Edit `monitoring/customer_monitor.go` to adjust risk patterns:

```go
// Example: Add new risk rule
{
    Name:        "crypto_targeting",
    Pattern:     `coinbase|binance|crypto|blockchain`,
    RiskWeight:  35,
    Threshold:   5,
    Description: "Targeting cryptocurrency platforms",
    Enabled:     true,
}
```

### Law Enforcement IP Ranges

Update the `isLawEnforcementIP()` function with current intelligence:

```go
leRanges := []string{
    "192.52.178.0/24", // FBI
    "149.101.0.0/16",  // DHS  
    "204.248.25.0/24", // DOJ
    // Add your intelligence here
}
```

## 🚨 Alert Types

The system generates these security alerts:

- **Government Targeting**: Customer attacking .gov domains
- **Law Enforcement Contact**: LE IP accessing customer URLs
- **Honeypot Detection**: Interaction with honeypot systems
- **Critical Customer**: Customer exceeds critical risk threshold

## 📈 Integration

The monitoring system automatically integrates with:

- **MITM Proxy**: Tracks all customer URL hits
- **Credential Handlers**: Records capture success rates
- **WebSocket System**: Real-time risk updates
- **Firestore**: Persistent risk data storage

## 🛡️ Automatic Protection

When customers are detected as high-risk:

1. **Automatic blocking** for critical risk customers
2. **Enhanced logging** of all their activities  
3. **Real-time alerts** to admin dashboard
4. **Decoy traffic** to confuse their analysis

## 💡 Best Practices

### Regular Monitoring

- Check dashboard daily for new alerts
- Review high-risk customers weekly
- Update IP intelligence monthly

### Investigation Workflow  

1. **Investigate** suspicious customers before blocking
2. **Document** reasons for all actions taken
3. **Whitelist** legitimate customers if false positive

### Security Operations

- Use secure admin tokens
- Monitor from secure networks only
- Keep IP intelligence updated
- Regular backup of customer data

## 🔒 OPSEC Considerations

- Admin dashboard accessible only from trusted IPs
- Use VPN when accessing admin functions
- Rotate admin tokens regularly
- Monitor for admin endpoint access attempts

This system helps protect your framework service while maintaining professional customer experience for legitimate users.
