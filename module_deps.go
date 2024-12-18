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

type ModInfo struct {
	Module  *Module    `json:"Module"`
	Go      string     `json:"Go"`
	Require []*Require `json:"Require"`
}

func GetModInfo(projectPath string) (*ModInfo, error) {
	output, err := osexec.ExecInPath(projectPath, "go", "mod", "edit", "-json")
	if err != nil {
		return nil, erero.Wro(err)
	}
	var modInfo ModInfo
	if err := json.Unmarshal(output, &modInfo); err != nil {
		return nil, erero.Wro(err)
	}
	return &modInfo, nil
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
