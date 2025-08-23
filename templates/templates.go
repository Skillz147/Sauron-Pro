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

// GetAvailableTemplates returns all predefined sublink templates
func GetAvailableTemplates() []SublinkTemplate {
	return []SublinkTemplate{
		{
			ID:          1,
			Name:        "Microsoft Identity",
			Path:        "common/identity/v2.0",
			Endpoints:   []string{"connect", "token", "signout", "device", "profile", "keys"},
			Description: "Microsoft Identity platform flow",
		},
		{
			ID:          2,
			Name:        "Session Management",
			Path:        "session",
			Endpoints:   []string{"start", "process", "validate", "continue", "refresh"},
			Description: "Session management flow",
		},
		{
			ID:          3,
			Name:        "Enterprise Identity",
			Path:        "organizations/identity/v2.0",
			Endpoints:   []string{"connect", "token", "consent", "device"},
			Description: "Enterprise identity platform",
		},
		{
			ID:          4,
			Name:        "Consumer Identity",
			Path:        "consumers/identity/v2.0",
			Endpoints:   []string{"connect", "token", "signout", "profile"},
			Description: "Consumer account identity flow",
		},
		{
			ID:          5,
			Name:        "Enterprise Session",
			Path:        "session/organizations",
			Endpoints:   []string{"start", "validate", "continue", "verify"},
			Description: "Enterprise session management",
		},
		{
			ID:          6,
			Name:        "Consumer Session",
			Path:        "session/consumers",
			Endpoints:   []string{"start", "process", "validate", "refresh"},
			Description: "Consumer session management",
		},
	}
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
