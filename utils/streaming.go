package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// StreamHandler defines a callback for each chunk received
type StreamHandler func(chunk map[string]interface{}) error

// StreamRequest makes a streaming request to the given URL with auth token
func StreamRequest(ctx context.Context, url, token string, payload interface{}, handleChunk StreamHandler) error {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if len(line) == 0 {
			continue
		}

		var chunk map[string]interface{}
		if err := json.Unmarshal(line, &chunk); err != nil {
			// skip invalid JSON lines
			continue
		}

		if err := handleChunk(chunk); err != nil {
			return err
		}
	}

	return nil
}
