# Basic Usage Examples

Common usage patterns for Carrot Agent.

## Example 1: Simple Chat

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello! What can you do?",
    "session_id": "example-1"
  }'
```

## Example 2: File Creation

Ask the agent to create a file:

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Create a Python script that prints Hello World",
    "session_id": "example-2"
  }'
```

The agent will use the `file_write` tool to create the script.

## Example 3: Memory Management

Save important information:

```bash
curl -X POST http://localhost:8080/api/memory \
  -H "Content-Type: application/json" \
  -d '{
    "type": "snapshot",
    "content": "User prefers Python over JavaScript",
    "metadata": "{\"category\": \"preference\"}"
  }'
```

## Example 4: Skill Listing

View available skills:

```bash
curl http://localhost:8080/api/skills
```

Response:
```json
{
  "skills": [
    {
      "id": "skill_123",
      "name": "csv_processor",
      "description": "Process CSV files and generate statistics"
    }
  ],
  "count": 1
}
```

## Example 5: Session Management

List all sessions:

```bash
curl http://localhost:8080/api/session/
```

Delete a session:

```bash
curl -X DELETE http://localhost:8080/api/session/example-1
```

## Example 6: Statistics

Get agent statistics:

```bash
curl http://localhost:8080/api/stats
```

Response:
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

## Python Example

Using Python requests library:

```python
import requests

url = "http://localhost:8080/api/chat"
payload = {
    "message": "What is the capital of France?",
    "session_id": "python-example"
}

response = requests.post(url, json=payload)
print(response.json()["message"])
```

## JavaScript Example

Using fetch API:

```javascript
const response = await fetch('http://localhost:8080/api/chat', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    message: 'Tell me a joke',
    session_id: 'js-example'
  })
});

const data = await response.json();
console.log(data.message);
```

## Next Steps

- [Memory Management Examples](/examples/memory)
- [Skill Creation Examples](/examples/skills)
- [File Operations Examples](/examples/files)
