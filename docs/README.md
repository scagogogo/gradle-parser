# Gradle Parser Documentation

This directory contains the documentation for the Gradle Parser project, built with [VitePress](https://vitepress.dev/).

## 📚 Documentation Structure

```
docs/
├── .vitepress/          # VitePress configuration
│   └── config.js        # Main configuration file
├── public/              # Static assets
│   ├── logo.svg         # Project logo
│   └── favicon.ico      # Site favicon
├── guide/               # English guides
│   ├── getting-started.md
│   ├── basic-usage.md
│   ├── advanced-features.md
│   ├── structured-editing.md
│   └── configuration.md
├── zh/                  # Chinese documentation
│   ├── guide/           # Chinese guides
│   └── api/             # Chinese API docs
├── api/                 # English API documentation
│   ├── index.md         # API overview
│   ├── core.md          # Core API
│   ├── models.md        # Data models
│   ├── editor.md        # Editor API
│   ├── parser.md        # Parser API
│   └── utilities.md     # Utility functions
├── examples/            # Code examples
│   ├── index.md         # Examples overview
│   ├── basic-parsing.md
│   ├── dependency-analysis.md
│   └── ...
├── index.md             # English homepage
└── package.json         # Node.js dependencies
```

## 🚀 Quick Start

### Prerequisites

- Node.js 18 or higher
- npm or yarn

### Local Development

1. **Install dependencies:**
   ```bash
   cd docs
   npm install
   ```

2. **Start development server:**
   ```bash
   npm run dev
   ```

3. **Open in browser:**
   Visit http://localhost:5173/gradle-parser/

### Build for Production

```bash
npm run build
```

The built files will be in `.vitepress/dist/`.

### Preview Production Build

```bash
npm run preview
```

## 📝 Writing Documentation

### Adding New Pages

1. Create a new `.md` file in the appropriate directory
2. Add the page to the sidebar configuration in `.vitepress/config.js`
3. Use proper frontmatter if needed

### Markdown Features

VitePress supports enhanced Markdown with:

- **Code syntax highlighting**
- **Custom containers** (tip, warning, danger, etc.)
- **Math expressions** with KaTeX
- **Mermaid diagrams**
- **Vue components** in Markdown

### Code Examples

Use proper language tags for syntax highlighting:

````markdown
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Gradle Parser!")
}
```
````

### Custom Containers

```markdown
::: tip
This is a tip container.
:::

::: warning
This is a warning container.
:::

::: danger
This is a danger container.
:::
```

## 🌐 Multi-language Support

The documentation supports both English and Chinese:

- **English**: Root level (`/`)
- **Chinese**: Under `/zh/` prefix

### Adding Translations

1. Create corresponding files in the `zh/` directory
2. Update the navigation and sidebar in `.vitepress/config.js`
3. Maintain the same structure as English docs

## 🔧 Configuration

### VitePress Configuration

Main configuration is in `.vitepress/config.js`:

- **Site metadata** (title, description)
- **Navigation menus**
- **Sidebar structure**
- **Theme customization**
- **Multi-language settings**

### Customizing Theme

You can customize the theme by:

1. Modifying `.vitepress/config.js`
2. Adding custom CSS in `.vitepress/theme/`
3. Creating custom Vue components

## 🚀 Deployment

### GitHub Pages (Automatic)

The documentation is automatically deployed to GitHub Pages when:

1. Changes are pushed to the `main` branch
2. Files in the `docs/` directory are modified
3. The GitHub Actions workflow completes successfully

### Manual Deployment

1. Build the documentation:
   ```bash
   npm run build
   ```

2. Deploy the `.vitepress/dist/` directory to your hosting provider

## 🔍 Link Validation

The project includes automatic link validation:

```bash
npm run lint-links
```

This checks for:
- Broken internal links
- Missing files
- Invalid references

## 📊 Analytics and SEO

The documentation includes:

- **SEO optimization** with proper meta tags
- **Open Graph** tags for social sharing
- **Structured data** for search engines
- **Sitemap generation** (automatic)

## 🛠️ Development Tips

### Hot Reload

The development server supports hot reload for:
- Markdown content changes
- Configuration updates
- Theme modifications

### Debugging

1. Check the browser console for errors
2. Verify file paths and links
3. Test with `npm run build` before deploying

### Performance

- Optimize images in the `public/` directory
- Use appropriate image formats (WebP, AVIF)
- Minimize custom CSS and JavaScript

## 📚 Resources

- [VitePress Documentation](https://vitepress.dev/)
- [Markdown Guide](https://www.markdownguide.org/)
- [Vue.js Documentation](https://vuejs.org/)

## 🤝 Contributing

When contributing to documentation:

1. Follow the existing structure and style
2. Test locally before submitting
3. Update both English and Chinese versions when applicable
4. Ensure all links work correctly
5. Add examples for new features

## 📄 License

This documentation is part of the Gradle Parser project and is licensed under the MIT License.
