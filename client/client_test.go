package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tigillo/githubmodels-go/models"
)

// Helper to create a test client pointing to a mock server
func newTestClient(ts *httptest.Server) *Client {
	return &Client{
		token:   "test-token",
		BaseURL: ts.URL,
		Client:  ts.Client(),
	}
}

func TestListModels(t *testing.T) {
	// Mock response
	mockModels := []models.Model{
		{ID: "github/code-chat", Name: "Code Chat", Description: "Chat with code model"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/catalog/models" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing or wrong Authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockModels)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	ctx := context.Background()

	modelsList, err := c.ListModels(ctx)
	if err != nil {
		t.Fatalf("ListModels failed: %v", err)
	}

	if len(modelsList) != 1 {
		t.Fatalf("expected 1 model, got %d", len(modelsList))
	}
	if modelsList[0].ID != "github/code-chat" {
		t.Errorf("unexpected model ID: %s", modelsList[0].ID)
	}
}

func TestChatCompletion(t *testing.T) {
	mockResponse := models.ChatResponse{
		ID: "chat-1",
		Choices: []models.Choice{
			{Message: models.Message{Role: "assistant", Content: "Hello!"}},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/inference/chat/completions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	ctx := context.Background()

	req := models.ChatRequest{
		Model: "github/code-chat",
		Messages: []models.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp, err := c.ChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("ChatCompletion failed: %v", err)
	}

	if resp.ID != "chat-1" {
		t.Errorf("unexpected response ID: %s", resp.ID)
	}
	if len(resp.Choices) != 1 {
		t.Errorf("expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "Hello!" {
		t.Errorf("unexpected response content: %s", resp.Choices[0].Message.Content)
	}
}
