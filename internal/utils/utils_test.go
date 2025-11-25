// Package utils: Unit tests of version comparison and Go version matching functions
// Tests semantic version comparison logic and toolchain version matching validation
//
// utils: 版本比较和 Go 版本匹配功能的单元测试
// 测试语义版本比较逻辑和工具链版本匹配验证
package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCompareVersions validates semantic version comparison logic
// Expects versions with "v" prefix from go list output
//
// TestCompareVersions 验证语义版本比较逻辑
// 期望带 "v" 前缀的版本来自 go list 输出
func TestCompareVersions(t *testing.T) {
	require.Equal(t, 1, CompareVersions("v1.0.1", "v1.0.0"))
	require.Equal(t, 0, CompareVersions("v1.0.0", "v1.0.0"))
	require.Equal(t, -1, CompareVersions("v1.0.0", "v1.0.1"))
}

// TestCanUseGoVersion validates Go version matching checks
// Tests toolchain version matching validation logic using go/version.Compare
//
// TestCanUseGoVersion 验证 Go 版本匹配检查
// 使用 go/version.Compare 测试工具链版本匹配验证逻辑
func TestCanUseGoVersion(t *testing.T) {
	// Basic version comparison // 基本版本比较
	require.True(t, CanUseGoVersion("1.20", "1.21"))
	require.False(t, CanUseGoVersion("1.21", "1.20"))
	require.True(t, CanUseGoVersion("1.22", "1.22.8"))

	// RC and patch versions // RC 和补丁版本
	require.True(t, CanUseGoVersion("1.22rc1", "1.22.0"))
	require.True(t, CanUseGoVersion("1.22.0", "1.22.1"))

	// Empty required version // 空的需求版本
	require.True(t, CanUseGoVersion("", "1.20"))
}

// TestIsStableVersion validates stable version detection logic
// Tests filtering of prerelease versions (preview, rc, beta, alpha) and +incompatible
// Expects version strings with "v" prefix from go list -m -versions output
//
// TestIsStableVersion 验证稳定版本检测逻辑
// 测试预发布版本（preview、rc、beta、alpha）和 +incompatible 的过滤
// 期望带 "v" 前缀的版本字符串来自 go list -m -versions 输出
func TestIsStableVersion(t *testing.T) {
	// Stable versions (from go list output) // 稳定版本（来自 go list 输出）
	require.True(t, IsStableVersion("v1.0.0"))
	require.True(t, IsStableVersion("v1.39.2"))
	require.True(t, IsStableVersion("v1.20.5"))

	// Unstable versions // 不稳定版本
	require.False(t, IsStableVersion("v2.0.0-preview.4"))
	require.False(t, IsStableVersion("v2.0.0-preview.4+incompatible"))
	require.False(t, IsStableVersion("v1.0.0-rc1"))
	require.False(t, IsStableVersion("v1.0.0-beta"))
	require.False(t, IsStableVersion("v1.0.0-alpha"))
	require.False(t, IsStableVersion("v1.0.0+incompatible"))
}
