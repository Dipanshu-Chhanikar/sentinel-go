package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/ai"
	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/recon"
	"github.com/spf13/cobra"
)

var reconCmd = &cobra.Command{
	Use:   "recon [domain]",
	Short: "Passively gather historical URLs and analyze them with AI",
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
		
		// Take a chunk of URLs so we don't blow up the AI's 2048 context window
		chunkSize := 50
		if len(urls) < chunkSize {
			chunkSize = len(urls)
		}
		
		var urlList strings.Builder
		for i := 0; i < chunkSize; i++ {
			urlList.WriteString(urls[i] + "\n")
		}

		fmt.Printf("\n[*] Passing top %d URLs to local AI for triage...\n", chunkSize)
		
		// The strict prompt engineered to force a pure JSON array output
		prompt := fmt.Sprintf(`You are an expert bug bounty data parser. 
Analyze the following list of raw URLs and extract ONLY the URLs that might be interesting for security testing (e.g., API endpoints, sensitive files like .csv or .json, or admin paths).
You MUST return the output as a valid JSON array of strings. 
DO NOT include any conversational text, markdown formatting (like backticks), or explanations. 

Example of required output format:
["http://example.com/api/v1/data.json", "https://example.com/admin/login"]

URLs to analyze:
%s`, urlList.String())

		// Using llama3.2:3b because it handles strict formatting prompts well
		response, err := ai.QueryOllama("llama3.2:3b", prompt)
		if err != nil {
			log.Fatalf("[-] AI Engine Failure: %v", err)
		}
		
		fmt.Printf("\n[+] AI Triage Results:\n%s\n", response)
	},
}

func init() {
	rootCmd.AddCommand(reconCmd)
}
