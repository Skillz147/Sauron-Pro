package utils

import (
	"net"
	"net/http"
	"strings"

	"o365/configdb"

	"github.com/oschwald/geoip2-golang"
)

var (
	geoDB *geoip2.Reader

	// ISO 3166-1 alpha-2 codes for allowed countries
	allowedCountries = map[string]bool{
		"US": true, "CA": true, "GB": true, "DE": true, "FR": true,
		"IT": true, "NL": true, "SE": true, "NO": true, "DK": true,
		"FI": true, "IE": true, "JP": true, "SG": true, "AU": true,
		"CH": true, "AT": true, "BE": true, "LU": true, "IS": true,
		"ES": true, "PT": true, "SM": true, "AD": true, "MT": true,
		"SI": true, "CY": true, "EE": true, "CZ": true,
	}
)

// InitGeoDB loads a MaxMind GeoLite2 database from the given path
func InitGeoDB(path string) error {
	var err error
	geoDB, err = geoip2.Open(path)
	return err
}

// GetRealIP attempts to determine the client's real IP address
func GetRealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// GetCountryCode resolves a country ISO code from an IP address
func GetCountryCode(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed == nil || geoDB == nil {
		return ""
	}
	record, err := geoDB.Country(parsed)
	if err != nil || record == nil {
		return ""
	}
	return record.Country.IsoCode
}

// CountryFilter blocks requests from disallowed countries or banned IPs
func CountryFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetRealIP(r)

		// 🧪 Allow local development
		if ip == "127.0.0.1" || ip == "::1" {
			next.ServeHTTP(w, r)
			return
		}

		// 🚫 IP already banned
		if configdb.IsIPBanned(ip) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		// 🌍 Geo filter
		country := GetCountryCode(ip)
		if !allowedCountries[country] {
			BotLogger.Warn().
				Str("ip", ip).
				Str("country", country).
				Str("ua", r.UserAgent()).
				Msg("Blocked disallowed country")

			configdb.SaveBannedIP(ip, country)
			http.Error(w, "Blocked by region", http.StatusForbidden)
			return
		}

		// ✅ Allow
		next.ServeHTTP(w, r)
	})
}
