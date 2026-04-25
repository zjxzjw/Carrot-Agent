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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type StreamHandler func(*ChatResponse) error

type Provider interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	ChatStream(ctx context.Context, req *ChatRequest, handler StreamHandler) error
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
	Role             string     `json:"role"`
	Content          string     `json:"content"`
	ReasoningContent string     `json:"reasoning_content,omitempty"` // DeepSeek thinking mode
	Name             string     `json:"name,omitempty"`
	ToolCallID       string     `json:"tool_call_id,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"` // DeepSeek puts tool_calls inside message
	Type             string     `json:"type"`                 // Required for assistant with tool_calls and tool messages
}

type ChatResponse struct {
	Model      string   `json:"model"`
	Choices    []Choice `json:"choices"`
	Usage      Usage    `json:"usage"`
	StopReason string   `json:"stop_reason,omitempty"`
	Done       bool     `json:"done"`
}

type Choice struct {
	Index        int        `json:"index"`
	Message      Message    `json:"message"`
	FinishReason string     `json:"finish_reason"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
}

type Delta struct {
	Role             string     `json:"role,omitempty"`
	Content          string     `json:"content,omitempty"`
	ReasoningContent string     `json:"reasoning_content,omitempty"` // DeepSeek thinking mode
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
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
	Type     string          `json:"type"`
	Function ToolFunctionDef `json:"function"`
}

type ToolFunctionDef struct {
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

	// Debug: Log raw response body to see actual structure
	fmt.Printf("[DEBUG] Raw API response (first 2000 chars):\n%s\n", string(body[:min(len(body), 2000)]))

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Debug: Log the parsed response structure for tool calls
	if len(chatResp.Choices) > 0 {
		choice := chatResp.Choices[0]
		// DeepSeek puts tool_calls inside message, not at choice level
		toolCallsCount := len(choice.Message.ToolCalls)
		fmt.Printf("[DEBUG] Finish reason: %s, ToolCalls count: %d (in message), Content length: %d\n",
			choice.FinishReason, toolCallsCount, len(choice.Message.Content))
		if toolCallsCount > 0 {
			for i, tc := range choice.Message.ToolCalls {
				fmt.Printf("[DEBUG] ToolCall[%d]: ID=%s, Function.Name=%s, Function.Arguments=%s\n",
					i, tc.ID, tc.Function.Name, tc.Function.Arguments)
			}
		}
		// Also log the raw Message field to see if tool_calls are there
		messageJSON, _ := json.Marshal(choice.Message)
		fmt.Printf("[DEBUG] Parsed Message field: %s\n", string(messageJSON))
	}

	return &chatResp, nil
}

func (p *OpenAIProvider) ChatStream(ctx context.Context, req *ChatRequest, handler StreamHandler) error {
	if req.Model == "" {
		req.Model = p.Model
	}
	req.Stream = true

	url := p.BaseURL + "/chat/completions"

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Cache-Control", "no-cache")
	httpReq.Header.Set("Connection", "keep-alive")
	httpReq.Header.Set("X-Request-ID", fmt.Sprintf("%d", time.Now().UnixNano()))

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	reader := resp.Body
	buffer := make([]byte, 0, 4096)
	lineBuffer := make([]byte, 0, 4096)

	accumulatedContent := ""
	var totalUsage Usage

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read response: %w", err)
		}

		buffer = append(buffer, buf[:n]...)

		for len(buffer) > 0 {
			idx := 0
			for i, b := range buffer {
				if b == '\n' {
					idx = i
					break
				}
			}

			if idx == 0 && len(buffer) > 0 && buffer[0] == '\n' {
				buffer = buffer[1:]
				continue
			}

			if idx == 0 {
				break
			}

			line := string(buffer[:idx])
			buffer = buffer[idx+1:]

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				chatResp := &ChatResponse{
					Model:   req.Model,
					Done:    true,
					Choices: []Choice{{Index: 0, FinishReason: "stop"}},
					Usage:   totalUsage,
				}
				if err := handler(chatResp); err != nil {
					return err
				}
				return nil
			}

			var chunk struct {
				ID      string `json:"id"`
				Object  string `json:"object"`
				Created int64  `json:"created"`
				Model   string `json:"model"`
				Choices []struct {
					Index int `json:"index"`
					Delta struct {
						Content          string `json:"content,omitempty"`
						ReasoningContent string `json:"reasoning_content,omitempty"` // DeepSeek thinking mode
						ToolCalls        []struct {
							ID       string `json:"id"`
							Type     string `json:"type"`
							Function struct {
								Name      string `json:"name"`
								Arguments string `json:"arguments"`
							} `json:"function"`
						} `json:"tool_calls"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
				Usage struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
				} `json:"usage"`
			}

			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) > 0 {
				delta := chunk.Choices[0].Delta
				accumulatedContent += delta.Content

				totalUsage = Usage{
					PromptTokens:     chunk.Usage.PromptTokens,
					CompletionTokens: chunk.Usage.CompletionTokens,
					TotalTokens:      chunk.Usage.TotalTokens,
				}

				chatResp := &ChatResponse{
					Model: chunk.Model,
					Choices: []Choice{{
						Index: chunk.Choices[0].Index,
						Message: Message{
							Role:             "assistant",
							Content:          delta.Content,
							ReasoningContent: delta.ReasoningContent,
						},
						FinishReason: chunk.Choices[0].FinishReason,
					}},
					Usage: totalUsage,
				}

				if len(delta.ToolCalls) > 0 {
					toolCalls := make([]ToolCall, len(delta.ToolCalls))
					for i, tc := range delta.ToolCalls {
						toolCalls[i] = ToolCall{
							ID:   tc.ID,
							Type: "function",
							Function: FunctionCall{
								Name:      tc.Function.Name,
								Arguments: tc.Function.Arguments,
							},
						}
					}
					chatResp.Choices[0].ToolCalls = toolCalls
				}

				if err := handler(chatResp); err != nil {
					return err
				}

				if chunk.Choices[0].FinishReason == "stop" {
					return nil
				}
			}

			_ = lineBuffer
			_ = accumulatedContent
		}
	}

	return nil
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
		"model":      p.Model,
		"messages":   req.Messages,
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
			TotalTokens:      claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
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

func (p *ClaudeProvider) ChatStream(ctx context.Context, req *ChatRequest, handler StreamHandler) error {
	return fmt.Errorf("Claude provider does not support streaming")
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
