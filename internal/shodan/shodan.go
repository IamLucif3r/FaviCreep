package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const baseURL = "https://api.shodan.io/shodan/host/search"

type ShodanResponse struct {
	Matches []struct {
		IPStr     string   `json:"ip_str"`
		Port      int      `json:"port"`
		Hostnames []string `json:"hostnames"`
	} `json:"matches"`
	Total int `json:"total"`
}

func SearchByFaviconHash(hash uint32) (*ShodanResponse, error) {
	apiKey := os.Getenv("SHODAN_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SHODAN_API_KEY environment variable is not set")
	}

	query := fmt.Sprintf("http.favicon.hash:%d", hash)
	url := fmt.Sprintf("%s?key=%s&query=%s", baseURL, apiKey, query)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("shodan API returned status %d", resp.StatusCode)
	}

	var result ShodanResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
