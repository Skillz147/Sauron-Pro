package decoy

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"o365/utils"
)

// AdvancedDecoyPatterns provides sophisticated traffic patterns that mimic real M365 usage
type AdvancedDecoyPatterns struct {
	manager *DecoyTrafficManager
}

// NewAdvancedDecoyPatterns creates an advanced pattern generator
func NewAdvancedDecoyPatterns(manager *DecoyTrafficManager) *AdvancedDecoyPatterns {
	return &AdvancedDecoyPatterns{manager: manager}
}

// StartAdvancedPatterns begins sophisticated decoy traffic generation
func (a *AdvancedDecoyPatterns) StartAdvancedPatterns() {
	// Start different types of realistic patterns
	go a.simulateLoginAttempts()
	go a.simulateEmailActivity()
	go a.simulateTeamsActivity()
	go a.simulateAdminActivity()
	go a.simulateDocumentAccess()
	go a.simulateAPITraffic()

	utils.SystemLogger.Info().Msg("🎭 Advanced decoy patterns started")
}

// simulateLoginAttempts creates realistic login flows (including failed attempts)
func (a *AdvancedDecoyPatterns) simulateLoginAttempts() {
	ticker := time.NewTicker(time.Duration(float64(120*time.Second) / a.manager.intensity))
	defer ticker.Stop()

	loginDomains := []string{
		"login.microsoftonline.com",
		"login.microsoft.com",
		"account.microsoft.com",
	}

	for {
		select {
		case <-a.manager.stopChan:
			return
		case <-ticker.C:
			a.simulateLoginFlow(loginDomains[rand.Intn(len(loginDomains))])
		}
	}
}

// simulateLoginFlow creates a realistic login sequence
func (a *AdvancedDecoyPatterns) simulateLoginFlow(domain string) {
	baseURL := "https://" + domain

	// Realistic login flow sequence
	sequence := []struct {
		path   string
		method string
		delay  time.Duration
	}{
		{"/common/oauth2/authorize", "GET", 1 * time.Second},
		{"/common/login", "GET", 2 * time.Second},
		{"/common/SAS/ProcessAuth", "POST", 3 * time.Second},
		{"/common/SAS/BeginAuth", "POST", 1 * time.Second},
		{"/common/SAS/EndAuth", "POST", 2 * time.Second},
	}

	for _, step := range sequence {
		time.Sleep(step.delay)

		headers := map[string]string{
			"Referer": baseURL,
			"Origin":  baseURL,
			"Accept":  "application/json, text/plain, */*",
		}

		// Add realistic form data for POST requests
		var body []byte
		if step.method == "POST" {
			body = []byte(a.generateLoginFormData())
			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}

		a.manager.makeDecoyRequestWithHeaders(step.method, baseURL+step.path, headers, body)

		select {
		case <-a.manager.stopChan:
			return
		default:
		}
	}
}

// generateLoginFormData creates realistic login form data
func (a *AdvancedDecoyPatterns) generateLoginFormData() string {
	// Use realistic but fake credentials
	emails := []string{
		"john.smith@contoso.com",
		"mary.johnson@fabrikam.com",
		"david.wilson@northwind.com",
		"sarah.brown@adventure-works.com",
	}

	email := emails[rand.Intn(len(emails))]

	return fmt.Sprintf("login=%s&passwd=TempPassword123&flowToken=%s&canary=%s",
		email,
		generateRandomToken(32),
		generateRandomToken(16))
}

// simulateEmailActivity creates realistic Outlook usage patterns
func (a *AdvancedDecoyPatterns) simulateEmailActivity() {
	ticker := time.NewTicker(time.Duration(float64(90*time.Second) / a.manager.intensity))
	defer ticker.Stop()

	outlookActions := []string{
		"/mail/inbox",
		"/mail/sent",
		"/mail/drafts",
		"/mail/search",
		"/mail/compose",
		"/calendar/view",
		"/calendar/events",
		"/people/contacts",
	}

	for {
		select {
		case <-a.manager.stopChan:
			return
		case <-ticker.C:
			action := outlookActions[rand.Intn(len(outlookActions))]
			a.manager.makeDecoyRequest("GET", "https://outlook.office.com"+action, nil)
		}
	}
}

// simulateTeamsActivity creates realistic Teams usage
func (a *AdvancedDecoyPatterns) simulateTeamsActivity() {
	ticker := time.NewTicker(time.Duration(float64(150*time.Second) / a.manager.intensity))
	defer ticker.Stop()

	teamsActions := []string{
		"/conversations",
		"/teams",
		"/calendar",
		"/calls",
		"/files",
		"/apps",
	}

	for {
		select {
		case <-a.manager.stopChan:
			return
		case <-ticker.C:
			action := teamsActions[rand.Intn(len(teamsActions))]
			a.manager.makeDecoyRequest("GET", "https://teams.microsoft.com"+action, nil)
		}
	}
}

// simulateAdminActivity creates admin portal usage patterns
func (a *AdvancedDecoyPatterns) simulateAdminActivity() {
	ticker := time.NewTicker(time.Duration(float64(300*time.Second) / a.manager.intensity))
	defer ticker.Stop()

	adminActions := []string{
		"/AdminPortal/Home",
		"/AdminPortal/Users/List",
		"/AdminPortal/Groups/List",
		"/AdminPortal/Reports/Usage",
		"/AdminPortal/Settings/Security",
		"/AdminPortal/Billing/Subscriptions",
	}

	for {
		select {
		case <-a.manager.stopChan:
			return
		case <-ticker.C:
			action := adminActions[rand.Intn(len(adminActions))]
			headers := map[string]string{
				"X-Requested-With": "XMLHttpRequest",
				"Accept":           "application/json",
			}
			a.manager.makeDecoyRequestWithHeaders("GET", "https://admin.microsoft.com"+action, headers, nil)
		}
	}
}

// simulateDocumentAccess creates SharePoint/OneDrive usage
func (a *AdvancedDecoyPatterns) simulateDocumentAccess() {
	ticker := time.NewTicker(time.Duration(float64(180*time.Second) / a.manager.intensity))
	defer ticker.Stop()

	documentSites := []string{
		"https://contoso.sharepoint.com",
		"https://fabrikam.sharepoint.com",
		"https://onedrive.live.com",
	}

	documentActions := []string{
		"/Documents/Forms/AllItems.aspx",
		"/Shared%20Documents",
		"/SitePages/Home.aspx",
		"/Lists/Tasks/AllItems.aspx",
		"/_layouts/15/storman.aspx",
	}

	for {
		select {
		case <-a.manager.stopChan:
			return
		case <-ticker.C:
			site := documentSites[rand.Intn(len(documentSites))]
			action := documentActions[rand.Intn(len(documentActions))]
			a.manager.makeDecoyRequest("GET", site+action, nil)
		}
	}
}

// simulateAPITraffic creates realistic Graph API and REST calls
func (a *AdvancedDecoyPatterns) simulateAPITraffic() {
	ticker := time.NewTicker(time.Duration(float64(60*time.Second) / a.manager.intensity))
	defer ticker.Stop()

	apiEndpoints := []string{
		"https://graph.microsoft.com/v1.0/me",
		"https://graph.microsoft.com/v1.0/me/messages",
		"https://graph.microsoft.com/v1.0/me/events",
		"https://graph.microsoft.com/v1.0/me/drive/root/children",
		"https://graph.microsoft.com/v1.0/users",
		"https://graph.microsoft.com/v1.0/groups",
		"https://management.azure.com/subscriptions",
		"https://api.office.com/discovery/v2.0/me/services",
	}

	for {
		select {
		case <-a.manager.stopChan:
			return
		case <-ticker.C:
			endpoint := apiEndpoints[rand.Intn(len(apiEndpoints))]
			headers := map[string]string{
				"Authorization": "Bearer " + generateRandomToken(64),
				"Accept":        "application/json",
				"Content-Type":  "application/json",
			}
			a.manager.makeDecoyRequestWithHeaders("GET", endpoint, headers, nil)
		}
	}
}

// generateRandomToken creates a realistic but fake token
func generateRandomToken(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// InjectDecoyIntoRealTraffic adds decoy requests around real traffic
func (a *AdvancedDecoyPatterns) InjectDecoyIntoRealTraffic(realRequest *http.Request) {
	if !a.manager.enabled {
		return
	}

	// Generate 1-3 decoy requests around the real one
	numDecoys := rand.Intn(3) + 1

	for i := 0; i < numDecoys; i++ {
		// Random delay before/after real request
		delay := time.Duration(rand.Intn(5)+1) * time.Second

		go func() {
			time.Sleep(delay)
			a.generateContextualDecoy(realRequest)
		}()
	}
}

// generateContextualDecoy creates decoy traffic similar to real request
func (a *AdvancedDecoyPatterns) generateContextualDecoy(realReq *http.Request) {
	host := realReq.Host

	// Generate similar but different URLs
	var decoyURL string

	if strings.Contains(host, "microsoft") || strings.Contains(host, "office") {
		// Microsoft-related decoys
		microsoftPaths := []string{
			"/login",
			"/auth",
			"/api/v1.0/me",
			"/mail/inbox",
			"/calendar",
			"/files",
		}
		path := microsoftPaths[rand.Intn(len(microsoftPaths))]
		decoyURL = "https://" + host + path
	} else {
		// Generic decoys
		genericPaths := []string{
			"/home",
			"/dashboard",
			"/profile",
			"/settings",
			"/api/status",
		}
		path := genericPaths[rand.Intn(len(genericPaths))]
		decoyURL = "https://" + host + path
	}

	// Use similar headers to real request
	headers := map[string]string{
		"User-Agent": realReq.UserAgent(),
		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Referer":    "https://" + host,
	}

	a.manager.makeDecoyRequestWithHeaders("GET", decoyURL, headers, nil)
}
