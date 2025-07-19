package favicon

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spaolacci/murmur3"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 {
			return fmt.Errorf("too many redirects")
		}
		return nil
	},
}

func HashFavicon(url string) (uint32, error) {
	favURL := normalizeFaviconURL(url)

	resp, err := client.Get(favURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch favicon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read favicon: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(body)

	hasher := murmur3.New32()
	_, err = hasher.Write([]byte(encoded))
	if err != nil {
		return 0, fmt.Errorf("failed to write to hash: %w", err)
	}

	return hasher.Sum32(), nil
}

func normalizeFaviconURL(url string) string {
	url = strings.TrimRight(url, "/")
	if !strings.HasSuffix(url, ".ico") {
		return url + "/favicon.ico"
	}
	return url
}
