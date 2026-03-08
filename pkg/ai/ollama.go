package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OllamaRequest structures the JSON payload sent to the API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse structures the JSON returned by the API
type OllamaResponse struct {
	Response string `json:"response"`
}

// QueryOllama sends a prompt to the local Ollama instance
func QueryOllama(model, prompt string) (string, error) {
	url := "http://localhost:11434/api/generate"

	reqBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false, // We want the full response at once, not streaming
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making request to Ollama: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var ollamaResp OllamaResponse
	err = json.Unmarshal(bodyBytes, &ollamaResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return ollamaResp.Response, nil
}
