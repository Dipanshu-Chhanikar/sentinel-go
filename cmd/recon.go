package cmd

import (
	"fmt"
	"log"

	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/recon"
	"github.com/spf13/cobra"
)

var reconCmd = &cobra.Command{
	Use:   "recon [domain]",
	Short: "Passively gather historical URLs for a domain",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		
		fmt.Printf("[*] Initiating passive recon for: %s\n", domain)
		fmt.Println("[*] Querying Wayback Machine CDX API...")
		
		urls, err := recon.FetchWaybackURLs(domain)
		if err != nil {
			log.Fatalf("[-] Recon failed: %v", err)
		}
		
		fmt.Printf("[+] Successfully retrieved %d historical endpoints.\n", len(urls))
		
		// Print the first 5 just to verify
		limit := 5
		if len(urls) < 5 {
			limit = len(urls)
		}
		fmt.Println("[*] Sample endpoints:")
		for i := 0; i < limit; i++ {
			fmt.Printf("    - %s\n", urls[i])
		}
	},
}

func init() {
	rootCmd.AddCommand(reconCmd)
}
