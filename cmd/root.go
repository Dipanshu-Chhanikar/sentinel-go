package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "Sentinel-Go: AI-Driven Bug Bounty Agent",
	Long: `Sentinel-Go is an autonomous, context-aware bug bounty agent.
It leverages local LLMs to analyze JavaScript, fuzz endpoints, and orchestrate security tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[+] Sentinel-Go initialized.")
		fmt.Println("[+] Waiting for target instructions...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
