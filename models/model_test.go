package models

import (
	"encoding/json"
	"testing"
)

func TestModelJSON(t *testing.T) {
	m := Model{
		ID:          "github/code-chat",
		Name:        "Code Chat",
		Description: "Chat with code model",
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal Model: %v", err)
	}

	var decoded Model
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal Model: %v", err)
	}

	if decoded.ID != m.ID || decoded.Name != m.Name || decoded.Description != m.Description {
		t.Errorf("decoded model does not match original")
	}
}

func TestChatRequestAndResponse(t *testing.T) {
	req := ChatRequest{
		Model: "github/code-chat",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal ChatRequest: %v", err)
	}

	var decoded ChatRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal ChatRequest: %v", err)
	}

	if decoded.Model != req.Model || len(decoded.Messages) != 1 || decoded.Messages[0].Content != "Hello" {
		t.Errorf("decoded ChatRequest does not match original")
	}

	resp := ChatResponse{
		ID: "chat-1",
		Choices: []Choice{
			{Message: Message{Role: "assistant", Content: "Hi!"}},
		},
	}

	dataResp, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal ChatResponse: %v", err)
	}

	var decodedResp ChatResponse
	if err := json.Unmarshal(dataResp, &decodedResp); err != nil {
		t.Fatalf("failed to unmarshal ChatResponse: %v", err)
	}

	if decodedResp.ID != resp.ID || decodedResp.Choices[0].Message.Content != "Hi!" {
		t.Errorf("decoded ChatResponse does not match original")
	}
}
