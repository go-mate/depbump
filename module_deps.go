// Package depbump: Go module dep management and upgrade automation
// Provides comprehensive tools that analyze, update, and manage Go module deps
// Supports both direct and indirect dep handling with version management integration
//
// depbump: Go 模块依赖管理和升级自动化
// 提供全面的工具来分析、更新和管理 Go 模块依赖
// 支持直接和间接依赖处理，集成版本控制
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

// DepCate defines the type of dependencies to be processed
// Supports filtering: direct, indirect, and complete package types
//
// DepCate 定义要处理的依赖类别
// 支持按直接、间接或所有依赖类型过滤
type DepCate string

const (
	DepCateDirect   DepCate = "DIRECT"   // Direct dependencies just // 仅直接依赖
	DepCateIndirect DepCate = "INDIRECT" // Indirect dependencies just // 仅间接依赖
	DepCateEveryone DepCate = "EVERYONE" // All dependencies // 所有依赖
)

// Module represents the main module information from go.mod
// Contains the module path and related metadata
//
// Module 代表 go.mod 中的主模块信息
// 包含模块路径和相关元数据
type Module struct {
	Path string `json:"Path"` // Module path // 模块路径
}

// Require represents a single dep requirement
// Contains dep path, version, and indirect status
//
// Require 代表单个依赖需求
// 包含依赖路径、版本和间接状态
type Require struct {
	Path     string `json:"Path"`     // Dep path // 依赖路径
	Version  string `json:"Version"`  // Current version // 当前版本
	Indirect bool   `json:"Indirect"` // Whether indirect dep // 是否为间接依赖
}

// ModuleInfo contains complete module dep information
// Parsed from go mod edit -json output with toolchain details
//
// ModuleInfo 包含完整的模块依赖信息
// 从 go mod edit -json 输出解析，包含工具链详情
type ModuleInfo struct {
	Module    *Module    `json:"Module"`    // Main module info // 主模块信息
	Go        string     `json:"Go"`        // Go version requirement // Go 版本需求
	Toolchain string     `json:"Toolchain"` // Toolchain specification // 工具链规范
	Require   []*Require `json:"Require"`   // Dependency list // 依赖列表
}

// GetToolchainVersion returns the effective Go toolchain version within this module
// Falls back to module Go version with "go" prefix when toolchain is not specified
// Used to ensure consistent Go version during package updates
//
// GetToolchainVersion 返回此模块的有效 Go 工具链版本
// 当未指定工具链时，回退到带有 "go" 前缀的模块 Go 版本
// 用于确保依赖更新期间的 Go 版本一致性
func (a *ModuleInfo) GetToolchainVersion() string {
	return zerotern.VF(a.Toolchain, func() string {
		return "go" + a.Go // Add "go" prefix since version format is "1.22.8" // 添加 "go" 前缀，因为版本格式为 "1.22.8"
	})
}

// GetDirectRequires filters and returns just direct (non-indirect) dependencies
// Useful when updating explicit dependencies
//
// GetDirectRequires 过滤并返回仅直接（非间接）依赖
// 适用于仅更新显式声明的依赖
func (a *ModuleInfo) GetDirectRequires() []*Require {
	var directs []*Require
	for _, require := range a.Require {
		if !require.Indirect {
			directs = append(directs, require)
		}
	}
	return directs
}

// GetScopedRequires returns dependencies filtered according to the specified type
// Supports filtering: direct, indirect, and complete package types
//
// GetScopedRequires 返回按指定类别过滤的依赖
// 支持按直接、间接或所有依赖类型过滤
func (a *ModuleInfo) GetScopedRequires(cate DepCate) []*Require {
	var results []*Require
	for _, dep := range a.Require {
		switch cate {
		case DepCateDirect:
			if !dep.Indirect {
				results = append(results, dep)
			}
		case DepCateIndirect:
			if dep.Indirect {
				results = append(results, dep)
			}
		default:
			results = append(results, dep)
		}
	}
	return results
}

// GetModuleInfo executes 'go mod edit -json' and parses module information
// Returns structured data about the module and its dependencies
//
// GetModuleInfo 执行 'go mod edit -json' 并解析模块信息
// 返回关于模块及其依赖的结构化数据
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

// ParseModuleFile reads and parses the go.mod file using golang.org/x/mod/modfile
// Returns the parsed module file structure enabling advanced manipulation
//
// ParseModuleFile 使用 golang.org/x/mod/modfile 读取和解析 go.mod 文件
// 返回解析后的模块文件结构，用于高级操作
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
