package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "favicreep",
	Short: "FaviCreep is used to discover forgotten infrastructure using favicon hashing and Shodan",
	Long: `FaviCreep is a recon tool that finds subdomains, hashes favicons (mmh3), 
clusters them, and searches Shodan for matching hashes to uncover exposed assets.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}
}
