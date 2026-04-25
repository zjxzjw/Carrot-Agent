# Introduction

Carrot Agent is an intelligent agent framework developed in Go, designed for containerized operation with persistent memory and automatic skill learning capabilities.

## 🌟 Why Carrot Agent?

Unlike traditional chatbots, Carrot Agent:

- **Remembers**: Maintains context across sessions with hierarchical memory
- **Learns**: Automatically creates reusable skills from complex workflows
- **Acts**: Executes real-world tasks through tool calling
- **Scales**: Runs efficiently in containers with minimal resource usage

## Key Features

### 🧠 Hierarchical Memory System

Three-tier memory architecture ensures your agent never forgets important information:

- **Snapshot Memory**: Short-term contextual information
- **Session Memory**: Conversation history and state
- **Long-term Memory**: Persistent knowledge across sessions

### 🎯 Automatic Skill Learning

After completing complex tasks (5+ tool calls), the agent automatically:
1. Analyzes the workflow
2. Generates a reusable skill
3. Saves it for future use

This means your agent gets smarter over time!

### 🔧 Powerful Tool Registry

Built-in tools include:
- File read/write operations
- HTTP requests
- Memory management
- Skill CRUD operations
- System information
- Time utilities

All tools run with security constraints to prevent unauthorized access.

### 🐳 Container-First Design

- Official Docker images
- Non-root user execution
- Volume-based data persistence
- Health checks and auto-restart
- One-command deployment with Docker Compose

### 🤖 Multi-Model Support

Works seamlessly with:
- OpenAI GPT models (GPT-4, GPT-3.5)
- Anthropic Claude
- Any OpenAI-compatible API (OpenRouter, etc.)

## Use Cases

- **Personal Assistant**: Remember preferences and maintain context
- **Development Helper**: Code generation, file operations, documentation
- **Research Agent**: Web scraping, data collection, analysis
- **Automation Workflows**: Complex multi-step task automation
- **Knowledge Management**: Organize and retrieve information

## Quick Example

```bash
# Start with Docker Compose
docker-compose up -d

# Chat with your agent
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Help me create a Python script that reads a CSV file",
    "session_id": "my-session"
  }'
```

The agent will:
1. Understand your request
2. Use appropriate tools (file operations, code generation)
3. Save the workflow as a skill for future use
4. Remember your preferences

## Next Steps

- [Quick Start](/guide/quick-start) - Get up and running in 5 minutes
- [Installation](/guide/installation) - Detailed installation guide
- [Architecture](/guide/architecture) - Understand how it works
