# 🥕 Carrot Agent

![Carrot Agent Logo](logo.png)

[中文版本 (Chinese Version)](README_zh.md)

An intelligent agent framework developed in Go, designed for containerized operation, providing core features such as persistent memory, skill learning, and tool calling.

## 🌟 Core Features

- **Intelligent Agent Functions**: Tool calling, hierarchical memory management, automatic skill learning, cross-session context maintenance
- **Multi-model Support**: OpenAI GPT, Claude, and other large language models
- **Container-first**: Official Docker images and Docker Compose configuration
- **Persistent Storage**: SQLite database for storing memories, skills, and sessions
- **Self-evolution**: Automatically generate reusable skills after completing complex tasks
- **High Performance**: Go language implementation with low resource consumption
- **Secure Isolation**: Runs as non-root user with least privilege principle
- **API Interface**: RESTful API for integration with other systems
- **Web UI**: Built-in web interface for easy management

## 📦 Quick Start

### 1. Docker Deployment (Recommended)

```bash
# Clone code
git clone https://github.com/your-org/carrot-agent.git
cd carrot-agent

# Configure environment variables
cp .env.example .env
# Edit .env file to add your API key
# Example: CARROT_API_KEY=your-api-key

# Start container
docker-compose up -d

# Access web interface
# Open http://localhost:8080 in your browser

# Enter container (for CLI access)
docker exec -it carrot-agent /bin/sh
/app/carrot-agent
```

### 2. Local Run

```bash
# Install dependencies
go mod tidy

# Create configuration directory
mkdir -p ~/.carrot

# Configure
cp config.yaml.example ~/.carrot/config.yaml
# Edit ~/.carrot/config.yaml to add your API key

# Run CLI
go run ./cmd/cli

# Or run API server
go run ./cmd/api
# Then access http://localhost:8080
```

## 🎯 Features

### Intelligent Agent

- **Tool Calling**: Execute various tool operations
- **Memory Management**: Hierarchical memory architecture (snapshot, skill, context, long-term)
- **Skill Learning**: Automatically generate and update skills
- **Session Management**: Maintain cross-session context

### Containerization

- **Docker Support**: Official Docker images
- **Docker Compose**: One-click deployment
- **Data Persistence**: Volume mounting ensures data is not lost
- **Secure Isolation**: Runs as non-root user

## 📁 Project Structure

```
carrot-agent/
├── cmd/                    # Command line tools
│   ├── api/main.go         # API server entry point
│   └── cli/main.go         # CLI entry point
├── config/                 # Configuration management
├── pkg/                    # Core packages
│   ├── agent/              # Core agent engine
│   │   ├── memory/         # Hierarchical memory management
│   │   ├── model/          # Model providers (OpenAI, Claude)
│   │   ├── skill/          # Skill system
│   │   ├── tool/           # Tool registry
│   │   ├── agent.go        # Agent core logic
│   │   └── agent_test.go   # Agent tests
│   ├── logger/             # Logging system
│   └── storage/            # Storage management (SQLite)
├── ui/                     # Web interface
├── website/                # Documentation website
├── Dockerfile              # Containerization build
├── docker-compose.yaml     # Docker Compose deployment
├── config.yaml.example     # Configuration example
├── ARCHITECTURE.md         # Architecture documentation
├── README.md               # English documentation
├── README_zh.md            # Chinese documentation
├── go.mod                  # Go module file
└── go.sum                  # Go module checksums
```

## 🛠️ Tech Stack

- **Language**: Go 1.22+
- **Storage**: SQLite (embedded database)
- **Configuration**: YAML
- **Container**: Docker, Docker Compose
- **Frontend**: React, TypeScript, Ant Design
- **API**: RESTful HTTP API
- **Models**: OpenAI GPT, Claude (Anthropic)
- **Logging**: Structured logging
- **Testing**: Go testing framework

## 📚 Configuration

### Environment Variables

| Variable Name            | Description      | Default Value                 |
| ----------------------- | --------------- | ----------------------------- |
| CARROT\_API\_KEY        | Model API Key   | -                             |
| CARROT\_MODEL\_PROVIDER | Model Provider  | openai                        |
| CARROT\_MODEL\_NAME     | Model Name      | gpt-4                         |
| CARROT\_BASE\_URL       | API Base URL    | <https://api.openai.com/v1>   |

### Configuration File

Refer to the `config.yaml.example` file for detailed configuration parameters.

## 🤖 Command Line Operations

| Command    | Description                |
| --------- | -------------------------- |
| help      | Display help information   |
| quit/exit | Exit the agent             |
| reset     | Reset conversation history |
| stats     | Display agent statistics   |
| skills    | List available skills      |

## 📈 Development Roadmap

1. **Core Agent Functionality**: Implement model calling and tool execution
2. **Memory System**: Implement hierarchical memory management
3. **Skill System**: Implement automatic skill generation
4. **Containerization**: Complete Docker deployment
5. **API Service**: Implement REST API

## 📄 License

MIT License

## 🤔 FAQ

### Q: How to configure the model?

A: Set the parameters in the `model` section of `config.yaml`, or configure via environment variables.

### Q: Where is the data stored?

A: Data is stored in the SQLite database under the `~/.carrot` directory.

### Q: How to add custom tools?

A: Register new tools in the `registerDefaultTools` method in `pkg/agent/agent.go`.

## 📞 Support

For questions or suggestions, please submit an Issue or contact us.
