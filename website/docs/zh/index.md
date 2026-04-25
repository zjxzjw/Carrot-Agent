---
layout: home

hero:
  name: "Carrot Agent"
  text: "智能代理框架"
  tagline: "持久记忆、自动技能学习和强大的工具调用 - 全部在容器化的 Go 应用中"
  image:
    src: /logo.png
    alt: Carrot Agent Logo
  actions:
    - theme: brand
      text: 开始使用
      link: /zh/guide/quick-start
    - theme: alt
      text: 在 GitHub 上查看
      link: https://github.com/zjxzjw/Carrot-Agent
    - theme: alt
      text: API 参考
      link: /zh/api/overview

features:
  - icon: 🧠
    title: 分层记忆
    details: 三层记忆系统（快照、会话、长期），实现跨对话的持久上下文
  - icon: 🎯
    title: 自动技能学习
    details: 完成任务后自动从复杂工作流中生成可重用技能
  - icon: 🔧
    title: 工具调用
    details: 内置文件操作、HTTP 请求、记忆管理、系统信息等工具
  - icon: 🐳
    title: 容器优先
    details: 官方 Docker 镜像，安全隔离，一键部署
  - icon: 🤖
    title: 多模型支持
    details: 通过统一接口支持 OpenAI GPT、Claude 等多种大语言模型
  - icon: ⚡
    title: 高性能
    details: 使用 Go 语言构建，资源占用低，执行速度快
  - icon: 🔒
    title: 安全操作
    details: 路径和 URL 验证，非 root 用户执行，最小权限原则
  - icon: 📊
    title: 会话管理
    details: 保存和加载对话，实现跨会话的连续上下文
  - icon: 🖥️
    title: Web 界面
    details: 内置基于 React 的 Web UI，方便管理和监控

---

## 🚀 为什么选择 Carrot Agent？

Carrot Agent 旨在成为当今最灵活、最强大的智能代理框架。无论您是构建个人助手、自动化复杂工作流，还是创建自定义 AI 解决方案，Carrot Agent 都能为您提供成功所需的工具和基础设施。

### 核心优势

- **持久记忆**：跨对话永不丢失上下文
- **自动技能学习**：持续提高性能
- **可扩展架构**：轻松添加自定义工具和集成
- **容器化部署**：简单安全的部署选项
- **多模型支持**：使用您首选的 LLM 提供商

## 🔧 使用场景

- **个人助手**：智能伴侣，记住您的偏好并从交互中学习
- **工作流自动化**：以最少的配置自动化复杂任务
- **知识库**：构建随时间增长的持久知识系统
- **API 集成**：与外部服务和系统连接
- **研究助手**：帮助数据收集、分析和综合

<style>
:root {
  --vp-home-hero-name-color: transparent;
  --vp-home-hero-name-background: -webkit-linear-gradient(120deg, #ff6b35 30%, #f7931e);
  --vp-home-hero-image-background-image: linear-gradient(-45deg, #ff6b3530 50%, #f7931e30 50%);
  --vp-home-hero-image-filter: blur(44px);
}

@media (min-width: 640px) {
  :root {
    --vp-home-hero-image-filter: blur(56px);
  }
}

@media (min-width: 960px) {
  :root {
    --vp-home-hero-image-filter: blur(68px);
  }
}

/* 添加自定义动画 */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.features {
  animation: fadeIn 0.6s ease-out;
}

.hero {
  animation: fadeIn 0.4s ease-out;
}
</style>
