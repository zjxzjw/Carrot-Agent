# Quick Start

Get Carrot Agent running in under 5 minutes!

## Prerequisites

- Docker and Docker Compose installed
- An API key from OpenAI, Anthropic, or another LLM provider

## Step 1: Clone the Repository

```bash
git clone https://github.com/zjxzjw/Carrot-Agent.git
cd carrot-agent
```

## Step 2: Configure Environment Variables

Create a `.env` file:

```bash
cp .env.example .env
```

Edit `.env` and add your API key:

```env
CARROT_API_KEY=your-api-key-here
CARROT_MODEL_PROVIDER=openai
CARROT_MODEL_NAME=gpt-4
CARROT_BASE_URL=https://api.openai.com/v1
```

### Supported Model Providers

| Provider | Model Examples | Base URL |
|----------|---------------|----------|
| OpenAI | gpt-4, gpt-3.5-turbo | `https://api.openai.com/v1` |
| Claude | claude-3-opus, claude-3-sonnet | `https://api.anthropic.com/v1` |
| OpenRouter | Various models | `https://openrouter.ai/api/v1` |

## Step 3: Start with Docker Compose

```bash
docker-compose up -d
```

This will:
- Build the Docker image
- Start the API server on port 8080
- Create persistent volumes for data storage

## Step 4: Verify Installation

Check if the service is running:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

## Step 5: Start Chatting

### Option A: Using cURL

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello! What can you help me with?",
    "session_id": "test-session"
  }'
```

### Option B: Using Web UI

Open your browser and navigate to:

```
http://localhost:8080
```

The web interface provides:
- 💬 Chat interface
- 📚 Skill management
- 🧠 Memory browser
- 📊 Statistics dashboard
- 📝 Session history

### Option C: Using CLI

Enter the container:

```bash
docker exec -it carrot-agent /bin/sh
/app/carrot-agent
```

## First Conversation

Try asking the agent to perform a task:

```
User: "Help me create a Python script that reads a CSV file and displays statistics"

Agent: I'll help you create that script. Let me use the file_write tool...
       [Creates the script]
       
       Would you like me to save this as a skill for future use?
```

After completing complex tasks, the agent will automatically suggest creating skills!

## Next Steps

- [Installation Guide](/guide/installation) - Detailed installation options
- [Configuration](/guide/configuration) - Customize your agent
- [Core Concepts](/guide/architecture) - Learn how it works
- [API Reference](/api/overview) - Explore the REST API

## Troubleshooting

### Container won't start

Check logs:
```bash
docker-compose logs -f
```

### API key error

Verify your `.env` file:
```bash
docker-compose config
```

### Port already in use

Change the port in `docker-compose.yaml`:
```yaml
ports:
  - "8081:8080"  # Use port 8081 instead
```

## Need Help?

- 📖 Read the [Documentation](/guide/introduction)
- 🐛 Report issues on [GitHub](https://github.com/zjxzjw/Carrot-Agent/issues)
- 💬 Join our community discussions
