package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"carrotagent/carrot-agent/config"
	"carrotagent/carrot-agent/pkg/agent"
	"carrotagent/carrot-agent/pkg/agent/memory"
	"carrotagent/carrot-agent/pkg/agent/model"
	"carrotagent/carrot-agent/pkg/agent/skill"
	"carrotagent/carrot-agent/pkg/storage"
)

var (
	aiAgent *agent.AIAgent
	store   *storage.Store
	appConfig *config.Config
)

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
				applyEnvOverrides(cfg)
				return cfg
			}
		}
	}

	cfg := config.Default()
	applyEnvOverrides(cfg)
	return cfg
}

func applyEnvOverrides(cfg *config.Config) {
	if apiKey := os.Getenv("CARROT_API_KEY"); apiKey != "" {
		cfg.Model.APIKey = apiKey
	}
	if modelName := os.Getenv("CARROT_MODEL_NAME"); modelName != "" {
		cfg.Model.ModelName = modelName
	}
	if baseURL := os.Getenv("CARROT_BASE_URL"); baseURL != "" {
		cfg.Model.BaseURL = baseURL
	}
	if provider := os.Getenv("CARROT_MODEL_PROVIDER"); provider != "" {
		cfg.Model.Provider = provider
	}
}

func initStorage(cfg *config.Config) (*storage.Store, error) {
	dbPath := os.ExpandEnv(cfg.Storage.DBPath)
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return storage.NewStore(dbPath)
}

func initAgent(cfg *config.Config, store *storage.Store) *agent.AIAgent {
	skillMgr := skill.NewSkillManager(store)
	memMgr := memory.NewMemoryManager(store)

	agentCfg := &agent.AgentConfig{
		Name:          cfg.Agent.Name,
		Version:       cfg.Agent.Version,
		DataDir:       cfg.Agent.DataDir,
		ModelProvider:  cfg.Model.Provider,
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

func main() {
	fmt.Printf("Carrot Agent API Server v%s\n", "0.1.0")

	cfg := loadConfig()
	appConfig = cfg

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	store, err = initStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	agentInstance := initAgent(cfg, store)
	if err := agentInstance.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	aiAgent = agentInstance

	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/api/chat", handleChat)
	http.HandleFunc("/api/skills", handleSkills)
	http.HandleFunc("/api/memory", handleMemory)
	http.HandleFunc("/api/stats", handleStats)
	http.HandleFunc("/api/session/", handleSession)
	http.HandleFunc("/api/config", handleConfig)
	http.HandleFunc("/api/models", handleModels)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: nil,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %d", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleSession(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	sessionID := strings.TrimPrefix(path, "/api/session/")

	if sessionID == "" {
		switch r.Method {
		case http.MethodGet:
			handleListSessions(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetSession(w, r, sessionID)
	case http.MethodDelete:
		handleDeleteSession(w, r, sessionID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleListSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := store.ListSessions("default", 100)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list sessions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"sessions": sessions, "count": len(sessions)})
}

func handleGetSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	session, err := store.GetSession(sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get session: %v", err), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func handleDeleteSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	if err := store.DeleteSession(sessionID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete session: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Session deleted successfully"})
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Message string `json:"message"`
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	resp, err := aiAgent.RunConversation(r.Context(), req.Message)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get response: %v", err), http.StatusInternalServerError)
		return
	}

	if len(resp.Choices) > 0 {
		response := struct {
			Message string `json:"message"`
			Usage model.Usage `json:"usage"`
		}{
			Message: resp.Choices[0].Message.Content,
			Usage: resp.Usage,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "No response from model", http.StatusInternalServerError)
	}
}

func handleSkills(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleListSkills(w, r)
	case http.MethodPost:
		handleCreateSkill(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleListSkills(w http.ResponseWriter, r *http.Request) {
	skills := aiAgent.SkillManager.List(100)

	response := struct {
		Skills []*storage.Skill `json:"skills"`
		Count int `json:"count"`
	}{
		Skills: skills,
		Count: len(skills),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleCreateSkill(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := aiAgent.SkillManager.Create(r.Context(), req.Name, req.Description, req.Content); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create skill: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Skill created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handleMemory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleListMemory(w, r)
	case http.MethodPost:
		handleAddMemory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleListMemory(w http.ResponseWriter, r *http.Request) {
	memType := r.URL.Query().Get("type")
	memories, err := aiAgent.Memory.List(memType, 100)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list memories: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Memories []*storage.Memory `json:"memories"`
		Count int `json:"count"`
	}{
		Memories: memories,
		Count: len(memories),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleAddMemory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type string `json:"type"`
		Content string `json:"content"`
		Metadata string `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := aiAgent.Memory.Add(r.Context(), req.Type, req.Content, req.Metadata); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add memory: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Memory added successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := aiAgent.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetConfig(w, r)
	case http.MethodPut:
		handleUpdateConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appConfig)
}

func handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Model struct {
			Provider    string  `json:"provider"`
			APIKey      string  `json:"api_key"`
			ModelName   string  `json:"model_name"`
			BaseURL     string  `json:"base_url"`
			Temperature float64 `json:"temperature"`
			MaxTokens   int     `json:"max_tokens"`
		} `json:"model"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update config
	if req.Model.Provider != "" {
		appConfig.Model.Provider = req.Model.Provider
	}
	if req.Model.APIKey != "" {
		appConfig.Model.APIKey = req.Model.APIKey
	}
	if req.Model.ModelName != "" {
		appConfig.Model.ModelName = req.Model.ModelName
	}
	if req.Model.BaseURL != "" {
		appConfig.Model.BaseURL = req.Model.BaseURL
	}
	if req.Model.Temperature > 0 {
		appConfig.Model.Temperature = req.Model.Temperature
	}
	if req.Model.MaxTokens > 0 {
		appConfig.Model.MaxTokens = req.Model.MaxTokens
	}

	// Validate config
	if errs := appConfig.Validate(); len(errs) > 0 {
		http.Error(w, fmt.Sprintf("Configuration validation failed: %v", errs), http.StatusBadRequest)
		return
	}

	// Save config to file
	configPath := os.ExpandEnv("~/.carrot/config.yaml")
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create config directory: %v", err), http.StatusInternalServerError)
		return
	}

	if err := appConfig.Save(configPath); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save config: %v", err), http.StatusInternalServerError)
		return
	}

	// Restart agent with new config
	ctx := context.Background()
	agentInstance := initAgent(appConfig, store)
	if err := agentInstance.Initialize(ctx); err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize agent: %v", err), http.StatusInternalServerError)
		return
	}

	aiAgent = agentInstance

	response := struct {
		Message string `json:"message"`
		Config *config.Config `json:"config"`
	}{
		Message: "Configuration updated successfully",
		Config: appConfig,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	models := []struct {
		Provider string `json:"provider"`
		Models []string `json:"models"`
	}{
		{
			Provider: "openai",
			Models: []string{
				"gpt-4o",
				"gpt-4-turbo",
				"gpt-4",
				"gpt-3.5-turbo",
				"gpt-3.5-turbo-16k",
			},
		},
		{
			Provider: "claude",
			Models: []string{
				"claude-3-opus-20240229",
				"claude-3-sonnet-20240229",
				"claude-3-haiku-20240307",
				"claude-2.1",
				"claude-2",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}