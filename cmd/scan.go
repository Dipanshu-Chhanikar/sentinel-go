package cmd

import (
	"fmt"
	"log"

	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/attack"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [url] [tech-tag]",
	Short: "Run a targeted Nuclei scan against a specific URL",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		techTag := args[1]
		
		fmt.Printf("[*] Initiating targeted Nuclei strike on: %s\n", target)
		fmt.Printf("[*] Technology Profile: %s\n", techTag)
		fmt.Println("[*] Executing subprocess... (this may take a minute)")
		
		findings, err := attack.RunNuclei(target, techTag)
		if err != nil {
			log.Fatalf("[-] Subprocess failure: %v", err)
		}
		
		if len(findings) == 0 {
			fmt.Println("[+] Scan complete. No vulnerabilities found for this tech stack.")
			return
		}

		fmt.Printf("\n[!] CRITICAL: Found %d potential vulnerabilities!\n", len(findings))
		for _, finding := range findings {
			fmt.Printf("    -> %s\n", finding)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
