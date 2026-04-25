# Architecture

Learn about Carrot Agent's internal architecture and design principles.

## Overview

Carrot Agent follows a modular architecture with clear separation of concerns.

## Core Components

### Agent Engine

The central component that orchestrates:
- Model interactions
- Tool execution
- Memory management
- Skill learning

### Memory System

Three-tier architecture:
1. **Snapshot Memory**: Short-term context
2. **Session Memory**: Conversation history
3. **Long-term Memory**: Persistent knowledge

### Skill Manager

Handles:
- Skill creation
- Skill updates
- Skill search
- Automatic skill generation

### Tool Registry

Manages available tools:
- File operations
- HTTP requests
- Memory tools
- System utilities

## Data Flow

```
User Input → Agent Engine → Model Provider
                ↓
         Tool Execution
                ↓
        Memory/Skill Update
                ↓
          Response to User
```

## Security Model

- Path whitelisting
- URL filtering
- Non-root execution
- Command blocking

## Scalability

- Stateless API design
- SQLite for single-instance
- Ready for horizontal scaling
