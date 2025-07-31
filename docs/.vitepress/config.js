import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Gradle Parser',
  description: 'A powerful Gradle build file parser for Go',
  
  // 多语言配置
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'Gradle Parser',
      description: 'A powerful Gradle build file parser for Go',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Guide', link: '/guide/getting-started' },
          { text: 'API Reference', link: '/api/' },
          { text: 'Examples', link: '/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/gradle-parser' }
        ],
        sidebar: {
          '/guide/': [
            {
              text: 'Guide',
              items: [
                { text: 'Getting Started', link: '/guide/getting-started' },
                { text: 'Basic Usage', link: '/guide/basic-usage' },
                { text: 'Advanced Features', link: '/guide/advanced-features' },
                { text: 'Structured Editing', link: '/guide/structured-editing' },
                { text: 'Configuration', link: '/guide/configuration' }
              ]
            }
          ],
          '/api/': [
            {
              text: 'API Reference',
              items: [
                { text: 'Overview', link: '/api/' },
                { text: 'Core API', link: '/api/core' },
                { text: 'Parser', link: '/api/parser' },
                { text: 'Models', link: '/api/models' },
                { text: 'Editor', link: '/api/editor' },
                { text: 'Utilities', link: '/api/utilities' }
              ]
            }
          ],
          '/examples/': [
            {
              text: 'Examples',
              items: [
                { text: 'Overview', link: '/examples/' },
                { text: 'Basic Parsing', link: '/examples/basic-parsing' },
                { text: 'Dependency Analysis', link: '/examples/dependency-analysis' },
                { text: 'Plugin Detection', link: '/examples/plugin-detection' },
                { text: 'Repository Parsing', link: '/examples/repository-parsing' },
                { text: 'Structured Editing', link: '/examples/structured-editing' },
                { text: 'Custom Parser', link: '/examples/custom-parser' }
              ]
            }
          ]
        },
        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/gradle-parser' }
        ],
        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2024 Gradle Parser Contributors'
        }
      }
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'Gradle Parser',
      description: '强大的 Go 语言 Gradle 构建文件解析器',
      themeConfig: {
        nav: [
          { text: '首页', link: '/zh/' },
          { text: '指南', link: '/zh/guide/getting-started' },
          { text: 'API 参考', link: '/zh/api/' },
          { text: '示例', link: '/zh/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/gradle-parser' }
        ],
        sidebar: {
          '/zh/guide/': [
            {
              text: '指南',
              items: [
                { text: '快速开始', link: '/zh/guide/getting-started' },
                { text: '基本用法', link: '/zh/guide/basic-usage' },
                { text: '高级功能', link: '/zh/guide/advanced-features' },
                { text: '结构化编辑', link: '/zh/guide/structured-editing' },
                { text: '配置选项', link: '/zh/guide/configuration' }
              ]
            }
          ],
          '/zh/api/': [
            {
              text: 'API 参考',
              items: [
                { text: '概览', link: '/zh/api/' },
                { text: '核心 API', link: '/zh/api/core' },
                { text: '解析器', link: '/zh/api/parser' },
                { text: '数据模型', link: '/zh/api/models' },
                { text: '编辑器', link: '/zh/api/editor' },
                { text: '工具函数', link: '/zh/api/utilities' }
              ]
            }
          ],
          '/zh/examples/': [
            {
              text: '示例',
              items: [
                { text: '概览', link: '/zh/examples/' },
                { text: '基础解析', link: '/zh/examples/basic-parsing' },
                { text: '依赖分析', link: '/zh/examples/dependency-analysis' },
                { text: '插件检测', link: '/zh/examples/plugin-detection' },
                { text: '仓库解析', link: '/zh/examples/repository-parsing' },
                { text: '结构化编辑', link: '/zh/examples/structured-editing' },
                { text: '自定义解析器', link: '/zh/examples/custom-parser' }
              ]
            }
          ]
        },
        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/gradle-parser' }
        ],
        footer: {
          message: '基于 MIT 许可证发布。',
          copyright: 'Copyright © 2024 Gradle Parser 贡献者'
        },
        docFooter: {
          prev: '上一页',
          next: '下一页'
        },
        outline: {
          label: '页面导航'
        },
        lastUpdated: {
          text: '最后更新于',
          formatOptions: {
            dateStyle: 'short',
            timeStyle: 'medium'
          }
        },
        langMenuLabel: '多语言',
        returnToTopLabel: '回到顶部',
        sidebarMenuLabel: '菜单',
        darkModeSwitchLabel: '主题',
        lightModeSwitchTitle: '切换到浅色模式',
        darkModeSwitchTitle: '切换到深色模式'
      }
    }
  },

  // 主题配置
  themeConfig: {
    logo: '/logo.svg',
    search: {
      provider: 'local'
    }
  },

  // 构建配置
  base: '/gradle-parser/',
  outDir: '.vitepress/dist',
  ignoreDeadLinks: true,
  
  // Markdown 配置
  markdown: {
    lineNumbers: true,
    theme: {
      light: 'github-light',
      dark: 'github-dark'
    }
  },

  // 头部配置
  head: [
    ['link', { rel: 'icon', href: '/gradle-parser/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#3c8772' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'en' }],
    ['meta', { property: 'og:title', content: 'Gradle Parser | A powerful Gradle build file parser for Go' }],
    ['meta', { property: 'og:site_name', content: 'Gradle Parser' }],
    ['meta', { property: 'og:image', content: 'https://scagogogo.github.io/gradle-parser/og-image.png' }],
    ['meta', { property: 'og:url', content: 'https://scagogogo.github.io/gradle-parser/' }]
  ]
})
