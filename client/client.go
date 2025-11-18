package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tigillo/githubmodels-go/models"
)

// Client is the main GitHub Models API client
type Client struct {
	token   string
	Client  *http.Client
	BaseURL string // exported so tests can override
}

// NewClient creates a new GitHub Models API client
func NewClient(token string) *Client {
	return &Client{
		token:   token,
		Client:  http.DefaultClient,
		BaseURL: "https://models.github.ai", // production default
	}
}

// Model represents a GitHub Models API model
type Model struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// ListModels returns all available models from the catalog
func (c *Client) ListModels(ctx context.Context) ([]Model, error) {
	url := fmt.Sprintf("%s/catalog/models", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var models []Model
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, err
	}

	return models, nil
}

// ChatCompletion sends a chat completion request to GitHub Models API
func (c *Client) ChatCompletion(ctx context.Context, reqData models.ChatRequest) (*models.ChatResponse, error) {
	url := fmt.Sprintf("%s/inference/chat/completions", c.BaseURL)

	bodyBytes, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var chatResp models.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, err
	}

	return &chatResp, nil
}
