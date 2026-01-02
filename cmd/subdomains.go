package cmd

import (
	"fmt"
	"log"

	"github.com/iamlucif3r/favicreep/internal/subdomain"
	"github.com/iamlucif3r/favicreep/internal/utils"
	"github.com/spf13/cobra"
)

var domain string

var subdomainsCmd = &cobra.Command{
	Use:   "subdomains",
	Short: "Enumerate subdomains for a given domain",
	Run: func(cmd *cobra.Command, args []string) {
		if domain == "" {
			log.Fatal("❌ Please provide a domain using --domain")
		}

		spin := utils.NewSpinner("🔍 Enumerating subdomains for: " + domain)
		spin.Start()

		subs, err := subdomain.Discover(domain)
		spin.Stop()

		if err != nil {
			log.Fatalf("❌ Failed to enumerate subdomains: %v", err)
		}

		fmt.Printf("✅ Found %d subdomains:\n", len(subs))
		for _, s := range subs {
			fmt.Println(" -", s)
		}
	},
}

func init() {
	subdomainsCmd.Flags().StringVar(&domain, "domain", "", "Target domain (e.g. example.com)")
	rootCmd.AddCommand(subdomainsCmd)
}
