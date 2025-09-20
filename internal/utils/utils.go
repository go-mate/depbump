// Package utils: Common functions for depbump package management
// Provides semantic version comparison and Go version matching checks
// Implements official Go module version comparison logic with pseudo-version support
// Optimized for package analysis and toolchain matching validation
//
// utils: depbump 包管理的通用工具函数
// 提供语义版本比较和 Go 版本匹配检查
// 实现官方 Go 模块版本比较逻辑，支持伪版本
// 为包分析和工具链匹配验证进行了优化
package utils

import (
	"strings"

	"golang.org/x/mod/semver"
)

// CompareVersions compares two version strings using official Go semantic versioning
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
// Handles version prefix normalization for consistent comparison
// Supports pseudo-versions, pre-release versions, and complex version formats
//
// CompareVersions 使用官方 Go 语义版本比较两个版本字符串
// 如果 v1 < v2 返回 -1，v1 == v2 返回 0，v1 > v2 返回 1
// 自动处理版本前缀标准化以实现一致比较
// 支持伪版本、预发布版本和复杂版本格式
func CompareVersions(v1, v2 string) int {
	// Ensure versions have "v" prefix for semver compatibility
	// 确保版本号带有 "v" 前缀以兼容 semver
	if !strings.HasPrefix(v1, "v") {
		v1 = "v" + v1
	}
	if !strings.HasPrefix(v2, "v") {
		v2 = "v" + v2
	}
	return semver.Compare(v1, v2)
}

// CanUseGoVersion validates if a required Go version can operate with target version
// Returns true when required version is less than or equal to target version
// Handles empty required version as universally compatible (returns true)
// Essential for preventing Go toolchain version conflicts in package upgrades
//
// CanUseGoVersion 验证所需 Go 版本是否能与目标版本配合工作
// 当所需版本小于或等于目标版本时返回 true
// 将空的所需版本处理为通用兼容（返回 true）
// 对于防止包升级中的 Go 工具链版本冲突至关重要
func CanUseGoVersion(required, target string) bool {
	if required == "" {
		return true
	}
	return CompareVersions(required, target) <= 0
}
