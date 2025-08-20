package honeypot

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// HoneypotTemplate defines the structure for serving realistic business websites as honeypots
type HoneypotTemplate struct {
	templates map[string]*template.Template
}

// NewHoneypotTemplate creates a new template handler for honeypot responses
func NewHoneypotTemplate() *HoneypotTemplate {
	ht := &HoneypotTemplate{
		templates: make(map[string]*template.Template),
	}

	// Load Microsoft-style error page templates
	ht.loadTemplates()
	return ht
}

// loadTemplates precompiles realistic business website templates for honeypot intelligence gathering
func (ht *HoneypotTemplate) loadTemplates() {
	// High-end Tech Startup Landing Page - appears legitimate but secretly logs attackers
	techStartupTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NeuralScale AI - Revolutionizing Enterprise Intelligence</title>
    <meta name="description" content="Leading AI-powered enterprise solutions. Transform your business with cutting-edge artificial intelligence and machine learning technologies.">
    <meta name="keywords" content="AI, artificial intelligence, machine learning, enterprise, automation, neural networks, data analytics">
    <meta property="og:title" content="NeuralScale AI - Enterprise AI Solutions">
    <meta property="og:description" content="Transform your enterprise with cutting-edge AI technology">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif; line-height: 1.6; color: #333; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 1rem 0; position: fixed; width: 100%; top: 0; z-index: 1000; }
        .nav { max-width: 1200px; margin: 0 auto; display: flex; justify-content: space-between; align-items: center; padding: 0 2rem; }
        .logo { font-size: 1.5rem; font-weight: 700; }
        .nav-links { display: flex; list-style: none; gap: 2rem; }
        .nav-links a { color: white; text-decoration: none; transition: opacity 0.3s; }
        .nav-links a:hover { opacity: 0.8; }
        .hero { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 120px 2rem 80px; text-align: center; }
        .hero h1 { font-size: 3.5rem; font-weight: 800; margin-bottom: 1.5rem; }
        .hero p { font-size: 1.25rem; margin-bottom: 2rem; opacity: 0.9; }
        .cta-btn { background: #ff6b6b; color: white; padding: 15px 30px; border: none; border-radius: 50px; font-size: 1.1rem; font-weight: 600; cursor: pointer; transition: transform 0.3s; }
        .cta-btn:hover { transform: translateY(-2px); }
        .features { padding: 80px 2rem; max-width: 1200px; margin: 0 auto; }
        .features-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 3rem; margin-top: 3rem; }
        .feature { text-align: center; padding: 2rem; border-radius: 10px; box-shadow: 0 10px 30px rgba(0,0,0,0.1); }
        .feature-icon { font-size: 3rem; margin-bottom: 1rem; }
        .feature h3 { font-size: 1.5rem; margin-bottom: 1rem; color: #333; }
        .footer { background: #2c3e50; color: white; text-align: center; padding: 2rem; }
    </style>
</head>
<body>
    <header class="header">
        <nav class="nav">
            <div class="logo">NeuralScale AI</div>
            <ul class="nav-links">
                <li><a href="#solutions">Solutions</a></li>
                <li><a href="#about">About</a></li>
                <li><a href="#contact">Contact</a></li>
                <li><a href="#demo">Demo</a></li>
            </ul>
        </nav>
    </header>
    
    <section class="hero">
        <h1>The Future of Enterprise AI</h1>
        <p>Unlock unprecedented insights with our next-generation artificial intelligence platform</p>
        <button class="cta-btn" onclick="trackClick('cta_hero')">Start Free Trial</button>
    </section>
    
    <section class="features">
        <h2 style="text-align: center; font-size: 2.5rem; margin-bottom: 1rem;">Cutting-Edge Solutions</h2>
        <div class="features-grid">
            <div class="feature">
                <div class="feature-icon">🧠</div>
                <h3>Neural Analytics</h3>
                <p>Advanced machine learning algorithms that adapt and evolve with your business needs</p>
            </div>
            <div class="feature">
                <div class="feature-icon">⚡</div>
                <h3>Real-time Processing</h3>
                <p>Lightning-fast data processing with sub-millisecond response times</p>
            </div>
            <div class="feature">
                <div class="feature-icon">🔒</div>
                <h3>Enterprise Security</h3>
                <p>Bank-grade encryption and compliance with global security standards</p>
            </div>
        </div>
    </section>
    
    <footer class="footer">
        <p>&copy; 2024 NeuralScale AI. All rights reserved. | Session: {{.CorrelationID}}</p>
    </footer>
    
    <script>
        function trackClick(element) {
            console.log('User interaction tracked:', element);
            // This secretly logs all user interactions
        }
        
        // Track page views and user behavior
        window.addEventListener('load', function() {
            console.log('Page loaded - visitor tracked');
        });
    </script>
</body>
</html>`

	// Professional Cybersecurity Consulting Firm - looks legitimate but captures threat intelligence
	cybersecTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CyberShield Pro - Elite Cybersecurity Consulting</title>
    <meta name="description" content="Leading cybersecurity consulting firm. Protect your enterprise with advanced threat detection, penetration testing, and security architecture solutions.">
    <meta name="keywords" content="cybersecurity, penetration testing, threat detection, security consulting, CISO, compliance, red team">
    <meta property="og:title" content="CyberShield Pro - Enterprise Cybersecurity">
    <meta property="og:description" content="Elite cybersecurity consulting and threat intelligence services">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Roboto', Arial, sans-serif; line-height: 1.6; color: #2c3e50; background: #f8f9fa; }
        .header { background: #1a252f; color: white; padding: 1rem 0; position: sticky; top: 0; z-index: 1000; }
        .nav { max-width: 1200px; margin: 0 auto; display: flex; justify-content: space-between; align-items: center; padding: 0 2rem; }
        .logo { font-size: 1.8rem; font-weight: 700; color: #00d4aa; }
        .nav-menu { display: flex; list-style: none; gap: 2.5rem; }
        .nav-menu a { color: white; text-decoration: none; font-weight: 500; transition: color 0.3s; }
        .nav-menu a:hover { color: #00d4aa; }
        .hero { background: linear-gradient(rgba(26,37,47,0.9), rgba(26,37,47,0.9)), url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 1000"><rect fill="%23000" width="1000" height="1000"/><circle fill="%2300d4aa" cx="200" cy="200" r="2"/></svg>'); color: white; padding: 120px 2rem 100px; text-align: center; }
        .hero h1 { font-size: 3.2rem; font-weight: 800; margin-bottom: 1.5rem; }
        .hero p { font-size: 1.3rem; margin-bottom: 2.5rem; opacity: 0.9; max-width: 600px; margin-left: auto; margin-right: auto; }
        .cta-group { display: flex; gap: 1rem; justify-content: center; flex-wrap: wrap; }
        .btn-primary { background: #00d4aa; color: white; padding: 15px 30px; border: none; border-radius: 5px; font-size: 1.1rem; font-weight: 600; cursor: pointer; transition: all 0.3s; }
        .btn-secondary { background: transparent; color: white; padding: 15px 30px; border: 2px solid white; border-radius: 5px; font-size: 1.1rem; font-weight: 600; cursor: pointer; transition: all 0.3s; }
        .btn-primary:hover { background: #00b894; transform: translateY(-2px); }
        .btn-secondary:hover { background: white; color: #1a252f; }
        .services { padding: 100px 2rem; max-width: 1200px; margin: 0 auto; }
        .section-title { text-align: center; font-size: 2.8rem; margin-bottom: 3rem; color: #1a252f; }
        .services-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(350px, 1fr)); gap: 2.5rem; }
        .service-card { background: white; padding: 2.5rem; border-radius: 10px; box-shadow: 0 15px 35px rgba(0,0,0,0.1); transition: transform 0.3s; }
        .service-card:hover { transform: translateY(-5px); }
        .service-icon { font-size: 3.5rem; margin-bottom: 1.5rem; color: #00d4aa; }
        .service-card h3 { font-size: 1.6rem; margin-bottom: 1rem; color: #1a252f; }
        .trust-badges { background: #1a252f; color: white; padding: 60px 2rem; text-align: center; }
        .badges-grid { display: flex; justify-content: center; gap: 3rem; flex-wrap: wrap; margin-top: 2rem; }
        .badge { padding: 1rem; background: rgba(255,255,255,0.1); border-radius: 10px; }
        .footer { background: #0d1117; color: white; text-align: center; padding: 3rem 2rem; }
    </style>
</head>
<body>
    <header class="header">
        <nav class="nav">
            <div class="logo">🛡️ CyberShield Pro</div>
            <ul class="nav-menu">
                <li><a href="#services">Services</a></li>
                <li><a href="#expertise">Expertise</a></li>
                <li><a href="#case-studies">Case Studies</a></li>
                <li><a href="#contact">Contact</a></li>
            </ul>
        </nav>
    </header>
    
    <section class="hero">
        <h1>Elite Cybersecurity Consulting</h1>
        <p>Protect your enterprise with world-class threat intelligence, penetration testing, and security architecture from industry veterans</p>
        <div class="cta-group">
            <button class="btn-primary" onclick="trackEngagement('free_assessment')">Free Security Assessment</button>
            <button class="btn-secondary" onclick="trackEngagement('emergency_response')">24/7 Incident Response</button>
        </div>
    </section>
    
    <section class="services">
        <h2 class="section-title">Comprehensive Security Solutions</h2>
        <div class="services-grid">
            <div class="service-card">
                <div class="service-icon">🎯</div>
                <h3>Advanced Penetration Testing</h3>
                <p>Red team exercises and comprehensive vulnerability assessments to identify and eliminate security gaps before attackers do</p>
            </div>
            <div class="service-card">
                <div class="service-icon">🔍</div>
                <h3>Threat Intelligence</h3>
                <p>Real-time threat monitoring and intelligence gathering to stay ahead of emerging cyber threats and attack vectors</p>
            </div>
            <div class="service-card">
                <div class="service-icon">⚡</div>
                <h3>Incident Response</h3>
                <p>24/7 emergency response team with forensic capabilities to contain breaches and minimize business impact</p>
            </div>
        </div>
    </section>
    
    <section class="trust-badges">
        <h2>Trusted by Global Enterprises</h2>
        <div class="badges-grid">
            <div class="badge">ISO 27001 Certified</div>
            <div class="badge">CISSP Experts</div>
            <div class="badge">Fortune 500 Clients</div>
            <div class="badge">99.9% Success Rate</div>
        </div>
    </section>
    
    <footer class="footer">
        <p>&copy; 2024 CyberShield Pro Consulting. All rights reserved. | Request ID: {{.CorrelationID}}</p>
        <p style="margin-top: 0.5rem; font-size: 0.9rem; opacity: 0.7;">Protecting enterprises worldwide since 2018</p>
    </footer>
    
    <script>
        function trackEngagement(action) {
            console.log('Security consultation requested:', action);
            // Honeypot intelligence gathering
        }
        
        // Advanced visitor tracking
        document.addEventListener('DOMContentLoaded', function() {
            console.log('CyberShield visitor logged:', new Date().toISOString());
        });
    </script>
</body>
</html>`

	// Premium Digital Marketing Agency - sophisticated design that secretly logs visitors
	marketingTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Elevate Digital - Premium Marketing Agency</title>
    <meta name="description" content="Award-winning digital marketing agency. Drive growth with data-driven campaigns, creative excellence, and innovative marketing strategies.">
    <meta name="keywords" content="digital marketing, SEO, social media marketing, PPC, content marketing, brand strategy, marketing automation">
    <meta property="og:title" content="Elevate Digital - Premium Marketing Solutions">
    <meta property="og:description" content="Transform your business with award-winning digital marketing strategies">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Poppins', Arial, sans-serif; line-height: 1.6; color: #2d3748; }
        .header { background: #ffffff; box-shadow: 0 2px 10px rgba(0,0,0,0.1); padding: 1rem 0; position: fixed; width: 100%; top: 0; z-index: 1000; }
        .nav { max-width: 1200px; margin: 0 auto; display: flex; justify-content: space-between; align-items: center; padding: 0 2rem; }
        .logo { font-size: 2rem; font-weight: 800; background: linear-gradient(45deg, #ff6b6b, #4ecdc4); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
        .nav-menu { display: flex; list-style: none; gap: 2rem; }
        .nav-menu a { color: #2d3748; text-decoration: none; font-weight: 500; transition: color 0.3s; }
        .nav-menu a:hover { color: #ff6b6b; }
        .hero { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 140px 2rem 100px; text-align: center; overflow: hidden; position: relative; }
        .hero::before { content: ''; position: absolute; top: 0; left: 0; right: 0; bottom: 0; background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><circle fill="rgba(255,255,255,0.1)" cx="20" cy="20" r="2"/><circle fill="rgba(255,255,255,0.1)" cx="80" cy="80" r="3"/></svg>'); }
        .hero-content { position: relative; z-index: 2; }
        .hero h1 { font-size: 4rem; font-weight: 800; margin-bottom: 1.5rem; }
        .hero p { font-size: 1.4rem; margin-bottom: 2.5rem; opacity: 0.9; max-width: 700px; margin-left: auto; margin-right: auto; }
        .cta-btn { background: linear-gradient(45deg, #ff6b6b, #ee5a52); color: white; padding: 18px 40px; border: none; border-radius: 50px; font-size: 1.2rem; font-weight: 700; cursor: pointer; transition: all 0.3s; box-shadow: 0 10px 30px rgba(255,107,107,0.3); }
        .cta-btn:hover { transform: translateY(-3px); box-shadow: 0 15px 40px rgba(255,107,107,0.4); }
        .features { padding: 120px 2rem; background: #f7fafc; }
        .container { max-width: 1200px; margin: 0 auto; }
        .section-title { text-align: center; font-size: 3rem; margin-bottom: 4rem; color: #2d3748; }
        .features-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 3rem; }
        .feature { background: white; padding: 3rem; border-radius: 20px; text-align: center; box-shadow: 0 20px 40px rgba(0,0,0,0.1); transition: transform 0.3s; }
        .feature:hover { transform: translateY(-10px); }
        .feature-icon { font-size: 4rem; margin-bottom: 2rem; }
        .feature h3 { font-size: 1.8rem; margin-bottom: 1.5rem; color: #2d3748; }
        .stats { background: linear-gradient(45deg, #ff6b6b, #4ecdc4); color: white; padding: 80px 2rem; text-align: center; }
        .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 2rem; margin-top: 3rem; }
        .stat { padding: 2rem; }
        .stat-number { font-size: 3.5rem; font-weight: 800; margin-bottom: 0.5rem; }
        .stat-label { font-size: 1.2rem; opacity: 0.9; }
        .footer { background: #1a202c; color: white; text-align: center; padding: 4rem 2rem; }
    </style>
</head>
<body>
    <header class="header">
        <nav class="nav">
            <div class="logo">Elevate Digital</div>
            <ul class="nav-menu">
                <li><a href="#services">Services</a></li>
                <li><a href="#portfolio">Portfolio</a></li>
                <li><a href="#about">About</a></li>
                <li><a href="#contact">Contact</a></li>
            </ul>
        </nav>
    </header>
    
    <section class="hero">
        <div class="hero-content">
            <h1>Marketing That Moves Mountains</h1>
            <p>Award-winning digital marketing strategies that transform brands and drive exponential growth through data-driven excellence</p>
            <button class="cta-btn" onclick="trackConversion('hero_consultation')">Get Free Marketing Audit</button>
        </div>
    </section>
    
    <section class="features">
        <div class="container">
            <h2 class="section-title">Results-Driven Solutions</h2>
            <div class="features-grid">
                <div class="feature">
                    <div class="feature-icon">📈</div>
                    <h3>Performance Marketing</h3>
                    <p>ROI-focused campaigns with advanced analytics and conversion optimization that deliver measurable business growth</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">🎨</div>
                    <h3>Creative Excellence</h3>
                    <p>Award-winning creative campaigns that capture attention, build brand loyalty, and drive customer engagement</p>
                </div>
                <div class="feature">
                    <div class="feature-icon">🔍</div>
                    <h3>Data Intelligence</h3>
                    <p>Advanced market research and consumer insights that inform strategic decisions and maximize campaign effectiveness</p>
                </div>
            </div>
        </div>
    </section>
    
    <section class="stats">
        <h2>Delivering Exceptional Results</h2>
        <div class="stats-grid">
            <div class="stat">
                <div class="stat-number">250%</div>
                <div class="stat-label">Average ROI Increase</div>
            </div>
            <div class="stat">
                <div class="stat-number">500+</div>
                <div class="stat-label">Successful Campaigns</div>
            </div>
            <div class="stat">
                <div class="stat-number">98%</div>
                <div class="stat-label">Client Retention Rate</div>
            </div>
            <div class="stat">
                <div class="stat-number">24/7</div>
                <div class="stat-label">Campaign Monitoring</div>
            </div>
        </div>
    </section>
    
    <footer class="footer">
        <p>&copy; 2024 Elevate Digital Marketing Agency. All rights reserved. | Campaign ID: {{.CorrelationID}}</p>
        <p style="margin-top: 1rem; font-size: 0.9rem; opacity: 0.7;">Driving business growth through innovative marketing since 2019</p>
    </footer>
    
    <script>
        function trackConversion(source) {
            console.log('Marketing consultation requested from:', source);
            // Secret visitor tracking for honeypot intelligence
        }
        
        // Track visitor engagement
        window.addEventListener('load', function() {
            console.log('Elevate Digital visitor tracked:', document.referrer);
        });
    </script>
</body>
</html>`

	// Parse templates
	var err error
	ht.templates["techstartup"], err = template.New("techstartup").Parse(techStartupTemplate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse tech startup template")
	}

	ht.templates["cybersec"], err = template.New("cybersec").Parse(cybersecTemplate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse cybersec template")
	}

	ht.templates["marketing"], err = template.New("marketing").Parse(marketingTemplate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse marketing template")
	}

	log.Info().Int("templates", len(ht.templates)).Msg("Honeypot templates loaded")
}

// ServeBusinessSite serves a realistic business website to the visitor while secretly logging them
func (ht *HoneypotTemplate) ServeBusinessSite(w http.ResponseWriter, r *http.Request, templateName string) {
	template, exists := ht.templates[templateName]
	if !exists {
		// Fallback to tech startup template
		template = ht.templates["techstartup"]
		if template == nil {
			http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
			return
		}
	}

	// Generate realistic business data
	data := map[string]interface{}{
		"CorrelationID": generateCorrelationID(),
		"Timestamp":     time.Now().Format("2006-01-02 15:04:05 UTC"),
		"UserAgent":     r.UserAgent(),
		"VisitorIP":     r.RemoteAddr,
	}

	// Set headers to appear as a legitimate business website
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("Server", "nginx/1.21.6")

	// Return 200 OK - appears as a normal legitimate website
	w.WriteHeader(http.StatusOK)

	err := template.Execute(w, data)
	if err != nil {
		log.Error().Err(err).Str("template", templateName).Msg("Failed to execute honeypot business template")
		http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
	}
}

// generateCorrelationID creates realistic tracking IDs for visitor logging
func generateCorrelationID() string {
	// Generate a realistic UUID-style tracking ID for the honeypot
	timestamp := time.Now().Unix()
	return generateUUIDFromTimestamp(timestamp)
}

// generateUUIDFromTimestamp creates a deterministic UUID from timestamp for tracking
func generateUUIDFromTimestamp(timestamp int64) string {
	// Create a deterministic but realistic UUID for visitor tracking
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		timestamp&0xffffffff,
		(timestamp>>32)&0xffff,
		((timestamp>>48)&0x0fff)|0x4000,
		((timestamp>>32)&0x3fff)|0x8000,
		timestamp&0xffffffffffff)
}
