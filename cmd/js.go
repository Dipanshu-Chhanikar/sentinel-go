package cmd

import (
	"fmt"
	"log"

	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/ai"
	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/recon"
	"github.com/spf13/cobra"
)

var jsCmd = &cobra.Command{
	Use:   "js [url]",
	Short: "Download a JS file and analyze it with DeepSeek for secrets",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jsURL := args[0]
		
		fmt.Printf("[*] Fetching JavaScript from: %s\n", jsURL)
		
		jsCode, err := recon.FetchJSFile(jsURL)
		if err != nil {
			log.Fatalf("[-] JS Fetch failed: %v", err)
		}
		
		fmt.Printf("[+] Successfully downloaded %d bytes of JavaScript.\n", len(jsCode))
		
		// Slice the JS to the first 1500 characters to fit the context window safely
		chunk := jsCode
		if len(chunk) > 1500 {
			chunk = chunk[:1500]
		}
		
		fmt.Println("[*] Passing code chunk to DeepSeek-Coder (GPU) for semantic analysis...")
		
		prompt := fmt.Sprintf(`You are a senior security researcher auditing front-end JavaScript.
Analyze the following JavaScript code snippet.
Look for hardcoded secrets, API keys, hidden endpoints (like /api/v1/internal), or developer comments.
If you find nothing interesting, respond exactly with "CLEAN".
If you find something, list it out clearly.

Code Snippet:
%s`, chunk)

		// Routing to deepseek-coder since it is specifically trained on code logic
		response, err := ai.QueryOllama("deepseek-coder:6.7b", prompt)
		if err != nil {
			log.Fatalf("[-] AI Engine Failure: %v", err)
		}
		
		fmt.Printf("\n[+] DeepSeek Analysis:\n%s\n", response)
	},
}

func init() {
	rootCmd.AddCommand(jsCmd)
}
