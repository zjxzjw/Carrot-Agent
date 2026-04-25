# 快速上手

在 5 分钟内启动并运行 Carrot Agent！

## 前置要求

- 已安装 Docker 和 Docker Compose
- 拥有 OpenAI、Anthropic 或其他 LLM 提供商的 API 密钥

## 步骤 1：克隆仓库

```bash
git clone https://github.com/zjxzjw/Carrot-Agent.git
cd carrot-agent
```

## 步骤 2：配置环境变量

创建 `.env` 文件：

```bash
cp .env.example .env
```

编辑 `.env` 并添加您的 API 密钥：

```env
CARROT_API_KEY=your-api-key-here
CARROT_MODEL_PROVIDER=openai
CARROT_MODEL_NAME=gpt-4
CARROT_BASE_URL=https://api.openai.com/v1
```

### 支持的模型提供商

| 提供商 | 模型示例 | Base URL |
|----------|---------------|----------|
| OpenAI | gpt-4, gpt-3.5-turbo | `https://api.openai.com/v1` |
| Claude | claude-3-opus, claude-3-sonnet | `https://api.anthropic.com/v1` |
| OpenRouter | 各种模型 | `https://openrouter.ai/api/v1` |

## 步骤 3：使用 Docker Compose 启动

```bash
docker-compose up -d
```

这将：
- 构建 Docker 镜像
- 在 8080 端口启动 API 服务器
- 创建持久化卷用于数据存储

## 步骤 4：验证安装

检查服务是否正在运行：

```bash
curl http://localhost:8080/health
```

预期响应：
```json
{"status":"ok"}
```

## 步骤 5：开始聊天

### 选项 A：使用 cURL

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好！你能帮我做什么？",
    "session_id": "test-session"
  }'
```

### 选项 B：使用 Web UI

在浏览器中打开：

```
http://localhost:8080
```

Web 界面提供：
- 💬 聊天界面
- 📚 技能管理
- 🧠 记忆浏览器
- 📊 统计仪表板
- 📝 会话历史

### 选项 C：使用 CLI

进入容器：

```bash
docker exec -it carrot-agent /bin/sh
/app/carrot-agent
```

## 第一次对话

尝试让代理执行任务：

```
用户："帮我创建一个读取 CSV 文件并显示统计信息的 Python 脚本"

代理：我将帮您创建该脚本。让我使用 file_write 工具...
       [创建脚本]
       
       您希望我将此保存为技能以备将来使用吗？
```

完成复杂任务后，代理会自动建议创建技能！

## 下一步

- [安装指南](/zh/guide/installation) - 详细的安装选项
- [配置说明](/zh/guide/configuration) - 自定义您的代理
- [核心概念](/zh/guide/architecture) - 了解工作原理
- [API 参考](/zh/api/overview) - 探索 REST API

## 故障排除

### 容器无法启动

检查日志：
```bash
docker-compose logs -f
```

### API 密钥错误

验证您的 `.env` 文件：
```bash
docker-compose config
```

### 端口已被占用

在 `docker-compose.yaml` 中更改端口：
```yaml
ports:
  - "8081:8080"  # 使用 8081 端口
```

## 需要帮助？

- 📖 阅读[文档](/zh/guide/introduction)
- 🐛 在 [GitHub](https://github.com/zjxzjw/Carrot-Agent/issues) 报告问题
- 💬 加入我们的社区讨论
