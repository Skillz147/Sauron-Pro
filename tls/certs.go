package tls

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"o365/utils"
)

const (
	certDir    = "tls"
	certPath   = "tls/cert.pem"
	keyPath    = "tls/key.pem"
	DomainRoot = "microsoftlogin.com"
)

func LoadDevCert() (tls.Certificate, string) {
	ensureCert()
	if err := ensureHosts(); err != nil {
		utils.SystemLogger.Warn().Err(err).Msg("Failed to update /etc/hosts")
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		utils.SystemLogger.Fatal().Err(err).Msg("Failed to load dev TLS cert")
	}
	return cert, DomainRoot
}

func ensureCert() {
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		_ = os.Mkdir(certDir, 0700)
	}
	if fileExists(certPath) && fileExists(keyPath) {
		return
	}

	utils.SystemLogger.Info().Msg("🔐 Generating dev certs with mkcert…")

	args := []string{
		"-cert-file", certPath,
		"-key-file", keyPath,
		"127.0.0.1",
		fmt.Sprintf("*.%s", DomainRoot),
		"admin." + DomainRoot,
	}

	seen := map[string]bool{
		"127.0.0.1":                     true,
		fmt.Sprintf("*.%s", DomainRoot): true,
	}

	for _, sub := range BaseSubdomains {
		parts := strings.Split(sub, ".")
		for i := 1; i <= len(parts); i++ {
			prefix := strings.Join(parts[:i], ".")
			wild := fmt.Sprintf("*.%s.%s", prefix, DomainRoot)
			if !seen[wild] {
				args = append(args, wild)
				seen[wild] = true
			}
		}

		host := fmt.Sprintf("%s.%s", sub, DomainRoot)
		if !seen[host] {
			args = append(args, host)
			seen[host] = true
		}

		www := fmt.Sprintf("www.%s.%s", sub, DomainRoot)
		if !seen[www] {
			args = append(args, www)
			seen[www] = true
		}
	}

	utils.SystemLogger.Info().Int("count", len(args)-4).Msg("mkcert domains prepared")
	cmd := exec.Command("mkcert", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		utils.SystemLogger.Fatal().Err(err).Msg("mkcert execution failed")
	}
}

func ensureHosts() error {
	var lines []string
	seen := map[string]bool{}

	for _, sub := range BaseSubdomains {
		for _, prefix := range []string{"", "www."} {
			entry := fmt.Sprintf("%s%s.%s", prefix, sub, DomainRoot)
			if !seen[entry] {
				lines = append(lines, fmt.Sprintf("127.0.0.1 %s", entry))
				seen[entry] = true
			}
		}
	}
	lines = append(lines, fmt.Sprintf("127.0.0.1 admin.%s", DomainRoot))

	content, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	var missing []string
	for _, line := range lines {
		if !strings.Contains(string(content), line) {
			missing = append(missing, line)
		}
	}
	if len(missing) == 0 {
		return nil
	}

	utils.SystemLogger.Info().Int("count", len(missing)).Msg("Appending entries to /etc/hosts")
	cmd := exec.Command("sudo", "tee", "-a", "/etc/hosts")
	cmd.Stdin = strings.NewReader("\n" + strings.Join(missing, "\n") + "\n")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func fileExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return err == nil
}
