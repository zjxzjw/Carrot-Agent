# API 概览

Carrot Agent 提供 RESTful API 以便以编程方式与代理交互。

## Base URL

```
http://localhost:8080/api
```

## 认证

目前 API 不需要认证。对于生产环境，建议添加带有认证的 API 网关或反向代理。

## 端点列表

| 方法 | 端点 | 描述 |
|--------|----------|-------------|
| GET | `/health` | 健康检查 |
| POST | `/chat` | 发送消息并获取响应 |
| GET | `/skills` | 列出所有技能 |
| POST | `/skills` | 创建新技能 |
| GET | `/memory` | 列出记忆 |
| POST | `/memory` | 添加记忆 |
| GET | `/session/` | 列出会话 |
| GET | `/session/{id}` | 获取会话详情 |
| DELETE | `/session/{id}` | 删除会话 |
| GET | `/stats` | 获取代理统计信息 |

## 响应格式

所有响应均为 JSON 格式：

```json
{
  "status": "success",
  "data": {...}
}
```

错误响应：

```json
{
  "error": "错误信息",
  "code": 400
}
```

## 限流

目前未实现限流。对于生产部署，建议在反向代理层添加限流。

## 示例

查看各个端点的详细文档：

- [聊天 API](/zh/api/chat)
- [技能 API](/zh/api/skills)
- [记忆 API](/zh/api/memory)
- [会话 API](/zh/api/sessions)
- [统计 API](/zh/api/stats)
