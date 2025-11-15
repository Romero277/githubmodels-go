package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tigillo/githubmodels-go/models"
)

func TestStreamRequest(t *testing.T) {
	chunks := []models.ChatResponse{
		{ID: "1", Choices: []models.Choice{{Message: models.Message{Role: "assistant", Content: "Hello"}}}},
		{ID: "1", Choices: []models.Choice{{Message: models.Message{Role: "assistant", Content: " World"}}}},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)
		for _, chunk := range chunks {
			line, _ := json.Marshal(chunk)
			w.Write(line)
			w.Write([]byte("\n"))
			flusher.Flush()
		}
	}))
	defer ts.Close()

	var result strings.Builder
	err := StreamRequest(context.Background(), ts.URL, "test-token", models.ChatRequest{}, func(chunk map[string]interface{}) error {
		data, _ := json.Marshal(chunk)
		var resp models.ChatResponse
		json.Unmarshal(data, &resp)
		for _, c := range resp.Choices {
			result.WriteString(c.Message.Content)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("StreamRequest failed: %v", err)
	}

	if result.String() != "Hello World" {
		t.Errorf("unexpected result: %q", result.String())
	}
}
