// Package depbump tests: Module dep parsing and analysis test suite
// Tests module information fetch, go.mod parsing, and dep filtering functions
// Validates module path detection, toolchain configuration, and scoped requirement filtering
//
// depbump 测试包：模块依赖解析和分析测试套件
// 测试模块信息检索、go.mod 解析和依赖过滤功能
// 验证模块路径检测、工具链配置和作用域需求过滤
package depbump

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
	"golang.org/x/mod/modfile"
)

// TestGetModuleInfo validates module information parsing from go mod edit -json
// Tests JSON unmarshaling, module path extraction, and Go version detection
//
// TestGetModuleInfo 验证从 go mod edit -json 解析模块信息
// 测试 JSON 解组、模块路径提取和 Go 版本检测
func TestGetModuleInfo(t *testing.T) {
	moduleInfo, err := GetModuleInfo(runpath.PARENT.Path())
	require.NoError(t, err)
	t.Log(neatjsons.S(moduleInfo))

	require.Equal(t, syntaxgo_reflect.GetPkgPathV2[Module](), moduleInfo.Module.Path)

	t.Log(moduleInfo.Module.Path)
	t.Log(moduleInfo.Go)
}

// TestParseModuleFileDemo demonstrates direct go.mod file parsing using modfile lib
// Tests file reading, parsing, and module path validation with explicit file operations
//
// TestParseModuleFileDemo 演示使用 modfile 库直接解析 go.mod 文件
// 测试文件读取、解析和模块路径验证，带显式文件操作
func TestParseModuleFileDemo(t *testing.T) {
	const fileName = "go.mod"
	modPath := osmustexist.FILE(runpath.PARENT.Join(fileName))
	t.Log(modPath)

	modData, err := os.ReadFile(modPath)
	require.NoError(t, err)

	modFile, err := modfile.Parse(fileName, modData, nil)
	require.NoError(t, err)

	t.Log(neatjsons.S(modFile))
	require.Equal(t, syntaxgo_reflect.GetPkgPathV2[Module](), modFile.Module.Mod.Path)
}

// TestParseModuleFile tests the ParseModuleFile utility function wrapper
// Validates encapsulated file path handling and module parsing function
//
// TestParseModuleFile 测试 ParseModuleFile 工具函数包装器
// 验证封装的文件路径处理和模块解析功能
func TestParseModuleFile(t *testing.T) {
	moduleFile, err := ParseModuleFile(runpath.PARENT.Path())
	require.NoError(t, err)

	t.Log(neatjsons.S(moduleFile))
	require.Equal(t, syntaxgo_reflect.GetPkgPathV2[Module](), moduleFile.Module.Mod.Path)
}

// TestModuleInfo_GetDirectRequires tests filtering for direct (non-indirect) dependencies
// Validates that indirect dependencies are properly excluded from the result set
//
// TestModuleInfo_GetDirectRequires 测试过滤直接（非间接）依赖
// 验证间接依赖被正确排除在结果集之外
func TestModuleInfo_GetDirectRequires(t *testing.T) {
	moduleInfo, err := GetModuleInfo(runpath.PARENT.Path())
	require.NoError(t, err)
	requires := moduleInfo.GetDirectRequires()
	t.Log(neatjsons.S(requires))
}

// TestModuleInfo_GetScopedRequires tests dep filtering by type scope
// Validates filtering logic for direct, indirect, and all dep types
//
// TestModuleInfo_GetScopedRequires 测试按类别范围过滤依赖
// 验证直接、间接和所有依赖类别的过滤逻辑
func TestModuleInfo_GetScopedRequires(t *testing.T) {
	moduleInfo, err := GetModuleInfo(runpath.PARENT.Path())
	require.NoError(t, err)
	requires := moduleInfo.GetScopedRequires(DepCateEveryone)
	t.Log(neatjsons.S(requires))
}
