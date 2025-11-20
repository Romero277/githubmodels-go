package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DoRequest is a helper to make HTTP requests to GitHub Models API
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}, result interface{}) (http.Header, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Extract only relevant headers
	headers := make(http.Header)
	for k, v := range resp.Header {
		if k == "X-RateLimit-Limit" || k == "X-RateLimit-Remaining" || k == "X-RateLimit-Reset" || k == "Retry-After" {
			headers[k] = v
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read response body for error message
		respBody, _ := io.ReadAll(resp.Body)
		return headers, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return headers, err
		}
	}

	return headers, nil
}
