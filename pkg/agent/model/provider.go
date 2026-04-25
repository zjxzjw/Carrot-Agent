package model

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Provider interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
}

type ChatRequest struct {
	Model       string           `json:"model"`
	Messages    []Message        `json:"messages"`
	Temperature float64          `json:"temperature"`
	MaxTokens   int              `json:"max_tokens"`
	Tools       []ToolDefinition `json:"tools,omitempty"`
	Stream      bool             `json:"stream,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type ChatResponse struct {
	Model      string    `json:"model"`
	Choices    []Choice  `json:"choices"`
	Usage      Usage     `json:"usage"`
	StopReason string    `json:"stop_reason,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string    `json:"id"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
}

type ToolParameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type OpenAIProvider struct {
	APIKey  string
	BaseURL string
	Model   string
	Client  *http.Client
}

func NewOpenAIProvider(apiKey, baseURL, model string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAIProvider{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
		Client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (p *OpenAIProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = p.Model
	}

	url := p.BaseURL + "/chat/completions"

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &chatResp, nil
}

type ClaudeProvider struct {
	APIKey  string
	BaseURL string
	Model   string
	Client  *http.Client
}

func NewClaudeProvider(apiKey, baseURL, model string) *ClaudeProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}

	return &ClaudeProvider{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
		Client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (p *ClaudeProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	claudeReq := map[string]interface{}{
		"model": p.Model,
		"messages": req.Messages,
		"max_tokens": req.MaxTokens,
	}

	if req.Temperature > 0 {
		claudeReq["temperature"] = req.Temperature
	}

	payload, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := p.BaseURL + "/messages"

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var claudeResp struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	chatResp := &ChatResponse{
		Model: p.Model,
		Usage: Usage{
			PromptTokens:     claudeResp.Usage.InputTokens,
			CompletionTokens: claudeResp.Usage.OutputTokens,
			TotalTokens:     claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		},
	}

	if len(claudeResp.Content) > 0 {
		chatResp.Choices = []Choice{
			{
				Message: Message{
					Role:    "assistant",
					Content: claudeResp.Content[0].Text,
				},
				FinishReason: "stop",
			},
		}
	}

	return chatResp, nil
}

type ProviderFactory struct {
	providers map[string]Provider
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[string]Provider),
	}
}

func (f *ProviderFactory) Register(name string, provider Provider) {
	f.providers[name] = provider
}

func (f *ProviderFactory) Get(name string) (Provider, bool) {
	p, ok := f.providers[name]
	return p, ok
}

func (f *ProviderFactory) CreateProvider(providerType, apiKey, baseURL, model string) Provider {
	switch strings.ToLower(providerType) {
	case "openai", "openrouter":
		return NewOpenAIProvider(apiKey, baseURL, model)
	case "claude", "anthropic":
		return NewClaudeProvider(apiKey, baseURL, model)
	default:
		return NewOpenAIProvider(apiKey, baseURL, model)
	}
}