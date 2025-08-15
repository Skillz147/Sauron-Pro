package decoy

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"o365/utils"

	"github.com/rs/zerolog/log"
)

// DecoyTrafficManager generates realistic background traffic to confuse analysis
type DecoyTrafficManager struct {
	mu         sync.RWMutex
	enabled    bool
	intensity  float64 // 0.0 to 1.0 (low to high traffic volume)
	domains    []string
	userAgents []string
	ipPool     []string
	patterns   []TrafficPattern
	stopChan   chan struct{}
	running    bool
}

// TrafficPattern defines a realistic user behavior pattern
type TrafficPattern struct {
	Name        string
	MinRequests int
	MaxRequests int
	Timing      TimingPattern
	URLs        []string
	Methods     []string
	Headers     map[string][]string
	Weight      float64 // Probability of this pattern being selected
}

// TimingPattern defines request timing behavior
type TimingPattern struct {
	MinDelay    time.Duration
	MaxDelay    time.Duration
	BurstChance float64 // Probability of burst requests
	BurstSize   int     // Number of requests in a burst
}

var (
	defaultManager *DecoyTrafficManager
	initOnce       sync.Once
)

// GetDecoyManager returns the global decoy traffic manager
func GetDecoyManager() *DecoyTrafficManager {
	initOnce.Do(func() {
		defaultManager = &DecoyTrafficManager{
			enabled:    false,
			intensity:  0.3, // Default medium intensity
			stopChan:   make(chan struct{}),
			domains:    getDefaultDomains(),
			userAgents: getRealisticUserAgents(),
			ipPool:     generateIPPool(),
			patterns:   getDefaultPatterns(),
		}
	})
	return defaultManager
}

// StartDecoyTraffic begins generating background traffic
func (d *DecoyTrafficManager) StartDecoyTraffic() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.running {
		return fmt.Errorf("decoy traffic already running")
	}

	d.enabled = true
	d.running = true

	// Start multiple goroutines for different traffic patterns
	go d.generateLegitimateTraffic()
	go d.generateSearchTraffic()
	go d.generateSocialMediaTraffic()
	go d.generateCorporateTraffic()
	go d.generateMaintenanceRequests()

	log.Info().
		Float64("intensity", d.intensity).
		Int("patterns", len(d.patterns)).
		Msg("🎭 Decoy traffic generation started")

	return nil
}

// StopDecoyTraffic stops background traffic generation
func (d *DecoyTrafficManager) StopDecoyTraffic() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.running {
		return
	}

	close(d.stopChan)
	d.enabled = false
	d.running = false

	log.Info().Msg("🛑 Decoy traffic generation stopped")
}

// IsRunning returns whether decoy traffic is currently running
func (d *DecoyTrafficManager) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.running
}

// GetIntensity returns the current traffic intensity
func (d *DecoyTrafficManager) GetIntensity() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.intensity
}

// GetPatternCount returns the number of configured traffic patterns
func (d *DecoyTrafficManager) GetPatternCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.patterns)
}

// SetIntensity adjusts the volume of decoy traffic (0.0 to 1.0)
func (d *DecoyTrafficManager) SetIntensity(intensity float64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if intensity < 0.0 {
		intensity = 0.0
	}
	if intensity > 1.0 {
		intensity = 1.0
	}

	d.intensity = intensity
	log.Info().Float64("intensity", intensity).Msg("🎭 Decoy traffic intensity adjusted")
}

// generateLegitimateTraffic creates normal web browsing patterns
func (d *DecoyTrafficManager) generateLegitimateTraffic() {
	ticker := time.NewTicker(time.Duration(float64(30*time.Second) / d.intensity))
	defer ticker.Stop()

	for {
		select {
		case <-d.stopChan:
			return
		case <-ticker.C:
			d.simulateWebBrowsing()
		}
	}
}

// generateSearchTraffic simulates search engine queries
func (d *DecoyTrafficManager) generateSearchTraffic() {
	ticker := time.NewTicker(time.Duration(float64(45*time.Second) / d.intensity))
	defer ticker.Stop()

	searchEngines := []string{
		"www.google.com",
		"www.bing.com",
		"duckduckgo.com",
	}

	queries := []string{
		"office 365 login",
		"microsoft teams",
		"outlook email",
		"sharepoint access",
		"azure portal",
		"microsoft authentication",
		"business productivity tools",
		"cloud collaboration",
	}

	for {
		select {
		case <-d.stopChan:
			return
		case <-ticker.C:
			engine := searchEngines[rand.Intn(len(searchEngines))]
			query := queries[rand.Intn(len(queries))]
			d.makeDecoyRequest("GET", "https://"+engine+"/search?q="+strings.ReplaceAll(query, " ", "+"), nil)
		}
	}
}

// generateSocialMediaTraffic simulates social media activity
func (d *DecoyTrafficManager) generateSocialMediaTraffic() {
	ticker := time.NewTicker(time.Duration(float64(60*time.Second) / d.intensity))
	defer ticker.Stop()

	socialSites := []string{
		"www.linkedin.com",
		"twitter.com",
		"www.facebook.com",
	}

	for {
		select {
		case <-d.stopChan:
			return
		case <-ticker.C:
			site := socialSites[rand.Intn(len(socialSites))]
			d.makeDecoyRequest("GET", "https://"+site+"/feed", nil)
		}
	}
}

// generateCorporateTraffic simulates corporate website visits
func (d *DecoyTrafficManager) generateCorporateTraffic() {
	ticker := time.NewTicker(time.Duration(float64(90*time.Second) / d.intensity))
	defer ticker.Stop()

	corporateSites := []string{
		"www.microsoft.com",
		"office.com",
		"portal.azure.com",
		"admin.microsoft.com",
		"security.microsoft.com",
		"compliance.microsoft.com",
	}

	for {
		select {
		case <-d.stopChan:
			return
		case <-ticker.C:
			site := corporateSites[rand.Intn(len(corporateSites))]
			d.makeDecoyRequest("GET", "https://"+site, nil)
		}
	}
}

// generateMaintenanceRequests simulates automated system requests
func (d *DecoyTrafficManager) generateMaintenanceRequests() {
	ticker := time.NewTicker(time.Duration(float64(120*time.Second) / d.intensity))
	defer ticker.Stop()

	for {
		select {
		case <-d.stopChan:
			return
		case <-ticker.C:
			// Simulate health checks, API calls, etc.
			d.makeDecoyRequest("GET", "https://status.office.com/api/v1.0/status", nil)
		}
	}
}

// simulateWebBrowsing creates realistic browsing session
func (d *DecoyTrafficManager) simulateWebBrowsing() {
	pattern := d.selectRandomPattern()
	numRequests := rand.Intn(pattern.MaxRequests-pattern.MinRequests+1) + pattern.MinRequests

	for i := 0; i < numRequests; i++ {
		// Random delay between requests
		delay := time.Duration(rand.Int63n(int64(pattern.Timing.MaxDelay-pattern.Timing.MinDelay))) + pattern.Timing.MinDelay
		time.Sleep(delay)

		// Select random URL and method
		url := pattern.URLs[rand.Intn(len(pattern.URLs))]
		method := pattern.Methods[rand.Intn(len(pattern.Methods))]

		// Add realistic headers
		headers := make(map[string]string)
		for k, v := range pattern.Headers {
			headers[k] = v[rand.Intn(len(v))]
		}

		d.makeDecoyRequestWithHeaders(method, url, headers, nil)

		select {
		case <-d.stopChan:
			return
		default:
		}
	}
}

// makeDecoyRequest performs a single decoy HTTP request
func (d *DecoyTrafficManager) makeDecoyRequest(method, url string, body []byte) {
	d.makeDecoyRequestWithHeaders(method, url, nil, body)
}

// makeDecoyRequestWithHeaders performs a decoy request with custom headers
func (d *DecoyTrafficManager) makeDecoyRequestWithHeaders(method, url string, headers map[string]string, body []byte) {
	if !d.enabled {
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return
	}

	// Add realistic headers
	req.Header.Set("User-Agent", d.getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Add custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Simulate different source IPs (for logging purposes)
	if fakeIP := d.getRandomIP(); fakeIP != "" {
		req.Header.Set("X-Forwarded-For", fakeIP)
	}

	// Make the request (ignore response for decoy traffic)
	resp, err := client.Do(req)
	if err == nil && resp != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	// Log decoy traffic (with special marker for filtering)
	utils.SystemLogger.Debug().
		Str("method", method).
		Str("url", url).
		Str("type", "decoy").
		Msg("🎭 Generated decoy request")
}

// selectRandomPattern chooses a traffic pattern based on weights
func (d *DecoyTrafficManager) selectRandomPattern() TrafficPattern {
	if len(d.patterns) == 0 {
		return getDefaultPatterns()[0]
	}

	totalWeight := 0.0
	for _, p := range d.patterns {
		totalWeight += p.Weight
	}

	r := rand.Float64() * totalWeight
	currentWeight := 0.0

	for _, p := range d.patterns {
		currentWeight += p.Weight
		if r <= currentWeight {
			return p
		}
	}

	return d.patterns[0]
}

// getRandomUserAgent returns a realistic user agent string
func (d *DecoyTrafficManager) getRandomUserAgent() string {
	if len(d.userAgents) == 0 {
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	}
	return d.userAgents[rand.Intn(len(d.userAgents))]
}

// getRandomIP returns a realistic IP address for header spoofing
func (d *DecoyTrafficManager) getRandomIP() string {
	if len(d.ipPool) == 0 {
		return ""
	}
	return d.ipPool[rand.Intn(len(d.ipPool))]
}

// Configuration functions
func getDefaultDomains() []string {
	return []string{
		"office.com",
		"microsoft.com",
		"outlook.com",
		"teams.microsoft.com",
		"portal.azure.com",
	}
}

func getRealisticUserAgents() []string {
	return []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/120.0.0.0 Safari/537.36",
	}
}

func generateIPPool() []string {
	// Generate realistic corporate IP ranges
	var ips []string

	// Common corporate IP ranges
	corporateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"203.0.113.0/24", // RFC5737 test range
	}

	// Generate sample IPs for logging
	for _, cidr := range corporateRanges {
		for i := 0; i < 50; i++ {
			ip := generateRandomIPInRange(cidr)
			if ip != "" {
				ips = append(ips, ip)
			}
		}
	}

	return ips
}

func generateRandomIPInRange(cidr string) string {
	// Simplified IP generation for decoy purposes
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return ""
	}

	baseIP := parts[0]
	octets := strings.Split(baseIP, ".")
	if len(octets) != 4 {
		return ""
	}

	// Randomize last octet for variety
	lastOctet := rand.Intn(254) + 1
	return fmt.Sprintf("%s.%s.%s.%d", octets[0], octets[1], octets[2], lastOctet)
}

func getDefaultPatterns() []TrafficPattern {
	return []TrafficPattern{
		{
			Name:        "CorporateBrowsing",
			MinRequests: 3,
			MaxRequests: 8,
			Timing: TimingPattern{
				MinDelay:    2 * time.Second,
				MaxDelay:    10 * time.Second,
				BurstChance: 0.2,
				BurstSize:   3,
			},
			URLs: []string{
				"https://portal.office.com",
				"https://outlook.office.com",
				"https://teams.microsoft.com",
				"https://admin.microsoft.com",
			},
			Methods: []string{"GET", "POST"},
			Headers: map[string][]string{
				"Referer": {
					"https://portal.office.com",
					"https://www.microsoft.com",
				},
			},
			Weight: 0.4,
		},
		{
			Name:        "EmailAccess",
			MinRequests: 2,
			MaxRequests: 5,
			Timing: TimingPattern{
				MinDelay:    1 * time.Second,
				MaxDelay:    5 * time.Second,
				BurstChance: 0.3,
				BurstSize:   2,
			},
			URLs: []string{
				"https://outlook.office.com/mail",
				"https://outlook.office.com/calendar",
				"https://outlook.office.com/people",
			},
			Methods: []string{"GET", "POST"},
			Headers: map[string][]string{
				"Accept": {
					"application/json",
					"text/html",
				},
			},
			Weight: 0.3,
		},
		{
			Name:        "DocumentAccess",
			MinRequests: 1,
			MaxRequests: 4,
			Timing: TimingPattern{
				MinDelay:    3 * time.Second,
				MaxDelay:    15 * time.Second,
				BurstChance: 0.1,
				BurstSize:   2,
			},
			URLs: []string{
				"https://onedrive.live.com",
				"https://sharepoint.com",
				"https://office.com/launch/word",
				"https://office.com/launch/excel",
			},
			Methods: []string{"GET"},
			Headers: map[string][]string{
				"Accept": {"text/html,application/xhtml+xml"},
			},
			Weight: 0.2,
		},
		{
			Name:        "APIRequests",
			MinRequests: 1,
			MaxRequests: 3,
			Timing: TimingPattern{
				MinDelay:    500 * time.Millisecond,
				MaxDelay:    2 * time.Second,
				BurstChance: 0.4,
				BurstSize:   4,
			},
			URLs: []string{
				"https://graph.microsoft.com/v1.0/me",
				"https://management.azure.com/subscriptions",
				"https://api.office.com/discovery/v2.0/me/services",
			},
			Methods: []string{"GET", "POST"},
			Headers: map[string][]string{
				"Accept":       {"application/json"},
				"Content-Type": {"application/json"},
			},
			Weight: 0.1,
		},
	}
}
