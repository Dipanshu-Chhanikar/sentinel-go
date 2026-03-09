package attack

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RunNuclei executes Nuclei against a target with specific tags
func RunNuclei(targetURL string, techTag string) ([]string, error) {
	// The command: nuclei -u target.com -tags techTag -silent
	cmd := exec.Command("nuclei", "-u", targetURL, "-tags", techTag, "-silent")
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Nuclei returns an error code if it finds nothing or fails, 
		// but we mostly care about the stdout string.
		fmt.Printf("[-] Nuclei execution note: %v\n", err)
	}

	// Parse the output into a slice of strings
	rawOutput := out.String()
	if rawOutput == "" {
		return nil, nil
	}

	var findings []string
	lines := strings.Split(strings.TrimSpace(rawOutput), "\n")
	for _, line := range lines {
		if line != "" {
			findings = append(findings, line)
		}
	}

	return findings, nil
}
