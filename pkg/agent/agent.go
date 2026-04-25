package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"carrotagent/carrot-agent/pkg/agent/memory"
	"carrotagent/carrot-agent/pkg/agent/model"
	"carrotagent/carrot-agent/pkg/agent/skill"
	"carrotagent/carrot-agent/pkg/agent/tool"
	"carrotagent/carrot-agent/pkg/logger"
	"carrotagent/carrot-agent/pkg/storage"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func buildToolParams(properties map[string]interface{}, required []string) map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": properties,
		"required":   required,
	}
}

var (
	allowedPaths    = []string{"~/.carrot", "/tmp"}
	blockListRegex  = regexp.MustCompile(`(?i)(proc|sys|etc/passwd|etc/shadow|\.ssh|\.aws|\.git/config)`)
	privateIPRanges = []*net.IPNet{
		parseCIDR("10.0.0.0/8"),
		parseCIDR("172.16.0.0/12"),
		parseCIDR("192.168.0.0/16"),
		parseCIDR("127.0.0.0/8"),
		parseCIDR("169.254.0.0/16"),
	}
)

func parseCIDR(cidr string) *net.IPNet {
	_, ipNet, _ := net.ParseCIDR(cidr)
	return ipNet
}

func isPathAllowed(path string) bool {
	absPath, err := filepath.Abs(os.ExpandEnv(path))
	if err != nil {
		return false
	}

	if blockListRegex.MatchString(absPath) {
		return false
	}

	for _, allowed := range allowedPaths {
		allowedAbs, _ := filepath.Abs(os.ExpandEnv(allowed))
		if strings.HasPrefix(absPath, allowedAbs) {
			return true
		}
	}

	return false
}

func isURLAllowed(rawURL string) bool {
	if blockListRegex.MatchString(rawURL) {
		return false
	}

	urlLower := strings.ToLower(rawURL)
	if strings.HasPrefix(urlLower, "http://") || strings.HasPrefix(urlLower, "https://") {
		return true
	}

	host, _, err := net.SplitHostPort(rawURL)
	if err != nil {
		host = rawURL
	}

	ip := net.ParseIP(host)
	if ip != nil {
		for _, privateRange := range privateIPRanges {
			if privateRange.Contains(ip) {
				return false
			}
		}
	}

	return true
}

type AIAgent struct {
	name          string
	version       string
	store         *storage.Store
	Memory        *memory.MemoryManager
	SkillManager  *skill.SkillManager
	modelProvider model.Provider
	toolRegistry  *tool.ToolRegistry
	config        *AgentConfig
	conversation  []model.Message
	skillNudgeInt int
	toolCallCount int
	mu            sync.RWMutex
}

type AgentConfig struct {
	Name          string
	Version       string
	DataDir       string
	ModelProvider string
	ModelName     string
	Temperature   float64
	MaxTokens     int
	EnableSkills  bool
	EnableMemory  bool
	SkillNudgeInt int
}

type AgentOption func(*AIAgent)

func WithModelProvider(provider model.Provider) AgentOption {
	return func(a *AIAgent) {
		a.modelProvider = provider
	}
}

func WithSkillManager(skillMgr *skill.SkillManager) AgentOption {
	return func(a *AIAgent) {
		a.SkillManager = skillMgr
	}
}

func WithMemoryManager(memMgr *memory.MemoryManager) AgentOption {
	return func(a *AIAgent) {
		a.Memory = memMgr
	}
}

func WithToolRegistry(registry *tool.ToolRegistry) AgentOption {
	return func(a *AIAgent) {
		a.toolRegistry = registry
	}
}

func (a *AIAgent) GetToolRegistry() *tool.ToolRegistry {
	return a.toolRegistry
}

func (a *AIAgent) GetModelProvider() model.Provider {
	return a.modelProvider
}

func (a *AIAgent) GetConfig() *AgentConfig {
	return a.config
}

func (a *AIAgent) GetMessages() []model.Message {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.conversation
}

func (a *AIAgent) AddMessage(msg model.Message) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.conversation = append(a.conversation, msg)
}

func (a *AIAgent) ClearMessages() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.conversation = make([]model.Message, 0)
}

func NewAIAgent(cfg *AgentConfig, store *storage.Store, opts ...AgentOption) *AIAgent {
	agent := &AIAgent{
		name:          "carrot-agent",
		version:       "0.1.0",
		store:         store,
		Memory:        memory.NewMemoryManager(store),
		SkillManager:  skill.NewSkillManager(store),
		toolRegistry:  tool.NewToolRegistry(),
		conversation:  make([]model.Message, 0),
		skillNudgeInt: 10,
		toolCallCount: 0,
		config:        cfg,
	}

	if cfg.SkillNudgeInt > 0 {
		agent.skillNudgeInt = cfg.SkillNudgeInt
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

func (a *AIAgent) Initialize(ctx context.Context) error {
	if err := a.Memory.Load(ctx); err != nil {
		return fmt.Errorf("failed to load memory: %w", err)
	}

	if err := a.SkillManager.Load(ctx); err != nil {
		return fmt.Errorf("failed to load skills: %w", err)
	}

	a.registerDefaultTools()

	return nil
}

func (a *AIAgent) registerDefaultTools() {
	a.toolRegistry.Register("memory_read",
		"Read from memory storage",
		buildToolParams(
			map[string]interface{}{
				"memory_id": map[string]interface{}{"type": "string", "description": "Memory ID to read"},
			},
			[]string{"memory_id"},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			id, _ := args["memory_id"].(string)
			mem, err := a.Memory.Get(id)
			if err != nil {
				return nil, err
			}
			if mem == nil {
				return "Memory not found", nil
			}
			return mem.Content, nil
		})

	a.toolRegistry.Register("memory_write",
		"Write to memory storage",
		buildToolParams(
			map[string]interface{}{
				"type":     map[string]interface{}{"type": "string", "description": "Memory type (snapshot, session, longterm)"},
				"content":  map[string]interface{}{"type": "string", "description": "Content to store"},
				"metadata": map[string]interface{}{"type": "string", "description": "Optional metadata as JSON"},
			},
			[]string{"type", "content"},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			memType, _ := args["type"].(string)
			content, _ := args["content"].(string)
			metadata, _ := args["metadata"].(string)
			if metadata == "" {
				metadata = "{}"
			}

			if err := a.Memory.Add(ctx, memType, content, metadata); err != nil {
				return nil, err
			}
			return "Memory saved successfully", nil
		})

	a.toolRegistry.Register("skill_create",
		"Create a new skill",
		buildToolParams(
			map[string]interface{}{
				"name":        map[string]interface{}{"type": "string", "description": "Skill name"},
				"description": map[string]interface{}{"type": "string", "description": "Skill description"},
				"content":     map[string]interface{}{"type": "string", "description": "Skill content in markdown format"},
			},
			[]string{"name", "description", "content"},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			name, _ := args["name"].(string)
			description, _ := args["description"].(string)
			content, _ := args["content"].(string)

			if err := a.SkillManager.Create(ctx, name, description, content); err != nil {
				return nil, err
			}
			return fmt.Sprintf("Skill '%s' created successfully", name), nil
		})

	a.toolRegistry.Register("skill_update",
		"Update an existing skill",
		buildToolParams(
			map[string]interface{}{
				"skill_id": map[string]interface{}{"type": "string", "description": "Skill ID to update"},
				"content":  map[string]interface{}{"type": "string", "description": "New skill content"},
			},
			[]string{"skill_id", "content"},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			skillID, _ := args["skill_id"].(string)
			content, _ := args["content"].(string)

			if err := a.SkillManager.Update(ctx, skillID, content); err != nil {
				return nil, err
			}
			return "Skill updated successfully", nil
		})

	a.toolRegistry.Register("skill_list",
		"List all available skills",
		buildToolParams(
			map[string]interface{}{},
			[]string{},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			skills := a.SkillManager.List(100)
			if len(skills) == 0 {
				return "No skills available.", nil
			}

			var lines []string
			lines = append(lines, "# Skills\n")
			for _, s := range skills {
				lines = append(lines, fmt.Sprintf("- **%s**: %s (v%s)", s.Name, s.Description, s.Version))
			}
			return strings.Join(lines, "\n"), nil
		})

	a.toolRegistry.Register("skill_search",
		"Search skills by keyword",
		buildToolParams(
			map[string]interface{}{
				"keyword": map[string]interface{}{"type": "string", "description": "Search keyword"},
			},
			[]string{"keyword"},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			keyword, _ := args["keyword"].(string)
			skills, err := a.SkillManager.Search(keyword, 50)
			if err != nil {
				return nil, err
			}

			if len(skills) == 0 {
				return fmt.Sprintf("No skills found matching '%s'", keyword), nil
			}

			var lines []string
			lines = append(lines, fmt.Sprintf("# Skills matching '%s'\n", keyword))
			for _, s := range skills {
				lines = append(lines, fmt.Sprintf("- **%s**: %s", s.Name, s.Description))
			}
			return strings.Join(lines, "\n"), nil
		})

	a.toolRegistry.Register("system_info",
		"Get system information",
		buildToolParams(
			map[string]interface{}{},
			[]string{},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			info := map[string]interface{}{
				"os":            runtime.GOOS,
				"arch":          runtime.GOARCH,
				"go_version":    runtime.Version(),
				"agent_version": a.version,
				"time":          time.Now().Format(time.RFC3339),
			}
			return info, nil
		})

	a.toolRegistry.Register("file_read",
		"Read file content",
		buildToolParams(
			map[string]interface{}{
				"file_path": map[string]interface{}{"type": "string", "description": "File path to read"},
			},
			[]string{"file_path"},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			filePath, _ := args["file_path"].(string)
			if !isPathAllowed(filePath) {
				return nil, fmt.Errorf("access denied: path %q is not allowed", filePath)
			}
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}
			return string(content), nil
		})

	a.toolRegistry.Register("file_write",
		"Write content to file",
		buildToolParams(
			map[string]interface{}{
				"file_path": map[string]interface{}{"type": "string", "description": "File path to write"},
				"content":   map[string]interface{}{"type": "string", "description": "Content to write"},
			},
			[]string{"file_path", "content"},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			filePath, _ := args["file_path"].(string)
			if !isPathAllowed(filePath) {
				return nil, fmt.Errorf("access denied: path %q is not allowed", filePath)
			}
			content, _ := args["content"].(string)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return nil, fmt.Errorf("failed to write file: %w", err)
			}
			return "File written successfully", nil
		})

	a.toolRegistry.Register("http_get",
		"Send HTTP GET request",
		buildToolParams(
			map[string]interface{}{
				"url": map[string]interface{}{"type": "string", "description": "URL to request"},
			},
			[]string{"url"},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			url, _ := args["url"].(string)
			if !isURLAllowed(url) {
				return nil, fmt.Errorf("access denied: URL %q is not allowed", url)
			}
			resp, err := http.Get(url)
			if err != nil {
				return nil, fmt.Errorf("failed to send request: %w", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response: %w", err)
			}

			return map[string]interface{}{
				"status":      resp.Status,
				"status_code": resp.StatusCode,
				"content":     string(body),
			}, nil
		})

	a.toolRegistry.Register("get_time",
		"Get current time",
		buildToolParams(
			map[string]interface{}{},
			[]string{},
		),
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			return map[string]interface{}{
				"current_time": time.Now().Format(time.RFC3339),
				"unix_time":    time.Now().Unix(),
			}, nil
		})
}

func (a *AIAgent) RunConversation(ctx context.Context, userInput string) (*model.ChatResponse, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	logger.Info("Received user input: %s", userInput)

	if userInput != "" {
		a.conversation = append(a.conversation, model.Message{
			Role:    "user",
			Content: userInput,
		})
	}

	systemPrompt := a.buildSystemPrompt()
	messages := append([]model.Message{{Role: "system", Content: systemPrompt}}, a.conversation...)

	logger.Info("MESSAGES_DEBUG: Built messages array, conversation length=%d", len(a.conversation))
	for i, m := range a.conversation {
		logger.Info("MESSAGES_DEBUG: conversation[%d] Role=%q Type=%q ToolCalls=%d", i, m.Role, m.Type, len(m.ToolCalls))
	}

	req := &model.ChatRequest{
		Model:       a.config.ModelName,
		Messages:    messages,
		Temperature: a.config.Temperature,
		MaxTokens:   a.config.MaxTokens,
	}

	// Ensure all messages have proper type field for DeepSeek compatibility
	for i := range req.Messages {
		if req.Messages[i].Role == "assistant" && len(req.Messages[i].ToolCalls) > 0 && req.Messages[i].Type == "" {
			req.Messages[i].Type = "function"
		}
		if req.Messages[i].Role == "tool" && req.Messages[i].Type == "" {
			req.Messages[i].Type = "function"
		}
	}

	tools := tool.ConvertToModelTools(a.toolRegistry)
	if len(tools) > 0 {
		req.Tools = make([]model.ToolDefinition, len(tools))
		for i, t := range tools {
			if fn, ok := t["function"].(map[string]interface{}); ok {
				functionDef := model.ToolFunctionDef{
					Name:        getString(fn, "name"),
					Description: getString(fn, "description"),
				}
				// Always set parameters, even if not present
				properties := make(map[string]model.Property)
				required := []string{}

				if params, ok := fn["parameters"].(map[string]interface{}); ok {
					if props, ok := params["properties"].(map[string]interface{}); ok {
						for name, prop := range props {
							if propMap, ok := prop.(map[string]interface{}); ok {
								properties[name] = model.Property{
									Type:        getString(propMap, "type"),
									Description: getString(propMap, "description"),
								}
							}
						}
					}

					if reqList, ok := params["required"].([]interface{}); ok {
						for _, req := range reqList {
							if reqStr, ok := req.(string); ok {
								required = append(required, reqStr)
							}
						}
					}
				}

				functionDef.Parameters = model.ToolParameters{
					Type:       "object",
					Properties: properties,
					Required:   required,
				}
				toolDef := model.ToolDefinition{
					Type:     "function",
					Function: functionDef,
				}
				req.Tools[i] = toolDef
			}
		}
		logger.Debug("Registered %d tools for model", len(tools))
		logger.Info("Tool names: %v", func() []string {
			names := make([]string, len(tools))
			for i, t := range tools {
				if fn, ok := t["function"].(map[string]interface{}); ok {
					names[i] = getString(fn, "name")
				}
			}
			return names
		}())
	} else {
		logger.Warn("No tools registered!")
	}

	if a.modelProvider == nil {
		logger.Error("No model provider configured")
		return nil, fmt.Errorf("no model provider configured")
	}

	logger.Info("Sending request to model: %s", a.config.ModelName)

	// Debug: log the request messages structure
	reqJSON, _ := json.Marshal(req)
	logger.Info("REQUEST_DEBUG: Full Request payload: %s", string(reqJSON))

	// Log individual messages for debugging
	for i, msg := range req.Messages {
		msgJSON, _ := json.Marshal(msg)
		logger.Info("REQUEST_DEBUG: Message[%d] role=%s type=%s: %s", i, msg.Role, msg.Type, string(msgJSON[:min(len(msgJSON), 500)]))
	}

	resp, err := a.modelProvider.Chat(ctx, req)
	if err != nil {
		logger.Error("Failed to get model response: %v", err)
		return nil, fmt.Errorf("failed to get model response: %w", err)
	}

	// Debug: log the raw response
	respJSON, _ := json.Marshal(resp)
	logger.Debug("Raw model response: %s", string(respJSON))

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		msg := choice.Message
		logger.Info("STORAGE_DEBUG: Before storing assistant message - Type=%q, Role=%q, ToolCalls=%d, ReasoningContent=%d", 
			msg.Type, msg.Role, len(msg.ToolCalls), len(msg.ReasoningContent))
		if len(msg.ToolCalls) > 0 {
			msg.Type = "function"
		}
		logger.Info("STORAGE_DEBUG: After fix before append - Type=%q", msg.Type)
		a.conversation = append(a.conversation, msg)
		logger.Info("Received model response with finish reason: %s", choice.FinishReason)
		logger.Info("Message content: %s", msg.Content)
		if msg.ReasoningContent != "" {
			logger.Info("Reasoning content present: %d chars", len(msg.ReasoningContent))
		}
		// DeepSeek puts tool_calls inside message
		if len(msg.ToolCalls) > 0 {
			logger.Info("Tool calls present: %d", len(msg.ToolCalls))
		}
	} else {
		logger.Warn("No choices returned from model")
	}

	return resp, nil
}

func (a *AIAgent) buildSystemPrompt() string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("You are %s, version %s.\n\n", a.name, a.version))

	if a.config.EnableMemory {
		snapshot := a.Memory.GetSnapshotContent()
		if snapshot != "" {
			prompt.WriteString("## Memory\n")
			prompt.WriteString(snapshot)
			prompt.WriteString("\n\n")
		}
	}

	if a.config.EnableSkills {
		skillsIndex := a.SkillManager.GetSkillsIndex()
		if skillsIndex != "" {
			prompt.WriteString(skillsIndex)
			prompt.WriteString("\n")
		}

		prompt.WriteString("\n## Skill Creation Guidance\n")
		prompt.WriteString("After completing a complex task (5+ tool calls), fixing a tricky error, or discovering a non-trivial workflow, save the approach as a skill using the skill_create tool.\n\n")
	}

	prompt.WriteString("## Memory Update Guidance\n")
	prompt.WriteString("During conversation, if you learn important information about the user or environment, save it using the memory_write tool with type 'snapshot'.\n\n")

	prompt.WriteString(a.toolRegistry.GetToolsForPrompt())

	return prompt.String()
}

func (a *AIAgent) ProcessToolCalls(ctx context.Context, toolCalls []map[string]interface{}) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0, len(toolCalls))

	for _, tc := range toolCalls {
		toolCallID, _ := tc["id"].(string)
		parsed, err := tool.ParseToolCall(tc)
		if err != nil {
			results = append(results, map[string]interface{}{
				"tool_call_id": toolCallID,
				"output":       fmt.Sprintf("Error parsing tool call: %v", err),
			})
			logger.Error("Failed to parse tool call: %v", err)
			continue
		}

		// Create a context with timeout for tool execution
		toolCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		result := a.toolRegistry.Execute(toolCtx, parsed.Name, parsed.Arguments)
		cancel()

		output := ""
		if result.Success {
			if outputBytes, err := json.Marshal(result.Output); err == nil {
				output = string(outputBytes)
			} else {
				output = fmt.Sprintf("%v", result.Output)
			}
			logger.Info("Tool %s executed successfully", parsed.Name)
		} else {
			output = result.Error
			logger.Error("Tool %s execution failed: %s", parsed.Name, result.Error)
		}

		results = append(results, map[string]interface{}{
			"tool_call_id": toolCallID,
			"output":       output,
		})

		a.toolCallCount++

		if a.toolCallCount > 0 && a.toolCallCount%a.skillNudgeInt == 0 {
			a.triggerSkillNudge(ctx)
		}
	}

	return results, nil
}

func (a *AIAgent) triggerSkillNudge(ctx context.Context) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		lastMessages := a.conversation
		if len(lastMessages) < 2 {
			return
		}

		logger.Info("Checking if recent work should be saved as a skill...")

		var recentContent strings.Builder
		recentContent.WriteString("# Recent Work\n\n")

		start := len(lastMessages) - 10
		if start < 0 {
			start = 0
		}

		for i := start; i < len(lastMessages); i++ {
			msg := lastMessages[i]
			recentContent.WriteString(fmt.Sprintf("## %s\n", msg.Role))
			recentContent.WriteString(msg.Content)
			recentContent.WriteString("\n\n")
		}

		skillName := fmt.Sprintf("workflow_%d", time.Now().Unix())
		skillDescription := "Automatically generated skill from recent workflow"

		skillContent := skill.GenerateSkillFile(
			skillName,
			skillDescription,
			recentContent.String(),
		)

		if err := a.SkillManager.Create(ctx, skillName, skillDescription, skillContent); err != nil {
			logger.Error("Failed to create skill: %v", err)
		} else {
			logger.Info("Skill '%s' created successfully", skillName)
		}
	}()
}

func (a *AIAgent) SaveConversation(sessionID string) error {
	msgs := make([]string, len(a.conversation))
	for i, msg := range a.conversation {
		msgBytes, _ := json.Marshal(msg)
		msgs[i] = string(msgBytes)
	}

	messagesJSON := "[" + strings.Join(msgs, ",") + "]"

	session := &storage.Session{
		ID:        sessionID,
		UserID:    "default",
		Messages:  messagesJSON,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return a.store.SaveSession(session)
}

func (a *AIAgent) LoadConversation(sessionID string) error {
	session, err := a.store.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		a.conversation = make([]model.Message, 0)
		return nil
	}

	var messages []model.Message
	if err := json.Unmarshal([]byte(session.Messages), &messages); err != nil {
		return fmt.Errorf("failed to unmarshal session messages: %w", err)
	}

	a.conversation = messages
	return nil
}

func (a *AIAgent) ResetConversation() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.conversation = make([]model.Message, 0)
}

func (a *AIAgent) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"tool_call_count":  a.toolCallCount,
		"skill_count":      a.SkillManager.GetSkillCount(),
		"memory_stats":     a.Memory.GetMemoryStats(),
		"conversation_len": len(a.conversation),
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
