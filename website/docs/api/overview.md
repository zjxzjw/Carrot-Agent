# API Overview

Carrot Agent provides a RESTful API for interacting with the agent programmatically.

## Base URL

```
http://localhost:8080/api
```

## Authentication

Currently, the API does not require authentication. For production use, consider adding an API gateway or reverse proxy with authentication.

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/chat` | Send message and get response |
| GET | `/skills` | List all skills |
| POST | `/skills` | Create a new skill |
| GET | `/memory` | List memories |
| POST | `/memory` | Add memory |
| GET | `/session/` | List sessions |
| GET | `/session/{id}` | Get session details |
| DELETE | `/session/{id}` | Delete session |
| GET | `/stats` | Get agent statistics |

## Response Format

All responses are in JSON format:

```json
{
  "status": "success",
  "data": {...}
}
```

Error responses:

```json
{
  "error": "Error message",
  "code": 400
}
```

## Rate Limiting

Currently no rate limiting is implemented. For production deployments, add rate limiting at the reverse proxy level.

## Examples

See individual endpoint documentation for detailed examples:

- [Chat API](/api/chat)
- [Skills API](/api/skills)
- [Memory API](/api/memory)
- [Sessions API](/api/sessions)
- [Stats API](/api/stats)
