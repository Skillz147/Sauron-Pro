//utils/geo.go

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GeoInfo struct {
	Query       string `json:"query"`
	Country     string `json:"country"`
	City        string `json:"city"`
	RegionName  string `json:"regionName"`
	ISP         string `json:"isp"`
	Org         string `json:"org"`
	Timezone    string `json:"timezone"`
	CountryCode string `json:"countryCode"`
}

// GeoLookup queries ip-api.com for IP location info
func GeoLookup(ip string) (*GeoInfo, error) {
	if ip == "127.0.0.1" || ip == "::1" {
		return &GeoInfo{Country: "Localhost"}, nil
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,city,regionName,isp,org,query,timezone,countryCode", ip)
	client := http.Client{Timeout: 4 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GeoInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
