# Advanced Features Example / 高级功能示例

This example demonstrates advanced features of the Gradle Parser, including source mapping, custom configurations, and performance optimization.

本示例演示了 Gradle Parser 的高级功能，包括源码映射、自定义配置和性能优化。

## Features Demonstrated / 演示功能

### English
- **Source-aware parsing** - Track exact source locations of parsed elements
- **Custom parser configuration** - Optimize parsing for specific use cases
- **Performance benchmarking** - Compare different parsing configurations
- **Multi-file analysis** - Analyze multiple Gradle files efficiently
- **Error handling strategies** - Robust error handling and recovery
- **Memory optimization** - Minimize memory usage for large projects

### 中文
- **源码感知解析** - 跟踪解析元素的精确源位置
- **自定义解析器配置** - 为特定用例优化解析
- **性能基准测试** - 比较不同的解析配置
- **多文件分析** - 高效分析多个 Gradle 文件
- **错误处理策略** - 健壮的错误处理和恢复
- **内存优化** - 为大型项目最小化内存使用

## Running the Example / 运行示例

```bash
cd examples/07_advanced
go run main.go
```

## Key Concepts / 关键概念

### Source Mapping / 源码映射

Source mapping allows you to track the exact location of every parsed element in the original file, enabling precise modifications and detailed analysis.

源码映射允许您跟踪原始文件中每个解析元素的确切位置，实现精确修改和详细分析。

### Performance Optimization / 性能优化

Different parser configurations can significantly impact performance. This example shows how to:
- Skip unnecessary parsing features
- Optimize memory usage
- Batch process multiple files

不同的解析器配置会显著影响性能。本示例展示如何：
- 跳过不必要的解析功能
- 优化内存使用
- 批量处理多个文件

### Custom Configuration / 自定义配置

The parser can be configured for specific use cases:
- **Fast parsing** - Skip comments and raw content collection
- **Memory efficient** - Minimize memory footprint
- **Complete analysis** - Parse all available information

解析器可以为特定用例进行配置：
- **快速解析** - 跳过注释和原始内容收集
- **内存高效** - 最小化内存占用
- **完整分析** - 解析所有可用信息

## Sample Output / 示例输出

```
🔍 Advanced Gradle Parser Features Demo
========================================

📊 Performance Benchmark:
  Fast Parser: 2.5ms (dependencies only)
  Standard Parser: 8.1ms (full parsing)
  Memory Optimized: 5.2ms (reduced memory)

📍 Source Mapping Analysis:
  Dependencies found at:
    - Line 32: org.springframework.boot:spring-boot-starter-web
    - Line 36: mysql:mysql-connector-java:8.0.29
    - Line 39: org.apache.commons:commons-lang3:3.12.0

🔧 Custom Analysis:
  Project Type: Spring Boot + Java
  Total Dependencies: 8
  Outdated Dependencies: 2
  Security Issues: 0

💾 Memory Usage:
  Standard Parsing: 2.1 MB
  Optimized Parsing: 0.8 MB
  Memory Saved: 62%
```

## Code Structure / 代码结构

The example is organized into several functions:

示例代码组织为几个函数：

- `demonstrateSourceMapping()` - Source location tracking / 源位置跟踪
- `benchmarkConfigurations()` - Performance comparison / 性能比较
- `analyzeMultipleFiles()` - Batch processing / 批量处理
- `customAnalysis()` - Advanced analysis / 高级分析
- `memoryOptimization()` - Memory usage optimization / 内存使用优化

## Learning Objectives / 学习目标

After running this example, you will understand:

运行此示例后，您将了解：

1. How to use source mapping for precise editing / 如何使用源码映射进行精确编辑
2. How to configure the parser for optimal performance / 如何配置解析器以获得最佳性能
3. How to handle large-scale Gradle file processing / 如何处理大规模 Gradle 文件处理
4. How to implement custom analysis logic / 如何实现自定义分析逻辑
5. How to optimize memory usage / 如何优化内存使用

## Next Steps / 下一步

- Explore the [API documentation](../../docs/api/) for detailed reference
- Check out [structured editing examples](../06_editor/) for file modification
- Review [performance optimization guide](../../docs/guide/configuration.md)

- 探索 [API 文档](../../docs/api/) 获取详细参考
- 查看 [结构化编辑示例](../06_editor/) 了解文件修改
- 查看 [性能优化指南](../../docs/guide/configuration.md)
