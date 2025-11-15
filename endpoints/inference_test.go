package endpoints

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tigillo/githubmodels-go/client"
	"github.com/tigillo/githubmodels-go/models"
)

// Minimal test for Inference wrapper
func TestCreateChatCompletion(t *testing.T) {
	mockResp := models.ChatResponse{
		ID: "chat-1",
		Choices: []models.Choice{
			{Message: models.Message{Role: "assistant", Content: "Hello"}},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResp)
	}))
	defer ts.Close()

	c := client.NewClient("test-token")
	c.BaseURL = ts.URL
	c.Client = ts.Client()

	ctx := context.Background()
	req := models.ChatRequest{
		Model: "github/code-chat",
		Messages: []models.Message{
			{Role: "user", Content: "Hi"},
		},
	}

	resp, err := ChatCompletion(ctx, c, req)
	if err != nil {
		t.Fatalf("CChatCompletion failed: %v", err)
	}

	if resp.ID != "chat-1" {
		t.Errorf("unexpected response ID: %s", resp.ID)
	}
}
