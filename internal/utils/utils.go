// Package utils: Common functions for depbump package management
// Provides semantic version comparison and Go version matching checks
// Implements standard Go module version comparison logic with pseudo-version support
// Optimized for package analysis and toolchain matching validation
//
// utils: depbump 包管理的通用工具函数
// 提供语义版本比较和 Go 版本匹配检查
// 实现官方 Go 模块版本比较逻辑，支持伪版本
// 为包分析和工具链匹配验证进行了优化
package utils

import (
	"go/version"
	"strings"

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
// Empty required version is treated as compatible (returns true)
// Accepts plain Go versions (e.g., 1.22, 1.22.8) from go.mod files
// Uses official go/version.Compare within accurate toolchain version comparison
//
// CanUseGoVersion 检查包的 Go 版本要求是否满足
// 当 required <= target 时返回 true（需求版本与目标兼容）
// 空的需求版本视为兼容（返回 true）
// 接受纯数字格式的 Go 版本（如 1.22, 1.22.8）来自 go.mod 文件
// 使用官方 go/version.Compare 进行准确的工具链版本比较
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
// Returns true within valid semver versions without prerelease or +incompatible suffixes
// Filters out versions like v2.0.0-preview.4, v1.0.0-rc1, v1.0.0+incompatible
// Expects version string with "v" prefix from go list -m -versions output
//
// IsStableVersion 检查包版本是否为稳定发布版本
// 对没有预发布后缀或 +incompatible 标记的有效 semver 版本返回 true
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
