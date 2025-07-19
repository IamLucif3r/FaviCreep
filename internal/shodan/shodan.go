package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const baseURL = "https://api.shodan.io/shodan/host/search"

type Match struct {
	IPStr     string   `json:"ip_str"`
	Port      int      `json:"port"`
	Hostnames []string `json:"hostnames"`
}

type ShodanResponse struct {
	Matches []Match `json:"matches"`
	Total   int     `json:"total"`
}

func SearchByFaviconHash(hash uint32) (*ShodanResponse, error) {
	apiKey := os.Getenv("SHODAN_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SHODAN_API_KEY environment variable not set")
	}

	query := fmt.Sprintf("http.favicon.hash:%d", hash)
	url := fmt.Sprintf("%s?key=%s&query=%s", baseURL, apiKey, query)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query Shodan: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("shodan API returned HTTP %d", resp.StatusCode)
	}

	var result ShodanResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Shodan response: %w", err)
	}

	return &result, nil
}
