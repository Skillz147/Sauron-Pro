package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

// LicenseProof represents the payload structure for the license token
type LicenseProof struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"` // Unix timestamp
}

// VerifyLicenseProofToken validates an HMAC-authenticated license token
func VerifyLicenseProofToken(token string) (*LicenseProof, error) {
	// Expected format: base64(payload).base64(hmac)
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}
	payloadB64 := parts[0]
	sigB64 := parts[1]

	// Decode base64 payload
	payloadJSON, err := base64.StdEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, errors.New("invalid base64 in payload")
	}

	// Verify HMAC signature
	expectedSig := computeHMAC(payloadB64)
	if !hmac.Equal([]byte(expectedSig), []byte(sigB64)) {
		return nil, errors.New("invalid HMAC signature")
	}

	// Parse payload into struct
	var proof LicenseProof
	if err := json.Unmarshal(payloadJSON, &proof); err != nil {
		return nil, errors.New("invalid JSON in payload")
	}
	if proof.UserID == "" || proof.Exp <= 0 {
		return nil, errors.New("missing user_id or expiry")
	}
	if time.Now().Unix() > proof.Exp {
		return nil, errors.New("token expired")
	}

	return &proof, nil
}

// computeHMAC returns the base64-encoded HMAC-SHA256 signature for input
func computeHMAC(input string) string {
	secret := os.Getenv("LICENSE_TOKEN_SECRET")
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
