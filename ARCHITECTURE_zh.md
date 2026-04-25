# Carrot Agent 架构文档

## 1. 项目架构

Carrot Agent 采用分层架构设计，主要由以下几个核心组件组成：

```
carrot-agent/
├── cmd/                  # 命令行和 API 入口
│   ├── cli/              # 命令行接口
│   └── api/              # REST API 接口
├── config/               # 配置管理
├── pkg/                  # 核心功能包
│   ├── agent/            # 智能代理引擎
│   │   ├── memory/       # 分层记忆管理
│   │   ├── model/        # 模型提供者
│   │   ├── skill/        # 技能系统
│   │   ├── tool/         # 工具注册表
│   │   └── agent.go      # 代理核心逻辑
│   └── storage/          # 存储管理
├── ui/                   # 前端界面
├── Dockerfile            # 容器化构建
└── docker-compose.yaml   # Docker Compose 部署
```

## 2. 核心组件

### 2.1 智能代理引擎 (AIAgent)

智能代理引擎是整个系统的核心，负责协调各个组件的工作，处理用户输入并生成响应。

**主要功能：**
- 处理用户输入并生成响应
- 管理对话历史
- 协调工具调用
- 管理技能和记忆
- 自动技能生成

**核心方法：**
- `RunConversation()`: 运行对话，处理用户输入并生成响应
- `ProcessToolCalls()`: 处理工具调用
- `buildSystemPrompt()`: 构建系统提示
- `triggerSkillNudge()`: 触发技能生成

### 2.2 分层记忆管理 (MemoryManager)

记忆管理系统负责存储和检索代理的记忆，支持不同类型的记忆存储。

**记忆类型：**
- **快照记忆 (snapshot)**: 存储重要的用户信息和环境信息
- **会话记忆 (session)**: 存储对话会话信息
- **长期记忆 (longterm)**: 存储长期有用的信息

**主要功能：**
- 添加、更新、删除记忆
- 按类型查询记忆
- 搜索记忆
- 管理记忆统计

### 2.3 技能系统 (SkillManager)

技能系统负责管理代理的技能，支持技能的创建、更新和查询。

**主要功能：**
- 创建和更新技能
- 列出和搜索技能
- 技能索引管理
- 自动技能生成

### 2.4 模型提供者 (Provider)

模型提供者负责与大语言模型进行交互，支持不同的模型提供商。

**支持的模型：**
- OpenAI GPT
- Claude (Anthropic)

**主要功能：**
- 发送聊天请求
- 处理模型响应
- 支持工具调用

### 2.5 工具注册表 (ToolRegistry)

工具注册表管理代理可以使用的工具，支持工具的注册和执行。

**内置工具：**
- `memory_read`: 读取记忆
- `memory_write`: 写入记忆
- `skill_create`: 创建技能
- `skill_update`: 更新技能
- `skill_list`: 列出技能
- `skill_search`: 搜索技能
- `system_info`: 获取系统信息
- `file_read`: 读取文件
- `file_write`: 写入文件
- `http_get`: 发送 HTTP GET 请求
- `get_time`: 获取当前时间

### 2.6 存储管理 (Store)

存储管理负责持久化数据，使用 SQLite 数据库存储记忆、技能和会话信息。

**主要功能：**
- 保存和读取记忆
- 保存和读取技能
- 保存和读取会话
- 搜索数据

## 3. 数据流

### 3.1 对话流程

1. 用户输入 → AIAgent.RunConversation()
2. 构建系统提示 → AIAgent.buildSystemPrompt()
3. 发送请求到模型 → Provider.Chat()
4. 处理模型响应 → AIAgent.ProcessToolCalls()
5. 执行工具调用 → ToolRegistry.Execute()
6. 生成最终响应 → 返回给用户

### 3.2 技能生成流程

1. 工具调用计数 → AIAgent.toolCallCount
2. 达到阈值 → AIAgent.triggerSkillNudge()
3. 提取对话内容 → 生成技能内容
4. 保存技能 → SkillManager.Create()

## 4. API 接口

### 4.1 聊天接口

- **URL**: `/api/chat`
- **方法**: POST
- **请求体**:
  ```json
  {
    "message": "用户消息",
    "session_id": "会话 ID"
  }
  ```
- **响应**:
  ```json
  {
    "message": "AI 响应",
    "usage": {
      "prompt_tokens": 100,
      "completion_tokens": 50,
      "total_tokens": 150
    }
  }
  ```

### 4.2 技能接口

- **URL**: `/api/skills`
- **方法**: GET (列出技能), POST (创建技能)
- **响应**:
  ```json
  {
    "skills": [
      {
        "id": "skill_123",
        "name": "测试技能",
        "description": "测试技能描述",
        "version": "1.0.0",
        "platforms": "[\"macos\",\"linux\"]",
        "content": "技能内容",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "count": 1
  }
  ```

### 4.3 记忆接口

- **URL**: `/api/memory`
- **方法**: GET (列出记忆), POST (添加记忆)
- **响应**:
  ```json
  {
    "memories": [
      {
        "id": "memory_123",
        "type": "snapshot",
        "content": "记忆内容",
        "metadata": "{}",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "count": 1
  }
  ```

### 4.4 统计接口

- **URL**: `/api/stats`
- **方法**: GET
- **响应**:
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

## 5. 部署架构

### 5.1 Docker 部署

Carrot Agent 支持 Docker 部署，使用 Docker Compose 进行容器编排。

**核心配置：**
- 容器镜像: `carrotagent/carrot-agent:latest`
- 端口映射: `8080:8080`
- 数据卷: `carrot-data:/home/carrot/.carrot`
- 环境变量: 配置模型 API Key 和其他参数

### 5.2 本地部署

1. 安装依赖: `go mod tidy`
2. 配置: `cp config.yaml.example ~/.carrot/config.yaml`
3. 运行 API 服务: `go run ./cmd/api`
4. 运行前端: `cd ui && npm install && npm run dev`

## 6. 扩展指南

### 6.1 添加新工具

1. 在 `pkg/agent/agent.go` 的 `registerDefaultTools()` 方法中注册新工具
2. 实现工具函数，处理工具调用逻辑
3. 配置工具参数和描述

### 6.2 添加新模型

1. 在 `pkg/agent/model/provider.go` 中实现新的模型提供者
2. 在 `ProviderFactory.CreateProvider()` 方法中添加模型类型
3. 配置模型参数和认证

### 6.3 扩展前端

1. 在 `ui/src/store.ts` 中添加新的状态和方法
2. 在 `ui/src/App.tsx` 中添加新的组件和页面
3. 实现与后端 API 的交互

## 7. 最佳实践

- **技能管理**: 定期检查和更新技能，确保技能的准确性和实用性
- **记忆管理**: 合理使用不同类型的记忆，避免记忆膨胀
- **工具使用**: 为常用功能创建工具，提高代理的效率
- **模型选择**: 根据任务需求选择合适的模型，平衡性能和成本
- **容器化**: 使用 Docker 部署，确保环境一致性和可移植性

## 8. 故障排查

- **API 连接失败**: 检查网络连接和模型 API Key 配置
- **工具执行失败**: 检查工具参数和权限设置
- **内存不足**: 清理不需要的记忆和技能
- **性能问题**: 优化模型参数和工具执行逻辑

## 9. 未来规划

- **多模态支持**: 添加对图像、语音等多模态输入的支持
- **插件系统**: 实现插件系统，支持第三方插件
- **知识库集成**: 集成外部知识库，增强代理的知识能力
- **安全增强**: 加强安全措施，防止恶意工具调用
- **性能优化**: 优化系统性能，提高响应速度