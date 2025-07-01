package tls

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"o365/configdb"
	"o365/utils"

	"github.com/caddyserver/certmagic"
	"github.com/libdns/cloudflare"
	"github.com/libdns/libdns"
)

// LoggingDNSProvider logs TXT records
type LoggingDNSProvider struct {
	Provider *cloudflare.Provider
}

func (p *LoggingDNSProvider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	for _, rec := range recs {
		utils.SystemLogger.Info().
			Str("zone", zone).
			Str("name", rec.Name).
			Str("value", rec.Value).
			Msg("Expected TXT record")
	}
	records, err := p.Provider.AppendRecords(ctx, zone, recs)
	if err != nil {
		utils.SystemLogger.Error().Err(err).Str("zone", zone).Msg("Failed to append TXT")
	}
	return records, err
}

func (p *LoggingDNSProvider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	for _, rec := range recs {
		utils.SystemLogger.Info().
			Str("zone", zone).
			Str("name", rec.Name).
			Str("value", rec.Value).
			Msg("Deleting TXT record")
	}
	records, err := p.Provider.DeleteRecords(ctx, zone, recs)
	if err != nil {
		utils.SystemLogger.Error().Err(err).Str("zone", zone).Msg("Failed to delete TXT")
	}
	return records, err
}

func (p *LoggingDNSProvider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	records, err := p.Provider.GetRecords(ctx, zone)
	if err != nil {
		utils.SystemLogger.Error().Err(err).Str("zone", zone).Msg("Failed to get DNS records")
	}
	return records, err
}

func LoadProdCert(ctx context.Context) (tls.Certificate, string) {
	domain := configdb.GetDomain()
	if domain == "" {
		utils.SystemLogger.Fatal().Msg("❌ No domain configured")
	}
	cfToken := configdb.GetCloudflareToken()
	if cfToken == "" {
		utils.SystemLogger.Fatal().Msg("❌ Cloudflare API token missing")
	}

	if err := os.Setenv("CLOUDFLARE_API_TOKEN", cfToken); err != nil {
		utils.SystemLogger.Fatal().Err(err).Msg("❌ Failed to set CLOUDFLARE_API_TOKEN")
	}

	cache := certmagic.NewCache(certmagic.CacheOptions{
		GetConfigForCert: func(cert certmagic.Certificate) (*certmagic.Config, error) {
			return certmagic.NewDefault(), nil
		},
	})

	cloudflareProvider := &cloudflare.Provider{
		APIToken: cfToken,
	}
	loggingProvider := &LoggingDNSProvider{Provider: cloudflareProvider}

	cfg := certmagic.NewDefault()
	issuer := certmagic.NewACMEIssuer(cfg, certmagic.ACMEIssuer{
		Email:  "admin@" + domain,
		CA:     certmagic.LetsEncryptProductionCA,
		Agreed: true,
		DNS01Solver: &certmagic.DNS01Solver{
			DNSManager: certmagic.DNSManager{
				DNSProvider: loggingProvider,
			},
		},
	})

	cm := certmagic.New(cache, certmagic.Config{
		Issuers: []certmagic.Issuer{issuer},
	})

	domains := []string{"*." + domain, domain, "admin." + domain}
	seen := map[string]bool{
		"*." + domain:     true,
		domain:            true,
		"admin." + domain: true,
	}

	for _, sub := range BaseSubdomains {
		parts := strings.Split(sub, ".")
		for i := 1; i <= len(parts); i++ {
			prefix := strings.Join(parts[:i], ".")
			wild := fmt.Sprintf("*.%s.%s", prefix, domain)
			if !seen[wild] {
				domains = append(domains, wild)
				seen[wild] = true
			}
		}
		host := fmt.Sprintf("%s.%s", sub, domain)
		if !seen[host] {
			domains = append(domains, host)
			seen[host] = true
		}
		www := fmt.Sprintf("www.%s.%s", sub, domain)
		if !seen[www] {
			domains = append(domains, www)
			seen[www] = true
		}
	}

	utils.SystemLogger.Info().Strs("domains", domains).Msg("🔐 Requesting certificate")

	if err := cm.ManageSync(ctx, domains); err != nil {
		utils.SystemLogger.Fatal().Err(err).Msg("❌ Failed to obtain certificate")
	}

	tlsCfg := cm.TLSConfig()
	if len(tlsCfg.Certificates) == 0 {
		utils.SystemLogger.Fatal().Msg("❌ No certificates returned")
	}
	return tlsCfg.Certificates[0], domain
}
