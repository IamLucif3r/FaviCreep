package favicon

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"net/url"

	"github.com/spaolacci/murmur3"
	"golang.org/x/net/html"
)

func HashFavicon(baseURL string) (uint32, error) {
	baseURL = strings.TrimSuffix(baseURL, "/")

	favURL, err := fetchFaviconURL(baseURL)
	if err != nil {
		return 0, fmt.Errorf("failed to find favicon: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(favURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch favicon at %s: %w", favURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("favicon fetch returned status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read favicon body: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	return murmur3.Sum32([]byte(encoded)), nil
}

func fetchFaviconURL(baseURL string) (string, error) {

	favURL := baseURL + "/favicon.ico"
	ok, err := urlExists(favURL)
	if err == nil && ok {
		return favURL, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var faviconHref string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			rel := ""
			href := ""
			for _, attr := range n.Attr {
				if attr.Key == "rel" {
					rel = attr.Val
				}
				if attr.Key == "href" {
					href = attr.Val
				}
			}
			if strings.Contains(strings.ToLower(rel), "icon") && href != "" && faviconHref == "" {
				faviconHref = href
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if faviconHref == "" {
		return "", fmt.Errorf("no favicon found")
	}

	u, err := url.Parse(faviconHref)
	if err != nil {
		return "", err
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	resolved := base.ResolveReference(u).String()
	return resolved, nil
}

func urlExists(u string) (bool, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(u)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}
