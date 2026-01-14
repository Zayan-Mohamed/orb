package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// createSession creates a new session with the relay server
func createSession(relayURL, sharedPath string) (string, string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	reqBody := map[string]string{
		"shared_path": sharedPath,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := client.Post(
		relayURL+"/session/create",
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to contact relay: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("relay error: %s", string(body))
	}

	var result struct {
		SessionID string `json:"session_id"`
		Passcode  string `json:"passcode"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.SessionID, result.Passcode, nil
}
