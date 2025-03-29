# Gradle-Parser 示例程序

本目录包含一系列示例程序，展示如何使用gradle-parser库的各种功能。

## 示例结构

- [01_basic](01_basic/) - 基本用法示例
- [02_dependencies](02_dependencies/) - 依赖提取和分析示例
- [03_plugins](03_plugins/) - 插件提取和分析示例
- [04_repositories](04_repositories/) - 仓库提取和分析示例
- [05_complete](05_complete/) - 完整功能和定制选项示例

每个示例目录中都包含一个README.md文件，详细介绍了该示例的使用方法。

## 示例文件

[sample_files](sample_files/) 目录包含了用于测试的示例Gradle文件：

### 基本文件
- `build.gradle` - 标准Groovy DSL格式的Gradle文件示例
- `build.gradle.kts` - Kotlin DSL格式的Gradle文件示例
- `settings.gradle` - Gradle设置文件示例

### 多模块项目
- `app/build.gradle` - 应用模块的构建文件（Groovy DSL）
- `common/build.gradle` - 通用库模块的构建文件（Groovy DSL）
- `data/build.gradle.kts` - 数据模块的构建文件（Kotlin DSL）

这些文件可以作为一个完整的多模块Gradle项目结构示例，也可以单独用于测试解析功能。

## 运行示例

所有示例现在都使用硬编码参数，使您可以直接运行示例而无需提供命令行参数。每个示例都已配置为使用相对路径（`../sample_files/build.gradle`）来引用示例文件：

```bash
cd 01_basic
go run main.go
```

如果您想修改示例的行为，请直接编辑各示例目录下的 `main.go` 文件中标记有 `// MODIFY HERE` 的参数。

### 示例列表

1. **01_basic**: 基本解析功能，展示如何解析Gradle文件并提取基本信息
2. **02_dependencies**: 专注于依赖提取和分析
3. **03_plugins**: 专注于插件提取和项目类型检测
4. **04_repositories**: 专注于仓库配置提取和分析
5. **05_complete**: 展示完整功能，包括项目分析和高级配置

### 示例文件

在 `sample_files` 目录中提供了多种Gradle文件用于测试：
- `build.gradle`: 标准Groovy DSL格式的Gradle构建文件
- `build.gradle.kts`: Kotlin DSL格式的Gradle构建文件
- `settings.gradle`: Gradle设置文件
- 多模块项目示例 (`app/`, `common/`, `data/`)

## 使用方法

要运行任意示例，进入对应目录，并按照其README.md中的指导运行：

```bash
cd 01_basic
go run main.go
```

测试多模块项目解析：

```bash
cd 05_complete
go run main.go
```

## 导入方式

所有示例代码现在都直接导入 `pkg/api` 包：

```go
import (
    "github.com/scagogogo/gradle-parser/pkg/api"
)

// 然后使用 api 包中的函数
result, err := api.ParseFile(filePath)
```

这样做使得代码结构更加清晰和模块化，建议在自己的项目中也使用这种导入方式。 