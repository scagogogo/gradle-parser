// Package util 提供工具函数
package util

import (
	"os"
	"path/filepath"
	"strings"
)

// IsBuildGradleFile 检查文件是否是Gradle构建文件
func IsBuildGradleFile(filePath string) bool {
	fileName := filepath.Base(filePath)
	return fileName == "build.gradle" || fileName == "build.gradle.kts"
}

// IsSettingsGradleFile 检查文件是否是Gradle设置文件
func IsSettingsGradleFile(filePath string) bool {
	fileName := filepath.Base(filePath)
	return fileName == "settings.gradle" || fileName == "settings.gradle.kts"
}

// IsKotlinDSL 检查文件是否使用Kotlin DSL
func IsKotlinDSL(filePath string) bool {
	return strings.HasSuffix(filePath, ".kts")
}

// FindGradleFiles 在指定目录中查找所有Gradle文件
func FindGradleFiles(rootDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (IsBuildGradleFile(path) || IsSettingsGradleFile(path)) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// FindProjectRoot 查找包含build.gradle的项目根目录
func FindProjectRoot(startDir string) (string, error) {
	currentDir := startDir
	for {
		// 检查当前目录是否有build.gradle文件
		buildGradle := filepath.Join(currentDir, "build.gradle")
		buildGradleKts := filepath.Join(currentDir, "build.gradle.kts")

		if fileExists(buildGradle) || fileExists(buildGradleKts) {
			return currentDir, nil
		}

		// 获取父目录
		parentDir := filepath.Dir(currentDir)
		// 检查是否已到达文件系统根目录
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", os.ErrNotExist
}

// fileExists 检查文件是否存在
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetFileContent 获取文件内容
func GetFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
