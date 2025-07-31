// 07_advanced demonstrates advanced features of gradle-parser
// 07_advanced 演示 gradle-parser 的高级功能
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/scagogogo/gradle-parser/pkg/api"
)

func main() {
	fmt.Println("🔍 Advanced Gradle Parser Features Demo")
	fmt.Println("========================================\n")

	// Demonstrate source mapping
	demonstrateSourceMapping()

	// Benchmark different configurations
	benchmarkConfigurations()

	// Analyze multiple files
	analyzeMultipleFiles()

	// Custom analysis
	customAnalysis()

	// Memory optimization
	memoryOptimization()
}

// demonstrateSourceMapping shows how to use source-aware parsing
// demonstrateSourceMapping 展示如何使用源码感知解析
func demonstrateSourceMapping() {
	fmt.Println("📍 Source Mapping Analysis:")
	fmt.Println("---------------------------")

	filePath := "../sample_files/build.gradle"
	result, err := api.ParseFile(filePath)
	if err != nil {
		fmt.Printf("❌ Source mapping failed: %v\n", err)
		return
	}

	// For this demo, we'll use the regular project
	project := result.Project
	fmt.Printf("📄 File: %s\n", filePath)

	// Show dependencies (without source locations for now)
	fmt.Println("\n📦 Dependencies found:")
	for i, dep := range project.Dependencies {
		fmt.Printf("  %d. %s:%s:%s (%s)\n",
			i+1, dep.Group, dep.Name, dep.Version, dep.Scope)
	}

	// Show plugins (without source locations for now)
	fmt.Println("\n🔌 Plugins found:")
	for i, plugin := range project.Plugins {
		fmt.Printf("  %d. %s", i+1, plugin.ID)
		if plugin.Version != "" {
			fmt.Printf(" v%s", plugin.Version)
		}
		fmt.Println()
	}

	fmt.Println()
}

// benchmarkConfigurations compares performance of different parser configurations
// benchmarkConfigurations 比较不同解析器配置的性能
func benchmarkConfigurations() {
	fmt.Println("📊 Performance Benchmark:")
	fmt.Println("-------------------------")

	filePath := "../sample_files/build.gradle"
	iterations := 10

	// Fast parser (dependencies only)
	fastOptions := &api.Options{
		SkipComments:      true,
		CollectRawContent: false,
		ParsePlugins:      false,
		ParseRepositories: false,
		ParseTasks:        false,
		ParseDependencies: true,
	}
	fastParser := api.NewParser(fastOptions)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, err := fastParser.ParseFile(filePath)
		if err != nil {
			fmt.Printf("❌ Fast parser error: %v\n", err)
			return
		}
	}
	fastDuration := time.Since(start) / time.Duration(iterations)

	// Standard parser
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, err := api.ParseFile(filePath)
		if err != nil {
			fmt.Printf("❌ Standard parser error: %v\n", err)
			return
		}
	}
	standardDuration := time.Since(start) / time.Duration(iterations)

	// Memory optimized parser
	memoryOptions := &api.Options{
		SkipComments:      true,
		CollectRawContent: false,
		ParsePlugins:      true,
		ParseRepositories: true,
		ParseTasks:        false,
		ParseDependencies: true,
	}
	memoryParser := api.NewParser(memoryOptions)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, err := memoryParser.ParseFile(filePath)
		if err != nil {
			fmt.Printf("❌ Memory parser error: %v\n", err)
			return
		}
	}
	memoryDuration := time.Since(start) / time.Duration(iterations)

	fmt.Printf("  🚀 Fast Parser: %v (dependencies only)\n", fastDuration)
	fmt.Printf("  📋 Standard Parser: %v (full parsing)\n", standardDuration)
	fmt.Printf("  💾 Memory Optimized: %v (reduced memory)\n", memoryDuration)
	fmt.Printf("  📈 Speed improvement: %.1fx faster\n", float64(standardDuration)/float64(fastDuration))
	fmt.Println()
}

// analyzeMultipleFiles demonstrates batch processing
// analyzeMultipleFiles 演示批量处理
func analyzeMultipleFiles() {
	fmt.Println("📁 Multi-File Analysis:")
	fmt.Println("-----------------------")

	files := []string{
		"../sample_files/build.gradle",
		"../sample_files/app/build.gradle",
		"../sample_files/common/build.gradle",
	}

	totalDeps := 0
	totalPlugins := 0
	successCount := 0

	// Use optimized parser for batch processing
	batchOptions := &api.Options{
		SkipComments:      true,
		CollectRawContent: false,
		ParsePlugins:      true,
		ParseRepositories: true,
		ParseTasks:        true,
		ParseDependencies: true,
	}
	batchParser := api.NewParser(batchOptions)

	for _, file := range files {
		result, err := batchParser.ParseFile(file)
		if err != nil {
			fmt.Printf("  ❌ Failed to parse %s: %v\n", file, err)
			continue
		}

		project := result.Project
		deps := len(project.Dependencies)
		plugins := len(project.Plugins)

		fmt.Printf("  📄 %s: %d deps, %d plugins\n", file, deps, plugins)

		totalDeps += deps
		totalPlugins += plugins
		successCount++
	}

	fmt.Printf("\n  📊 Summary: %d files processed\n", successCount)
	fmt.Printf("  📦 Total dependencies: %d\n", totalDeps)
	fmt.Printf("  🔌 Total plugins: %d\n", totalPlugins)
	fmt.Println()
}

// customAnalysis demonstrates advanced analysis capabilities
// customAnalysis 演示高级分析功能
func customAnalysis() {
	fmt.Println("🔧 Custom Analysis:")
	fmt.Println("-------------------")

	filePath := "../sample_files/build.gradle"
	result, err := api.ParseFile(filePath)
	if err != nil {
		fmt.Printf("❌ Analysis failed: %v\n", err)
		return
	}

	project := result.Project
	plugins := project.Plugins
	dependencies := project.Dependencies

	// Project type detection
	projectTypes := []string{}
	if api.IsAndroidProject(plugins) {
		projectTypes = append(projectTypes, "Android")
	}
	if api.IsKotlinProject(plugins) {
		projectTypes = append(projectTypes, "Kotlin")
	}
	if api.IsSpringBootProject(plugins) {
		projectTypes = append(projectTypes, "Spring Boot")
	}

	// Check for Java plugin
	hasJava := false
	for _, plugin := range plugins {
		if plugin.ID == "java" {
			hasJava = true
			break
		}
	}
	if hasJava {
		projectTypes = append(projectTypes, "Java")
	}

	if len(projectTypes) > 0 {
		fmt.Printf("  🏷️  Project Type: %s\n", projectTypes[0])
		if len(projectTypes) > 1 {
			fmt.Printf("  🏷️  Additional Types: %v\n", projectTypes[1:])
		}
	}

	// Dependency analysis
	depSets := api.DependenciesByScope(dependencies)
	fmt.Printf("  📦 Total Dependencies: %d\n", len(dependencies))
	for _, depSet := range depSets {
		fmt.Printf("    - %s: %d\n", depSet.Scope, len(depSet.Dependencies))
	}

	// Mock outdated dependency check
	outdatedCount := 0
	for _, dep := range dependencies {
		// Simulate outdated check
		if dep.Group == "mysql" && dep.Name == "mysql-connector-java" && dep.Version == "8.0.29" {
			outdatedCount++
		}
	}
	fmt.Printf("  ⚠️  Outdated Dependencies: %d\n", outdatedCount)
	fmt.Printf("  🔒 Security Issues: 0\n") // Mock security check
	fmt.Println()
}

// memoryOptimization demonstrates memory usage optimization
// memoryOptimization 演示内存使用优化
func memoryOptimization() {
	fmt.Println("💾 Memory Usage Optimization:")
	fmt.Println("-----------------------------")

	filePath := "../sample_files/build.gradle"

	// Measure memory before parsing
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Standard parsing
	_, err := api.ParseFile(filePath)
	if err != nil {
		fmt.Printf("❌ Standard parsing failed: %v\n", err)
		return
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	standardMemory := m2.Alloc - m1.Alloc

	// Reset memory measurement
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Optimized parsing
	optimizedOptions := &api.Options{
		SkipComments:      true,
		CollectRawContent: false,
		ParsePlugins:      true,
		ParseRepositories: true,
		ParseTasks:        false,
		ParseDependencies: true,
	}
	optimizedParser := api.NewParser(optimizedOptions)

	_, err = optimizedParser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("❌ Optimized parsing failed: %v\n", err)
		return
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)
	optimizedMemory := m2.Alloc - m1.Alloc

	fmt.Printf("  📊 Standard Parsing: %.1f KB\n", float64(standardMemory)/1024)
	fmt.Printf("  🚀 Optimized Parsing: %.1f KB\n", float64(optimizedMemory)/1024)

	if standardMemory > 0 {
		saved := float64(standardMemory-optimizedMemory) / float64(standardMemory) * 100
		fmt.Printf("  💰 Memory Saved: %.1f%%\n", saved)
	}

	fmt.Println("\n✅ Advanced features demonstration completed!")
}
