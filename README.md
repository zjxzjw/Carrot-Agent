# 🥕 Carrot Agent

[中文版本 (Chinese Version)](readme_zh.md)

An intelligent agent framework developed in Go, designed for containerized operation, providing core features such as persistent memory, skill learning, and tool calling.

## 🌟 Core Features

- **Intelligent Agent Functions**: Tool calling, hierarchical memory management, automatic skill learning, cross-session context maintenance
- **Multi-model Support**: OpenAI GPT, Claude, and other large language models
- **Container-first**: Official Docker images and Docker Compose configuration
- **Persistent Storage**: SQLite database for storing memories, skills, and sessions
- **Self-evolution**: Automatically generate reusable skills after completing complex tasks
- **High Performance**: Go language implementation with low resource consumption
- **Secure Isolation**: Runs as non-root user with least privilege principle

## 📦 Quick Start

### 1. Docker Deployment (Recommended)

```bash
# Clone code
git clone https://github.com/your-org/carrot-agent.git
cd carrot-agent

# Configure environment variables
cp .env.example .env
vim .env  # Fill in your API Key

# Start container
docker-compose up -d

# Enter container
docker exec -it carrot-agent /bin/sh
/app/carrot-agent
```

### 2. Local Run

```bash
# Install dependencies
go mod tidy

# Configure
cp config.yaml.example ~/.carrot/config.yaml
vim ~/.carrot/config.yaml

# Run
go run ./cmd/cli
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
├── cmd/cli/main.go           # CLI entry point
├── pkg/
│   ├── agent/               # Core agent engine
│   │   ├── agent.go         # Agent core logic
│   │   ├── memory/          # Hierarchical memory management
│   │   ├── skill/           # Skill system
│   │   ├── model/           # Model providers
│   │   └── tool/            # Tool registry
│   └── storage/            # Storage management
├── config/                # Configuration management
├── Dockerfile             # Containerization build
├── docker-compose.yaml    # Docker Compose deployment
├── config.yaml.example    # Configuration example
└── ARCHITECTURE.md       # Architecture documentation
```

## 🛠️ Tech Stack

- **Language**: Go 1.22+
- **Storage**: SQLite
- **Configuration**: YAML
- **Container**: Docker
- **Models**: OpenAI GPT, Claude

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
