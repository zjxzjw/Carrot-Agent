package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"carrotagent/carrot-agent/config"
	"carrotagent/carrot-agent/pkg/agent"
	"carrotagent/carrot-agent/pkg/agent/memory"
	"carrotagent/carrot-agent/pkg/agent/model"
	"carrotagent/carrot-agent/pkg/agent/skill"
	"carrotagent/carrot-agent/pkg/storage"
)

var aiAgent *agent.AIAgent

func main() {
	// 加载配置
	cfg := config.Default()
	// 从环境变量覆盖配置
	if apiKey := os.Getenv("CARROT_API_KEY"); apiKey != "" {
		cfg.Model.APIKey = apiKey
	}
	if modelName := os.Getenv("CARROT_MODEL_NAME"); modelName != "" {
		cfg.Model.ModelName = modelName
	}
	if baseURL := os.Getenv("CARROT_BASE_URL"); baseURL != "" {
		cfg.Model.BaseURL = baseURL
	}

	// 初始化存储
	store, err := storage.NewStore(cfg.Agent.DataDir)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// 初始化内存管理器
	memoryManager := memory.NewMemoryManager(store)

	// 初始化技能管理器
	skillManager := skill.NewSkillManager(store)

	// 初始化模型提供者
	provider := model.NewOpenAIProvider(
		cfg.Model.APIKey,
		cfg.Model.BaseURL,
		cfg.Model.ModelName,
	)

	// 初始化代理
	aiAgent = agent.NewAIAgent(
		&agent.AgentConfig{
			Name:           "carrot-agent",
			Version:        "0.1.0",
			DataDir:        cfg.Agent.DataDir,
			ModelProvider:  cfg.Model.Provider,
			ModelName:      cfg.Model.ModelName,
			Temperature:    cfg.Model.Temperature,
			MaxTokens:      cfg.Model.MaxTokens,
			EnableSkills:   true,
			EnableMemory:   true,
			SkillNudgeInt:  cfg.Agent.SkillNudgeInt,
		},
		store,
		agent.WithMemoryManager(memoryManager),
		agent.WithSkillManager(skillManager),
		agent.WithModelProvider(provider),
	)

	// 初始化代理
	if err := aiAgent.Initialize(context.Background()); err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	// 设置路由
	http.HandleFunc("/api/chat", handleChat)
	http.HandleFunc("/api/skills", handleSkills)
	http.HandleFunc("/api/memory", handleMemory)
	http.HandleFunc("/api/stats", handleStats)

	// 启动服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: nil,
	}

	// 优雅关闭
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited")
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