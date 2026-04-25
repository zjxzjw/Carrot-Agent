# 基础使用示例

Carrot Agent 的常见使用模式。

## 示例 1：简单聊天

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好！你能做什么？",
    "session_id": "example-1"
  }'
```

## 示例 2：文件创建

让代理创建文件：

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "创建一个打印 Hello World 的 Python 脚本",
    "session_id": "example-2"
  }'
```

代理将使用 `file_write` 工具创建脚本。

## 示例 3：记忆管理

保存重要信息：

```bash
curl -X POST http://localhost:8080/api/memory \
  -H "Content-Type: application/json" \
  -d '{
    "type": "snapshot",
    "content": "用户更喜欢 Python 而不是 JavaScript",
    "metadata": "{\"category\": \"preference\"}"
  }'
```

## 示例 4：技能列表

查看可用技能：

```bash
curl http://localhost:8080/api/skills
```

响应：
```json
{
  "skills": [
    {
      "id": "skill_123",
      "name": "csv_processor",
      "description": "处理 CSV 文件并生成统计信息"
    }
  ],
  "count": 1
}
```

## 示例 5：会话管理

列出所有会话：

```bash
curl http://localhost:8080/api/session/
```

删除会话：

```bash
curl -X DELETE http://localhost:8080/api/session/example-1
```

## 示例 6：统计信息

获取代理统计：

```bash
curl http://localhost:8080/api/stats
```

响应：
```json
{
  "tool_call_count": 42,
  "skill_count": 5,
  "memory_stats": {
    "snapshot": 10,
    "session": 15,
    "longterm": 3
  },
  "conversation_len": 20
}
```

## Python 示例

使用 Python requests 库：

```python
import requests

url = "http://localhost:8080/api/chat"
payload = {
    "message": "法国的首都是什么？",
    "session_id": "python-example"
}

response = requests.post(url, json=payload)
print(response.json()["message"])
```

## JavaScript 示例

使用 fetch API：

```javascript
const response = await fetch('http://localhost:8080/api/chat', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    message: '给我讲个笑话',
    session_id: 'js-example'
  })
});

const data = await response.json();
console.log(data.message);
```

## 下一步

- [记忆管理示例](/zh/examples/memory)
- [技能创建示例](/zh/examples/skills)
- [文件操作示例](/zh/examples/files)
