package depbump

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
)

func TestUpdateModule(t *testing.T) {
	projectPath := runpath.PARENT.Path()
	t.Log(projectPath)

	moduleInfo, err := GetModuleInfo(projectPath)
	require.NoError(t, err)
	require.Equal(t, syntaxgo_reflect.GetPkgPathV2[Module](), moduleInfo.Module.Path)

	execConfig := osexec.NewExecConfig().WithDebug().WithPath(projectPath)
	for _, dep := range moduleInfo.Require {
		if !dep.Indirect {
			require.NoError(t, UpdateModule(execConfig, dep.Path, moduleInfo.GetToolchainVersion()))
			return // once is enough
		}
	}
}
