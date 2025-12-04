// Package depbumpsubcmd: Command-line interface to bump deps
// Provides Cobra-based CLI commands to handle direct and comprehensive dep updates
// Supports workspace operations with configurable filtering and update strategies
//
// depbumpsubcmd: 包升级操作的命令行接口
// 提供基于 Cobra 的 CLI 命令，用于直接和全面的依赖更新
// 支持带有可配置过滤和更新策略的工作区操作
package depbumpsubcmd

import (
	"fmt"

	"github.com/go-mate/depbump"
	"github.com/go-mate/go-work/workspath"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

// NewUpdateCmd creates update command with D/E/R subcommands
//
// NewUpdateCmd 创建 update 命令，包含 D/E/R 子命令
func NewUpdateCmd(execConfig *osexec.ExecConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update dependencies",
		Long:  "Update dependencies with various strategies and filtering options.",
	}
	cmd.AddCommand(NewDirectCmd(execConfig))
	cmd.AddCommand(NewEveryoneCmd(execConfig))
	cmd.AddCommand(NewRecursiveCmd(execConfig))
	return cmd
}

// ============================================================================
// Direct Commands: D, D L, D R
// 直接依赖命令：D, D L, D R
// ============================================================================

// NewDirectCmd creates command to update direct dependencies
// Contains subcommands: latest (L), recursive (R)
// Uses PersistentFlags so children inherit config
//
// NewDirectCmd 创建更新直接依赖的命令
// 包含子命令：latest (L), recursive (R)
// 使用 PersistentFlags 让子命令继承配置
func NewDirectCmd(execConfig *osexec.ExecConfig) *cobra.Command {
	config := &depbump.UpdateDepsConfig{
		Cate: depbump.DepCateDirect,
		Mode: depbump.GetModeUpdate,
	}
	cmd := &cobra.Command{
		Use:     "direct",
		Aliases: []string{"D", "directs"},
		Short:   "Update direct dependencies",
		Long:    "Update direct dependencies to stable versions.",
		Run: func(cmd *cobra.Command, args []string) {
			updateDeps(execConfig, config)
		},
	}
	setPersistentFlags(cmd, config)

	// D L: direct latest // D L：直接依赖最新版
	{
		cmd.AddCommand(&cobra.Command{
			Use:     "latest",
			Aliases: []string{"L"},
			Short:   "Update direct dependencies to latest versions",
			Long:    "Update direct dependencies to latest versions (including prerelease).",
			Args:    cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				config.Cate = depbump.DepCateDirect
				config.Mode = depbump.GetModeLatest
				updateDeps(execConfig, config)
			},
		})
	}

	// D R: direct recursive // D R：直接依赖递归
	{
		cmd.AddCommand(&cobra.Command{
			Use:     "recursive",
			Aliases: []string{"R"},
			Short:   "Update direct dependencies across workspace",
			Long:    "Update direct dependencies across all modules in the workspace.",
			Args:    cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				config.Cate = depbump.DepCateDirect
				config.Mode = depbump.GetModeUpdate
				updateDepsRecursive(execConfig, config)
			},
		})
	}

	return cmd
}

// ============================================================================
// Everyone Commands: E, E L, E R
// 全部依赖命令：E, E L, E R
// ============================================================================

// NewEveryoneCmd creates command to update all dependencies
// Contains subcommands: latest (L), recursive (R)
// Uses PersistentFlags so children inherit config
//
// NewEveryoneCmd 创建更新全部依赖的命令
// 包含子命令：latest (L), recursive (R)
// 使用 PersistentFlags 让子命令继承配置
func NewEveryoneCmd(execConfig *osexec.ExecConfig) *cobra.Command {
	config := &depbump.UpdateDepsConfig{
		Cate: depbump.DepCateEveryone,
		Mode: depbump.GetModeUpdate,
	}
	cmd := &cobra.Command{
		Use:     "everyone",
		Aliases: []string{"E", "each"},
		Short:   "Update all dependencies",
		Long:    "Update all dependencies (direct + indirect) to stable versions.",
		Run: func(cmd *cobra.Command, args []string) {
			updateDeps(execConfig, config)
		},
	}
	setPersistentFlags(cmd, config)

	// E L: everyone latest // E L：全部依赖最新版
	{
		cmd.AddCommand(&cobra.Command{
			Use:     "latest",
			Aliases: []string{"L"},
			Short:   "Update all dependencies to latest versions",
			Long:    "Update all dependencies (direct + indirect) to latest versions (including prerelease).",
			Args:    cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				config.Cate = depbump.DepCateEveryone
				config.Mode = depbump.GetModeLatest
				updateDeps(execConfig, config)
			},
		})
	}

	// E R: everyone recursive // E R：全部依赖递归
	{
		cmd.AddCommand(&cobra.Command{
			Use:     "recursive",
			Aliases: []string{"R"},
			Short:   "Update all dependencies across workspace",
			Long:    "Update all dependencies (direct + indirect) across all modules in the workspace.",
			Args:    cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				config.Cate = depbump.DepCateEveryone
				config.Mode = depbump.GetModeUpdate
				updateDepsRecursive(execConfig, config)
			},
		})
	}

	return cmd
}

// ============================================================================
// Recursive Commands: R, R D, R E
// 递归命令：R, R D, R E
// ============================================================================

// NewRecursiveCmd creates recursive command to update dependencies across workspace
// Contains subcommands: direct (D), everyone (E)
// Uses PersistentFlags so children inherit config
//
// NewRecursiveCmd 创建递归命令，在工作区中更新依赖
// 包含子命令：direct (D), everyone (E)
// 使用 PersistentFlags 让子命令继承配置
func NewRecursiveCmd(execConfig *osexec.ExecConfig) *cobra.Command {
	config := &depbump.UpdateDepsConfig{
		Cate: depbump.DepCateDirect,
		Mode: depbump.GetModeUpdate,
	}
	cmd := &cobra.Command{
		Use:     "recursive",
		Aliases: []string{"R"},
		Short:   "Update direct dependencies across workspace",
		Long:    "Update direct dependencies across all modules in the workspace. Executes 'depbump D' on each module.",
		Run: func(cmd *cobra.Command, args []string) {
			updateDepsRecursive(execConfig, config)
		},
	}
	setPersistentFlags(cmd, config)

	// R D: recursive direct // R D：递归直接依赖更新
	{
		cmd.AddCommand(&cobra.Command{
			Use:     "direct",
			Aliases: []string{"D", "directs"},
			Short:   "Update direct dependencies across workspace",
			Long:    "Update direct dependencies across all modules in the workspace.",
			Args:    cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				config.Cate = depbump.DepCateDirect
				config.Mode = depbump.GetModeUpdate
				updateDepsRecursive(execConfig, config)
			},
		})
	}

	// R E: recursive everyone // R E：递归全部依赖更新
	{
		cmd.AddCommand(&cobra.Command{
			Use:     "everyone",
			Aliases: []string{"E", "each"},
			Short:   "Update all dependencies across workspace",
			Long:    "Update all dependencies (direct + indirect) across all modules in the workspace.",
			Args:    cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				config.Cate = depbump.DepCateEveryone
				config.Mode = depbump.GetModeUpdate
				updateDepsRecursive(execConfig, config)
			},
		})
	}

	return cmd
}

// ============================================================================
// Core Functions
// 核心函数
// ============================================================================

// setPersistentFlags configures persistent command-line flags inherited by children
//
// setPersistentFlags 配置持久化命令行标志，子命令自动继承
func setPersistentFlags(cmd *cobra.Command, config *depbump.UpdateDepsConfig) {
	cmd.PersistentFlags().BoolVarP(&config.GitlabOnly, "gitlab-only", "", false, "Update gitlab dependencies")
	cmd.PersistentFlags().BoolVarP(&config.SkipGitlab, "skip-gitlab", "", false, "Skip gitlab dependencies")
	cmd.PersistentFlags().BoolVarP(&config.GithubOnly, "github-only", "", false, "Update github dependencies")
	cmd.PersistentFlags().BoolVarP(&config.SkipGithub, "skip-github", "", false, "Skip github dependencies")
}

// updateDeps executes package updates with specified configuration
//
// updateDeps 使用指定配置执行包更新
func updateDeps(execConfig *osexec.ExecConfig, config *depbump.UpdateDepsConfig) {
	zaplog.SUG.Debugln(neatjsons.S(config))

	projectDIR := osmustexist.ROOT(execConfig.Path)
	depbump.UpdateDeps(execConfig, rese.P1(depbump.GetModuleInfo(projectDIR)), config)
	rese.V1(execConfig.Exec("go", "mod", "tidy", "-e"))
}

// updateDepsRecursive executes package updates across workspace modules
//
// updateDepsRecursive 在工作区模块中执行包更新
func updateDepsRecursive(execConfig *osexec.ExecConfig, config *depbump.UpdateDepsConfig) {
	workPath := osmustexist.ROOT(execConfig.Path)

	options := workspath.NewOptions().
		WithIncludeCurrentProject(true).
		WithIncludeCurrentPackage(false).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(false)

	moduleRoots := workspath.GetModulePaths(workPath, options)

	zaplog.SUG.Infoln("Recursive mode: found", eroticgo.CYAN.Sprint(len(moduleRoots)), "modules")

	for idx, modulePath := range moduleRoots {
		zaplog.SUG.Infoln("Module", eroticgo.GREEN.Sprint(fmt.Sprintf("(%d/%d)", idx+1, len(moduleRoots))), "Processing:", eroticgo.CYAN.Sprint(modulePath))

		moduleExecConfig := execConfig.NewConfig().WithPath(modulePath)
		updateDeps(moduleExecConfig, config)
	}

	zaplog.SUG.Infoln("✅ Recursive updates completed!")
}
