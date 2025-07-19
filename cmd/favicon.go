package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/IamLucif3r/favicreep/internal/cluster"
	"github.com/IamLucif3r/favicreep/internal/favicon"
	"github.com/IamLucif3r/favicreep/internal/shodan"
	"github.com/IamLucif3r/favicreep/internal/subdomain"
	"github.com/IamLucif3r/favicreep/internal/utils"
	"github.com/spf13/cobra"
)

var (
	concurrency int
	outputFile  string
	withShodan  bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform subdomain discovery, hash favicons, cluster and search via Shodan",
	Run: func(cmd *cobra.Command, args []string) {
		if domain == "" {
			log.Fatal("[ERROR] Please provide a domain using --domain")
		}

		spin := utils.NewSpinner("🔍 Enumerating subdomains for " + domain)
		spin.Start()
		subs, err := subdomain.Discover(domain)
		spin.Stop()

		if err != nil {
			log.Fatalf("[ERROR] Failed to discover subdomains: %v", err)
		}

		if len(subs) == 0 {
			log.Println("[INFO]  No subdomains found.")
			return
		}

		fmt.Printf("[RESULT] Found %d subdomains\n", len(subs))

		spin = utils.NewSpinner("🎨 Fetching and hashing favicons...")
		spin.Start()

		results := make(map[string]uint32)
		var mu sync.Mutex
		var wg sync.WaitGroup
		sem := make(chan struct{}, concurrency)

		for _, sub := range subs {
			wg.Add(1)
			sem <- struct{}{}

			go func(s string) {
				defer wg.Done()
				defer func() { <-sem }()

				url := ensureURL(s)
				hash, err := favicon.HashFavicon(url)
				if err != nil {
					log.Printf("[ERROR] [%s] Favicon error: %v", url, err)
					return
				}

				mu.Lock()
				results[s] = hash
				mu.Unlock()
			}(sub)
		}

		wg.Wait()
		spin.Stop()

		clusters := cluster.Cluster(results)
		fmt.Printf("[RESULT] Found %d favicon hash clusters\n", len(clusters))

		shodanResults := make(map[uint32][]string)
		if withShodan {
			fmt.Println("[INFO] Performing Shodan lookups...")
			for hash := range clusters {
				res, err := shodan.SearchByFaviconHash(hash)
				if err != nil {
					log.Printf("[ERROR] Shodan search error for hash %d: %v", hash, err)
					continue
				}
				var hosts []string
				for _, match := range res.Matches {
					hosts = append(hosts, fmt.Sprintf("%s:%d", match.IPStr, match.Port))
				}
				shodanResults[hash] = hosts
			}
		}

		if outputFile != "" {
			data := map[string]any{
				"domain":       domain,
				"discovered":   subs,
				"clusters":     clusters,
				"shodan":       shodanResults,
				"generated_at": time.Now().Format(time.RFC3339),
			}

			bytes, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				log.Fatalf("[ERROR] Failed to marshal output: %v", err)
			}
			if err := os.WriteFile(outputFile, bytes, 0644); err != nil {
				log.Fatalf("[ERROR] Failed to write output file: %v", err)
			}
			fmt.Printf("[INFO] Results saved to %s\n", outputFile)
		}
	},
}

func ensureURL(sub string) string {
	if !strings.HasPrefix(sub, "http://") && !strings.HasPrefix(sub, "https://") {
		return "https://" + sub
	}
	return sub
}

func init() {
	scanCmd.Flags().StringVar(&domain, "domain", "", "Target domain (required)")
	scanCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent favicon fetches")
	scanCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Save results to a JSON file")
	scanCmd.Flags().BoolVar(&withShodan, "shodan", false, "Enable Shodan lookup for each favicon hash")
	rootCmd.AddCommand(scanCmd)
}
