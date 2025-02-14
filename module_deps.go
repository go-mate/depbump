package depbump

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/yyle88/erero"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath"
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
	Module  *Module    `json:"Module"`
	Go      string     `json:"Go"`
	Require []*Require `json:"Require"`
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
