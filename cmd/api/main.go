package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"carrotagent/carrot-agent/config"
	"carrotagent/carrot-agent/pkg/agent"
	"carrotagent/carrot-agent/pkg/agent/memory"
	"carrotagent/carrot-agent/pkg/agent/model"
	"carrotagent/carrot-agent/pkg/agent/skill"
	"carrotagent/carrot-agent/pkg/storage"

	"github.com/joho/godotenv"
)

var (
	aiAgent        *agent.AIAgent
	store          *storage.Store
	appConfig      *config.Config
	sessions       = make(map[string]*sessionInfo) // 存储会话信息
	sessionMu      sync.RWMutex                    // 保护sessions map的并发访问
	sessionTimeout = 24 * time.Hour                // 会话超时时间
)

type sessionInfo struct {
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type rateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*clientRate
	maxReqs int
	window  time.Duration
}

type clientRate struct {
	count   int
	resetAt time.Time
}

var globalRateLimiter *rateLimiter

func init() {
	cwd, err := os.Getwd()
	if err == nil {
		envPath := filepath.Join(cwd, ".env")
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: .env file not found at %s: %v", envPath, err)
		} else {
			log.Println(".env file loaded successfully from:", envPath)
		}
	}
	if username := os.Getenv("CARROT_AUTH_USERNAME"); username == "" {
		log.Println("Warning: CARROT_AUTH_USERNAME is not set")
	}
	if password := os.Getenv("CARROT_AUTH_PASSWORD"); password == "" {
		log.Println("Warning: CARROT_AUTH_PASSWORD is not set")
	}
}

func newRateLimiter(maxReqs int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		clients: make(map[string]*clientRate),
		maxReqs: maxReqs,
		window:  window,
	}
}

func (rl *rateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[ip]

	if !exists || now.After(client.resetAt) {
		rl.clients[ip] = &clientRate{
			count:   1,
			resetAt: now.Add(rl.window),
		}
		return true
	}

	if client.count >= rl.maxReqs {
		return false
	}

	client.count++
	return true
}

func (rl *rateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, client := range rl.clients {
		if now.After(client.resetAt) {
			delete(rl.clients, ip)
		}
	}
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
	// Auth configuration from environment
	if username := os.Getenv("CARROT_AUTH_USERNAME"); username != "" {
		cfg.Auth.Username = username
	}
	if password := os.Getenv("CARROT_AUTH_PASSWORD"); password != "" {
		cfg.Auth.Password = password
	}
	// Server configuration from environment
	if port := os.Getenv("CARROT_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil && p > 0 {
			cfg.Server.Port = p
		}
	}
	if host := os.Getenv("CARROT_SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func cleanupExpiredSessions() {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	now := time.Now()
	for sessionID, info := range sessions {
		if now.After(info.ExpiresAt) {
			delete(sessions, sessionID)
		}
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 清理过期会话
		cleanupExpiredSessions()

		// 跳过登录和健康检查路径的鉴权
		if r.URL.Path == "/api/login" || r.URL.Path == "/health" {
			next(w, r)
			return
		}

		// 从请求头获取会话ID
		sessionID := r.Header.Get("Authorization")
		if sessionID == "" {
			http.Error(w, `{"error":"Unauthorized: No session ID provided"}`, http.StatusUnauthorized)
			return
		}

		// 验证会话是否有效
		sessionMu.RLock()
		info, ok := sessions[sessionID]
		sessionMu.RUnlock()

		if !ok || time.Now().After(info.ExpiresAt) {
			http.Error(w, `{"error":"Unauthorized: Invalid or expired session"}`, http.StatusUnauthorized)
			return
		}

		// 延长会话过期时间
		sessionMu.Lock()
		if info, exists := sessions[sessionID]; exists {
			info.ExpiresAt = time.Now().Add(sessionTimeout)
		}
		sessionMu.Unlock()

		next(w, r)
	}
}

func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取客户端IP
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = strings.Split(forwarded, ",")[0]
		}

		if !globalRateLimiter.Allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		next(w, r)
	}
}

func securityHeadersMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置安全头
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("Pragma", "no-cache")

		next(w, r)
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

func main() {
	fmt.Printf("Carrot Agent API Server v%s\n", "0.1.0")

	cfg := loadConfig()
	appConfig = cfg

	// 初始化速率限制器（每分钟最多100个请求）
	globalRateLimiter = newRateLimiter(100, time.Minute)

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

	http.HandleFunc("/health", securityHeadersMiddleware(handleHealth))
	http.HandleFunc("/api/login", securityHeadersMiddleware(rateLimitMiddleware(handleLogin)))
	http.HandleFunc("/api/chat", authMiddleware(securityHeadersMiddleware(rateLimitMiddleware(handleChat))))
	http.HandleFunc("/api/chat/stream", authMiddleware(securityHeadersMiddleware(rateLimitMiddleware(handleChatStream))))
	http.HandleFunc("/api/skills", authMiddleware(securityHeadersMiddleware(rateLimitMiddleware(handleSkills))))
	http.HandleFunc("/api/memory", authMiddleware(securityHeadersMiddleware(rateLimitMiddleware(handleMemory))))
	http.HandleFunc("/api/stats", authMiddleware(securityHeadersMiddleware(rateLimitMiddleware(handleStats))))
	http.HandleFunc("/api/session/", authMiddleware(securityHeadersMiddleware(rateLimitMiddleware(handleSession))))
	http.HandleFunc("/api/config", authMiddleware(securityHeadersMiddleware(rateLimitMiddleware(handleConfig))))
	http.HandleFunc("/api/models", authMiddleware(securityHeadersMiddleware(rateLimitMiddleware(handleModels))))
	http.HandleFunc("/api/tools", authMiddleware(securityHeadersMiddleware(rateLimitMiddleware(handleTools))))

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        nil,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 启动后台清理任务
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				globalRateLimiter.Cleanup()
			case <-ctx.Done():
				return
			}
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

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	// 验证用户名和密码
	if appConfig.Auth.Username == "" || appConfig.Auth.Password == "" {
		http.Error(w, `{"error":"Authentication not configured. Please set username and password in config."}`, http.StatusServiceUnavailable)
		return
	}

	if req.Username != appConfig.Auth.Username || req.Password != appConfig.Auth.Password {
		http.Error(w, `{"error":"Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	// 生成安全的会话ID
	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, `{"error":"Failed to generate session"}`, http.StatusInternalServerError)
		return
	}

	now := time.Now()
	sessionMu.Lock()
	sessions[sessionID] = &sessionInfo{
		UserID:    req.Username,
		CreatedAt: now,
		ExpiresAt: now.Add(sessionTimeout),
	}
	sessionMu.Unlock()

	response := map[string]interface{}{
		"session_id": sessionID,
		"message":    "Login successful",
		"expires_in": int(sessionTimeout.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Message   string `json:"message"`
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// 输入验证
	if req.Message == "" {
		http.Error(w, `{"error":"Message is required"}`, http.StatusBadRequest)
		return
	}

	if len(req.Message) > 10000 {
		http.Error(w, `{"error":"Message too long (max 10000 characters)"}`, http.StatusBadRequest)
		return
	}

	// 第一次调用模型
	resp, err := aiAgent.RunConversation(r.Context(), req.Message)
	if err != nil {
		log.Printf("Error in chat: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Failed to get response: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// 检查是否需要执行工具调用
	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		log.Printf("Received model response with finish reason: %s", choice.FinishReason)
		// DeepSeek puts tool_calls inside message
		log.Printf("Number of tool calls: %d", len(choice.Message.ToolCalls))

		// 如果模型决定调用工具
		if choice.FinishReason == "tool_calls" {
			if len(choice.Message.ToolCalls) > 0 {
				log.Printf("Processing tool calls...")
				// 提取工具调用信息
				toolCalls := make([]map[string]interface{}, 0, len(choice.Message.ToolCalls))
				for _, toolCall := range choice.Message.ToolCalls {
					log.Printf("Tool call: %s, arguments: %s", toolCall.Function.Name, toolCall.Function.Arguments)
					toolCallMap := map[string]interface{}{
						"id":   toolCall.ID,
						"type": "function",
						"function": map[string]interface{}{
							"name":      toolCall.Function.Name,
							"arguments": toolCall.Function.Arguments,
						},
					}
					toolCalls = append(toolCalls, toolCallMap)
				}

				// 执行工具调用
				log.Printf("Executing tool calls...")
				for _, tc := range toolCalls {
					toolCallID, _ := tc["id"].(string)
					toolName, _ := tc["function"].(map[string]interface{})["name"].(string)
					argsStr, _ := tc["function"].(map[string]interface{})["arguments"].(string)
					var args map[string]interface{}
					if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
						log.Printf("Error parsing tool arguments: %v", err)
						// Add error as tool response instead of returning
						toolMessage := model.Message{
							Role:       "tool",
							Content:    fmt.Sprintf("Error parsing arguments: %v", err),
							Name:       toolName,
							ToolCallID: toolCallID,
							Type:       "function",
						}
						aiAgent.AddMessage(toolMessage)
						continue
					}
					
					toolResult := aiAgent.GetToolRegistry().Execute(r.Context(), toolName, args)
					
					// Always add tool response message, even if execution failed
					output := ""
					if toolResult.Success {
						if outputBytes, err := json.Marshal(toolResult.Output); err == nil {
							output = string(outputBytes)
						} else {
							output = fmt.Sprintf("%v", toolResult.Output)
						}
						log.Printf("Tool %s executed successfully", toolName)
					} else {
						output = toolResult.Error
						log.Printf("Tool %s execution failed: %s", toolName, toolResult.Error)
					}
					
					// Add tool response to conversation history
					toolMessage := model.Message{
						Role:       "tool",
						Content:    output,
						Name:       toolName,
						ToolCallID: toolCallID,
						Type:       "function",
					}
					aiAgent.AddMessage(toolMessage)
					log.Printf("Added tool message for %s (success=%v)", toolName, toolResult.Success)
				}

				// 再次调用模型，获取最终响应
				log.Printf("Getting final response from model...")
				finalResp, err := aiAgent.RunConversation(r.Context(), "")
				if err != nil {
					log.Printf("Error getting final response: %v", err)
					http.Error(w, fmt.Sprintf(`{"error":"Failed to get final response: %s"}`, err.Error()), http.StatusInternalServerError)
					return
				}
				log.Printf("Final response received")

				// 返回最终响应
				if len(finalResp.Choices) > 0 {
					log.Printf("Final response content: %s", finalResp.Choices[0].Message.Content)
					response := map[string]interface{}{
						"message": finalResp.Choices[0].Message.Content,
						"usage":   finalResp.Usage,
					}

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				} else {
					http.Error(w, `{"error":"No final response from model"}`, http.StatusInternalServerError)
				}
			} else if choice.Message.Content != "" {
				// 模型返回了 tool_calls 但没有实际工具调用，但有内容
				// 将内容作为普通响应返回
				log.Printf("Model returned tool_calls but no actual tool calls, returning content as normal response")
				response := map[string]interface{}{
					"message": choice.Message.Content,
					"usage":   resp.Usage,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else {
				// 处理边缘情况：模型返回了tool_calls但没有实际工具调用，也没有内容
				// 这可能是模型的异常行为，我们尝试重新请求或者返回友好提示
				log.Printf("WARNING: Model returned finish_reason=tool_calls but no tool calls and no content")
				log.Printf("This might be a model compatibility issue. Attempting to get a direct response...")

				// 尝试再次调用模型，但不带工具定义，强制模型直接回复
				log.Printf("Retrying without tools to get a direct response...")
				retryResp, retryErr := aiAgent.RunConversation(r.Context(), "")
				if retryErr != nil {
					log.Printf("Retry also failed: %v", retryErr)
					http.Error(w, `{"error":"The model encountered an issue processing your request. Please try rephrasing your question."}`, http.StatusInternalServerError)
					return
				}

				if len(retryResp.Choices) > 0 && retryResp.Choices[0].Message.Content != "" {
					log.Printf("Retry successful, content: %s", retryResp.Choices[0].Message.Content)
					response := map[string]interface{}{
						"message": retryResp.Choices[0].Message.Content,
						"usage":   retryResp.Usage,
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				} else {
					log.Printf("Retry also returned empty response")
					http.Error(w, `{"error":"The model could not process your request. Please try again with a different question."}`, http.StatusInternalServerError)
				}
			}
		} else {
			// 直接返回模型响应
			log.Printf("Direct response: %s", choice.Message.Content)
			response := map[string]interface{}{
				"message": choice.Message.Content,
				"usage":   resp.Usage,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	} else {
		http.Error(w, `{"error":"No response from model"}`, http.StatusInternalServerError)
	}
}

func handleChatStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Message   string `json:"message"`
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, `{"error":"Message is required"}`, http.StatusBadRequest)
		return
	}

	if len(req.Message) > 10000 {
		http.Error(w, `{"error":"Message too long (max 10000 characters)"}`, http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, `{"error":"Streaming not supported"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	aiAgent.AddMessage(model.Message{
		Role:    "user",
		Content: req.Message,
	})

	modelProvider := aiAgent.GetModelProvider()
	if modelProvider == nil {
		sendSSEvent(w, flusher, "error", map[string]string{"error": "Model provider not configured"})
		return
	}

	agentConfig := aiAgent.GetConfig()
	toolDefs := aiAgent.GetToolRegistry().ListEnabled()

	tools := make([]model.ToolDefinition, 0, len(toolDefs))
	for _, def := range toolDefs {
		properties := make(map[string]model.Property)
		required := []string{}

		for paramName, paramDef := range def.Parameters {
			paramMap, ok := paramDef.(map[string]interface{})
			if !ok {
				continue
			}
			paramType, _ := paramMap["type"].(string)
			paramDesc, _ := paramMap["description"].(string)
			properties[paramName] = model.Property{
				Type:        paramType,
				Description: paramDesc,
			}
			if requiredFlag, ok := paramMap["required"].(bool); ok && requiredFlag {
				required = append(required, paramName)
			}
		}

		tools = append(tools, model.ToolDefinition{
			Type: "function",
			Function: model.ToolFunctionDef{
				Name:        def.Name,
				Description: def.Description,
				Parameters: model.ToolParameters{
					Type:       "object",
					Properties: properties,
					Required:   required,
				},
			},
		})
	}

	chatReq := &model.ChatRequest{
		Model:       agentConfig.ModelName,
		Messages:    aiAgent.GetMessages(),
		Temperature: agentConfig.Temperature,
		MaxTokens:   agentConfig.MaxTokens,
		Tools:       tools,
		Stream:      true,
	}

	accumulatedContent := ""

	err := modelProvider.ChatStream(r.Context(), chatReq, func(chunk *model.ChatResponse) error {
		if len(chunk.Choices) > 0 {
			choice := chunk.Choices[0]
			accumulatedContent += choice.Message.Content

			event := map[string]interface{}{
				"content": choice.Message.Content,
				"done":    chunk.Done,
			}

			if chunk.Done {
				event["usage"] = chunk.Usage
			}

			if err := sendSSEvent(w, flusher, "chunk", event); err != nil {
				return err
			}

			if choice.FinishReason == "stop" {
				return nil
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Stream error: %v", err)
		sendSSEvent(w, flusher, "error", map[string]string{"error": err.Error()})
		return
	}

	aiAgent.AddMessage(model.Message{
		Role:    "assistant",
		Content: accumulatedContent,
	})
}

func sendSSEvent(w http.ResponseWriter, flusher http.Flusher, event string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "data: %s\n\nevent: %s\n\n", jsonData, event)
	if err != nil {
		return err
	}
	flusher.Flush()
	return nil
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
		Count  int              `json:"count"`
	}{
		Skills: skills,
		Count:  len(skills),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleCreateSkill(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// 输入验证
	if req.Name == "" || req.Description == "" || req.Content == "" {
		http.Error(w, `{"error":"Name, description and content are required"}`, http.StatusBadRequest)
		return
	}

	if len(req.Name) > 100 {
		http.Error(w, `{"error":"Name too long (max 100 characters)"}`, http.StatusBadRequest)
		return
	}

	if len(req.Description) > 500 {
		http.Error(w, `{"error":"Description too long (max 500 characters)"}`, http.StatusBadRequest)
		return
	}

	if err := aiAgent.SkillManager.Create(r.Context(), req.Name, req.Description, req.Content); err != nil {
		log.Printf("Error creating skill: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Failed to create skill: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Skill created successfully",
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
		Count    int               `json:"count"`
	}{
		Memories: memories,
		Count:    len(memories),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleAddMemory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type     string `json:"type"`
		Content  string `json:"content"`
		Metadata string `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// 输入验证
	if req.Type == "" || req.Content == "" {
		http.Error(w, `{"error":"Type and content are required"}`, http.StatusBadRequest)
		return
	}

	validTypes := map[string]bool{
		"snapshot": true,
		"session":  true,
		"longterm": true,
	}
	if !validTypes[req.Type] {
		http.Error(w, `{"error":"Invalid memory type. Must be one of: snapshot, session, longterm"}`, http.StatusBadRequest)
		return
	}

	if len(req.Content) > 50000 {
		http.Error(w, `{"error":"Content too long (max 50000 characters)"}`, http.StatusBadRequest)
		return
	}

	if err := aiAgent.Memory.Add(r.Context(), req.Type, req.Content, req.Metadata); err != nil {
		log.Printf("Error adding memory: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Failed to add memory: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Memory added successfully",
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
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Update config with validation
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
	if req.Model.Temperature > 0 && req.Model.Temperature <= 2.0 {
		appConfig.Model.Temperature = req.Model.Temperature
	}
	if req.Model.MaxTokens > 0 && req.Model.MaxTokens <= 128000 {
		appConfig.Model.MaxTokens = req.Model.MaxTokens
	}

	// Validate config
	if errs := appConfig.Validate(); len(errs) > 0 {
		http.Error(w, fmt.Sprintf(`{"error":"Configuration validation failed: %v"}`, errs), http.StatusBadRequest)
		return
	}

	// Save config to file
	configPath := os.ExpandEnv("~/.carrot/config.yaml")
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating config directory: %v", err)
		http.Error(w, `{"error":"Failed to create config directory"}`, http.StatusInternalServerError)
		return
	}

	if err := appConfig.Save(configPath); err != nil {
		log.Printf("Error saving config: %v", err)
		http.Error(w, `{"error":"Failed to save config"}`, http.StatusInternalServerError)
		return
	}

	// Restart agent with new config
	ctx := context.Background()
	agentInstance := initAgent(appConfig, store)
	if err := agentInstance.Initialize(ctx); err != nil {
		log.Printf("Error initializing agent: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Failed to initialize agent: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	aiAgent = agentInstance

	response := map[string]interface{}{
		"message": "Configuration updated successfully",
		"config":  appConfig,
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
		Provider string   `json:"provider"`
		Models   []string `json:"models"`
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

func handleTools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取所有工具
	toolRegistry := aiAgent.GetToolRegistry()
	toolInfos := toolRegistry.GetToolInfo()
	toolsetInfo := toolRegistry.GetToolsetInfo()
	toolsets := toolRegistry.GetToolsets()

	// 构建响应
	response := map[string]interface{}{
		"tools":        toolInfos,
		"toolsets":     toolsets,
		"toolset_info": toolsetInfo,
		"count":        len(toolInfos),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
