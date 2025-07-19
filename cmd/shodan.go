package cmd

import (
	"fmt"
	"log"

	"github.com/IamLucif3r/favicreep/internal/shodan"
	"github.com/spf13/cobra"
)

var shodanHash uint32

var shodanCmd = &cobra.Command{
	Use:   "shodan",
	Short: "Search Shodan for hosts matching a favicon hash",
	Run: func(cmd *cobra.Command, args []string) {
		if shodanHash == 0 {
			log.Fatal("[ERROR] Please provide a hash using --hash")
		}

		fmt.Printf("🔎 Searching Shodan for favicon hash %d...\n", shodanHash)

		resp, err := shodan.SearchByFaviconHash(shodanHash)
		if err != nil {
			log.Fatalf("[ERROR] Shodan search failed: %v", err)
		}

		fmt.Printf("Found %d matching hosts\n", resp.Total)
		for _, match := range resp.Matches {
			fmt.Printf(" - %s:%d (Hostnames: %v)\n", match.IPStr, match.Port, match.Hostnames)
		}
	},
}

func init() {
	shodanCmd.Flags().Uint32Var(&shodanHash, "hash", 0, "Favicon mmh3 hash to search")
	rootCmd.AddCommand(shodanCmd)
}
