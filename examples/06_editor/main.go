// 示例：使用Gradle结构化编辑器进行精确修改
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/gradle-parser/pkg/api"
	"github.com/scagogogo/gradle-parser/pkg/editor"
)

func main() {
	fmt.Println("=== Gradle结构化编辑器示例 ===")

	// 使用示例文件
	testFilePath := filepath.Join("..", "sample_files", "build.gradle")

	// 方法1：使用便捷API进行单个修改
	fmt.Println("\n1. 使用便捷API更新依赖版本")
	_, err := api.UpdateDependencyVersion(testFilePath, "mysql", "mysql-connector-java", "8.0.31")
	if err != nil {
		log.Printf("更新失败: %v", err)
	} else {
		fmt.Println("✅ 成功更新mysql版本到8.0.31")
		// 可以将返回的文本写入文件
	}

	// 方法2：使用编辑器进行批量修改
	fmt.Println("\n2. 使用编辑器进行批量修改")

	// 创建编辑器
	gradleEditor, err := api.CreateGradleEditor(testFilePath)
	if err != nil {
		log.Fatalf("创建编辑器失败: %v", err)
	}

	// 批量修改
	modifications := []struct {
		description string
		action      func() error
	}{
		{
			description: "更新项目版本",
			action: func() error {
				return gradleEditor.UpdateProperty("version", "1.0.0")
			},
		},
		{
			description: "更新Spring Boot插件版本",
			action: func() error {
				return gradleEditor.UpdatePluginVersion("org.springframework.boot", "2.7.2")
			},
		},
		{
			description: "更新Guava依赖版本",
			action: func() error {
				return gradleEditor.UpdateDependencyVersion("com.google.guava", "guava", "31.1-jre")
			},
		},
		{
			description: "添加新依赖",
			action: func() error {
				return gradleEditor.AddDependency("org.apache.commons", "commons-text", "1.9", "implementation")
			},
		},
	}

	// 执行所有修改
	for _, mod := range modifications {
		if err := mod.action(); err != nil {
			fmt.Printf("❌ %s失败: %v\n", mod.description, err)
		} else {
			fmt.Printf("✅ %s成功\n", mod.description)
		}
	}

	// 获取修改摘要
	fmt.Println("\n3. 修改摘要")
	allModifications := gradleEditor.GetModifications()
	serializer := editor.NewGradleSerializer(gradleEditor.GetSourceMappedProject().OriginalText)
	summary := serializer.GetModificationSummary(allModifications)

	fmt.Printf("总修改数: %d\n", summary.TotalModifications)
	fmt.Printf("修改类型分布:\n")
	for modType, count := range summary.ByType {
		fmt.Printf("  %s: %d\n", modType, count)
	}

	fmt.Println("\n修改详情:")
	for i, desc := range summary.Descriptions {
		fmt.Printf("  %d. %s\n", i+1, desc)
	}

	// 应用修改
	fmt.Println("\n4. 应用修改")
	finalText, err := serializer.ApplyModifications(allModifications)
	if err != nil {
		log.Fatalf("应用修改失败: %v", err)
	}

	// 生成diff
	fmt.Println("\n5. 修改Diff")
	diffLines := serializer.GenerateDiff(allModifications)
	for _, diffLine := range diffLines {
		fmt.Println(diffLine.String())
	}

	// 可选：将修改后的内容写入新文件
	outputPath := "build.gradle.new"
	if err := os.WriteFile(outputPath, []byte(finalText), 0644); err != nil {
		log.Printf("写入文件失败: %v", err)
	} else {
		fmt.Printf("\n✅ 修改后的文件已保存到: %s\n", outputPath)
	}

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n主要特性:")
	fmt.Println("✅ 源码位置精确追踪")
	fmt.Println("✅ 最小化diff修改")
	fmt.Println("✅ 批量修改支持")
	fmt.Println("✅ 修改验证和回滚")
	fmt.Println("✅ 保持原始格式和注释")
}
