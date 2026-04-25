---
layout: home

hero:
  name: "Carrot Agent"
  text: "智能代理框架"
  tagline: "持久化记忆、自动技能学习和强大的工具调用 - 全部集成在容器化的 Go 应用中"
  image:
    src: /logo.png
    alt: Carrot Agent Logo
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/quick-start
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/zjxzjw/Carrot-Agent

features:
  - icon: 🧠
    title: 分层记忆系统
    details: 三层记忆架构（快照、会话、长期）确保跨对话的持久化上下文
  - icon: 🎯
    title: 自动技能学习
    details: 完成任务后自动生成可复用的技能，让代理越来越智能
  - icon: 🔧
    title: 工具调用
    details: 内置文件操作、HTTP 请求、记忆管理、系统信息等强大工具
  - icon: 🐳
    title: 容器优先设计
    details: 官方 Docker 镜像，安全隔离，一键部署
  - icon: 🤖
    title: 多模型支持
    details: 支持 OpenAI GPT、Claude 等 LLM 提供商的统一接口
  - icon: ⚡
    title: 高性能
    details: 基于 Go 语言构建，低资源消耗，快速执行
  - icon: 🔒
    title: 安全操作
    details: 路径和 URL 验证，非 root 用户执行，最小权限原则
  - icon: 📊
    title: 会话管理
    details: 保存和加载对话，确保跨会话的连续上下文
  - icon: 🖥️
    title: Web 界面
    details: 内置基于 React 的 Web UI，方便管理和监控
---

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
</style>
