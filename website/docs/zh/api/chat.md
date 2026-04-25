# 聊天 API

向代理发送消息并接收 AI 驱动的响应。

## 端点

```
POST /api/chat
```

## 请求体

```json
{
  "message": "您的消息",
  "session_id": "可选的会话ID"
}
```

### 参数

| 字段 | 类型 | 必填 | 描述 |
|-------|------|----------|-------------|
| `message` | string | 是 | 用户的消息 |
| `session_id` | string | 否 | 用于上下文持久化的会话标识符 |

## 响应

```json
{
  "message": "代理的响应",
  "usage": {
    "prompt_tokens": 150,
    "completion_tokens": 80,
    "total_tokens": 230
  }
}
```

## 示例

### 基础聊天

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "今天天气怎么样？",
    "session_id": "user-123"
  }'
```

### 带会话上下文

```bash
# 第一条消息
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我叫 Alice，我喜欢 Python",
    "session_id": "alice-session"
  }'

# 第二条消息（代理会记住）
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我喜欢什么编程语言？",
    "session_id": "alice-session"
  }'

# 响应："您提到过您喜欢 Python！"
```

### 任务执行

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "创建一个名为 hello.txt 的文件，内容为 Hello World",
    "session_id": "task-session"
  }'
```

代理将：
1. 解析请求
2. 使用 `file_write` 工具
3. 返回确认信息

## 错误响应

### 缺少消息

```json
{
  "error": "无效请求",
  "code": 400
}
```

### 模型错误

```json
{
  "error": "获取响应失败：API 密钥无效",
  "code": 500
}
```

## 最佳实践

1. **使用会话 ID**：始终提供 `session_id` 以实现上下文持久化
2. **处理错误**：为瞬时故障实现重试逻辑
3. **监控 Token 使用**：跟踪 `usage` 字段以管理成本
4. **超时处理**：设置适当的客户端超时（建议：60秒）

## 限流

考虑在生产环境中实现限流：

```nginx
# Nginx 示例
limit_req_zone $binary_remote_addr zone=chat:10m rate=10r/m;

location /api/chat {
    limit_req zone=chat burst=5;
    proxy_pass http://carrot-agent:8080;
}
```
