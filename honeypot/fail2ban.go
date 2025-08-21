package honeypot

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Fail2BanManager handles integration with fail2ban
type Fail2BanManager struct {
	logFile  string
	jailName string
	enabled  bool
}

// NewFail2BanManager creates a new fail2ban manager
func NewFail2BanManager() *Fail2BanManager {
	manager := &Fail2BanManager{
		logFile:  "logs/sauron-security.log",
		jailName: "sauron-honeypot",
		enabled:  true,
	}

	// Check if fail2ban is available
	if _, err := exec.LookPath("fail2ban-client"); err != nil {
		log.Warn().Msg("fail2ban-client not found - disabling fail2ban integration")
		manager.enabled = false
	}

	return manager
}

// TriggerBan logs an event that fail2ban can parse and optionally bans directly
func (f *Fail2BanManager) TriggerBan(clientIP, attackType, reason string) {
	// Log in fail2ban parseable format
	logEntry := fmt.Sprintf("[%s] SECURITY_VIOLATION: IP=%s TYPE=%s REASON=%s",
		time.Now().Format("2006-01-02 15:04:05"),
		clientIP,
		attackType,
		reason)

	// Write to security log file
	f.writeToSecurityLog(logEntry)

	// Also log via zerolog for our own monitoring
	log.Error().
		Str("attack_type", attackType).
		Str("client_ip", clientIP).
		Str("reason", reason).
		Msg("FAIL2BAN_TRIGGER: Security violation logged")

	// If fail2ban is enabled, try direct ban
	if f.enabled {
		f.directBan(clientIP, reason)
	}
}

// writeToSecurityLog writes to the security log file that fail2ban monitors
func (f *Fail2BanManager) writeToSecurityLog(entry string) {
	file, err := os.OpenFile(f.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Str("file", f.logFile).Msg("Failed to open security log file")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(entry + "\n")
	if err != nil {
		log.Error().Err(err).Msg("Failed to write to security log")
		return
	}
	writer.Flush()
}

// directBan attempts to ban IP directly via fail2ban-client
func (f *Fail2BanManager) directBan(clientIP, reason string) {
	if !f.enabled {
		return
	}

	// Try to ban IP directly
	cmd := exec.Command("fail2ban-client", "set", f.jailName, "banip", clientIP)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Warn().
			Err(err).
			Str("ip", clientIP).
			Str("output", string(output)).
			Msg("Failed to ban IP via fail2ban-client")
		return
	}

	log.Info().
		Str("ip", clientIP).
		Str("jail", f.jailName).
		Str("reason", reason).
		Msg("IP banned via fail2ban")
}

// UnbanIP removes an IP from fail2ban jail
func (f *Fail2BanManager) UnbanIP(clientIP string) error {
	if !f.enabled {
		return fmt.Errorf("fail2ban not enabled")
	}

	cmd := exec.Command("fail2ban-client", "set", f.jailName, "unbanip", clientIP)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to unban IP %s: %v - output: %s", clientIP, err, string(output))
	}

	log.Info().
		Str("ip", clientIP).
		Str("jail", f.jailName).
		Msg("IP unbanned via fail2ban")

	return nil
}

// GetBannedIPs returns list of currently banned IPs
func (f *Fail2BanManager) GetBannedIPs() ([]string, error) {
	if !f.enabled {
		return nil, fmt.Errorf("fail2ban not enabled")
	}

	cmd := exec.Command("fail2ban-client", "status", f.jailName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("failed to get banned IPs: %v - output: %s", err, string(output))
	}

	// Parse output to extract banned IPs
	lines := strings.Split(string(output), "\n")
	var bannedIPs []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Currently banned:") {
			// Extract IPs from line like "Currently banned: 2     192.168.1.100 10.0.0.1"
			parts := strings.Fields(line)
			if len(parts) > 2 {
				bannedIPs = parts[2:] // Skip "Currently" and "banned:" and count
			}
			break
		}
	}

	return bannedIPs, nil
}

// IsEnabled returns whether fail2ban integration is enabled
func (f *Fail2BanManager) IsEnabled() bool {
	return f.enabled
}

// GetLogFile returns the path to the security log file
func (f *Fail2BanManager) GetLogFile() string {
	return f.logFile
}
