package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/ai"
	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/attack"
	"github.com/Dipanshu-Chhanikar/sentinel-go/pkg/recon"
	"github.com/spf13/cobra"
)

var targetFile string

var autoCmd = &cobra.Command{
	Use:   "auto [domain]",
	Short: "Run the complete autonomous Sentinel pipeline on a domain or list of domains",
	Run: func(cmd *cobra.Command, args []string) {
		var targets []string
		if targetFile != "" {
			file, err := os.Open(targetFile)
			if err != nil {
				log.Fatalf("[-] Failed to open target file: %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if domain := scanner.Text(); domain != "" {
					targets = append(targets, domain)
				}
			}
		} else if len(args) > 0 {
			targets = append(targets, args[0])
		} else {
			log.Fatal("[-] Please provide a domain or a target file using -f")
		}

		for _, domain := range targets {
			fmt.Printf("\n================================================\n")
			fmt.Printf("[$$$] INITIATING FULL PIPELINE FOR: %s\n", domain)
			fmt.Printf("================================================\n")

			// Step 1: PASSIVE RECON
			fmt.Println("[*] Step 1: Executing Passive Recon...")
			waybackURLs, err := recon.FetchWaybackURLs(domain)
			if err != nil {
				log.Printf("[-] Error fetching Wayback URLs: %v\n", err)
				continue
			}

			if len(waybackURLs) > 50 {
				waybackURLs = waybackURLs[:50]
			}

			// Using our strict JSON prompt!
			prompt := fmt.Sprintf(`You are an expert bug bounty data parser. 
Analyze the following list of raw URLs and extract ONLY the URLs that might be interesting for security testing (e.g., API endpoints, sensitive files like .csv or .json, or admin paths).
You MUST return the output as a valid JSON array of strings. 
DO NOT include any conversational text, markdown formatting, or explanations. 

URLs to analyze:
%s`, strings.Join(waybackURLs, "\n"))

			interestingEndpoints, err := ai.QueryOllama("llama3.2:3b", prompt)
			if err != nil {
				log.Printf("[-] Error querying AI: %v\n", err)
				continue
			}

			var parsedEndpoints []string
			if err := json.Unmarshal([]byte(interestingEndpoints), &parsedEndpoints); err != nil {
				log.Printf("[-] Error parsing JSON (AI hallucinated): %v\n", err)
				continue
			}
			fmt.Printf("[+] AI extracted %d high-value endpoints.\n", len(parsedEndpoints))

			// Step 2: INTELLIGENT ROUTING
			fmt.Println("\n[*] Step 2: AI Triage & Fuzzing...")
			for _, url := range parsedEndpoints {
				if strings.HasSuffix(url, ".js") {
					fmt.Printf("[*] Analyzing JS file: %s\n", url)
					jsContent, err := recon.FetchJSFile(url)
					if err != nil {
						continue
					}

					if len(jsContent) > 1500 {
						jsContent = jsContent[:1500]
					}

					secretPrompt := fmt.Sprintf(`Analyze the following JavaScript code snippet.
Look for hardcoded secrets, API keys, hidden endpoints, or developer comments.
If you find nothing, respond exactly with "CLEAN".
Code:
%s`, jsContent)

					secretData, err := ai.QueryOllama("deepseek-coder:6.7b", secretPrompt)
					if err == nil && !strings.Contains(secretData, "CLEAN") {
						fmt.Printf("[!] Secrets found in %s:\n%s\n", url, secretData)
					}
				} else if strings.Contains(url, "api") || strings.Contains(url, "admin") {
					statusCode, err := attack.CheckEndpoint(url, nil)
					if err == nil && (statusCode == 403 || statusCode == 401) {
						fmt.Printf("[!] Access Denied (%d) detected on %s\n", statusCode, url)
					}
				}
			}

			// Step 3: NUCLEI ORCHESTRATION
			fmt.Println("\n[*] Step 3: Executing Nuclei Strike...")
			findings, err := attack.RunNuclei(domain, "all")
			if err != nil {
				log.Printf("[-] Error running Nuclei: %v\n", err)
				continue
			}

			if len(findings) > 0 {
				fmt.Printf("[!] CRITICAL: Found %d Nuclei vulnerabilities!\n", len(findings))
				for _, finding := range findings {
					fmt.Printf("    -> %s\n", finding)
				}
			} else {
				fmt.Println("[+] No specific CVEs found by Nuclei.")
			}
		}

		fmt.Println("\n[+] Sentinel-Go Autonomous Run Complete.")
	},
}

func init() {
	autoCmd.Flags().StringVarP(&targetFile, "file", "f", "", "File containing a list of target domains")
	rootCmd.AddCommand(autoCmd)
}
