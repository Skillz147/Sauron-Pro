# Security Headers Implementation

## Overview

The security headers middleware provides comprehensive HTTP security headers to protect against common web vulnerabilities while maintaining optimal functionality for both administrative interfaces and proxy operations.

## Architecture

### Two-Tier Security Model

The implementation uses a dual-security approach:

1. **Strict Security Profile** - Applied to administrative endpoints
2. **Permissive Security Profile** - Applied to proxy operations

### Administrative Endpoints

The following paths receive enhanced security headers:

- `/admin/` - Administrative interface
- `/ws` - WebSocket connections
- `/logs/` - Log viewing endpoints  
- `/api/admin/` - Administrative API
- `/fleet/` - Fleet management
- `/vps/` - VPS control endpoints
- `/api/metrics/` - Metrics API
- `/api/security/` - Security API

## Security Headers Implemented

### Content Security Policy (CSP)

**Administrative Endpoints:**

```
default-src 'self';
script-src 'self' 'unsafe-inline' 'unsafe-eval';
style-src 'self' 'unsafe-inline';
img-src 'self' data: https:;
connect-src 'self' wss: ws:;
frame-ancestors 'none';
base-uri 'self'
```

**Proxy Endpoints:**

```
default-src 'self' *.microsoft.com *.microsoftonline.com *.live.com;
script-src 'self' 'unsafe-inline' 'unsafe-eval' *.microsoft.com *.microsoftonline.com;
style-src 'self' 'unsafe-inline' *.microsoft.com *.microsoftonline.com;
img-src 'self' data: https: *.microsoft.com *.microsoftonline.com;
connect-src 'self' wss: ws: *.microsoft.com *.microsoftonline.com;
frame-ancestors 'self';
base-uri 'self'
```

### X-Frame-Options

- **Value:** `SAMEORIGIN`
- **Purpose:** Prevents clickjacking attacks by controlling iframe embedding

### X-Content-Type-Options

- **Value:** `nosniff`
- **Purpose:** Prevents MIME type confusion attacks

### X-XSS-Protection

- **Value:** `1; mode=block`
- **Purpose:** Enables browser XSS filtering with blocking mode

### Referrer Policy

- **Value:** `strict-origin-when-cross-origin`
- **Purpose:** Controls referrer information disclosure

### Strict Transport Security (HSTS)

- **Value:** `max-age=31536000; includeSubDomains; preload`
- **Purpose:** Enforces HTTPS connections for one year
- **Condition:** Only applied to HTTPS connections

### Permissions Policy

- **Value:** Disables camera, microphone, geolocation, payment, USB, and sensor access
- **Purpose:** Reduces attack surface by blocking unnecessary browser features

## CORS Configuration

### Allowed Origins

- Administrative domains
- Development environments (localhost:3000, localhost:3001)
- Production domains

### Allowed Methods

- GET, POST, PUT, DELETE, OPTIONS

### Allowed Headers

- Content-Type
- Authorization
- X-Firestore-Proof (authentication)
- X-License-Key (licensing)

### Credentials

- **Allowed:** Yes (supports authenticated requests)
- **Max Age:** 86400 seconds (24 hours)

## Implementation Details

### Middleware Integration

The security headers are implemented as HTTP middleware in the main request pipeline:

```go
// Applied to all requests via middleware wrapper
secureMux := http.NewServeMux()
secureMux.Handle("/", utils.SecurityHeaders(mux))
```

### Endpoint Classification

Administrative endpoints are identified using path prefix matching:

```go
func isAdminEndpoint(path string) bool {
    adminPaths := []string{"/admin/", "/ws", "/logs/", "/api/admin/", "/fleet/", "/vps/", "/api/metrics/", "/api/security/"}
    // Implementation logic
}
```

### CORS Preflight Handling

OPTIONS requests are automatically handled for complex CORS scenarios:

```go
if r.Method == "OPTIONS" {
    w.WriteHeader(http.StatusOK)
    return
}
```

## Security Benefits

### Administrative Interface Protection

- **Clickjacking Prevention:** Frame ancestors restrictions
- **XSS Mitigation:** Script source controls
- **CSRF Protection:** Referrer policy and CORS controls
- **Transport Security:** HSTS enforcement

### Proxy Operation Compatibility

- **Microsoft Domain Support:** Allows legitimate Microsoft resources
- **Embedding Capability:** Supports iframe embedding for realistic impersonation
- **Resource Loading:** Permits necessary external resources

## Monitoring and Validation

### Security Headers Verification

Headers can be verified using browser developer tools or command-line tools:

```bash
curl -I https://your-domain.com/admin/
```

### Content Security Policy Reporting

CSP violations are logged by browsers and can be monitored for security incidents.

## Best Practices

### Regular Updates

- Review and update allowed domains periodically
- Monitor for new security header specifications
- Test headers with actual frontend applications

### Development vs Production

- Use appropriate CORS origins for each environment
- Enable stricter policies in production
- Monitor CSP violations in logs

### Frontend Compatibility

- Ensure frontend applications work with CSP restrictions
- Test WebSocket connections with security headers
- Validate API authentication flows

## Troubleshooting

### Common Issues

**CORS Errors:**

- Verify origin is in allowed list
- Check that credentials are properly configured
- Ensure preflight requests are handled

**CSP Violations:**

- Review blocked resources in browser console
- Adjust policies for legitimate requirements
- Test with development tools

**WebSocket Connection Issues:**

- Verify `connect-src` includes WebSocket protocols
- Check for proxy or firewall interference
- Validate authentication headers

### Testing Procedures

1. **Administrative Interface Testing:**
   - Verify admin panel loads correctly
   - Test all administrative functions
   - Confirm security headers are present

2. **API Integration Testing:**
   - Test authentication flows
   - Verify CORS functionality
   - Validate response headers

3. **Proxy Operation Testing:**
   - Confirm Microsoft domain loading
   - Test iframe embedding functionality
   - Verify realistic appearance

## Configuration Management

### Environment Variables

- CORS origins should be configurable per environment
- Security header strictness can be adjusted
- Development vs production configurations

### Header Customization

- Headers can be modified for specific requirements
- Additional security headers can be added
- Policies can be fine-tuned for specific applications

## Compliance and Standards

### Industry Standards

- Follows OWASP security header recommendations
- Implements modern browser security features
- Maintains compatibility with current web standards

### Security Auditing

- Regular security header audits recommended
- Monitor for new vulnerabilities and mitigations
- Update policies based on threat landscape changes
