# Contributing to Gradle Parser | 贡献指南

Thank you for your interest in contributing to Gradle Parser! This document provides guidelines and information for contributors.

感谢您对 Gradle Parser 项目的贡献兴趣！本文档为贡献者提供指南和信息。

## Table of Contents | 目录

- [Code of Conduct | 行为准则](#code-of-conduct--行为准则)
- [Getting Started | 开始贡献](#getting-started--开始贡献)
- [Development Setup | 开发环境设置](#development-setup--开发环境设置)
- [Contributing Process | 贡献流程](#contributing-process--贡献流程)
- [Coding Standards | 编码标准](#coding-standards--编码标准)
- [Testing Guidelines | 测试指南](#testing-guidelines--测试指南)
- [Documentation | 文档](#documentation--文档)
- [Submitting Changes | 提交更改](#submitting-changes--提交更改)

## Code of Conduct | 行为准则

This project adheres to a code of conduct. By participating, you are expected to uphold this code.

本项目遵循行为准则。参与项目时，您需要遵守这些准则。

- Be respectful and inclusive | 保持尊重和包容
- Focus on constructive feedback | 专注于建设性反馈
- Help others learn and grow | 帮助他人学习和成长
- Maintain professionalism | 保持专业性

## Getting Started | 开始贡献

### Prerequisites | 前置要求

- Go 1.19 or later | Go 1.19 或更高版本
- Git
- Basic understanding of Gradle build files | 对 Gradle 构建文件的基本了解

### Types of Contributions | 贡献类型

We welcome various types of contributions:

我们欢迎各种类型的贡献：

- 🐛 **Bug fixes** | 错误修复
- ✨ **New features** | 新功能
- 📚 **Documentation improvements** | 文档改进
- 🧪 **Test additions** | 测试添加
- 🎨 **Code quality improvements** | 代码质量改进
- 💡 **Examples and tutorials** | 示例和教程

## Development Setup | 开发环境设置

### 1. Fork and Clone | 分叉和克隆

```bash
# Fork the repository on GitHub
# 在 GitHub 上分叉仓库

# Clone your fork
# 克隆您的分叉
git clone https://github.com/YOUR_USERNAME/gradle-parser.git
cd gradle-parser

# Add upstream remote
# 添加上游远程仓库
git remote add upstream https://github.com/scagogogo/gradle-parser.git
```

### 2. Install Dependencies | 安装依赖

```bash
# Download Go modules
# 下载 Go 模块
go mod download

# Verify dependencies
# 验证依赖
go mod verify
```

### 3. Verify Setup | 验证设置

```bash
# Run tests
# 运行测试
go test ./...

# Run examples
# 运行示例
cd examples/01_basic && go run main.go

# Run linter
# 运行代码检查
golangci-lint run
```

## Contributing Process | 贡献流程

### 1. Create an Issue | 创建问题

Before starting work, create an issue to discuss:

开始工作前，创建问题进行讨论：

- Bug reports | 错误报告
- Feature requests | 功能请求
- Questions | 问题

### 2. Create a Branch | 创建分支

```bash
# Update main branch
# 更新主分支
git checkout main
git pull upstream main

# Create feature branch
# 创建功能分支
git checkout -b feature/your-feature-name
```

### 3. Make Changes | 进行更改

Follow the coding standards and make your changes:

遵循编码标准并进行更改：

- Write clean, readable code | 编写清洁、可读的代码
- Add tests for new functionality | 为新功能添加测试
- Update documentation | 更新文档
- Follow existing patterns | 遵循现有模式

### 4. Test Your Changes | 测试更改

```bash
# Run all tests
# 运行所有测试
go test ./...

# Run comprehensive test suite
# 运行综合测试套件
cd test && ./scripts/run-tests.sh

# Test examples
# 测试示例
cd examples && ./run-all-examples.sh

# Check code quality
# 检查代码质量
golangci-lint run
```

### 5. Commit Changes | 提交更改

Use conventional commit messages:

使用约定式提交消息：

```bash
# Format: type(scope): description
# 格式：类型(范围): 描述

git commit -m "feat(parser): add support for Kotlin DSL"
git commit -m "fix(api): handle empty dependency blocks"
git commit -m "docs(readme): update installation instructions"
git commit -m "test(integration): add multi-module parsing tests"
```

**Commit Types | 提交类型:**
- `feat`: New feature | 新功能
- `fix`: Bug fix | 错误修复
- `docs`: Documentation | 文档
- `test`: Tests | 测试
- `refactor`: Code refactoring | 代码重构
- `style`: Code style changes | 代码风格更改
- `perf`: Performance improvements | 性能改进

## Coding Standards | 编码标准

### Go Code Style | Go 代码风格

- Follow `gofmt` and `goimports` | 遵循 `gofmt` 和 `goimports`
- Use meaningful variable names | 使用有意义的变量名
- Write clear comments for exported functions | 为导出函数编写清晰注释
- Keep functions focused and small | 保持函数专注和简洁
- Handle errors appropriately | 适当处理错误

### Example | 示例

```go
// ParseGradleFile parses a Gradle build file and returns project information.
// ParseGradleFile 解析 Gradle 构建文件并返回项目信息。
func ParseGradleFile(filePath string) (*model.Project, error) {
    if filePath == "" {
        return nil, fmt.Errorf("file path cannot be empty")
    }
    
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
    }
    
    return parseContent(string(content))
}
```

### Documentation | 文档注释

- Document all exported functions and types | 为所有导出的函数和类型编写文档
- Use English for code comments | 代码注释使用英文
- Include examples in documentation | 在文档中包含示例

## Testing Guidelines | 测试指南

### Test Structure | 测试结构

```go
func TestFunctionName(t *testing.T) {
    // Arrange | 准备
    input := "test input"
    expected := "expected output"
    
    // Act | 执行
    result, err := FunctionToTest(input)
    
    // Assert | 断言
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test Categories | 测试类别

- **Unit tests**: Test individual functions | 单元测试：测试单个函数
- **Integration tests**: Test component interactions | 集成测试：测试组件交互
- **Example tests**: Verify examples work | 示例测试：验证示例工作

### Coverage Requirements | 覆盖率要求

- Aim for >90% test coverage | 目标 >90% 测试覆盖率
- Test both success and error cases | 测试成功和错误情况
- Include edge cases | 包含边缘情况

## Documentation | 文档

### Types of Documentation | 文档类型

1. **Code comments** | 代码注释
2. **API documentation** | API 文档
3. **User guides** | 用户指南
4. **Examples** | 示例

### Documentation Standards | 文档标准

- Write clear, concise documentation | 编写清晰、简洁的文档
- Include code examples | 包含代码示例
- Update documentation with code changes | 随代码更改更新文档
- Support both English and Chinese | 支持英文和中文

## Submitting Changes | 提交更改

### Pull Request Process | 拉取请求流程

1. **Push your branch** | 推送分支
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create Pull Request** | 创建拉取请求
   - Use descriptive title | 使用描述性标题
   - Reference related issues | 引用相关问题
   - Describe changes made | 描述所做更改
   - Include testing information | 包含测试信息

3. **PR Template** | PR 模板
   ```markdown
   ## Description | 描述
   Brief description of changes
   
   ## Type of Change | 更改类型
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Documentation update
   - [ ] Test addition
   
   ## Testing | 测试
   - [ ] Tests pass locally
   - [ ] Added new tests
   - [ ] Updated documentation
   
   ## Checklist | 检查清单
   - [ ] Code follows style guidelines
   - [ ] Self-review completed
   - [ ] Documentation updated
   ```

### Review Process | 审查流程

- Maintainers will review your PR | 维护者将审查您的 PR
- Address feedback promptly | 及时处理反馈
- Keep PR focused and small | 保持 PR 专注和小型化
- Be patient and respectful | 保持耐心和尊重

## Getting Help | 获取帮助

If you need help:

如果您需要帮助：

- 💬 **Discussions**: [GitHub Discussions](https://github.com/scagogogo/gradle-parser/discussions)
- 🐛 **Issues**: [GitHub Issues](https://github.com/scagogogo/gradle-parser/issues)
- 📧 **Email**: Contact maintainers directly

## Recognition | 致谢

Contributors will be recognized in:

贡献者将在以下地方得到认可：

- README.md contributors section | README.md 贡献者部分
- Release notes | 发布说明
- Project documentation | 项目文档

Thank you for contributing to Gradle Parser! 🎉

感谢您对 Gradle Parser 的贡献！🎉
