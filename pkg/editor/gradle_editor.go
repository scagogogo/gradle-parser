// Package editor 提供Gradle文件的结构化编辑功能。
package editor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

// GradleEditor 结构化Gradle编辑器。
type GradleEditor struct {
	sourceMappedProject *model.SourceMappedProject
	modifications       []Modification
}

// Modification 表示一个修改操作。
type Modification struct {
	Type        ModificationType  `json:"type"`
	SourceRange model.SourceRange `json:"sourceRange"`
	OldText     string            `json:"oldText"`
	NewText     string            `json:"newText"`
	Description string            `json:"description"`
}

// ModificationType 修改类型。
type ModificationType string

const (
	ModificationTypeReplace ModificationType = "replace"
	ModificationTypeInsert  ModificationType = "insert"
	ModificationTypeDelete  ModificationType = "delete"
)

// NewGradleEditor 创建新的Gradle编辑器。
func NewGradleEditor(sourceMappedProject *model.SourceMappedProject) *GradleEditor {
	return &GradleEditor{
		sourceMappedProject: sourceMappedProject,
		modifications:       make([]Modification, 0),
	}
}

// UpdateDependencyVersion 更新依赖版本。
func (ge *GradleEditor) UpdateDependencyVersion(group, name, newVersion string) error {
	// 检查项目是否为nil。
	if ge.sourceMappedProject == nil {
		return fmt.Errorf("source mapped project is nil")
	}

	// 查找匹配的依赖。
	var targetDep *model.SourceMappedDependency
	for _, dep := range ge.sourceMappedProject.SourceMappedDependencies {
		if dep.Group == group && dep.Name == name {
			targetDep = dep
			break
		}
	}

	if targetDep == nil {
		return fmt.Errorf("dependency %s:%s not found", group, name)
	}

	// 如果当前版本和新版本相同，不需要修改。
	if targetDep.Version == newVersion {
		return nil
	}

	// 生成新的依赖声明。
	var newText string
	if targetDep.Version == "" {
		// 原来没有版本号，需要添加版本号。
		if strings.Contains(targetDep.RawText, "'") {
			newText = fmt.Sprintf("'%s:%s:%s'", group, name, newVersion)
		} else {
			newText = fmt.Sprintf("\"%s:%s:%s\"", group, name, newVersion)
		}
	} else {
		// 替换现有版本号。
		oldVersionPattern := regexp.QuoteMeta(targetDep.Version)
		re := regexp.MustCompile(oldVersionPattern)
		newText = re.ReplaceAllString(targetDep.RawText, newVersion)
	}

	// 创建修改操作。
	modification := Modification{
		Type:        ModificationTypeReplace,
		SourceRange: targetDep.SourceRange,
		OldText:     targetDep.RawText,
		NewText:     newText,
		Description: fmt.Sprintf("Update %s:%s version from '%s' to '%s'", group, name, targetDep.Version, newVersion),
	}

	ge.modifications = append(ge.modifications, modification)

	// 更新内存中的依赖信息。
	targetDep.Version = newVersion
	targetDep.RawText = newText

	return nil
}

// UpdatePluginVersion 更新插件版本。
func (ge *GradleEditor) UpdatePluginVersion(pluginId, newVersion string) error {
	// 检查项目是否为nil。
	if ge.sourceMappedProject == nil {
		return fmt.Errorf("source mapped project is nil")
	}

	// 查找匹配的插件。
	var targetPlugin *model.SourceMappedPlugin
	for _, plugin := range ge.sourceMappedProject.SourceMappedPlugins {
		if plugin.ID == pluginId {
			targetPlugin = plugin
			break
		}
	}

	if targetPlugin == nil {
		return fmt.Errorf("plugin %s not found", pluginId)
	}

	// 如果当前版本和新版本相同，不需要修改。
	if targetPlugin.Version == newVersion {
		return nil
	}

	// 生成新的插件声明。
	var newText string
	if targetPlugin.Version == "" {
		// 原来没有版本号，需要添加版本号。
		if strings.Contains(targetPlugin.RawText, "'") {
			newText = fmt.Sprintf("id '%s' version '%s'", pluginId, newVersion)
		} else {
			newText = fmt.Sprintf("id \"%s\" version \"%s\"", pluginId, newVersion)
		}
	} else {
		// 替换现有版本号。
		oldVersionPattern := regexp.QuoteMeta(targetPlugin.Version)
		re := regexp.MustCompile(oldVersionPattern)
		newText = re.ReplaceAllString(targetPlugin.RawText, newVersion)
	}

	// 创建修改操作。
	modification := Modification{
		Type:        ModificationTypeReplace,
		SourceRange: targetPlugin.SourceRange,
		OldText:     targetPlugin.RawText,
		NewText:     newText,
		Description: fmt.Sprintf("Update plugin %s version from '%s' to '%s'", pluginId, targetPlugin.Version, newVersion),
	}

	ge.modifications = append(ge.modifications, modification)

	// 更新内存中的插件信息。
	targetPlugin.Version = newVersion
	targetPlugin.RawText = newText

	return nil
}

// UpdateProperty 更新项目属性。
func (ge *GradleEditor) UpdateProperty(key, newValue string) error {
	// 检查项目是否为nil。
	if ge.sourceMappedProject == nil {
		return fmt.Errorf("source mapped project is nil")
	}

	// 查找匹配的属性。
	var targetProperty *model.SourceMappedProperty
	for _, prop := range ge.sourceMappedProject.SourceMappedProperties {
		if prop.Key == key {
			targetProperty = prop
			break
		}
	}

	if targetProperty == nil {
		return fmt.Errorf("property %s not found", key)
	}

	// 如果当前值和新值相同，不需要修改。
	if targetProperty.Value == newValue {
		return nil
	}

	// 生成新的属性声明。
	var newText string
	if strings.Contains(targetProperty.RawText, "'") {
		newText = fmt.Sprintf("%s = '%s'", key, newValue)
	} else {
		newText = fmt.Sprintf("%s = \"%s\"", key, newValue)
	}

	// 创建修改操作。
	modification := Modification{
		Type:        ModificationTypeReplace,
		SourceRange: targetProperty.SourceRange,
		OldText:     targetProperty.RawText,
		NewText:     newText,
		Description: fmt.Sprintf("Update property %s from '%s' to '%s'", key, targetProperty.Value, newValue),
	}

	ge.modifications = append(ge.modifications, modification)

	// 更新内存中的属性信息。
	targetProperty.Value = newValue
	targetProperty.RawText = newText

	return nil
}

// AddDependency 添加新依赖。
func (ge *GradleEditor) AddDependency(group, name, version, scope string) error {
	// 检查项目是否为nil。
	if ge.sourceMappedProject == nil {
		return fmt.Errorf("source mapped project is nil")
	}

	// 查找dependencies块的位置。
	dependenciesBlockLine := ge.findDependenciesBlock()
	if dependenciesBlockLine == -1 {
		return fmt.Errorf("dependencies block not found")
	}

	// 生成新的依赖声明。
	var newText string
	if scope == "" {
		scope = "implementation"
	}

	if version != "" {
		newText = fmt.Sprintf("    %s '%s:%s:%s'", scope, group, name, version)
	} else {
		newText = fmt.Sprintf("    %s '%s:%s'", scope, group, name)
	}

	// 找到插入位置（dependencies块的最后一行之前）。
	insertLine := ge.findDependenciesBlockEnd(dependenciesBlockLine)
	if insertLine == -1 {
		return fmt.Errorf("could not find dependencies block end")
	}

	// 计算插入位置。
	insertPos := 0
	for i := 0; i < insertLine-1; i++ {
		insertPos += len(ge.sourceMappedProject.Lines[i]) + 1 // +1 for newline。
	}

	// 创建插入操作。
	modification := Modification{
		Type: ModificationTypeInsert,
		SourceRange: model.SourceRange{
			Start: model.SourcePosition{
				Line:     insertLine,
				Column:   1,
				StartPos: insertPos,
				EndPos:   insertPos,
				Length:   0,
			},
			End: model.SourcePosition{
				Line:     insertLine,
				Column:   1,
				StartPos: insertPos,
				EndPos:   insertPos,
				Length:   0,
			},
		},
		OldText:     "",
		NewText:     newText + "\n",
		Description: fmt.Sprintf("Add dependency %s:%s:%s with scope %s", group, name, version, scope),
	}

	ge.modifications = append(ge.modifications, modification)

	return nil
}

// GetModifications 获取所有修改操作。
func (ge *GradleEditor) GetModifications() []Modification {
	return ge.modifications
}

// GetSourceMappedProject 获取源码映射项目。
func (ge *GradleEditor) GetSourceMappedProject() *model.SourceMappedProject {
	return ge.sourceMappedProject
}

// ClearModifications 清除所有修改操作。
func (ge *GradleEditor) ClearModifications() {
	ge.modifications = make([]Modification, 0)
}

// findDependenciesBlock 查找dependencies块的起始行。
func (ge *GradleEditor) findDependenciesBlock() int {
	if ge.sourceMappedProject == nil {
		return -1
	}

	for i, line := range ge.sourceMappedProject.Lines {
		if strings.Contains(strings.TrimSpace(line), "dependencies") && strings.Contains(line, "{") {
			return i + 1 // 返回1-based行号。
		}
	}
	return -1
}

// findDependenciesBlockEnd 查找dependencies块的结束行。
func (ge *GradleEditor) findDependenciesBlockEnd(startLine int) int {
	if ge.sourceMappedProject == nil {
		return -1
	}

	braceCount := 0
	started := false

	for i := startLine - 1; i < len(ge.sourceMappedProject.Lines); i++ {
		line := strings.TrimSpace(ge.sourceMappedProject.Lines[i])

		if strings.Contains(line, "{") {
			braceCount++
			started = true
		}

		if strings.Contains(line, "}") {
			braceCount--
			if started && braceCount == 0 {
				return i + 1 // 返回1-based行号。
			}
		}
	}

	return -1
}
