# 🥕 Carrot Agent

![Carrot Agent Logo](logo.png)

[English Version](README.md)

基于 Go 语言开发的智能代理框架，专为容器化运行而设计，提供持久记忆、技能学习和工具调用等核心功能。

## 📋 版本信息

当前版本：0.1.0

## 🌟 核心特性

- **智能代理功能**：工具调用、分层记忆管理、自动技能学习、跨会话上下文保持
- **多模型支持**：OpenAI GPT、Claude 等多种大语言模型
- **容器优先**：官方 Docker 镜像和 Docker Compose 配置
- **持久化存储**：SQLite 数据库存储记忆、技能和会话
- **自进化能力**：完成复杂任务后自动生成可复用技能
- **高性能**：Go 语言实现，低资源占用
- **安全隔离**：非 root 用户运行，最小权限原则
- **API 接口**：RESTful API 用于与其他系统集成
- **Web 界面**：内置 Web 界面，方便管理

## 📦 快速开始

### 1. Docker 部署 (推荐)

```bash
# 克隆代码
git clone https://github.com/your-org/carrot-agent.git
cd carrot-agent

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件添加你的 API key
# 示例：CARROT_API_KEY=your-api-key

# 启动容器
docker-compose up -d

# 访问 Web 界面
# 在浏览器中打开 http://localhost:8080

# 进入容器（用于 CLI 访问）
docker exec -it carrot-agent /bin/sh
/app/carrot-agent
```

### 2. 本地运行

```bash
# 安装依赖
go mod tidy

# 创建配置目录
mkdir -p ~/.carrot

# 配置
cp config.yaml.example ~/.carrot/config.yaml
# 编辑 ~/.carrot/config.yaml 添加你的 API key

# 运行 CLI
go run ./cmd/cli

# 或运行 API 服务器
go run ./cmd/api
# 然后访问 http://localhost:8080
```

## 🎯 功能特性

### 智能代理

- **工具调用**：执行各种工具操作
- **记忆管理**：分层记忆架构（快照、技能、情景、长期）
- **技能学习**：自动生成和更新技能
- **会话管理**：保持跨会话上下文

### 容器化

- **Docker 支持**：官方 Docker 镜像
- **Docker Compose**：一键部署
- **数据持久化**：卷挂载确保数据不丢失
- **安全隔离**：非 root 用户运行

## 📁 项目结构

```
carrot-agent/
├── cmd/                    # 命令行工具
│   ├── api/main.go         # API 服务器入口
│   └── cli/main.go         # CLI 入口
├── config/                 # 配置管理
├── pkg/                    # 核心包
│   ├── agent/              # 核心代理引擎
│   │   ├── memory/         # 分层记忆管理
│   │   ├── model/          # 模型提供者 (OpenAI, Claude)
│   │   ├── skill/          # 技能系统
│   │   ├── tool/           # 工具注册表
│   │   ├── agent.go        # 代理核心逻辑
│   │   └── agent_test.go   # 代理测试
│   ├── logger/             # 日志系统
│   └── storage/            # 存储管理 (SQLite)
├── ui/                     # Web 界面
├── website/                # 文档网站
├── Dockerfile              # 容器化构建
├── docker-compose.yaml     # Docker Compose 部署
├── config.yaml.example     # 配置示例
├── ARCHITECTURE.md         # 架构文档
├── README.md               # 英文文档
├── README_zh.md            # 中文文档
├── go.mod                  # Go 模块文件
└── go.sum                  # Go 模块校验和
```

## 🛠️ 技术栈

- **语言**：Go 1.22+
- **存储**：SQLite (嵌入式数据库)
- **配置**：YAML
- **容器**：Docker、Docker Compose
- **前端**：React、TypeScript、Ant Design
- **API**：RESTful HTTP API
- **模型**：OpenAI GPT、Claude (Anthropic)
- **日志**：结构化日志
- **测试**：Go 测试框架

## 📚 配置说明

### 环境变量

| 变量名                     | 描述         | 默认值                         |
| ----------------------- | ---------- | --------------------------- |
| CARROT\_API\_KEY        | 模型 API Key | -                           |
| CARROT\_MODEL\_PROVIDER | 模型提供者      | openai                      |
| CARROT\_MODEL\_NAME     | 模型名称       | gpt-4                       |
| CARROT\_BASE\_URL       | API 基础 URL | <https://api.openai.com/v1> |

### 配置文件

参考 `config.yaml.example` 文件配置详细参数。

## 🤖 命令行操作

| 命令        | 描述       |
| --------- | -------- |
| help      | 显示帮助信息   |
| quit/exit | 退出代理     |
| reset     | 重置对话历史   |
| stats     | 显示代理统计信息 |
| skills    | 列出可用技能   |

## 📈 开发路线

1. **核心代理功能**：实现模型调用、工具执行
2. **记忆系统**：实现分层记忆管理
3. **技能系统**：实现自动技能生成
4. **容器化**：完善 Docker 部署
5. **API 服务**：实现 REST API

## 📄 许可证

MIT License

## 🤔 常见问题

### Q: 如何配置模型？

A: 在 `config.yaml` 中设置 `model` 部分的参数，或通过环境变量配置。

### Q: 数据存储在哪里？

A: 数据存储在 `~/.carrot` 目录下的 SQLite 数据库中。

### Q: 如何添加自定义工具？

A: 在 `pkg/agent/agent.go` 的 `registerDefaultTools` 方法中注册新工具。

## 📞 支持

如有问题或建议，请提交 Issue 或联系我们。