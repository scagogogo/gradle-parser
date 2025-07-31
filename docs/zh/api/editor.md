# 编辑器 API

编辑器 API 为 Gradle 构建文件提供结构化编辑功能。它允许您以编程方式修改依赖、插件和属性，同时保持原始格式并最小化差异。

## 包导入

```go
import "github.com/scagogogo/gradle-parser/pkg/api"
```

## 快速编辑函数

### UpdateDependencyVersion

更新 Gradle 文件中的依赖版本。

```go
func UpdateDependencyVersion(filePath, group, name, newVersion string) (string, error)
```

**参数:**
- `filePath` (string): Gradle 文件路径
- `group` (string): 依赖组 ID
- `name` (string): 依赖构件名称
- `newVersion` (string): 要设置的新版本

**返回值:**
- `string`: 更新后的文件内容
- `error`: 更新失败时的错误

**示例:**
```go
// 更新 MySQL 连接器版本
newContent, err := api.UpdateDependencyVersion(
    "build.gradle", 
    "mysql", 
    "mysql-connector-java", 
    "8.0.31"
)
if err != nil {
    log.Fatal(err)
}

// 写回文件
err = os.WriteFile("build.gradle", []byte(newContent), 0644)
if err != nil {
    log.Fatal(err)
}
```

### UpdatePluginVersion

更新 Gradle 文件中的插件版本。

```go
func UpdatePluginVersion(filePath, pluginId, newVersion string) (string, error)
```

**参数:**
- `filePath` (string): Gradle 文件路径
- `pluginId` (string): 插件标识符
- `newVersion` (string): 要设置的新版本

**返回值:**
- `string`: 更新后的文件内容
- `error`: 更新失败时的错误

## 高级编辑器

### CreateGradleEditor

为更复杂的修改创建结构化编辑器。

```go
func CreateGradleEditor(filePath string) (*editor.GradleEditor, error)
```

**参数:**
- `filePath` (string): Gradle 文件路径

**返回值:**
- `*editor.GradleEditor`: 文件的编辑器实例
- `error`: 创建失败时的错误

**示例:**
```go
editor, err := api.CreateGradleEditor("build.gradle")
if err != nil {
    log.Fatal(err)
}

// 执行多个编辑
err = editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
if err != nil {
    log.Fatal(err)
}

err = editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
if err != nil {
    log.Fatal(err)
}
```

## GradleEditor 方法

### UpdateDependencyVersion

使用编辑器更新依赖版本。

```go
func (ge *GradleEditor) UpdateDependencyVersion(group, name, newVersion string) error
```

### AddDependency

向项目添加新依赖。

```go
func (ge *GradleEditor) AddDependency(group, name, version, scope string) error
```

**示例:**
```go
err = editor.AddDependency(
    "org.apache.commons",
    "commons-text", 
    "1.9",
    "implementation"
)
if err != nil {
    log.Fatal(err)
}
```

### ApplyModifications

应用所有修改并返回更新的内容。

```go
func (ge *GradleEditor) ApplyModifications() (string, error)
```

**示例:**
```go
// 进行多个更改
editor.UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31")
editor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
editor.AddDependency("org.apache.commons", "commons-text", "1.9", "implementation")

// 应用所有更改
newContent, err := editor.ApplyModifications()
if err != nil {
    log.Fatal(err)
}

// 写入文件
err = os.WriteFile("build.gradle", []byte(newContent), 0644)
if err != nil {
    log.Fatal(err)
}
```

## 最小差异保证

编辑器设计为进行最小更改：

- **保持格式**: 保持原始缩进和间距
- **保持注释**: 注释不会被修改或删除
- **最小更改**: 只更改正在更新的特定值
- **逐行**: 逐行进行更改以最小化差异大小

**示例:**
```gradle
// 原始
implementation 'mysql:mysql-connector-java:8.0.28'

// 在 UpdateDependencyVersion("mysql", "mysql-connector-java", "8.0.31") 之后
implementation 'mysql:mysql-connector-java:8.0.31'
```

只有版本号被更改，其他所有内容保持相同。
