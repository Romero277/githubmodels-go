package endpoints

import (
	"context"

	"github.com/tigillo/githubmodels-go/client"
)

// ListModels fetches all available models from GitHub Models catalog
func ListModels(ctx context.Context, c *client.Client) ([]client.Model, error) {
	var models []client.Model
	_, err := c.DoRequest(ctx, "GET", "/catalog/models", nil, &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}
