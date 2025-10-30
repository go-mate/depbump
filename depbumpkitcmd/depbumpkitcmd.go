// Package depbumpkitcmd: Package checking and synchronization within Go modules
// Provides intelligent package upgrade tools that prevent Go toolchain contagion
// Implements version analysis and selective upgrades while maintaining Go version matching
// Supports upgrade-first method preventing package downgrades in production systems
//
// depbumpkitcmd: Go 模块的包兼容性检查和同步
// 提供智能包升级工具，防止 Go 工具链传染
// 实现版本分析和选择性升级，同时保持 Go 版本兼容性
// 支持仅升级策略，防止生产系统中的包降级
package depbumpkitcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-mate/depbump"
	"github.com/go-mate/depbump/internal/utils"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"golang.org/x/mod/modfile"
)

// SetupBumpCmd creates and configures the bump command handling package management
// Provides intelligent package analysis and upgrade capabilities
//
// SetupBumpCmd 创建并配置用于包管理的 bump 命令
// 提供智能包分析和升级功能
func SetupBumpCmd(rootCmd *cobra.Command, execConfig *osexec.ExecConfig) {
	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Bump dependencies to stable versions with Go version matching",
		Run: func(cmd *cobra.Command, args []string) {
			kit := NewBumpKit(execConfig)
			kit.SyncDependencies(&BumpDepsConfig{
				Cate: depbump.DepCateDirect,
				Mode: depbump.GetModeUpdate, // Default: stable versions within // 默认：仅稳定版本
			})
		},
	}

	// Add direct subcommand with latest sub-subcommand
	directCmd := &cobra.Command{
		Use:     "direct",
		Aliases: []string{"directs"},
		Short:   "Bump direct dependencies to stable versions",
		Run: func(cmd *cobra.Command, args []string) {
			kit := NewBumpKit(execConfig)
			kit.SyncDependencies(&BumpDepsConfig{
				Cate: depbump.DepCateDirect,
				Mode: depbump.GetModeUpdate, // Stable versions within // 仅稳定版本
			})
		},
	}
	directCmd.AddCommand(&cobra.Command{
		Use:   "latest",
		Short: "Bump direct dependencies to latest versions (including prerelease)",
		Run: func(cmd *cobra.Command, args []string) {
			kit := NewBumpKit(execConfig)
			kit.SyncDependencies(&BumpDepsConfig{
				Cate: depbump.DepCateDirect,
				Mode: depbump.GetModeLatest, // All versions // 所有版本
			})
		},
	})
	cmd.AddCommand(directCmd)

	// Add everyone subcommand with latest sub-subcommand
	everyoneCmd := &cobra.Command{
		Use:     "everyone",
		Aliases: []string{"require", "requires"},
		Short:   "Bump each dependencies to stable versions",
		Run: func(cmd *cobra.Command, args []string) {
			kit := NewBumpKit(execConfig)
			kit.SyncDependencies(&BumpDepsConfig{
				Cate: depbump.DepCateEveryone,
				Mode: depbump.GetModeUpdate, // Stable versions within // 仅稳定版本
			})
		},
	}
	everyoneCmd.AddCommand(&cobra.Command{
		Use:   "latest",
		Short: "Bump each dependencies to latest versions (including prerelease)",
		Run: func(cmd *cobra.Command, args []string) {
			kit := NewBumpKit(execConfig)
			kit.SyncDependencies(&BumpDepsConfig{
				Cate: depbump.DepCateEveryone,
				Mode: depbump.GetModeLatest, // All versions // 所有版本
			})
		},
	})
	cmd.AddCommand(everyoneCmd)

	rootCmd.AddCommand(cmd)
}

// BumpDepsConfig provides configuration needed in intelligent package bump operations
// Controls package types and upgrade actions with Go version matching
//
// BumpDepsConfig 为智能包升级操作提供配置
// 控制包类别和带 Go 版本匹配的升级行为
type BumpDepsConfig struct {
	Cate depbump.DepCate // Package type used in bump operations // 升级操作的包类型
	Mode depbump.GetMode // Version selection mode // 版本选择模式
}

// BumpKit handles package matching validation and intelligent upgrades
// Manages Go version requirements and package version resolution
// Implements caching mechanisms enabling efficient package analysis
//
// BumpKit 处理包兼容性验证和智能升级
// 管理 Go 版本要求和包版本解析
// 实现缓存机制以提高包分析效率
type BumpKit struct {
	TargetGoVersion string                // Target Go version during matching checks // 目标 Go 版本用于匹配检查
	MapDepGoVersion map[string]string     // Cache containing package Go version requirements // 包 Go 版本要求的缓存
	execConfig      *osexec.CommandConfig // Execution configuration handling command operations // 命令操作的执行配置
}

// NewBumpKit creates a new package matching engine with toolchain analysis
// Extracts target Go version from module toolchain configuration
// Initializes caching system enabling efficient package analysis
//
// NewBumpKit 创建新的包兼容性验证器，带有工具链分析
// 从模块工具链配置中提取目标 Go 版本
// 初始化缓存系统以实现高效的包分析
func NewBumpKit(execConfig *osexec.ExecConfig) *BumpKit {
	projectDIR := osmustexist.ROOT(execConfig.Path)

	moduleInfo := rese.P1(depbump.GetModuleInfo(projectDIR))
	// 获取有效的工具链版本，考虑 toolchain 字段
	toolchainVersion := moduleInfo.GetToolchainVersion()
	// 去掉 "go" 前缀，只保留版本号用于比较
	targetGoVersion := strings.TrimPrefix(toolchainVersion, "go")

	return &BumpKit{
		TargetGoVersion: targetGoVersion,
		MapDepGoVersion: make(map[string]string),
		execConfig:      execConfig,
	}
}

// SyncDependencies performs package analysis and applies intelligent upgrades
// Analyzes packages based on configuration during matching and version optimization
// Applies matching upgrades to prevent toolchain version conflicts
//
// SyncDependencies 执行包分析并应用智能升级
// 根据配置分析包的兼容性和版本优化
// 仅应用兼容的升级以防止工具链版本冲突
func (c *BumpKit) SyncDependencies(config *BumpDepsConfig) {
	zaplog.SUG.Infoln("starting", string(config.Cate), "dependencies analysis - Go", eroticgo.CYAN.Sprint(c.TargetGoVersion))
	deps := c.AnalyzeDependencies(config.Cate, config.Mode)
	zaplog.SUG.Debugln("analysis result", neatjsons.S(deps))

	zaplog.SUG.Infoln("🔧 Applying", string(config.Cate), "updates...")
	must.Done(c.ApplyUpdates(deps))
	zaplog.SUG.Infoln("✅", string(config.Cate), "updates success!")
}

// DependencyInfo represents comprehensive information about a package upgrade
// Contains version transition details and Go version requirements
// Used during analysis reporting and upgrade decision making
//
// DependencyInfo 表示包升级的全面信息
// 包含版本转换详情和 Go 版本要求
// 用于分析报告和升级决策制定
type DependencyInfo struct {
	Package       string
	OldDepVersion string
	NewDepVersion string
	NewGoVersion  string // Go version required in new package version // 新包版本需要的 Go 版本
}

// AnalyzeDependencies performs comprehensive analysis of dependencies according to type
// Evaluates each package during upgrades within Go version constraints
// Returns detailed upgrade recommendations with version matching information
//
// AnalyzeDependencies 根据类别对包执行全面分析
// 在 Go 版本约束内评估每个包的潜在升级
// 返回带有版本兼容性信息的详细升级建议
func (c *BumpKit) AnalyzeDependencies(cate depbump.DepCate, mode depbump.GetMode) []*DependencyInfo {
	projectDIR := osmustexist.ROOT(c.execConfig.Path)

	moduleInfo := rese.P1(depbump.GetModuleInfo(projectDIR))
	requires := moduleInfo.GetScopedRequires(cate)

	deps := make([]*DependencyInfo, 0, len(requires))
	zaplog.SUG.Infoln("analyzing", eroticgo.CYAN.Sprint(len(requires)), string(cate), "dependencies")

	for idx, req := range requires {
		// 显示进度以提升用户体验
		progress := fmt.Sprintf("(%d/%d)", idx+1, len(requires))
		zaplog.SUG.Infoln(progress, "analyzing", eroticgo.GREEN.Sprint(req.Path))

		versions := c.GetVersionList(req.Path)
		if len(versions) == 0 {
			continue
		}

		packageVersion := c.SelectBestPackageVersion(req.Path, versions, req.Version, mode)

		dep := &DependencyInfo{
			Package:       req.Path,
			OldDepVersion: req.Version,
			NewDepVersion: packageVersion.Version,
			NewGoVersion:  packageVersion.GoVersion,
		}

		if dep.OldDepVersion != dep.NewDepVersion {
			zaplog.SUG.Debugln("update recommended", eroticgo.GREEN.Sprint(neatjsons.S(dep)))
		}

		deps = append(deps, dep)
	}

	return deps
}

// BestPackageVersion contains the result of intelligent version selection
// Represents the best version choice with associated Go version requirements
// Used to communicate version analysis results between functions
//
// BestPackageVersion 包含智能版本选择的结果
// 表示最优版本选择及其关联的 Go 版本要求
// 用于在函数之间传递版本分析结果
type BestPackageVersion struct {
	Version   string // Selected version // 选中的版本
	GoVersion string // Required Go version // 需要的 Go 版本
}

// SelectBestPackageVersion finds the best matching version within a given package
// Implements upgrade-first method while respecting Go version matching constraints
// Returns the best available version, maintains current version if no upgrade possible
//
// SelectBestPackageVersion 为给定包找到最优兼容版本
// 实现仅升级方式，同时遵守 Go 版本兼容性约束
// 返回最佳可用版本，如果无法升级则保持当前版本
func (c *BumpKit) SelectBestPackageVersion(pkg string, versions []string, currentVersion string, mode depbump.GetMode) *BestPackageVersion {
	osmustexist.ROOT(c.execConfig.Path)

	// 找到当前版本在列表中的位置
	currentIndex := -1
	for i, version := range versions {
		if version == currentVersion {
			currentIndex = i
			break
		}
	}

	// 从当前版本开始，向上寻找兼容的版本（只升级，不降级）
	for i := 0; i <= currentIndex; i++ {
		version := versions[i]
		zaplog.SUG.Debugln("checking", eroticgo.CYAN.Sprint(version))

		// Skip unstable versions when mode is UPDATE
		// 当模式为 UPDATE 时跳过不稳定版本
		if mode == depbump.GetModeUpdate && !utils.IsStableVersion(version) {
			zaplog.SUG.Debugln("skip unstable version", eroticgo.YELLOW.Sprint(version))
			continue
		}

		goReq := c.GetPackageGoRequirement(pkg, version)
		if utils.CanUseGoVersion(goReq, c.TargetGoVersion) {
			// 只有当找到的版本比当前版本新或相等时才返回
			if utils.CompareVersions(version, currentVersion) >= 0 {
				packageVersion := &BestPackageVersion{
					Version:   version,
					GoVersion: goReq,
				}
				zaplog.SUG.Debugln("found best version", eroticgo.GREEN.Sprint(neatjsons.S(packageVersion)))
				return packageVersion
			}
		}
	}

	// 如果没有找到兼容的更新版本，保持当前版本，返回当前版本的 Go 要求
	packageVersion := &BestPackageVersion{
		Version:   currentVersion,
		GoVersion: c.GetPackageGoRequirement(pkg, currentVersion),
	}
	zaplog.SUG.Debugln("keep current version", eroticgo.YELLOW.Sprint(neatjsons.S(packageVersion)))
	return packageVersion
}

// GetVersionList retrieves and sorts available versions within a package
// Uses Go module system to fetch version information from package repositories
// Returns versions sorted in descending sequence enabling efficient newest-first processing
//
// GetVersionList 检索并排序包的所有可用版本
// 使用 Go 模块系统从包仓库获取版本信息
// 返回按降序排列的版本，以实现高效的最新版本优先处理
func (c *BumpKit) GetVersionList(pkg string) []string {
	osmustexist.ROOT(c.execConfig.Path)

	zaplog.SUG.Debugln("fetching versions", eroticgo.CYAN.Sprint(pkg))

	output, err := c.execConfig.Exec("go", "list", "-m", "-versions", pkg)
	if err != nil {
		zaplog.SUG.Warnln("failed to get versions", eroticgo.RED.Sprint(pkg), err.Error())
		return nil
	}

	parts := strings.Fields(string(output))
	if len(parts) <= 1 {
		return nil
	}

	versions := parts[1:]
	sort.Slice(versions, func(i, j int) bool {
		return utils.CompareVersions(versions[i], versions[j]) > 0
	})
	return versions
}

// GetPackageGoRequirement determines the Go version requirement within a specific package version
// Downloads and analyzes go.mod files to extract toolchain and Go version constraints
// Implements intelligent caching to minimize redundant package downloads
// Handles old packages without go.mod files with sensible defaults
//
// GetPackageGoRequirement 确定特定包版本的 Go 版本要求
// 下载并分析 go.mod 文件以提取工具链和 Go 版本约束
// 实现智能缓存以最小化冗余包下载
// 优雅处理没有 go.mod 文件的旧版包，提供合理的默认值
func (c *BumpKit) GetPackageGoRequirement(pkgPath, version string) string {
	osmustexist.ROOT(c.execConfig.Path)

	cacheKey := fmt.Sprintf("%s@%s", pkgPath, version)
	if cached, exists := c.MapDepGoVersion[cacheKey]; exists {
		return cached
	}

	zaplog.SUG.Debugln("downloading", eroticgo.CYAN.Sprint(pkgPath+"@"+version))

	// 直接获取模块的 go.mod 信息
	output, err := c.execConfig.Exec("go", "mod", "download", "-json", pkgPath+"@"+version)
	if err != nil {
		zaplog.SUG.Warnln("download failed", eroticgo.RED.Sprint(pkgPath+"@"+version), err.Error())
		return ""
	}

	var modInfo struct {
		GoMod string `json:"GoMod"`
	}
	must.Done(json.Unmarshal(output, &modInfo))

	var goReq string
	const defaultVersion = "1.0.0"
	if modInfo.GoMod == "" {
		// No go.mod file, use default version // 没有 go.mod 文件，使用默认版本
		goReq = defaultVersion
	} else {
		// Parse downloaded go.mod file // 解析下载的 go.mod 文件
		modData := rese.A1(os.ReadFile(modInfo.GoMod))
		modFile := rese.P1(modfile.Parse("go.mod", modData, nil))

		// Get effective toolchain version, considering toolchain field
		// 获取有效的工具链版本，考虑 toolchain 传染
		if modFile.Toolchain != nil {
			goReq = strings.TrimPrefix(modFile.Toolchain.Name, "go")
		} else if modFile.Go != nil {
			goReq = must.Nice(modFile.Go.Version)
		} else {
			// No go directive in go.mod, use default version // go.mod 中没有 go 指令，使用默认版本
			goReq = defaultVersion
		}
	}
	c.MapDepGoVersion[cacheKey] = goReq
	return goReq
}

// ApplyUpdates applies validated package updates to the current module
// Executes go get commands during each approved package upgrade
// Performs module cleanup to ensure consistent package state
//
// ApplyUpdates 将已验证的包更新应用到当前模块
// 为每个批准的依赖升级执行 go get 命令
// 执行模块清理以确保一致的依赖状态
func (c *BumpKit) ApplyUpdates(deps []*DependencyInfo) error {
	osmustexist.ROOT(c.execConfig.Path)

	for _, dep := range deps {
		if dep.OldDepVersion != dep.NewDepVersion {
			zaplog.SUG.Debugln("updating", eroticgo.GREEN.Sprint(dep.Package))

			_, err := c.execConfig.Exec("go", "get", dep.Package+"@"+dep.NewDepVersion)
			if err != nil {
				zaplog.SUG.Warnln("update failed", eroticgo.RED.Sprint(dep.Package))
				continue
			}
		}
	}

	zaplog.SUG.Infoln("cleaning up module dependencies")
	rese.V1(c.execConfig.Exec("go", "mod", "tidy"))
	return nil
}
