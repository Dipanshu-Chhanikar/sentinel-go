package recon

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// FetchJSFile downloads the raw text of a JavaScript file
func FetchJSFile(jsURL string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	
	req, err := http.NewRequest("GET", jsURL, nil)
	if err != nil {
		return "", err
	}
	
	// Masquerade as a normal browser to avoid WAF blocks
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch JS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read JS body: %v", err)
	}

	return string(bodyBytes), nil
}
