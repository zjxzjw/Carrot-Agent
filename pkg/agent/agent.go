package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"carrotagent/carrot-agent/pkg/agent/memory"
	"carrotagent/carrot-agent/pkg/agent/model"
	"carrotagent/carrot-agent/pkg/agent/skill"
	"carrotagent/carrot-agent/pkg/agent/tool"
	"carrotagent/carrot-agent/pkg/storage"
)

type AIAgent struct {
	name          string
	version       string
	store         *storage.Store
	memory        *memory.MemoryManager
	skillManager  *skill.SkillManager
	modelProvider model.Provider
	toolRegistry  *tool.ToolRegistry
	config        *AgentConfig
	conversation  []model.Message
	skillNudgeInt int
	toolCallCount int
	mu            sync.RWMutex
}

type AgentConfig struct {
	Name           string
	Version        string
	DataDir        string
	ModelProvider  string
	ModelName      string
	Temperature    float64
	MaxTokens      int
	EnableSkills   bool
	EnableMemory    bool
	SkillNudgeInt  int
}

type AgentOption func(*AIAgent)

func WithModelProvider(provider model.Provider) AgentOption {
	return func(a *AIAgent) {
		a.modelProvider = provider
	}
}

func WithSkillManager(skillMgr *skill.SkillManager) AgentOption {
	return func(a *AIAgent) {
		a.skillManager = skillMgr
	}
}

func WithMemoryManager(memMgr *memory.MemoryManager) AgentOption {
	return func(a *AIAgent) {
		a.memory = memMgr
	}
}

func WithToolRegistry(registry *tool.ToolRegistry) AgentOption {
	return func(a *AIAgent) {
		a.toolRegistry = registry
	}
}

func NewAIAgent(cfg *AgentConfig, store *storage.Store, opts ...AgentOption) *AIAgent {
	agent := &AIAgent{
		name:          "carrot-agent",
		version:       "0.1.0",
		store:         store,
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
	if err := a.memory.Load(ctx); err != nil {
		return fmt.Errorf("failed to load memory: %w", err)
	}

	if err := a.skillManager.Load(ctx); err != nil {
		return fmt.Errorf("failed to load skills: %w", err)
	}

	a.registerDefaultTools()

	return nil
}

func (a *AIAgent) registerDefaultTools() {
	a.toolRegistry.Register("memory_read",
		"Read from memory storage",
		map[string]interface{}{
			"memory_id": map[string]interface{}{"type": "string", "description": "Memory ID to read", "required": true},
		},
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			id, _ := args["memory_id"].(string)
			mem, err := a.memory.Get(id)
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
		map[string]interface{}{
			"type":     map[string]interface{}{"type": "string", "description": "Memory type (snapshot, session, longterm)", "required": true},
			"content":  map[string]interface{}{"type": "string", "description": "Content to store", "required": true},
			"metadata": map[string]interface{}{"type": "string", "description": "Optional metadata as JSON", "required": false},
		},
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			memType, _ := args["type"].(string)
			content, _ := args["content"].(string)
			metadata, _ := args["metadata"].(string)
			if metadata == "" {
				metadata = "{}"
			}

			if err := a.memory.Add(ctx, memType, content, metadata); err != nil {
				return nil, err
			}
			return "Memory saved successfully", nil
		})

	a.toolRegistry.Register("skill_create",
		"Create a new skill",
		map[string]interface{}{
			"name":        map[string]interface{}{"type": "string", "description": "Skill name", "required": true},
			"description": map[string]interface{}{"type": "string", "description": "Skill description", "required": true},
			"content":      map[string]interface{}{"type": "string", "description": "Skill content in markdown format", "required": true},
		},
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			name, _ := args["name"].(string)
			description, _ := args["description"].(string)
			content, _ := args["content"].(string)

			if err := a.skillManager.Create(ctx, name, description, content); err != nil {
				return nil, err
			}
			return fmt.Sprintf("Skill '%s' created successfully", name), nil
		})

	a.toolRegistry.Register("skill_update",
		"Update an existing skill",
		map[string]interface{}{
			"skill_id": map[string]interface{}{"type": "string", "description": "Skill ID to update", "required": true},
			"content":   map[string]interface{}{"type": "string", "description": "New skill content", "required": true},
		},
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			skillID, _ := args["skill_id"].(string)
			content, _ := args["content"].(string)

			if err := a.skillManager.Update(ctx, skillID, content); err != nil {
				return nil, err
			}
			return "Skill updated successfully", nil
		})

	a.toolRegistry.Register("skill_list",
		"List all available skills",
		map[string]interface{}{},
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			skills := a.skillManager.List(100)
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
		map[string]interface{}{
			"keyword": map[string]interface{}{"type": "string", "description": "Search keyword", "required": true},
		},
		func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			keyword, _ := args["keyword"].(string)
			skills, err := a.skillManager.Search(keyword, 50)
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
}

func (a *AIAgent) RunConversation(ctx context.Context, userInput string) (*model.ChatResponse, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.conversation = append(a.conversation, model.Message{
		Role:    "user",
		Content: userInput,
	})

	systemPrompt := a.buildSystemPrompt()
	messages := append([]model.Message{{Role: "system", Content: systemPrompt}}, a.conversation...)

	req := &model.ChatRequest{
		Model:       a.config.ModelName,
		Messages:    messages,
		Temperature: a.config.Temperature,
		MaxTokens:   a.config.MaxTokens,
	}

	tools := tool.ConvertToModelTools(a.toolRegistry)
	if len(tools) > 0 {
		req.Tools = make([]model.ToolDefinition, len(tools))
		for i, t := range tools {
			if fn, ok := t["function"].(map[string]interface{}); ok {
				req.Tools[i] = model.ToolDefinition{
					Name:        getString(fn, "name"),
					Description: getString(fn, "description"),
				}
			}
		}
	}

	if a.modelProvider == nil {
		return nil, fmt.Errorf("no model provider configured")
	}

	resp, err := a.modelProvider.Chat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get model response: %w", err)
	}

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		a.conversation = append(a.conversation, choice.Message)
	}

	return resp, nil
}

func (a *AIAgent) buildSystemPrompt() string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("You are %s, version %s.\n\n", a.name, a.version))

	if a.config.EnableMemory {
		snapshot := a.memory.GetSnapshotContent()
		if snapshot != "" {
			prompt.WriteString("## Memory\n")
			prompt.WriteString(snapshot)
			prompt.WriteString("\n\n")
		}
	}

	if a.config.EnableSkills {
		skillsIndex := a.skillManager.GetSkillsIndex()
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
		parsed, err := tool.ParseToolCall(tc)
		if err != nil {
			results = append(results, map[string]interface{}{
				"tool_call_id": tc["id"],
				"output":       fmt.Sprintf("Error parsing tool call: %v", err),
			})
			continue
		}

		result := a.toolRegistry.Execute(ctx, parsed.Name, parsed.Arguments)

		output := ""
		if result.Success {
			if outputBytes, err := json.Marshal(result.Output); err == nil {
				output = string(outputBytes)
			} else {
				output = fmt.Sprintf("%v", result.Output)
			}
		} else {
			output = result.Error
		}

		results = append(results, map[string]interface{}{
			"tool_call_id": tc["id"],
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
		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		lastMessages := a.conversation
		if len(lastMessages) < 2 {
			return
		}

		fmt.Println("[Background] Checking if recent work should be saved as a skill...")
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
		return err
	}
	if session == nil {
		return nil
	}

	a.conversation = make([]model.Message, 0)

	return nil
}

func (a *AIAgent) ResetConversation() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.conversation = make([]model.Message, 0)
}

func (a *AIAgent) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"tool_call_count":   a.toolCallCount,
		"skill_count":      a.skillManager.GetSkillCount(),
		"memory_stats":     a.memory.GetMemoryStats(),
		"conversation_len": len(a.conversation),
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}