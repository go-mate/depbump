// Package utils: Common functions supporting depbump package management
// Provides semantic version comparison, Go version matching checks, and workspace iteration
// Implements standard Go module version comparison logic with pseudo-version support
// Enables efficient package analysis, toolchain matching validation, and recursive module processing
//
// utils: depbump 包管理的通用工具函数
// 提供语义版本比较、Go 版本匹配检查和工作区遍历
// 实现官方 Go 模块版本比较逻辑，支持伪版本
// 用于高效的包分析、工具链匹配验证和递归模块处理
package utils

import (
	"fmt"
	"go/version"
	"strings"

	"github.com/go-mate/go-work/workspath"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/zaplog"
	"golang.org/x/mod/semver"
)

// CompareVersions compares package version strings using semver
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
// Expects versions with "v" prefix (e.g., v1.2.3) from go list -m -versions output
//
// CompareVersions 使用 semver 比较包版本字符串
// 如果 v1 < v2 返回 -1，v1 == v2 返回 0，v1 > v2 返回 1
// 期望带 "v" 前缀的版本（如 v1.2.3）来自 go list -m -versions 输出
func CompareVersions(v1, v2 string) int {
	return semver.Compare(v1, v2)
}

// CanUseGoVersion checks if a package's Go version requirement is satisfied
// Returns true when required <= target (required version is compatible with target)
// Missing required version is treated as compatible (returns true)
// Accepts plain Go versions (e.g., 1.22, 1.22.8) from go.mod files
// Uses standard go/version.Compare to enable accurate toolchain version comparison
//
// CanUseGoVersion 检查包的 Go 版本要求是否满足
// 当 required <= target 时返回 true（需求版本与目标兼容）
// 缺失的需求版本视作兼容（返回 true）
// 接受纯数字格式的 Go 版本（如 1.22, 1.22.8）来自 go.mod 文件
// 使用标准 go/version.Compare 进行准确的工具链版本比较
func CanUseGoVersion(required, target string) bool {
	// Add "go" prefix to match go/version.Compare format
	// 添加 "go" 前缀以匹配 go/version.Compare 格式
	if !strings.HasPrefix(required, "go") {
		required = "go" + required
	}
	if !strings.HasPrefix(target, "go") {
		target = "go" + target
	}
	return version.Compare(required, target) <= 0
}

// IsStableVersion checks if a package version is a stable release
// Returns true when version is valid semver without prerelease and +incompatible suffixes
// Filters out versions like v2.0.0-preview.4, v1.0.0-rc1, v1.0.0+incompatible
// Expects version string with "v" prefix from go list -m -versions output
//
// IsStableVersion 检查包版本是否是稳定发布版本
// 对没有预发布后缀和 +incompatible 标记的有效 semver 版本返回 true
// 过滤掉如 v2.0.0-preview.4, v1.0.0-rc1, v1.0.0+incompatible 等版本
// 期望带 "v" 前缀的版本字符串来自 go list -m -versions 输出
func IsStableVersion(version string) bool {
	// Reject +incompatible versions // 拒绝 +incompatible 版本
	if strings.Contains(version, "+incompatible") {
		return false
	}

	// Check valid semver with no prerelease suffix
	// 检查有效 semver 且无预发布后缀
	return semver.IsValid(version) && semver.Prerelease(version) == ""
}

// UIProgress formats progress display as "(<current>/total)" with 1-based index
// Example: UIProgress(0, 10) returns "(<1>/10)"
//
// UIProgress 格式化进度显示，使用 1-based 索引
// 示例：UIProgress(0, 10) 返回 "(<1>/10)"
func UIProgress(idx, cnt int) string {
	return fmt.Sprintf("(<%d>/%d)", idx+1, cnt)
}

// ForeachModule iterates over workspace modules and executes callback
// Scans workspace using workspath configuration and processes each module
//
// ForeachModule 遍历工作区模块并执行回调
// 使用 workspath 配置扫描工作区并处理每个模块
func ForeachModule(execConfig *osexec.ExecConfig, fn func(*osexec.ExecConfig)) {
	workPath := osmustexist.ROOT(execConfig.Path)

	options := workspath.NewOptions().
		WithIncludeCurrentProject(true).
		WithIncludeCurrentPackage(false).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(false)

	moduleRoots := workspath.GetModulePaths(workPath, options)

	zaplog.SUG.Infoln("Recursive mode: found", eroticgo.CYAN.Sprint(len(moduleRoots)), "modules")

	for idx, modulePath := range moduleRoots {
		zaplog.SUG.Infoln("Module", eroticgo.GREEN.Sprint(UIProgress(idx, len(moduleRoots))), "Processing:", eroticgo.CYAN.Sprint(modulePath))
		fn(execConfig.NewConfig().WithPath(modulePath))
	}

	zaplog.SUG.Infoln("✅ Recursive updates completed!")
}
