package honeypot

import (
	"sync"
	"time"
)

// IPAttempt tracks attempts per IP address
type IPAttempt struct {
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
	Banned    bool
}

// RateLimiter tracks IP attempts with sliding window
type RateLimiter struct {
	attempts map[string]*IPAttempt
	mutex    sync.RWMutex
	window   time.Duration
	maxTries int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(window time.Duration, maxTries int) *RateLimiter {
	rl := &RateLimiter{
		attempts: make(map[string]*IPAttempt),
		window:   window,
		maxTries: maxTries,
	}

	// Cleanup goroutine
	go rl.cleanup()

	return rl
}

// CheckAndRecord checks if IP should be banned and records attempt
func (rl *RateLimiter) CheckAndRecord(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	attempt, exists := rl.attempts[ip]
	if !exists {
		rl.attempts[ip] = &IPAttempt{
			Count:     1,
			FirstSeen: now,
			LastSeen:  now,
			Banned:    false,
		}
		return false
	}

	// If already banned, return true
	if attempt.Banned {
		return true
	}

	// Check if outside window - reset if so
	if now.Sub(attempt.FirstSeen) > rl.window {
		attempt.Count = 1
		attempt.FirstSeen = now
		attempt.LastSeen = now
		return false
	}

	// Increment count
	attempt.Count++
	attempt.LastSeen = now

	// Check if should be banned
	if attempt.Count > rl.maxTries {
		attempt.Banned = true
		return true
	}

	return false
}

// IsBanned checks if IP is currently banned
func (rl *RateLimiter) IsBanned(ip string) bool {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	attempt, exists := rl.attempts[ip]
	return exists && attempt.Banned
}

// GetAttemptCount returns current attempt count for IP
func (rl *RateLimiter) GetAttemptCount(ip string) int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	attempt, exists := rl.attempts[ip]
	if !exists {
		return 0
	}

	return attempt.Count
}

// cleanup removes old entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()

		for ip, attempt := range rl.attempts {
			// Remove entries older than 24 hours
			if now.Sub(attempt.LastSeen) > 24*time.Hour {
				delete(rl.attempts, ip)
			}
		}

		rl.mutex.Unlock()
	}
}
