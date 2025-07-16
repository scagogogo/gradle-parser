// Package model 提供源码位置追踪相关的数据结构
package model

import "fmt"

// SourcePosition 表示源码中的位置信息
type SourcePosition struct {
	Line      int `json:"line"`      // 行号（1-based）
	Column    int `json:"column"`    // 列号（1-based）
	StartPos  int `json:"startPos"`  // 在原始文本中的起始位置（0-based）
	EndPos    int `json:"endPos"`    // 在原始文本中的结束位置（0-based）
	Length    int `json:"length"`    // 文本长度
}

// SourceRange 表示源码中的范围
type SourceRange struct {
	Start SourcePosition `json:"start"`
	End   SourcePosition `json:"end"`
}

// String 返回位置的字符串表示
func (sp SourcePosition) String() string {
	return fmt.Sprintf("line %d, col %d", sp.Line, sp.Column)
}

// String 返回范围的字符串表示
func (sr SourceRange) String() string {
	return fmt.Sprintf("%s - %s", sr.Start.String(), sr.End.String())
}

// SourceMappedDependency 带源码位置信息的依赖
type SourceMappedDependency struct {
	*Dependency
	SourceRange SourceRange `json:"sourceRange"`
	RawText     string      `json:"rawText"` // 原始文本片段
}

// SourceMappedPlugin 带源码位置信息的插件
type SourceMappedPlugin struct {
	*Plugin
	SourceRange SourceRange `json:"sourceRange"`
	RawText     string      `json:"rawText"`
}

// SourceMappedRepository 带源码位置信息的仓库
type SourceMappedRepository struct {
	*Repository
	SourceRange SourceRange `json:"sourceRange"`
	RawText     string      `json:"rawText"`
}

// SourceMappedProperty 带源码位置信息的属性
type SourceMappedProperty struct {
	Key         string      `json:"key"`
	Value       string      `json:"value"`
	SourceRange SourceRange `json:"sourceRange"`
	RawText     string      `json:"rawText"`
}

// SourceMappedProject 带源码位置信息的项目
type SourceMappedProject struct {
	*Project
	
	// 带位置信息的组件
	SourceMappedDependencies []*SourceMappedDependency  `json:"sourceMappedDependencies"`
	SourceMappedPlugins      []*SourceMappedPlugin      `json:"sourceMappedPlugins"`
	SourceMappedRepositories []*SourceMappedRepository  `json:"sourceMappedRepositories"`
	SourceMappedProperties   []*SourceMappedProperty    `json:"sourceMappedProperties"`
	
	// 原始文本信息
	OriginalText string   `json:"originalText"`
	Lines        []string `json:"lines"` // 按行分割的原始文本
}

// SourceMappedParseResult 带源码位置信息的解析结果
type SourceMappedParseResult struct {
	*ParseResult
	SourceMappedProject *SourceMappedProject `json:"sourceMappedProject"`
}

// GetLineText 获取指定行的文本
func (smp *SourceMappedProject) GetLineText(lineNumber int) string {
	if lineNumber < 1 || lineNumber > len(smp.Lines) {
		return ""
	}
	return smp.Lines[lineNumber-1]
}

// GetTextRange 获取指定范围的文本
func (smp *SourceMappedProject) GetTextRange(sourceRange SourceRange) string {
	if sourceRange.Start.StartPos < 0 || sourceRange.End.EndPos > len(smp.OriginalText) {
		return ""
	}
	return smp.OriginalText[sourceRange.Start.StartPos:sourceRange.End.EndPos]
}

// FindDependencyByPosition 根据位置查找依赖
func (smp *SourceMappedProject) FindDependencyByPosition(line, column int) *SourceMappedDependency {
	for _, dep := range smp.SourceMappedDependencies {
		if line >= dep.SourceRange.Start.Line && line <= dep.SourceRange.End.Line {
			if line == dep.SourceRange.Start.Line && column < dep.SourceRange.Start.Column {
				continue
			}
			if line == dep.SourceRange.End.Line && column > dep.SourceRange.End.Column {
				continue
			}
			return dep
		}
	}
	return nil
}

// FindPluginByPosition 根据位置查找插件
func (smp *SourceMappedProject) FindPluginByPosition(line, column int) *SourceMappedPlugin {
	for _, plugin := range smp.SourceMappedPlugins {
		if line >= plugin.SourceRange.Start.Line && line <= plugin.SourceRange.End.Line {
			if line == plugin.SourceRange.Start.Line && column < plugin.SourceRange.Start.Column {
				continue
			}
			if line == plugin.SourceRange.End.Line && column > plugin.SourceRange.End.Column {
				continue
			}
			return plugin
		}
	}
	return nil
}

// FindPropertyByKey 根据键查找属性
func (smp *SourceMappedProject) FindPropertyByKey(key string) *SourceMappedProperty {
	for _, prop := range smp.SourceMappedProperties {
		if prop.Key == key {
			return prop
		}
	}
	return nil
}
