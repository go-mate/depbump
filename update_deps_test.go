package depbump

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/runpath"
	"github.com/yyle88/syntaxgo/syntaxgo_reflect"
)

func TestUpdateModule(t *testing.T) {
	projectPath := runpath.PARENT.Path()
	t.Log(projectPath)

	modInfo, err := GetModInfo(projectPath)
	require.NoError(t, err)
	require.Equal(t, syntaxgo_reflect.GetPkgPathV2[Module](), modInfo.Module.Path)

	for _, dep := range modInfo.Require {
		if !dep.Indirect {
			require.NoError(t, UpdateModule(projectPath, dep.Path))
			return // once is enough
		}
	}
}
