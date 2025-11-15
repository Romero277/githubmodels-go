package endpoints

import (
	"context"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/tigillo/githubmodels-go/client"
	"github.com/tigillo/githubmodels-go/models"
)

// Minimal test for ListModels wrapper
func TestListModels(t *testing.T) {
	mockModels := []models.Model{
		{ID: "github/code-chat", Name: "Code Chat", Description: "Chat with code model"},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockModels)
	}))
	defer ts.Close()

	c := client.NewClient("test-token")
	c.BaseURL = ts.URL
	c.Client = ts.Client()

	ctx := context.Background()
	modelsList, err := ListModels(ctx, c)
	if err != nil {
		t.Fatalf("ListModels failed: %v", err)
	}

	if len(modelsList) != 1 || modelsList[0].ID != "github/code-chat" {
		t.Errorf("unexpected models list: %+v", modelsList)
	}
}
