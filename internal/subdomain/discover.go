package subdomain

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Discover(domain string) ([]string, error) {

	cmd := exec.Command("subfinder", "-d", domain, "-silent")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("subfinder failed: %v\n%s", err, out.String())
	}

	lines := strings.Split(out.String(), "\n")
	var results []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			results = append(results, line)
		}
	}
	return results, nil
}
