package models

// Message represents a single message in a chat request
type Message struct {
	Role    string `json:"role"`    // "user", "system", "assistant"
	Content string `json:"content"` // Message content
}

// ChatRequest represents a request to the chat completion endpoint
type ChatRequest struct {
	Model    string    `json:"model"`    // Model ID, e.g., "github/code-chat"
	Messages []Message `json:"messages"` // Conversation messages
}

// Choice represents a single choice in the chat response
type Choice struct {
	Message Message `json:"message"` // The generated message from the model
}

// ChatResponse represents the response from the chat completion endpoint
type ChatResponse struct {
	ID      string   `json:"id"`      // Response ID
	Object  string   `json:"object"`  // Type of object, e.g., "chat.completion"
	Choices []Choice `json:"choices"` // List of choices
}
