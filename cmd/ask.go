package cmd

import (
	"fmt"
	"log"

	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/ai"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [prompt]",
	Short: "Ask the local LLM a security question",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		
		// Hardcoding phi4-mini for the test, we'll make this dynamic later
		model := "phi4-mini" 
		
		fmt.Printf("[*] Routing prompt to %s via GPU...\n", model)
		
		response, err := ai.QueryOllama(model, prompt)
		if err != nil {
			log.Fatalf("[-] AI Engine Failure: %v", err)
		}
		
		fmt.Printf("\n[+] Sentinel AI Response:\n%s\n", response)
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}
