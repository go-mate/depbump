// Package depbump tests: Dependency update functionality test suite
// Tests module update operations, toolchain version handling, and configuration validation
// Validates go get command execution and dependency upgrade behavior
//
// depbump 测试包：依赖更新功能测试套件
// 测试模块更新操作、工具链版本处理和配置验证
// 验证 go get 命令执行和依赖升级行为
package depbump

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
)

// TestUpdateModule validates single module dependency update functionality
// Tests toolchain configuration, update mode handling, and direct dependency processing
//
// TestUpdateModule 验证单个模块依赖更新功能
// 测试工具链配置、更新模式处理和直接依赖处理
func TestUpdateModule(t *testing.T) {
	projectPath := runpath.PARENT.Path()
	t.Log(projectPath)

	moduleInfo, err := GetModuleInfo(projectPath)
	require.NoError(t, err)
	require.Equal(t, syntaxgo_reflect.GetPkgPathV2[Module](), moduleInfo.Module.Path)

	execConfig := osexec.NewExecConfig().WithDebug().WithPath(projectPath)
	for _, dep := range moduleInfo.Require {
		if !dep.Indirect {
			require.NoError(t, UpdateModule(execConfig, dep.Path, &UpdateConfig{
				Toolchain: moduleInfo.GetToolchainVersion(),
				Mode:      GetModeUpdate,
			}))
			return // once is enough // 一次就足够了
		}
	}
}
