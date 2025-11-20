package endpoints

import (
	"context"

	"github.com/tigillo/githubmodels-go/client"
	"github.com/tigillo/githubmodels-go/models"
)

// ChatCompletion sends a chat request to the GitHub Models API
func ChatCompletion(ctx context.Context, c *client.Client, req models.ChatRequest) (*models.ChatResponse, error) {
	var resp models.ChatResponse
	headers, err := c.DoRequest(ctx, "POST", "/inference/chat/completions", req, &resp)

	// Always attach headers if available, even on error
	if headers != nil {
		resp.RateLimit = client.ParseRateLimitHeaders(headers)
	}

	if err != nil {
		// If we have headers (rate limits), return the partial response with the error
		if headers != nil {
			return &resp, err
		}
		return nil, err
	}
	return &resp, nil
}

// OrgChatCompletion sends a chat request to an organization-scoped endpoint
func OrgChatCompletion(ctx context.Context, c *client.Client, org string, req models.ChatRequest) (*models.ChatResponse, error) {
	path := "/orgs/" + org + "/inference/chat/completions"
	var resp models.ChatResponse
	headers, err := c.DoRequest(ctx, "POST", path, req, &resp)

	if headers != nil {
		resp.RateLimit = client.ParseRateLimitHeaders(headers)
	}

	if err != nil {
		if headers != nil {
			return &resp, err
		}
		return nil, err
	}
	return &resp, nil
}
