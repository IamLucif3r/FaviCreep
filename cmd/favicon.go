package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/IamLucif3r/favicreep/internal/cluster"
	"github.com/IamLucif3r/favicreep/internal/favicon"
	"github.com/IamLucif3r/favicreep/internal/subdomain"
	"github.com/IamLucif3r/favicreep/internal/utils"
	"github.com/spf13/cobra"
)

var (
	scanDomain  string
	concurrency int
	outputFile  string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Find subdomains and hash their favicons with concurrency and clustering",
	Run: func(cmd *cobra.Command, args []string) {
		if scanDomain == "" {
			log.Fatal("❌ Please provide a domain using --domain")
		}

		spin := utils.NewSpinner("Enumerating subdomains for " + scanDomain)
		spin.Start()
		subs, err := subdomain.Discover(scanDomain)
		spin.Stop()
		if err != nil {
			log.Fatalf("❌ Error discovering subdomains: %v", err)
		}

		if len(subs) == 0 {
			fmt.Println("⚠️  No subdomains found.")
			return
		}

		fmt.Printf("🔍 Found %d subdomains\n", len(subs))
		fmt.Println("🔑 Hashing favicons concurrently...")

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
					fmt.Printf("❌ [%s] Error: %v\n", url, err)
					return
				}

				mu.Lock()
				results[s] = hash
				mu.Unlock()

				fmt.Printf("✅ [%s] mmh3: %d\n", url, hash)
			}(sub)
		}

		wg.Wait()

		clusters := cluster.Cluster(results)

		fmt.Println("\n📊 Clustering results:")
		for hash, domains := range clusters {
			fmt.Printf("Hash: %d\n", hash)
			for _, d := range domains {
				fmt.Printf(" - %s\n", d)
			}
		}

		if outputFile != "" {
			data, err := json.MarshalIndent(clusters, "", "  ")
			if err != nil {
				log.Fatalf("❌ Failed to marshal clusters to JSON: %v", err)
			}
			err = os.WriteFile(outputFile, data, 0644)
			if err != nil {
				log.Fatalf("❌ Failed to write output file: %v", err)
			}
			fmt.Printf("✅ Results saved to %s\n", outputFile)
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
	scanCmd.Flags().StringVar(&scanDomain, "domain", "", "Domain to scan (e.g. example.com)")
	scanCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent favicon fetches")
	scanCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write clustering results to JSON file")
	rootCmd.AddCommand(scanCmd)
}
