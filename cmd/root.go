package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "favicreep",
	Short: "favicreep - Discover shadow assets using favicon hashes",
	Long: `favicreep is a fast, modular tool that finds forgotten systems
by fingerprinting favicons across domains.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}
}
