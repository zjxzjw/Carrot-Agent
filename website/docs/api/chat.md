# Chat API

Send messages to the agent and receive AI-powered responses.

## Endpoint

```
POST /api/chat
```

## Request Body

```json
{
  "message": "Your message here",
  "session_id": "optional-session-id"
}
```

### Parameters

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `message` | string | Yes | The user's message |
| `session_id` | string | No | Session identifier for context persistence |

## Response

```json
{
  "message": "Agent's response",
  "usage": {
    "prompt_tokens": 150,
    "completion_tokens": 80,
    "total_tokens": 230
  }
}
```

## Examples

### Basic Chat

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What is the weather like today?",
    "session_id": "user-123"
  }'
```

### With Session Context

```bash
# First message
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "My name is Alice and I like Python",
    "session_id": "alice-session"
  }'

# Second message (agent remembers)
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What programming language do I like?",
    "session_id": "alice-session"
  }'

# Response: "You mentioned that you like Python!"
```

### Task Execution

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Create a file called hello.txt with content Hello World",
    "session_id": "task-session"
  }'
```

The agent will:
1. Parse the request
2. Use the `file_write` tool
3. Return confirmation

## Error Responses

### Missing Message

```json
{
  "error": "Invalid request",
  "code": 400
}
```

### Model Error

```json
{
  "error": "Failed to get response: API key invalid",
  "code": 500
}
```

## Best Practices

1. **Use Session IDs**: Always provide a `session_id` for context persistence
2. **Handle Errors**: Implement retry logic for transient failures
3. **Monitor Token Usage**: Track `usage` field to manage costs
4. **Timeout Handling**: Set appropriate client-side timeouts (recommended: 60s)

## Rate Limiting

Consider implementing rate limiting in production:

```nginx
# Nginx example
limit_req_zone $binary_remote_addr zone=chat:10m rate=10r/m;

location /api/chat {
    limit_req zone=chat burst=5;
    proxy_pass http://carrot-agent:8080;
}
```
