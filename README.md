# Carrot Agent

![Carrot Agent Logo](logo.png)

[![GitHub stars](https://img.shields.io/github/stars/zjxzjw/Carrot-Agent.svg)](https://github.com/zjxzjw/Carrot-Agent/stargazers)
[![GitHub license](https://img.shields.io/github/license/zjxzjw/Carrot-Agent.svg)](https://github.com/zjxzjw/Carrot-Agent/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/zjxzjw/Carrot-Agent)](https://goreportcard.com/report/github.com/zjxzjw/Carrot-Agent)

[中文版本 (Chinese Version)](README_zh.md)

An intelligent agent framework developed in Go, designed for containerized operation, providing core features such as persistent memory, skill learning, and tool calling.

## 📑 Table of Contents

- [📋 Version](#-version)
- [🌟 Core Features](#-core-features)
- [📦 Quick Start](#-quick-start)
- [🎯 Features](#-features)
- [📁 Project Structure](#-project-structure)
- [🛠️ Tech Stack](#-tech-stack)
- [📚 Configuration](#-configuration)
- [🤖 Command Line Operations](#-command-line-operations)
- [📖 Usage Examples](#-usage-examples)
- [📈 Development Roadmap](#-development-roadmap)
- [📄 License](#-license)
- [🤔 FAQ](#-faq)
- [📞 Support](#-support)
- [🌐 Documentation](#-documentation)

## 📋 Version

Current Version: 0.1.0

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
- **Multi-language Support**: English and Chinese language options
- **File System Access**: Secure file operations with path restrictions
- **HTTP Network Access**: Secure HTTP requests with URL validation

## 📦 Quick Start

### 1. Docker Deployment (Recommended)

```bash
# Clone code
git clone https://github.com/zjxzjw/Carrot-Agent.git
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

- **Tool Calling**: Execute various tool operations including file operations, HTTP requests, and system commands
- **Memory Management**: Hierarchical memory architecture (snapshot, skill, context, long-term)
- **Skill Learning**: Automatically generate and update skills from completed tasks
- **Session Management**: Maintain cross-session context and conversation history
- **Self-evolution**: Automatically create skills from repeated workflows

### Containerization

- **Docker Support**: Official Docker images
- **Docker Compose**: One-click deployment with all dependencies
- **Data Persistence**: Volume mounting ensures data is not lost between container restarts
- **Secure Isolation**: Runs as non-root user with restricted file system access

### Web Interface

- **Modern UI**: Built with React, TypeScript, and Ant Design
- **Multi-language**: English and Chinese language support
- **Responsive Design**: Works on desktop and mobile devices
- **Real-time Updates**: Live chat interface with tool execution results

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
│   ├── public/             # Static assets
│   ├── src/                # React source code
│   │   ├── components/     # UI components
│   │   ├── pages/          # Application pages
│   │   ├── services/       # API services
│   │   └── store/          # Redux store
│   └── package.json        # Frontend dependencies
├── website/                # Documentation website
│   ├── docs/               # Markdown documentation
│   └── package.json        # Documentation site dependencies
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

- **Language**: Go 1.26.2+
- **Storage**: SQLite (embedded database)
- **Configuration**: YAML
- **Container**: Docker, Docker Compose
- **Frontend**: React 18.2.0, TypeScript 5.2.2, Ant Design 5.12.8
- **State Management**: Redux Toolkit 2.0.1
- **API**: RESTful HTTP API
- **Models**: OpenAI GPT, Claude (Anthropic)
- **Logging**: Structured logging
- **Testing**: Go testing framework
- **Build Tool**: Vite 5.0.8 (frontend)

## 📚 Configuration

### Environment Variables

| Variable Name           | Description    | Default Value               |
| ----------------------- | -------------- | --------------------------- |
| CARROT\_API\_KEY        | Model API Key  | -                           |
| CARROT\_MODEL\_PROVIDER | Model Provider | openai                      |
| CARROT\_MODEL\_NAME     | Model Name     | gpt-4                       |
| CARROT\_BASE\_URL       | API Base URL   | <https://api.openai.com/v1> |

### Authentication

Carrot Agent includes authentication functionality. The default credentials are:
- Username: `admin`
- Password: `admin123`

You can configure these credentials in the `auth` section of the config.yaml file:

```yaml
auth:
  username: your-username
  password: your-password
```

When accessing the web interface, you will be redirected to the login page first. After successful login, you will be able to access all features.

### Configuration File

Refer to the `config.yaml.example` file for detailed configuration parameters.

## 🤖 Command Line Operations

| Command   | Description                |
| --------- | -------------------------- |
| help      | Display help information   |
| quit/exit | Exit the agent             |
| reset     | Reset conversation history |
| stats     | Display agent statistics   |
| skills    | List available skills      |

## 📖 Usage Examples

### Basic Conversation

```bash
$ go run ./cmd/cli

Carrot Agent v0.1.0
Type 'help' for available commands

> Hello, what can you do?
I'm Carrot Agent, an intelligent assistant with memory and skill learning capabilities. I can:
- Execute tool operations
- Manage hierarchical memory
- Learn and generate skills
- Maintain cross-session context
- Provide system information
- Access files securely
- Make HTTP requests

How can I assist you today?
```

### Using Tools

```bash
> What's the current time?
{
  "current_time": "2024-01-01T12:00:00Z",
  "unix_time": 1704067200
}

> Read the config file
{
  "toolcall": {
    "thought": "I need to read the config file to see the current configuration",
    "name": "file_read",
    "args": {
      "file_path": "~/.carrot/config.yaml"
    }
  }
}

> Get system information
{
  "toolcall": {
    "thought": "Getting system information",
    "name": "system_info",
    "args": {}
  }
}
```

### Creating Skills

```bash
> Create a skill for generating daily reports
{
  "toolcall": {
    "thought": "Creating a skill for generating daily reports",
    "name": "skill_create",
    "args": {
      "name": "daily_report",
      "description": "Generate daily activity report",
      "content": "# Daily Report Generator\n\nThis skill generates a daily report based on recent activities.\n\n## Usage\n1. Collect activity data\n2. Analyze patterns\n3. Generate summary\n4. Save report to file"
    }
  }
}
```

### Using HTTP Tools

```bash
> Get the current weather in Beijing
{
  "toolcall": {
    "thought": "I need to make an HTTP request to get weather information",
    "name": "http_get",
    "args": {
      "url": "https://api.weatherapi.com/v1/current.json?key=YOUR_API_KEY&q=Beijing"
    }
  }
}
```

## 📈 Development Roadmap

1. **Core Agent Functionality**: ✅ Implement model calling and tool execution
2. **Memory System**: ✅ Implement hierarchical memory management
3. **Skill System**: ✅ Implement automatic skill generation
4. **Containerization**: ✅ Complete Docker deployment
5. **API Service**: ✅ Implement REST API
6. **Web Interface**: ✅ Build modern React UI
7. **Multi-language Support**: ✅ Add English and Chinese localization
8. **Advanced Tool Integration**: Enhance tool capabilities and security
9. **Performance Optimization**: Improve response times and resource usage
10. **Extended Model Support**: Add more model providers

## 📄 License

MIT License

## 🤔 FAQ

### Q: How to configure the model?

A: Set the parameters in the `model` section of `config.yaml`, or configure via environment variables.

### Q: Where is the data stored?

A: Data is stored in the SQLite database under the `~/.carrot` directory.

### Q: How to add custom tools?

A: Register new tools in the `registerDefaultTools` method in `pkg/agent/agent.go`.

### Q: Is the file system access secure?

A: Yes, Carrot Agent has path restrictions and a blocklist to prevent unauthorized access to sensitive files and directories.

### Q: Can I use custom model providers?

A: Yes, you can implement custom model providers by extending the `model.Provider` interface.

## 📞 Support

For questions or suggestions, please submit an Issue or contact us.

## 🌐 Documentation

For detailed documentation, please visit our official website:

- [Official Documentation](https://zjxzjw.github.io/Carrot-Agent/)
- [API Reference](https://zjxzjw.github.io/Carrot-Agent/api/overview)
- [Usage Examples](https://zjxzjw.github.io/Carrot-Agent/examples/basic)
- [Architecture Guide](https://zjxzjw.github.io/Carrot-Agent/guide/architecture)
