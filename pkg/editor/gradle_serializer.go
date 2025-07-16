// Package editor 提供Gradle文件的序列化功能
package editor

import (
	"fmt"
	"sort"
	"strings"
)

// GradleSerializer 最小diff序列化器
type GradleSerializer struct {
	originalText string
	lines        []string
}

// NewGradleSerializer 创建新的序列化器
func NewGradleSerializer(originalText string) *GradleSerializer {
	return &GradleSerializer{
		originalText: originalText,
		lines:        strings.Split(originalText, "\n"),
	}
}

// ApplyModifications 应用修改操作并返回新的文本
func (gs *GradleSerializer) ApplyModifications(modifications []Modification) (string, error) {
	if len(modifications) == 0 {
		return gs.originalText, nil
	}

	// 按位置排序修改操作（从后往前，避免位置偏移）
	sortedMods := make([]Modification, len(modifications))
	copy(sortedMods, modifications)
	sort.Slice(sortedMods, func(i, j int) bool {
		return sortedMods[i].SourceRange.Start.StartPos > sortedMods[j].SourceRange.Start.StartPos
	})

	// 应用修改
	result := gs.originalText
	for _, mod := range sortedMods {
		var err error
		result, err = gs.applyModification(result, mod)
		if err != nil {
			return "", fmt.Errorf("failed to apply modification: %w", err)
		}
	}

	return result, nil
}

// applyModification 应用单个修改操作
func (gs *GradleSerializer) applyModification(text string, mod Modification) (string, error) {
	switch mod.Type {
	case ModificationTypeReplace:
		return gs.applyReplace(text, mod)
	case ModificationTypeInsert:
		return gs.applyInsert(text, mod)
	case ModificationTypeDelete:
		return gs.applyDelete(text, mod)
	default:
		return "", fmt.Errorf("unknown modification type: %s", mod.Type)
	}
}

// applyReplace 应用替换操作
func (gs *GradleSerializer) applyReplace(text string, mod Modification) (string, error) {
	startPos := mod.SourceRange.Start.StartPos
	endPos := mod.SourceRange.End.StartPos

	if startPos < 0 || endPos > len(text) || startPos > endPos {
		return "", fmt.Errorf("invalid source range for replace operation")
	}

	// 验证要替换的文本是否匹配
	actualText := text[startPos:endPos]
	if actualText != mod.OldText {
		// 尝试在行内查找匹配的文本
		line := gs.getLineFromPosition(text, startPos)
		if strings.Contains(line, mod.OldText) {
			// 在行内查找精确位置
			lineStart := gs.getLineStartPosition(text, startPos)
			relativePos := strings.Index(line, mod.OldText)
			if relativePos != -1 {
				actualStartPos := lineStart + relativePos
				actualEndPos := actualStartPos + len(mod.OldText)
				return text[:actualStartPos] + mod.NewText + text[actualEndPos:], nil
			}
		}
		return "", fmt.Errorf("text mismatch: expected '%s', got '%s'", mod.OldText, actualText)
	}

	return text[:startPos] + mod.NewText + text[endPos:], nil
}

// applyInsert 应用插入操作
func (gs *GradleSerializer) applyInsert(text string, mod Modification) (string, error) {
	insertPos := mod.SourceRange.Start.StartPos

	if insertPos < 0 || insertPos > len(text) {
		return "", fmt.Errorf("invalid insert position")
	}

	return text[:insertPos] + mod.NewText + text[insertPos:], nil
}

// applyDelete 应用删除操作
func (gs *GradleSerializer) applyDelete(text string, mod Modification) (string, error) {
	startPos := mod.SourceRange.Start.StartPos
	endPos := mod.SourceRange.End.StartPos

	if startPos < 0 || endPos > len(text) || startPos > endPos {
		return "", fmt.Errorf("invalid source range for delete operation")
	}

	return text[:startPos] + text[endPos:], nil
}

// getLineFromPosition 根据位置获取所在行的文本
func (gs *GradleSerializer) getLineFromPosition(text string, pos int) string {
	lines := strings.Split(text, "\n")
	currentPos := 0

	for _, line := range lines {
		lineEnd := currentPos + len(line)
		if pos >= currentPos && pos <= lineEnd {
			return line
		}
		currentPos = lineEnd + 1 // +1 for newline
	}

	return ""
}

// getLineStartPosition 根据位置获取所在行的起始位置
func (gs *GradleSerializer) getLineStartPosition(text string, pos int) int {
	lines := strings.Split(text, "\n")
	currentPos := 0

	for _, line := range lines {
		lineEnd := currentPos + len(line)
		if pos >= currentPos && pos <= lineEnd {
			return currentPos
		}
		currentPos = lineEnd + 1 // +1 for newline
	}

	return 0
}

// GenerateDiff 生成修改的diff信息
func (gs *GradleSerializer) GenerateDiff(modifications []Modification) []DiffLine {
	diffLines := make([]DiffLine, 0)

	for _, mod := range modifications {
		switch mod.Type {
		case ModificationTypeReplace:
			diffLines = append(diffLines, DiffLine{
				Type:        DiffTypeRemove,
				LineNumber:  mod.SourceRange.Start.Line,
				Content:     mod.OldText,
				Description: mod.Description,
			})
			diffLines = append(diffLines, DiffLine{
				Type:        DiffTypeAdd,
				LineNumber:  mod.SourceRange.Start.Line,
				Content:     mod.NewText,
				Description: mod.Description,
			})
		case ModificationTypeInsert:
			diffLines = append(diffLines, DiffLine{
				Type:        DiffTypeAdd,
				LineNumber:  mod.SourceRange.Start.Line,
				Content:     mod.NewText,
				Description: mod.Description,
			})
		case ModificationTypeDelete:
			diffLines = append(diffLines, DiffLine{
				Type:        DiffTypeRemove,
				LineNumber:  mod.SourceRange.Start.Line,
				Content:     mod.OldText,
				Description: mod.Description,
			})
		}
	}

	return diffLines
}

// DiffLine 表示diff中的一行
type DiffLine struct {
	Type        DiffType `json:"type"`
	LineNumber  int      `json:"lineNumber"`
	Content     string   `json:"content"`
	Description string   `json:"description"`
}

// DiffType diff类型
type DiffType string

const (
	DiffTypeAdd    DiffType = "add"
	DiffTypeRemove DiffType = "remove"
	DiffTypeModify DiffType = "modify"
)

// String 返回diff行的字符串表示
func (dl DiffLine) String() string {
	prefix := " "
	switch dl.Type {
	case DiffTypeAdd:
		prefix = "+"
	case DiffTypeRemove:
		prefix = "-"
	case DiffTypeModify:
		prefix = "~"
	}

	return fmt.Sprintf("%s %d: %s", prefix, dl.LineNumber, dl.Content)
}

// ValidateModifications 验证修改操作的有效性
func (gs *GradleSerializer) ValidateModifications(modifications []Modification) []error {
	errors := make([]error, 0)

	for i, mod := range modifications {
		// 检查位置范围
		if mod.SourceRange.Start.StartPos < 0 {
			errors = append(errors, fmt.Errorf("modification %d: invalid start position %d", i, mod.SourceRange.Start.StartPos))
		}

		if mod.SourceRange.End.StartPos > len(gs.originalText) {
			errors = append(errors, fmt.Errorf("modification %d: end position %d exceeds text length %d", i, mod.SourceRange.End.StartPos, len(gs.originalText)))
		}

		if mod.SourceRange.Start.StartPos > mod.SourceRange.End.StartPos {
			errors = append(errors, fmt.Errorf("modification %d: start position %d > end position %d", i, mod.SourceRange.Start.StartPos, mod.SourceRange.End.StartPos))
		}

		// 检查替换操作的文本匹配
		if mod.Type == ModificationTypeReplace {
			startPos := mod.SourceRange.Start.StartPos
			endPos := mod.SourceRange.End.StartPos

			if startPos >= 0 && endPos <= len(gs.originalText) && startPos <= endPos {
				actualText := gs.originalText[startPos:endPos]
				if actualText != mod.OldText {
					errors = append(errors, fmt.Errorf("modification %d: text mismatch, expected '%s', got '%s'", i, mod.OldText, actualText))
				}
			}
		}
	}

	return errors
}

// GetModificationSummary 获取修改操作的摘要
func (gs *GradleSerializer) GetModificationSummary(modifications []Modification) ModificationSummary {
	summary := ModificationSummary{
		TotalModifications: len(modifications),
		ByType:             make(map[ModificationType]int),
		Descriptions:       make([]string, 0),
	}

	for _, mod := range modifications {
		summary.ByType[mod.Type]++
		summary.Descriptions = append(summary.Descriptions, mod.Description)
	}

	return summary
}

// ModificationSummary 修改操作摘要
type ModificationSummary struct {
	TotalModifications int                      `json:"totalModifications"`
	ByType             map[ModificationType]int `json:"byType"`
	Descriptions       []string                 `json:"descriptions"`
}
