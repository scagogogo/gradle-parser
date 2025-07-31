// Package dependency 提供Gradle依赖解析功能。
package dependency

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scagogogo/gradle-parser/pkg/model"
)

// 常见的依赖声明正则表达式。
var (
	// 格式: group:name:version。
	// 例如: org.springframework:spring-core:5.3.10。
	gavRegex = regexp.MustCompile(`^(['"]?)([^:'"]+):([^:'"]+):([^'"]+)(['"]?)$`)

	// 格式: group:name (没有版本号)。
	// 例如: org.springframework.boot:spring-boot-starter-web。
	gaRegex = regexp.MustCompile(`^(['"]?)([^:'"]+):([^:'"]+)(['"]?)$`)

	// 格式: group.name:name:version。
	// 例如: org.springframework.boot:spring-boot-starter:2.5.5。
	dotNameRegex = regexp.MustCompile(`^(['"]?)([^:'"]+)\.([^:'"]+):([^:'"]+):([^'"]+)(['"]?)$`)

	// 格式: project(":name")。
	// 例如: project(":app")。
	projectRefRegex = regexp.MustCompile(`^project\(['"]:(.*)['"]\)$`)
)

// 依赖配置范围。
var commonScopes = []string{
	"implementation", "api", "compile", "compileOnly", "runtime", "runtimeOnly",
	"testImplementation", "testApi", "testCompile", "testCompileOnly", "testRuntime", "testRuntimeOnly",
	"androidTestImplementation", "androidTestApi", "androidTestCompile",
	"debugImplementation", "releaseImplementation",
}

// DependencyParser 处理Gradle依赖解析。
type DependencyParser struct{}

// NewDependencyParser 创建新的依赖解析器。
func NewDependencyParser() *DependencyParser {
	return &DependencyParser{}
}

// ParseDependencyBlock 解析依赖块。
func (dp *DependencyParser) ParseDependencyBlock(block *model.ScriptBlock) ([]*model.Dependency, error) {
	if block == nil {
		return nil, fmt.Errorf("依赖块为空")
	}

	deps := make([]*model.Dependency, 0)

	// 遍历所有可能的依赖配置范围。
	for _, scope := range commonScopes {
		scopeDeps, err := dp.parseScopedDependencies(block, scope)
		if err != nil {
			return nil, err
		}
		deps = append(deps, scopeDeps...)
	}

	// 处理任何自定义范围的依赖。
	for methodName, closures := range block.Closures {
		if !contains(commonScopes, methodName) {
			// 这可能是自定义范围。
			for _, closure := range closures {
				customDeps, err := dp.parseCustomDependencies(closure, methodName)
				if err != nil {
					return nil, err
				}
				deps = append(deps, customDeps...)
			}
		}
	}

	return deps, nil
}

// parseScopedDependencies 解析指定范围的依赖项。
func (dp *DependencyParser) parseScopedDependencies(block *model.ScriptBlock, scope string) ([]*model.Dependency, error) {
	deps := make([]*model.Dependency, 0)

	// 检查是否有该范围的依赖项方法调用。
	if closures, ok := block.Closures[scope]; ok {
		for _, closure := range closures {
			for _, value := range closure.Values {
				if dep, ok := dp.parseDependencyString(fmt.Sprintf("%v", value), scope); ok {
					deps = append(deps, dep)
				}
			}
		}
	}

	return deps, nil
}

// parseCustomDependencies 解析自定义范围的依赖项。
func (dp *DependencyParser) parseCustomDependencies(block *model.ScriptBlock, scope string) ([]*model.Dependency, error) {
	deps := make([]*model.Dependency, 0)

	for _, value := range block.Values {
		if dep, ok := dp.parseDependencyString(fmt.Sprintf("%v", value), scope); ok {
			deps = append(deps, dep)
		}
	}

	return deps, nil
}

// parseDependencyString 从字符串解析依赖项。
func (dp *DependencyParser) parseDependencyString(depStr string, scope string) (*model.Dependency, bool) {
	// 清理字符串。
	depStr = strings.TrimSpace(depStr)

	// 项目依赖。
	if match := projectRefRegex.FindStringSubmatch(depStr); len(match) > 1 {
		return &model.Dependency{
			Name:  match[1],
			Scope: scope,
			Raw:   depStr,
		}, true
	}

	// 标准GAV格式: group:name:version。
	if match := gavRegex.FindStringSubmatch(depStr); len(match) > 4 {
		return &model.Dependency{
			Group:   match[2],
			Name:    match[3],
			Version: match[4],
			Scope:   scope,
			Raw:     depStr,
		}, true
	}

	// GA格式: group:name (没有版本号)。
	if match := gaRegex.FindStringSubmatch(depStr); len(match) > 3 {
		return &model.Dependency{
			Group:   match[2],
			Name:    match[3],
			Version: "", // 版本号为空，可能由dependency-management管理。
			Scope:   scope,
			Raw:     depStr,
		}, true
	}

	// 带命名空间的格式: group.name:name:version。
	if match := dotNameRegex.FindStringSubmatch(depStr); len(match) > 5 {
		group := match[2] + "." + match[3]
		return &model.Dependency{
			Group:   group,
			Name:    match[4],
			Version: match[5],
			Scope:   scope,
			Raw:     depStr,
		}, true
	}

	// 未识别的依赖格式。
	return nil, false
}

// ExtractDependenciesFromText 从原始文本中提取依赖。
func (dp *DependencyParser) ExtractDependenciesFromText(text string) []*model.Dependency {
	deps := make([]*model.Dependency, 0)

	// 分析文本中的依赖声明。
	lines := strings.Split(text, "\n")
	currentScope := ""

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查是否是依赖范围声明行。
		for _, scope := range commonScopes {
			if strings.HasPrefix(trimmedLine, scope) && strings.Contains(trimmedLine, "(") {
				currentScope = scope
				break
			}
		}

		// 尝试解析为依赖项。
		for _, pattern := range []string{
			`['"]([^'"]+):([^'"]+):([^'"]+)['"]`,           // "group:name:version"。
			`['"]([^'"]+):([^'"]+)['"]`,                    // "group:name" (没有版本号)。
			`['"]([^'"]+)\.([^'"]+):([^'"]+):([^'"]+)['"]`, // "group.name:name:version"。
			`project\(['"]:(.*)['"]\)`,                     // project(":name")。
		} {
			re := regexp.MustCompile(pattern)
			matches := re.FindAllStringSubmatch(trimmedLine, -1)

			for _, match := range matches {
				if len(match) > 0 {
					rawDep := match[0]
					// filter out some unwanted deps。
					if strings.Contains(rawDep, "https://github.com") ||
						strings.Contains(rawDep, "https://central.sonatype.com/repository/maven-snapshots") ||
						strings.Contains(rawDep, "https://ossrh-staging-api.central.sonatype.com/service/local/") {
						continue
					}
					if dep, ok := dp.parseDependencyString(rawDep, currentScope); ok {
						deps = append(deps, dep)
					}
				}
			}
		}
	}

	return deps
}

// GroupDependenciesByScope 按范围对依赖进行分组。
func (dp *DependencyParser) GroupDependenciesByScope(deps []*model.Dependency) []*model.DependencySet {
	// 使用map收集按范围分组的依赖。
	scopeMap := make(map[string][]*model.Dependency)

	for _, dep := range deps {
		if dep.Scope != "" {
			scopeMap[dep.Scope] = append(scopeMap[dep.Scope], dep)
		} else {
			// 默认范围。
			scopeMap["implementation"] = append(scopeMap["implementation"], dep)
		}
	}

	// 转换为DependencySet列表。
	sets := make([]*model.DependencySet, 0, len(scopeMap))
	for scope, scopeDeps := range scopeMap {
		sets = append(sets, &model.DependencySet{
			Scope:        scope,
			Dependencies: scopeDeps,
		})
	}

	return sets
}

// 辅助函数: 检查字符串是否在切片中。
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
