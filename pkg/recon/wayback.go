package recon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FetchWaybackURLs queries the CDX API for a target domain
func FetchWaybackURLs(domain string) ([]string, error) {
	// Added &limit=5000 to prevent massive memory spikes and API timeouts
	apiURL := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=json&fl=original&collapse=urlkey&limit=5000", domain)

	// Bumped timeout to 60 seconds for slower archive responses
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to reach Wayback Machine: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var records [][]string
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	var urls []string
	for i, row := range records {
		if i > 0 && len(row) > 0 {
			urls = append(urls, row[0])
		}
	}

	return urls, nil
}
