package agent

import (
	"context"
	"testing"
	"time"

	"carrotagent/carrot-agent/pkg/storage"
)

func TestNewAIAgent(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cfg := &AgentConfig{
		Name:          "test-agent",
		Version:       "0.1.0",
		DataDir:       ".",
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		Temperature:   0.7,
		MaxTokens:     4096,
		EnableSkills:  true,
		EnableMemory:  true,
		SkillNudgeInt: 10,
	}

	agent := NewAIAgent(cfg, store)

	if agent == nil {
		t.Fatal("Failed to create agent")
	}
}

func TestAgentInitialize(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cfg := &AgentConfig{
		Name:          "test-agent",
		Version:       "0.1.0",
		DataDir:       ".",
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		Temperature:   0.7,
		MaxTokens:     4096,
		EnableSkills:  true,
		EnableMemory:  true,
		SkillNudgeInt: 10,
	}

	agent := NewAIAgent(cfg, store)

	err = agent.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize agent: %v", err)
	}
}

func TestAgentGetStats(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cfg := &AgentConfig{
		Name:          "test-agent",
		Version:       "0.1.0",
		DataDir:       ".",
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		Temperature:   0.7,
		MaxTokens:     4096,
		EnableSkills:  true,
		EnableMemory:  true,
		SkillNudgeInt: 10,
	}

	agent := NewAIAgent(cfg, store)

	err = agent.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize agent: %v", err)
	}

	stats := agent.GetStats()

	if stats["tool_call_count"] == nil {
		t.Error("Expected tool_call_count in stats")
	}

	if stats["skill_count"] == nil {
		t.Error("Expected skill_count in stats")
	}

	if stats["memory_stats"] == nil {
		t.Error("Expected memory_stats in stats")
	}

	if stats["conversation_len"] == nil {
		t.Error("Expected conversation_len in stats")
	}
}

func TestSaveAndLoadConversation(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cfg := &AgentConfig{
		Name:          "test-agent",
		Version:       "0.1.0",
		DataDir:       ".",
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		Temperature:   0.7,
		MaxTokens:     4096,
		EnableSkills:  true,
		EnableMemory:  true,
		SkillNudgeInt: 10,
	}

	agent := NewAIAgent(cfg, store)

	err = agent.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize agent: %v", err)
	}

	sessionID := "test_session_1"
	err = agent.SaveConversation(sessionID)
	if err != nil {
		t.Fatalf("Failed to save conversation: %v", err)
	}

	err = agent.LoadConversation(sessionID)
	if err != nil {
		t.Fatalf("Failed to load conversation: %v", err)
	}

	if len(agent.conversation) != 0 {
		t.Errorf("Expected empty conversation after load, got %d messages", len(agent.conversation))
	}
}

func TestLoadConversationNonExistent(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cfg := &AgentConfig{
		Name:          "test-agent",
		Version:       "0.1.0",
		DataDir:       ".",
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		Temperature:   0.7,
		MaxTokens:     4096,
		EnableSkills:  true,
		EnableMemory:  true,
		SkillNudgeInt: 10,
	}

	agent := NewAIAgent(cfg, store)

	err = agent.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize agent: %v", err)
	}

	err = agent.LoadConversation("non_existent_session")
	if err != nil {
		t.Fatalf("Failed to load non-existent session: %v", err)
	}

	if len(agent.conversation) != 0 {
		t.Errorf("Expected empty conversation for non-existent session, got %d messages", len(agent.conversation))
	}
}

func TestResetConversation(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	cfg := &AgentConfig{
		Name:          "test-agent",
		Version:       "0.1.0",
		DataDir:       ".",
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		Temperature:   0.7,
		MaxTokens:     4096,
		EnableSkills:  true,
		EnableMemory:  true,
		SkillNudgeInt: 10,
	}

	agent := NewAIAgent(cfg, store)

	err = agent.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize agent: %v", err)
	}

	agent.ResetConversation()

	if len(agent.conversation) != 0 {
		t.Errorf("Expected empty conversation after reset, got %d messages", len(agent.conversation))
	}
}

func TestSaveSession(t *testing.T) {
	store, err := storage.NewStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	session := &storage.Session{
		ID:        "test_session",
		UserID:    "test_user",
		Messages:  "[]",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = store.SaveSession(session)
	if err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	retrieved, err := store.GetSession("test_session")
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected to retrieve session, got nil")
	}

	if retrieved.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, retrieved.ID)
	}
}
