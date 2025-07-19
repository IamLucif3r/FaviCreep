package subdomain

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Discover(domain string) ([]string, error) {
	var out bytes.Buffer
	cmd := exec.Command("subfinder", "-silent", "-d", domain)

	cmd.Stdout = &out
	cmd.Stderr = nil

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run subfinder: %w", err)
	}

	var subdomains []string
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		sub := strings.TrimSpace(scanner.Text())
		if sub != "" {
			subdomains = append(subdomains, sub)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading subfinder output: %w", err)
	}

	return subdomains, nil
}
