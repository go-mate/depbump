package depbump

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/yyle88/erero"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath"
	"github.com/yyle88/tern/zerotern"
	"golang.org/x/mod/modfile"
)

type Module struct {
	Path string `json:"Path"`
}

type Require struct {
	Path     string `json:"Path"`
	Version  string `json:"Version"`
	Indirect bool   `json:"Indirect"`
}

type ModuleInfo struct {
	Module    *Module    `json:"Module"`
	Go        string     `json:"Go"`
	Toolchain string     `json:"Toolchain"`
	Require   []*Require `json:"Require"`
}

// GetToolchainVersion 获取当前使用的工具链的go版本号，当没有配置工具链时返回模块的go版本号
func (a *ModuleInfo) GetToolchainVersion() string {
	return zerotern.VF(a.Toolchain, func() string {
		return "go" + a.Go // 因为这里的版本不带go前缀，只是 1.22.8 这样的，需要拼接信息
	})
}

func (a *ModuleInfo) GetDirectModules() []*Require {
	var directModules []*Require
	for _, require := range a.Require {
		if !require.Indirect {
			directModules = append(directModules, require)
		}
	}
	return directModules
}

func GetModuleInfo(projectPath string) (*ModuleInfo, error) {
	output, err := osexec.ExecInPath(projectPath, "go", "mod", "edit", "-json")
	if err != nil {
		return nil, erero.Wro(err)
	}
	var moduleInfo ModuleInfo
	if err := json.Unmarshal(output, &moduleInfo); err != nil {
		return nil, erero.Wro(err)
	}
	return &moduleInfo, nil
}

func ParseModuleFile(projectPath string) (*modfile.File, error) {
	const fileName = "go.mod"
	modPath, err := osexistpath.FILE(filepath.Join(projectPath, fileName))
	if err != nil {
		return nil, erero.Wro(err)
	}
	modData, err := os.ReadFile(modPath)
	if err != nil {
		return nil, erero.Wro(err)
	}
	modFile, err := modfile.Parse(fileName, modData, nil)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return modFile, nil
}
