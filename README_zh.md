# Carrot Agent

![Carrot Agent Logo](logo.png)

[![GitHub stars](https://img.shields.io/github/stars/zjxzjw/Carrot-Agent.svg)](https://github.com/zjxzjw/Carrot-Agent/stargazers)
[![GitHub license](https://img.shields.io/github/license/zjxzjw/Carrot-Agent.svg)](https://github.com/zjxzjw/Carrot-Agent/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/zjxzjw/Carrot-Agent)](https://goreportcard.com/report/github.com/zjxzjw/Carrot-Agent)

[English Version](README.md)

基于 Go 语言开发的智能代理框架，专为容器化运行而设计，提供持久记忆、技能学习和工具调用等核心功能。

## 📑 目录

- [📋 版本信息](#-版本信息)
- [🌟 核心特性](#-核心特性)
- [📦 快速开始](#-快速开始)
- [🎯 功能特性](#-功能特性)
- [📁 项目结构](#-项目结构)
- [🛠️ 技术栈](#-技术栈)
- [📚 配置说明](#-配置说明)
- [🤖 命令行操作](#-命令行操作)
- [📖 使用示例](#-使用示例)
- [📈 开发路线](#-开发路线)
- [📄 许可证](#-许可证)
- [🤔 常见问题](#-常见问题)
- [📞 支持](#-支持)
- [🌐 文档](#-文档)

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
- **多语言支持**：英文和中文语言选项
- **文件系统访问**：安全的文件操作，带路径限制
- **HTTP 网络访问**：安全的 HTTP 请求，带 URL 验证

## 📦 快速开始

### 1. Docker 部署 (推荐)

```bash
# 克隆代码
git clone https://github.com/zjxzjw/Carrot-Agent.git
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

- **工具调用**：执行各种工具操作，包括文件操作、HTTP 请求和系统命令
- **记忆管理**：分层记忆架构（快照、技能、情景、长期）
- **技能学习**：从完成的任务中自动生成和更新技能
- **会话管理**：保持跨会话上下文和对话历史
- **自进化**：从重复的工作流程中自动创建技能

### 容器化

- **Docker 支持**：官方 Docker 镜像
- **Docker Compose**：一键部署，包含所有依赖
- **数据持久化**：卷挂载确保容器重启后数据不丢失
- **安全隔离**：非 root 用户运行，带有限制的文件系统访问

### Web 界面

- **现代 UI**：使用 React、TypeScript 和 Ant Design 构建
- **多语言**：支持英文和中文
- **响应式设计**：适用于桌面和移动设备
- **实时更新**：带工具执行结果的实时聊天界面

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
│   ├── public/             # 静态资源
│   ├── src/                # React 源代码
│   │   ├── components/     # UI 组件
│   │   ├── pages/          # 应用页面
│   │   ├── services/       # API 服务
│   │   └── store/          # Redux 存储
│   └── package.json        # 前端依赖
├── website/                # 文档网站
│   ├── docs/               # Markdown 文档
│   └── package.json        # 文档站点依赖
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

- **语言**：Go 1.26.2+
- **存储**：SQLite (嵌入式数据库)
- **配置**：YAML
- **容器**：Docker、Docker Compose
- **前端**：React 18.2.0、TypeScript 5.2.2、Ant Design 5.12.8
- **状态管理**：Redux Toolkit 2.0.1
- **API**：RESTful HTTP API
- **模型**：OpenAI GPT、Claude (Anthropic)
- **日志**：结构化日志
- **测试**：Go 测试框架
- **构建工具**：Vite 5.0.8 (前端)

## 📚 配置说明

### 环境变量

| 变量名                     | 描述         | 默认值                         |
| ----------------------- | ---------- | --------------------------- |
| CARROT\_API\_KEY        | 模型 API Key | -                           |
| CARROT\_MODEL\_PROVIDER | 模型提供者      | openai                      |
| CARROT\_MODEL\_NAME     | 模型名称       | gpt-4                       |
| CARROT\_BASE\_URL       | API 基础 URL | <https://api.openai.com/v1> |

### 认证

Carrot Agent 包含认证功能。默认凭据为：
- 用户名：`admin`
- 密码：`admin123`

你可以在 config.yaml 文件的 `auth` 部分配置这些凭据：

```yaml
auth:
  username: your-username
  password: your-password
```

访问 Web 界面时，你将首先被重定向到登录页面。成功登录后，你将能够访问所有功能。

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

## 📖 使用示例

### 基本对话

```bash
$ go run ./cmd/cli

Carrot Agent v0.1.0
Type 'help' for available commands

> 你好，你能做什么？
我是 Carrot Agent，一个具有记忆和技能学习能力的智能助手。我可以：
- 执行工具操作
- 管理分层记忆
- 学习和生成技能
- 保持跨会话上下文
- 提供系统信息
- 安全访问文件
- 发送 HTTP 请求

今天我能为您提供什么帮助？
```

### 使用工具

```bash
> 当前时间是什么？
{
  "current_time": "2024-01-01T12:00:00Z",
  "unix_time": 1704067200
}

> 读取配置文件
{
  "toolcall": {
    "thought": "我需要读取配置文件来查看当前配置",
    "name": "file_read",
    "args": {
      "file_path": "~/.carrot/config.yaml"
    }
  }
}

> 获取系统信息
{
  "toolcall": {
    "thought": "获取系统信息",
    "name": "system_info",
    "args": {}
  }
}
```

### 创建技能

```bash
> 创建一个生成每日报告的技能
{
  "toolcall": {
    "thought": "创建一个生成每日报告的技能",
    "name": "skill_create",
    "args": {
      "name": "daily_report",
      "description": "生成每日活动报告",
      "content": "# 每日报告生成器\n\n此技能根据最近的活动生成每日报告。\n\n## 使用方法\n1. 收集活动数据\n2. 分析模式\n3. 生成摘要\n4. 将报告保存到文件"
    }
  }
}
```

### 使用 HTTP 工具

```bash
> 获取北京的当前天气
{
  "toolcall": {
    "thought": "我需要发送 HTTP 请求来获取天气信息",
    "name": "http_get",
    "args": {
      "url": "https://api.weatherapi.com/v1/current.json?key=YOUR_API_KEY&q=Beijing"
    }
  }
}
```

## 📈 开发路线

1. **核心代理功能**：✅ 实现模型调用、工具执行
2. **记忆系统**：✅ 实现分层记忆管理
3. **技能系统**：✅ 实现自动技能生成
4. **容器化**：✅ 完善 Docker 部署
5. **API 服务**：✅ 实现 REST API
6. **Web 界面**：✅ 构建现代 React UI
7. **多语言支持**：✅ 添加英文和中文本地化
8. **高级工具集成**：增强工具能力和安全性
9. **性能优化**：提高响应时间和资源使用效率
10. **扩展模型支持**：添加更多模型提供者

## 📄 许可证

MIT License

## 🤔 常见问题

### Q: 如何配置模型？

A: 在 `config.yaml` 中设置 `model` 部分的参数，或通过环境变量配置。

### Q: 数据存储在哪里？

A: 数据存储在 `~/.carrot` 目录下的 SQLite 数据库中。

### Q: 如何添加自定义工具？

A: 在 `pkg/agent/agent.go` 的 `registerDefaultTools` 方法中注册新工具。

### Q: 文件系统访问安全吗？

A: 是的，Carrot Agent 有路径限制和黑名单，防止未授权访问敏感文件和目录。

### Q: 我可以使用自定义模型提供者吗？

A: 是的，你可以通过扩展 `model.Provider` 接口来实现自定义模型提供者。

## 📞 支持

如有问题或建议，请提交 Issue 或联系我们。

## 🌐 文档

详细文档请访问我们的官方网站：

- [官方文档](https://zjxzjw.github.io/Carrot-Agent/zh/)
- [API 参考](https://zjxzjw.github.io/Carrot-Agent/zh/api/overview)
- [使用示例](https://zjxzjw.github.io/Carrot-Agent/zh/examples/basic)
- [架构指南](https://zjxzjw.github.io/Carrot-Agent/zh/guide/architecture)
