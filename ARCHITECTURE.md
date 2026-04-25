# Carrot Agent Architecture Documentation

[中文版本 (Chinese Version)](architecture_zh.md)

## 1. Project Architecture

Carrot Agent adopts a layered architecture design, consisting of the following core components:

```
carrot-agent/
├── cmd/                  # Command line and API entry points
│   ├── cli/              # Command line interface
│   └── api/              # REST API interface
├── config/               # Configuration management
├── pkg/                  # Core functionality packages
│   ├── agent/            # Intelligent agent engine
│   │   ├── memory/       # Hierarchical memory management
│   │   ├── model/        # Model providers
│   │   ├── skill/        # Skill system
│   │   ├── tool/         # Tool registry
│   │   └── agent.go      # Agent core logic
│   └── storage/          # Storage management
├── ui/                   # Frontend interface
├── Dockerfile            # Containerization build
└── docker-compose.yaml   # Docker Compose deployment
```

## 2. Core Components

### 2.1 Intelligent Agent Engine (AIAgent)

The intelligent agent engine is the core of the entire system, responsible for coordinating the work of various components, processing user input, and generating responses.

**Main Functions:**
- Process user input and generate responses
- Manage conversation history
- Coordinate tool calls
- Manage skills and memory
- Automatic skill generation

**Core Methods:**
- `RunConversation()`: Run conversation, process user input and generate responses
- `ProcessToolCalls()`: Process tool calls
- `buildSystemPrompt()`: Build system prompt
- `triggerSkillNudge()`: Trigger skill generation

### 2.2 Hierarchical Memory Management (MemoryManager)

The memory management system is responsible for storing and retrieving agent memories, supporting different types of memory storage.

**Memory Types:**
- **Snapshot Memory**: Store important user information and environmental information
- **Session Memory**: Store conversation session information
- **Long-term Memory**: Store long-term useful information

**Main Functions:**
- Add, update, delete memories
- Query memories by type
- Search memories
- Manage memory statistics

### 2.3 Skill System (SkillManager)

The skill system is responsible for managing agent skills, supporting skill creation, update, and query.

**Main Functions:**
- Create and update skills
- List and search skills
- Skill index management
- Automatic skill generation

### 2.4 Model Provider (Provider)

The model provider is responsible for interacting with large language models, supporting different model providers.

**Supported Models:**
- OpenAI GPT
- Claude (Anthropic)

**Main Functions:**
- Send chat requests
- Process model responses
- Support tool calls

### 2.5 Tool Registry (ToolRegistry)

The tool registry manages the tools that the agent can use, supporting tool registration and execution.

**Built-in Tools:**
- `memory_read`: Read memory
- `memory_write`: Write memory
- `skill_create`: Create skill
- `skill_update`: Update skill
- `skill_list`: List skills
- `skill_search`: Search skills
- `system_info`: Get system information
- `file_read`: Read file
- `file_write`: Write file
- `http_get`: Send HTTP GET request
- `get_time`: Get current time

### 2.6 Storage Management (Store)

Storage management is responsible for data persistence, using SQLite database to store memories, skills, and session information.

**Main Functions:**
- Save and read memories
- Save and read skills
- Save and read sessions
- Search data

## 3. Data Flow

### 3.1 Conversation Flow

1. User input → AIAgent.RunConversation()
2. Build system prompt → AIAgent.buildSystemPrompt()
3. Send request to model → Provider.Chat()
4. Process model response → AIAgent.ProcessToolCalls()
5. Execute tool calls → ToolRegistry.Execute()
6. Generate final response → Return to user

### 3.2 Skill Generation Flow

1. Tool call count → AIAgent.toolCallCount
2. Reach threshold → AIAgent.triggerSkillNudge()
3. Extract conversation content → Generate skill content
4. Save skill → SkillManager.Create()

## 4. API Interfaces

### 4.1 Chat Interface

- **URL**: `/api/chat`
- **Method**: POST
- **Request Body**:
  ```json
  {
    "message": "User message",
    "session_id": "Session ID"
  }
  ```
- **Response**:
  ```json
  {
    "message": "AI response",
    "usage": {
      "prompt_tokens": 100,
      "completion_tokens": 50,
      "total_tokens": 150
    }
  }
  ```

### 4.2 Skill Interface

- **URL**: `/api/skills`
- **Method**: GET (list skills), POST (create skill)
- **Response**:
  ```json
  {
    "skills": [
      {
        "id": "skill_123",
        "name": "Test Skill",
        "description": "Test skill description",
        "version": "1.0.0",
        "platforms": "[\"macos\",\"linux\"]",
        "content": "Skill content",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "count": 1
  }
  ```

### 4.3 Memory Interface

- **URL**: `/api/memory`
- **Method**: GET (list memories), POST (add memory)
- **Response**:
  ```json
  {
    "memories": [
      {
        "id": "memory_123",
        "type": "snapshot",
        "content": "Memory content",
        "metadata": "{}",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "count": 1
  }
  ```

### 4.4 Statistics Interface

- **URL**: `/api/stats`
- **Method**: GET
- **Response**:
  ```json
  {
    "tool_call_count": 10,
    "skill_count": 5,
    "memory_stats": {
      "snapshot": 2,
      "session": 3,
      "longterm": 1
    },
    "conversation_len": 5
  }
  ```

## 5. Deployment Architecture

### 5.1 Docker Deployment

Carrot Agent supports Docker deployment, using Docker Compose for container orchestration.

**Core Configuration:**
- Container image: `carrotagent/carrot-agent:latest`
- Port mapping: `8080:8080`
- Data volume: `carrot-data:/home/carrot/.carrot`
- Environment variables: Configure model API Key and other parameters

### 5.2 Local Deployment

1. Install dependencies: `go mod tidy`
2. Configure: `cp config.yaml.example ~/.carrot/config.yaml`
3. Run API service: `go run ./cmd/api`
4. Run frontend: `cd ui && npm install && npm run dev`

## 6. Extension Guide

### 6.1 Adding New Tools

1. Register new tools in the `registerDefaultTools()` method in `pkg/agent/agent.go`
2. Implement tool functions to handle tool call logic
3. Configure tool parameters and descriptions

### 6.2 Adding New Models

1. Implement new model providers in `pkg/agent/model/provider.go`
2. Add model types in the `ProviderFactory.CreateProvider()` method
3. Configure model parameters and authentication

### 6.3 Extending Frontend

1. Add new states and methods in `ui/src/store.ts`
2. Add new components and pages in `ui/src/App.tsx`
3. Implement interaction with backend API

## 7. Best Practices

- **Skill Management**: Regularly check and update skills to ensure their accuracy and usefulness
- **Memory Management**: Use different types of memory appropriately to avoid memory bloat
- **Tool Usage**: Create tools for common functions to improve agent efficiency
- **Model Selection**: Choose the appropriate model based on task requirements, balancing performance and cost
- **Containerization**: Use Docker deployment to ensure environment consistency and portability

## 8. Troubleshooting

- **API Connection Failure**: Check network connection and model API Key configuration
- **Tool Execution Failure**: Check tool parameters and permission settings
- **Insufficient Memory**: Clean up unnecessary memories and skills
- **Performance Issues**: Optimize model parameters and tool execution logic

## 9. Future Plans

- **Multimodal Support**: Add support for multimodal inputs such as images and voice
- **Plugin System**: Implement a plugin system to support third-party plugins
- **Knowledge Base Integration**: Integrate external knowledge bases to enhance agent knowledge capabilities
- **Security Enhancement**: Strengthen security measures to prevent malicious tool calls
- **Performance Optimization**: Optimize system performance and improve response speed
