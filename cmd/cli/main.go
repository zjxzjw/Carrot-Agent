package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"carrotagent/carrot-agent/config"
	"carrotagent/carrot-agent/pkg/agent"
	"carrotagent/carrot-agent/pkg/agent/memory"
	"carrotagent/carrot-agent/pkg/agent/model"
	"carrotagent/carrot-agent/pkg/agent/skill"
	"carrotagent/carrot-agent/pkg/storage"
)

var (
	Version   = "0.1.0"
	BuildTime = "unknown"
)

func main() {
	fmt.Printf("Carrot Agent v%s (build: %s)\n", Version, BuildTime)
	fmt.Println("Type 'help' for available commands, 'quit' to exit.")

	cfg := loadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store, err := initStorage(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	agentInstance := initAgent(cfg, store)
	if err := agentInstance.Initialize(ctx); err != nil {
		fmt.Printf("Failed to initialize agent: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Agent initialized successfully!")
	fmt.Printf("Model: %s (%s) | Skills: %v | Memory: %v\n\n",
		cfg.Model.Provider, cfg.Model.ModelName, cfg.Agent.SkillNudgeInt > 0, cfg.Agent.DataDir != "")

	scanner := bufio.NewScanner(os.Stdin)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	interactiveLoop(ctx, agentInstance, scanner)
}

func loadConfig() *config.Config {
	configPaths := []string{
		"~/.carrot/config.yaml",
		"./config.yaml",
		"/etc/carrot/config.yaml",
	}

	for _, path := range configPaths {
		expandedPath := os.ExpandEnv(path)
		if _, err := os.Stat(expandedPath); err == nil {
			cfg, err := config.Load(expandedPath)
			if err == nil {
				fmt.Printf("Loaded config from: %s\n", expandedPath)
				return cfg
			}
		}
	}

	fmt.Println("Using default config (no config file found)")
	return config.Default()
}

func initStorage(cfg *config.Config) (*storage.Store, error) {
	dbPath := os.ExpandEnv(cfg.Storage.DBPath)
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	store, err := storage.NewStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	return store, nil
}

func initAgent(cfg *config.Config, store *storage.Store) *agent.AIAgent {
	skillMgr := skill.NewSkillManager(store)
	memMgr := memory.NewMemoryManager(store)

	agentCfg := &agent.AgentConfig{
		Name:          cfg.Agent.Name,
		Version:       cfg.Agent.Version,
		DataDir:       cfg.Agent.DataDir,
		ModelProvider: cfg.Model.Provider,
		ModelName:     cfg.Model.ModelName,
		Temperature:   cfg.Model.Temperature,
		MaxTokens:     cfg.Model.MaxTokens,
		EnableSkills:  true,
		EnableMemory:  true,
		SkillNudgeInt: cfg.Agent.SkillNudgeInt,
	}

	var modelProvider model.Provider
	if cfg.Model.APIKey != "" {
		modelProvider = model.NewProviderFactory().CreateProvider(
			cfg.Model.Provider,
			cfg.Model.APIKey,
			cfg.Model.BaseURL,
			cfg.Model.ModelName,
		)
	}

	agentOpts := []agent.AgentOption{
		agent.WithSkillManager(skillMgr),
		agent.WithMemoryManager(memMgr),
	}

	if modelProvider != nil {
		agentOpts = append(agentOpts, agent.WithModelProvider(modelProvider))
	}

	return agent.NewAIAgent(agentCfg, store, agentOpts...)
}

func interactiveLoop(ctx context.Context, ag *agent.AIAgent, scanner *bufio.Scanner) {
	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "quit" || input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "help" {
			printHelp()
			continue
		}

		if input == "stats" {
			printStats(ag)
			continue
		}

		if input == "reset" {
			ag.ResetConversation()
			fmt.Println("Conversation reset.")
			continue
		}

		if input == "skills" {
			stats := ag.GetStats()
			fmt.Printf("Skills count: %v\n", stats["skill_count"])
			continue
		}

		if err := handleConversationTurn(ctx, ag, input); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func handleConversationTurn(ctx context.Context, ag *agent.AIAgent, input string) error {
	maxIterations := 10
	currentInput := input

	for i := 0; i < maxIterations; i++ {
		resp, err := ag.RunConversation(ctx, currentInput)
		if err != nil {
			return err
		}

		if resp == nil || len(resp.Choices) == 0 {
			return fmt.Errorf("no response from model")
		}

		choice := resp.Choices[0]

		if choice.FinishReason == "tool_calls" && len(choice.Message.ToolCalls) > 0 {
			fmt.Println("Tool calls detected, processing...")
			toolCallMaps := make([]map[string]interface{}, 0, len(choice.Message.ToolCalls))
			for _, tc := range choice.Message.ToolCalls {
				toolCallMaps = append(toolCallMaps, map[string]interface{}{
					"id":   tc.ID,
					"type": "function",
					"function": map[string]interface{}{
						"name":      tc.Function.Name,
						"arguments": tc.Function.Arguments,
					},
				})
			}

			toolMessages, err := ag.ProcessToolCalls(ctx, toolCallMaps)
			if err != nil {
				return fmt.Errorf("failed to process tool calls: %w", err)
			}

			// Add all tool response messages to conversation
			for i, toolMsg := range toolMessages {
				ag.AddMessage(toolMsg)
				fmt.Printf("Tool %d result: %s\n", i+1, toolMsg.Content)
			}

			currentInput = "Continue"
			continue
		}

		fmt.Printf("Response: %s\n", choice.Message.Content)
		return nil
	}

	return fmt.Errorf("max iterations reached during tool call processing")
}

func printHelp() {
	fmt.Print(`
Available Commands:
  help       - Show this help message
  quit/exit  - Exit the agent
  reset      - Reset conversation history
  stats      - Show agent statistics
  skills     - List available skills

Examples:
  > Hello, how are you?
  > skills
`)
}

func printStats(ag *agent.AIAgent) {
	stats := ag.GetStats()
	fmt.Println("Agent Statistics:")
	for k, v := range stats {
		fmt.Printf("  %s: %v\n", k, v)
	}
}
