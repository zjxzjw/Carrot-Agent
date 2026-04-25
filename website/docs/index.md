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
      link: https://github.com/zjxzjw/Carrot-Agent
    - theme: alt
      text: API Reference
      link: /api/overview

features:
  - icon: 🧠
    title: Hierarchical Memory
    details: Three-tier memory system (snapshot, session, long-term) for persistent context across conversations
  - icon: 🎯
    title: Automatic Skill Learning
    details: Automatically generates reusable skills from complex workflows after completing tasks
  - icon: 🔧
    title: Tool Calling
    details: Built-in tools for file operations, HTTP requests, memory management, system information, and more
  - icon: 🐳
    title: Container-First
    details: Official Docker images with secure isolation and one-click deployment
  - icon: 🤖
    title: Multi-Model Support
    details: Works with OpenAI GPT, Claude, and other LLM providers through unified interface
  - icon: ⚡
    title: High Performance
    details: Built with Go for low resource consumption and fast execution
  - icon: 🔒
    title: Secure Operation
    details: Path and URL validation, non-root user execution, and least privilege principle
  - icon: 📊
    title: Session Management
    details: Save and load conversations for continuous context across sessions
  - icon: 🖥️
    title: Web Interface
    details: Built-in React-based web UI for easy management and monitoring

---

## 🚀 Why Carrot Agent?

Carrot Agent is designed to be the most flexible and powerful intelligent agent framework available today. Whether you're building a personal assistant, automating complex workflows, or creating a custom AI solution, Carrot Agent provides the tools and infrastructure you need to succeed.

### Key Benefits

- **Persistent Memory**: Never lose context across conversations
- **Automatic Skill Learning**: Continuously improve performance
- **Extensible Architecture**: Easy to add custom tools and integrations
- **Containerized Deployment**: Simple and secure deployment options
- **Multi-Model Support**: Work with your preferred LLM provider

## 🔧 Use Cases

- **Personal Assistant**: Intelligent companion that remembers your preferences and learns from interactions
- **Workflow Automation**: Automate complex tasks with minimal configuration
- **Knowledge Base**: Build a persistent knowledge system that grows over time
- **API Integration**: Connect with external services and systems
- **Research Assistant**: Help with data collection, analysis, and synthesis

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

/* Add custom animations */
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
