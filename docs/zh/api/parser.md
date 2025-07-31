# 解析器 API

解析器 API 为高级用例提供低级解析接口和实现。虽然大多数用户应该使用[核心 API](./core.md)，但解析器 API 提供对解析行为的细粒度控制。

## 包导入

```go
import "github.com/scagogogo/gradle-parser/pkg/parser"
```

## Parser 接口

### Parser

所有实现必须满足的主要解析器接口。

```go
type Parser interface {
    Parse(content string) (*model.ParseResult, error)
    ParseFile(filePath string) (*model.ParseResult, error)
}
```

## GradleParser

Parser 接口的默认实现。

### 构造函数

```go
func NewParser() *GradleParser
```

**返回值:**
- `*GradleParser`: 具有默认配置的新解析器实例

**示例:**
```go
parser := parser.NewParser()
result, err := parser.Parse(gradleContent)
if err != nil {
    log.Fatal(err)
}
```

### 配置方法

GradleParser 支持方法链式配置：

#### WithSkipComments

控制解析期间是否处理注释。

```go
func (p *GradleParser) WithSkipComments(skip bool) *GradleParser
```

**参数:**
- `skip` (bool): 如果为 true，解析期间忽略注释

**示例:**
```go
parser := parser.NewParser().WithSkipComments(false)
```

## 源码感知解析器

对于需要精确源位置跟踪的应用程序：

### SourceAwareParser

跟踪所有解析元素源位置的扩展解析器。

```go
type SourceAwareParser struct {
    *GradleParser
}
```

### 构造函数

```go
func NewSourceAwareParser() *SourceAwareParser
```

**返回值:**
- `*SourceAwareParser`: 新的源码感知解析器实例

## 最佳实践

1. **使用适当的解析级别**: 根据需要在基础解析和源码感知解析之间选择
2. **为性能配置**: 禁用不必要的解析功能以获得更好的性能
3. **优雅地处理错误**: 始终检查错误和警告
4. **重用解析器实例**: 创建一次解析器并重用于多个文件
5. **内存管理**: 对于大规模处理，禁用原始内容收集
