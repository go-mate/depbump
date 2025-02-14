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

func TestGetModuleInfo(t *testing.T) {
	moduleInfo, err := GetModuleInfo(runpath.PARENT.Path())
	require.NoError(t, err)
	t.Log(neatjsons.S(moduleInfo))

	require.Equal(t, syntaxgo_reflect.GetPkgPathV2[Module](), moduleInfo.Module.Path)
}

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

func TestParseModuleFile(t *testing.T) {
	modFile, err := ParseModuleFile(runpath.PARENT.Path())
	require.NoError(t, err)

	t.Log(neatjsons.S(modFile))
	require.Equal(t, syntaxgo_reflect.GetPkgPathV2[Module](), modFile.Module.Mod.Path)
}
