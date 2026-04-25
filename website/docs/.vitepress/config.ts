import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Carrot Agent',
  description: 'An intelligent agent framework with persistent memory and skill learning',
  
  head: [
    ['link', { rel: 'icon', href: '/favicon.ico' }],
    ['meta', { name: 'viewport', content: 'width=device-width, initial-scale=1.0' }],
  ],

  locales: {
    root: {
      label: 'English',
      lang: 'en-US',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Guide', link: '/guide/introduction' },
          { text: 'API Reference', link: '/api/overview' },
          { text: 'Examples', link: '/examples/basic' },
          { 
            text: 'Links',
            items: [
              { text: 'GitHub', link: 'https://github.com/your-org/carrot-agent' },
              { text: 'Docker Hub', link: 'https://hub.docker.com/r/carrotagent/carrot-agent' },
            ]
          }
        ],
        sidebar: {
          '/guide/': [
            {
              text: 'Getting Started',
              items: [
                { text: 'Introduction', link: '/guide/introduction' },
                { text: 'Quick Start', link: '/guide/quick-start' },
                { text: 'Installation', link: '/guide/installation' },
              ]
            },
            {
              text: 'Core Concepts',
              items: [
                { text: 'Agent Architecture', link: '/guide/architecture' },
                { text: 'Memory System', link: '/guide/memory' },
                { text: 'Skill System', link: '/guide/skills' },
                { text: 'Tool Registry', link: '/guide/tools' },
              ]
            },
            {
              text: 'Deployment',
              items: [
                { text: 'Docker Deployment', link: '/guide/docker' },
                { text: 'Configuration', link: '/guide/configuration' },
                { text: 'Security', link: '/guide/security' },
              ]
            },
            {
              text: 'Advanced',
              items: [
                { text: 'Custom Tools', link: '/guide/custom-tools' },
                { text: 'Model Providers', link: '/guide/models' },
                { text: 'Best Practices', link: '/guide/best-practices' },
              ]
            }
          ],
          '/api/': [
            {
              text: 'API Reference',
              items: [
                { text: 'Overview', link: '/api/overview' },
                { text: 'Chat API', link: '/api/chat' },
                { text: 'Skills API', link: '/api/skills' },
                { text: 'Memory API', link: '/api/memory' },
                { text: 'Sessions API', link: '/api/sessions' },
                { text: 'Stats API', link: '/api/stats' },
              ]
            }
          ],
          '/examples/': [
            {
              text: 'Examples',
              items: [
                { text: 'Basic Usage', link: '/examples/basic' },
                { text: 'Memory Management', link: '/examples/memory' },
                { text: 'Skill Creation', link: '/examples/skills' },
                { text: 'File Operations', link: '/examples/files' },
              ]
            }
          ]
        },
        editLink: {
          pattern: 'https://github.com/your-org/carrot-agent/edit/main/website/docs/:path',
          text: 'Edit this page on GitHub'
        },
        footer: {
          message: 'Released under the MIT License.'
        }
      }
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      link: '/zh/',
      themeConfig: {
        nav: [
          { text: '首页', link: '/zh/' },
          { text: '指南', link: '/zh/guide/introduction' },
          { text: 'API 参考', link: '/zh/api/overview' },
          { text: '示例', link: '/zh/examples/basic' },
          { 
            text: '链接',
            items: [
              { text: 'GitHub', link: 'https://github.com/your-org/carrot-agent' },
              { text: 'Docker Hub', link: 'https://hub.docker.com/r/carrotagent/carrot-agent' },
            ]
          }
        ],
        sidebar: {
          '/zh/guide/': [
            {
              text: '快速开始',
              items: [
                { text: '介绍', link: '/zh/guide/introduction' },
                { text: '快速上手', link: '/zh/guide/quick-start' },
                { text: '安装指南', link: '/zh/guide/installation' },
              ]
            },
            {
              text: '核心概念',
              items: [
                { text: '代理架构', link: '/zh/guide/architecture' },
                { text: '记忆系统', link: '/zh/guide/memory' },
                { text: '技能系统', link: '/zh/guide/skills' },
                { text: '工具注册表', link: '/zh/guide/tools' },
              ]
            },
            {
              text: '部署',
              items: [
                { text: 'Docker 部署', link: '/zh/guide/docker' },
                { text: '配置说明', link: '/zh/guide/configuration' },
                { text: '安全指南', link: '/zh/guide/security' },
              ]
            },
            {
              text: '进阶',
              items: [
                { text: '自定义工具', link: '/zh/guide/custom-tools' },
                { text: '模型提供商', link: '/zh/guide/models' },
                { text: '最佳实践', link: '/zh/guide/best-practices' },
              ]
            }
          ],
          '/zh/api/': [
            {
              text: 'API 参考',
              items: [
                { text: '概览', link: '/zh/api/overview' },
                { text: '聊天 API', link: '/zh/api/chat' },
                { text: '技能 API', link: '/zh/api/skills' },
                { text: '记忆 API', link: '/zh/api/memory' },
                { text: '会话 API', link: '/zh/api/sessions' },
                { text: '统计 API', link: '/zh/api/stats' },
              ]
            }
          ],
          '/zh/examples/': [
            {
              text: '示例',
              items: [
                { text: '基础用法', link: '/zh/examples/basic' },
                { text: '记忆管理', link: '/zh/examples/memory' },
                { text: '技能创建', link: '/zh/examples/skills' },
                { text: '文件操作', link: '/zh/examples/files' },
              ]
            }
          ]
        },
        editLink: {
          pattern: 'https://github.com/your-org/carrot-agent/edit/main/website/docs/:path',
          text: '在 GitHub 上编辑此页'
        },
        footer: {
          message: '基于 MIT 许可证发布'
        }
      }
    }
  },

  themeConfig: {
    logo: '/logo.png',
    
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/introduction' },
      { text: 'API Reference', link: '/api/overview' },
      { text: 'Examples', link: '/examples/basic' },
      { 
        text: 'Links',
        items: [
          { text: 'GitHub', link: 'https://github.com/your-org/carrot-agent' },
          { text: 'Docker Hub', link: 'https://hub.docker.com/r/carrotagent/carrot-agent' },
        ]
      }
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Getting Started',
          items: [
            { text: 'Introduction', link: '/guide/introduction' },
            { text: 'Quick Start', link: '/guide/quick-start' },
            { text: 'Installation', link: '/guide/installation' },
          ]
        },
        {
          text: 'Core Concepts',
          items: [
            { text: 'Agent Architecture', link: '/guide/architecture' },
            { text: 'Memory System', link: '/guide/memory' },
            { text: 'Skill System', link: '/guide/skills' },
            { text: 'Tool Registry', link: '/guide/tools' },
          ]
        },
        {
          text: 'Deployment',
          items: [
            { text: 'Docker Deployment', link: '/guide/docker' },
            { text: 'Configuration', link: '/guide/configuration' },
            { text: 'Security', link: '/guide/security' },
          ]
        },
        {
          text: 'Advanced',
          items: [
            { text: 'Custom Tools', link: '/guide/custom-tools' },
            { text: 'Model Providers', link: '/guide/models' },
            { text: 'Best Practices', link: '/guide/best-practices' },
          ]
        }
      ],
      '/api/': [
        {
          text: 'API Reference',
          items: [
            { text: 'Overview', link: '/api/overview' },
            { text: 'Chat API', link: '/api/chat' },
            { text: 'Skills API', link: '/api/skills' },
            { text: 'Memory API', link: '/api/memory' },
            { text: 'Sessions API', link: '/api/sessions' },
            { text: 'Stats API', link: '/api/stats' },
          ]
        }
      ],
      '/examples/': [
        {
          text: 'Examples',
          items: [
            { text: 'Basic Usage', link: '/examples/basic' },
            { text: 'Memory Management', link: '/examples/memory' },
            { text: 'Skill Creation', link: '/examples/skills' },
            { text: 'File Operations', link: '/examples/files' },
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/your-org/carrot-agent' }
    ],

    search: {
      provider: 'local'
    }
  },

  markdown: {
    lineNumbers: true,
  }
})
