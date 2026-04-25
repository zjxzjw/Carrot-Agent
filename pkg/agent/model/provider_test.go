package model

import (
	"context"
	"testing"
)

func TestNewOpenAIProvider(t *testing.T) {
	provider := NewOpenAIProvider("test-api-key", "https://api.openai.com/v1", "gpt-4")

	if provider.APIKey != "test-api-key" {
		t.Errorf("Expected APIKey 'test-api-key', got '%s'", provider.APIKey)
	}

	if provider.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("Expected BaseURL 'https://api.openai.com/v1', got '%s'", provider.BaseURL)
	}

	if provider.Model != "gpt-4" {
		t.Errorf("Expected Model 'gpt-4', got '%s'", provider.Model)
	}
}

func TestNewOpenAIProviderWithEmptyBaseURL(t *testing.T) {
	provider := NewOpenAIProvider("test-api-key", "", "gpt-4")

	if provider.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("Expected default BaseURL 'https://api.openai.com/v1', got '%s'", provider.BaseURL)
	}
}

func TestNewClaudeProvider(t *testing.T) {
	provider := NewClaudeProvider("test-api-key", "https://api.anthropic.com/v1", "claude-3-opus")

	if provider.APIKey != "test-api-key" {
		t.Errorf("Expected APIKey 'test-api-key', got '%s'", provider.APIKey)
	}

	if provider.BaseURL != "https://api.anthropic.com/v1" {
		t.Errorf("Expected BaseURL 'https://api.anthropic.com/v1', got '%s'", provider.BaseURL)
	}

	if provider.Model != "claude-3-opus" {
		t.Errorf("Expected Model 'claude-3-opus', got '%s'", provider.Model)
	}
}

func TestNewClaudeProviderWithEmptyBaseURL(t *testing.T) {
	provider := NewClaudeProvider("test-api-key", "", "claude-3-opus")

	if provider.BaseURL != "https://api.anthropic.com/v1" {
		t.Errorf("Expected default BaseURL 'https://api.anthropic.com/v1', got '%s'", provider.BaseURL)
	}
}

func TestOpenAIProviderChatRequestValidation(t *testing.T) {
	req := &ChatRequest{
		Model:       "gpt-4",
		Messages:    []Message{{Role: "user", Content: "Hello"}},
		Temperature: 0.7,
		MaxTokens:   100,
	}

	if req.Model == "" {
		t.Error("Expected Model to be set")
	}

	if len(req.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(req.Messages))
	}
}

func TestClaudeProviderChatRequestValidation(t *testing.T) {
	req := map[string]interface{}{
		"model":      "claude-3-opus",
		"messages":   []Message{{Role: "user", Content: "Hello"}},
		"max_tokens": 100,
	}

	if req["model"] != "claude-3-opus" {
		t.Errorf("Expected model 'claude-3-opus', got '%v'", req["model"])
	}
}

func TestProviderFactory(t *testing.T) {
	factory := NewProviderFactory()

	if factory == nil {
		t.Fatal("Expected factory, got nil")
	}

	if factory.providers == nil {
		t.Error("Expected providers map to be initialized")
	}
}

func TestProviderFactoryRegister(t *testing.T) {
	factory := NewProviderFactory()

	provider := NewOpenAIProvider("test-key", "", "gpt-4")
	factory.Register("test", provider)

	p, ok := factory.Get("test")
	if !ok {
		t.Error("Expected to get registered provider")
	}

	if p != provider {
		t.Error("Expected to get same provider that was registered")
	}
}

func TestProviderFactoryGetNonExistent(t *testing.T) {
	factory := NewProviderFactory()

	_, ok := factory.Get("non-existent")
	if ok {
		t.Error("Expected ok to be false for non-existent provider")
	}
}

func TestProviderFactoryCreateProviderOpenAI(t *testing.T) {
	factory := NewProviderFactory()

	provider := factory.CreateProvider("openai", "test-key", "", "gpt-4")
	if provider == nil {
		t.Fatal("Expected provider, got nil")
	}

	openAIProvider, ok := provider.(*OpenAIProvider)
	if !ok {
		t.Error("Expected OpenAIProvider")
	}

	if openAIProvider.APIKey != "test-key" {
		t.Errorf("Expected APIKey 'test-key', got '%s'", openAIProvider.APIKey)
	}
}

func TestProviderFactoryCreateProviderClaude(t *testing.T) {
	factory := NewProviderFactory()

	provider := factory.CreateProvider("claude", "test-key", "", "claude-3-opus")
	if provider == nil {
		t.Fatal("Expected provider, got nil")
	}

	claudeProvider, ok := provider.(*ClaudeProvider)
	if !ok {
		t.Error("Expected ClaudeProvider")
	}

	if claudeProvider.APIKey != "test-key" {
		t.Errorf("Expected APIKey 'test-key', got '%s'", claudeProvider.APIKey)
	}
}

func TestProviderFactoryCreateProviderAnthropic(t *testing.T) {
	factory := NewProviderFactory()

	provider := factory.CreateProvider("anthropic", "test-key", "", "claude-3-sonnet")
	if provider == nil {
		t.Fatal("Expected provider, got nil")
	}

	claudeProvider, ok := provider.(*ClaudeProvider)
	if !ok {
		t.Error("Expected ClaudeProvider")
	}

	if claudeProvider.Model != "claude-3-sonnet" {
		t.Errorf("Expected Model 'claude-3-sonnet', got '%s'", claudeProvider.Model)
	}
}

func TestProviderFactoryCreateProviderOpenRouter(t *testing.T) {
	factory := NewProviderFactory()

	provider := factory.CreateProvider("openrouter", "test-key", "https://openrouter.ai/api", "anthropic/claude-3-opus")
	if provider == nil {
		t.Fatal("Expected provider, got nil")
	}

	openAIProvider, ok := provider.(*OpenAIProvider)
	if !ok {
		t.Error("Expected OpenAIProvider for openrouter")
	}

	if openAIProvider.BaseURL != "https://openrouter.ai/api" {
		t.Errorf("Expected BaseURL 'https://openrouter.ai/api', got '%s'", openAIProvider.BaseURL)
	}
}

func TestProviderFactoryCreateProviderDefault(t *testing.T) {
	factory := NewProviderFactory()

	provider := factory.CreateProvider("unknown", "test-key", "", "gpt-4")
	if provider == nil {
		t.Fatal("Expected provider, got nil")
	}

	_, ok := provider.(*OpenAIProvider)
	if !ok {
		t.Error("Expected OpenAIProvider as default")
	}
}

func TestChatRequestStruct(t *testing.T) {
	req := ChatRequest{
		Model:       "gpt-4",
		Messages:    []Message{{Role: "user", Content: "Hello"}, {Role: "assistant", Content: "Hi there"}},
		Temperature: 0.7,
		MaxTokens:   100,
		Tools:       []ToolDefinition{{Name: "test-tool", Description: "A test tool"}},
		Stream:      false,
	}

	if req.Model != "gpt-4" {
		t.Errorf("Expected Model 'gpt-4', got '%s'", req.Model)
	}

	if len(req.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(req.Messages))
	}

	if len(req.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(req.Tools))
	}

	if req.Stream != false {
		t.Error("Expected Stream to be false")
	}
}

func TestMessageStruct(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Hello, how are you?",
		Name:    "user123",
	}

	if msg.Role != "user" {
		t.Errorf("Expected Role 'user', got '%s'", msg.Role)
	}

	if msg.Content != "Hello, how are you?" {
		t.Errorf("Expected Content 'Hello, how are you?', got '%s'", msg.Content)
	}

	if msg.Name != "user123" {
		t.Errorf("Expected Name 'user123', got '%s'", msg.Name)
	}
}

func TestChatResponseStruct(t *testing.T) {
	resp := ChatResponse{
		Model: "gpt-4",
		Choices: []Choice{
			{
				Index:        0,
				Message:      Message{Role: "assistant", Content: "Hello!"},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
		StopReason: "stop",
	}

	if resp.Model != "gpt-4" {
		t.Errorf("Expected Model 'gpt-4', got '%s'", resp.Model)
	}

	if len(resp.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Usage.TotalTokens != 15 {
		t.Errorf("Expected TotalTokens 15, got %d", resp.Usage.TotalTokens)
	}
}

func TestToolDefinitionStruct(t *testing.T) {
	tool := ToolDefinition{
		Name:        "get_weather",
		Description: "Get current weather for a location",
		Parameters: ToolParameters{
			Type:       "object",
			Properties: map[string]Property{"location": {Type: "string", Description: "City name"}},
			Required:   []string{"location"},
		},
	}

	if tool.Name != "get_weather" {
		t.Errorf("Expected Name 'get_weather', got '%s'", tool.Name)
	}

	if tool.Parameters.Type != "object" {
		t.Errorf("Expected Type 'object', got '%s'", tool.Parameters.Type)
	}

	if len(tool.Parameters.Required) != 1 {
		t.Errorf("Expected 1 required param, got %d", len(tool.Parameters.Required))
	}
}

func TestOpenAIProviderContextCancellation(t *testing.T) {
	provider := NewOpenAIProvider("test-key", "", "gpt-4")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := &ChatRequest{
		Model:    "gpt-4",
		Messages: []Message{{Role: "user", Content: "Hello"}},
	}

	_, err := provider.Chat(ctx, req)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}
}