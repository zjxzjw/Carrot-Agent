---
layout: home

hero:
  name: "Carrot Agent"
  text: "Intelligent Agent Framework"
  tagline: "Persistent memory, automatic skill learning, and powerful tool calling - all in a containerized Go application"
  image:
    src: /logo.png
    alt: Carrot Agent Logo
  actions:
    - theme: brand
      text: Get Started
      link: /guide/quick-start
    - theme: alt
      text: View on GitHub
      link: https://github.com/your-org/carrot-agent

features:
  - icon: 🧠
    title: Hierarchical Memory
    details: Three-tier memory system (snapshot, session, long-term) for persistent context across conversations
  - icon: 🎯
    title: Automatic Skill Learning
    details: Automatically generates reusable skills from complex workflows after completing tasks
  - icon: 🔧
    title: Tool Calling
    details: Built-in tools for file operations, HTTP requests, memory management, and more
  - icon: 🐳
    title: Container-First
    details: Official Docker images with secure isolation and one-click deployment
  - icon: 🤖
    title: Multi-Model Support
    details: Works with OpenAI GPT, Claude, and other LLM providers through unified interface
  - icon: ⚡
    title: High Performance
    details: Built with Go for low resource consumption and fast execution
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
