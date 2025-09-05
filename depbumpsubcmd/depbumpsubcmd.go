// Package depbumpsubcmd: Command-line interface for dependency bump operations
// Provides Cobra-based CLI commands for module, direct, and comprehensive dependency updates
// Supports workspace operations with configurable filtering and update strategies
//
// depbumpsubcmd: 依赖升级操作的命令行接口
// 提供基于 Cobra 的 CLI 命令，用于模块、直接和全面的依赖更新
// 支持带有可配置过滤和更新策略的工作区操作
package depbumpsubcmd

import (
	"github.com/go-mate/depbump"
	"github.com/go-mate/go-work/worksexec"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/tern"
	"github.com/yyle88/zaplog"
)

var aliasesMap = map[string][]string{
	"direct":   {"directs"},
	"everyone": {"require", "requires"},
	"module":   {"modules"},
}

// NewUpdateCmd creates update commands and adds them to root command
// Provides both structured (update subcommand) and direct access commands
//
// NewUpdateCmd 创建更新命令并添加到根命令
// 提供结构化（update 子命令）和直接访问命令
func NewUpdateCmd(rootCmd *cobra.Command, config *worksexec.WorksExec) {
	// Create update subcommand group
	// 创建更新子命令组
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update dependencies",
		Long:  "Update dependencies with various strategies and filtering options.",
	}
	updateCmd.AddCommand(NewUpdateModuleCmd(config, "module"))
	updateCmd.AddCommand(NewUpdateDirectCmd(config, "direct"))
	updateCmd.AddCommand(NewUpdateEveryoneCmd(config, "everyone")) // Use "everyone" to avoid confusion with "all" // 使用 "everyone" 避免与 "all" 混淆

	// Add structured command
	// 添加结构化命令
	rootCmd.AddCommand(updateCmd)

	// Add direct access commands
	// 添加直接访问命令
	rootCmd.AddCommand(NewUpdateModuleCmd(config, "module"))
	rootCmd.AddCommand(NewUpdateDirectCmd(config, "direct"))
	rootCmd.AddCommand(NewUpdateEveryoneCmd(config, "everyone"))
}

// NewUpdateModuleCmd creates a command for updating Go modules in workspace
// Provides module-specific update functionality with configurable usage name
//
// NewUpdateModuleCmd 创建用于更新工作区中 Go 模块的命令
// 提供特定于模块的更新功能，带可配置的用法名称
func NewUpdateModuleCmd(config *worksexec.WorksExec, usageName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     usageName,
		Aliases: aliasesMap[usageName],
		Short:   "depbump module",
		Long:    "depbump module",
		Run: func(cmd *cobra.Command, args []string) {
			UpdateModules(config)
		},
	}
	return cmd
}

// UpdateModules performs comprehensive module updates across all workspaces
// Handles module info retrieval, toolchain detection, and cleanup operations
//
// UpdateModules 在所有工作区中执行全面的模块更新
// 处理模块信息检索、工具链检测和清理操作
func UpdateModules(config *worksexec.WorksExec) {
	for _, workspace := range config.GetWorkspaces() {
		for _, projectPath := range workspace.Projects {
			moduleInfo := rese.P1(depbump.GetModuleInfo(projectPath))
			updateModule(config.GetSubCommand(projectPath), projectPath, moduleInfo.GetToolchainVersion())
			must.Done(GoModTide(config.GetSubCommand(projectPath)))
		}
		if workspace.WorkRoot != "" {
			must.Done(GoWorkSync(config.GetSubCommand(workspace.WorkRoot)))
		}
	}
}

// updateModule executes go get -u for a single module with toolchain management
// Handles environment setup and output processing with success logging
//
// updateModule 为单个模块执行 go get -u，带工具链管理
// 处理环境设置和输出处理，带成功日志记录
func updateModule(execConfig *osexec.ExecConfig, projectPath string, toolchain string) {
	var success = true
	output := rese.V1(execConfig.NewConfig().
		WithEnvs([]string{"GOTOOLCHAIN=" + toolchain}). // Use project Go version to suppress dependency Go version requirements // 在升级时需要用项目的go版本号压制住依赖的go版本号
		WithPath(projectPath).
		WithMatchMore(true).
		WithMatchPipe(func(line string) bool {
			if upgradeInfo, matched := depbump.MatchUpgrade(line); matched {
				zaplog.SUG.Debugln("match-upgrade-output-message:", eroticgo.GREEN.Sprint(neatjsons.S(upgradeInfo)))
				return true
			}
			if warnMessage, matched := depbump.MatchToolchainVersionMismatch(line); matched {
				zaplog.SUG.Debugln("go-toolchain-mismatch-output:", eroticgo.RED.Sprint(neatjsons.S(warnMessage)))
				success = false
				return true
			}
			return false
		}).ExecInPipe("go", "get", "-u", "./..."))
	if success {
		zaplog.SUG.Debugln(string(output))
		zaplog.SUG.Infoln("success", eroticgo.RED.Sprint("success"))
	} else {
		zaplog.SUG.Warnln(string(output))
		zaplog.SUG.Warnln("warning", eroticgo.RED.Sprint("warning"))
	}
}

// NewUpdateDirectCmd creates a command for updating just direct dependencies
// Filters out indirect dependencies and provides selective update control
//
// NewUpdateDirectCmd 创建仅更新直接依赖的命令
// 过滤掉间接依赖并提供选择性更新控制
func NewUpdateDirectCmd(config *worksexec.WorksExec, usageName string) *cobra.Command {
	const usageNameLatest = "latest"

	updateDepsConfig := &depbump.UpdateDepsConfig{
		Cate: depbump.DepCateDirect,
		Mode: tern.BVV(usageName == usageNameLatest, depbump.GetModeLatest, depbump.GetModeUpdate),
	}
	cmd := &cobra.Command{
		Use:     usageName,
		Aliases: aliasesMap[usageName],
		Short:   "depbump direct (latest)",
		Long:    "depbump direct (latest)",
		Run: func(cmd *cobra.Command, args []string) {
			updateDeps(config, updateDepsConfig)
		},
	}
	setFlags(cmd, updateDepsConfig)

	if usageName != usageNameLatest {
		cmd.AddCommand(NewUpdateDirectCmd(config, usageNameLatest))
	}
	return cmd
}

// NewUpdateEveryoneCmd creates a command for updating all dependencies
// Updates both direct and indirect dependencies with comprehensive filtering options
//
// NewUpdateEveryoneCmd 创建用于更新所有依赖的命令
// 更新直接和间接依赖，带全面的过滤选项
func NewUpdateEveryoneCmd(config *worksexec.WorksExec, usageName string) *cobra.Command {
	const usageNameLatest = "latest"

	updateDepsConfig := &depbump.UpdateDepsConfig{
		Cate: depbump.DepCateEveryone,
		Mode: tern.BVV(usageName == usageNameLatest, depbump.GetModeLatest, depbump.GetModeUpdate),
	}
	cmd := &cobra.Command{
		Use:     usageName,
		Aliases: aliasesMap[usageName],
		Short:   "depbump require (latest)",
		Long:    "depbump require (latest)",
		Run: func(cmd *cobra.Command, args []string) {
			updateDeps(config, updateDepsConfig)
		},
	}
	setFlags(cmd, updateDepsConfig)

	if usageName != usageNameLatest {
		cmd.AddCommand(NewUpdateEveryoneCmd(config, usageNameLatest))
	}
	return cmd
}

// setFlags configures command-line flags for dependency filtering and source control
// Provides flags for GitLab/GitHub filtering and skip options
//
// setFlags 为依赖过滤和源代码控制配置命令行标志
// 提供 GitLab/GitHub 过滤和跳过选项的标志
func setFlags(cmd *cobra.Command, config *depbump.UpdateDepsConfig) {
	cmd.Flags().BoolVarP(&config.GitlabOnly, "gitlab-only", "", false, "gitlab only: only update gitlab dependencies")
	cmd.Flags().BoolVarP(&config.SkipGitlab, "skip-gitlab", "", false, "skip gitlab: skip update gitlab dependencies")
	cmd.Flags().BoolVarP(&config.GithubOnly, "github-only", "", false, "github only: only update github dependencies")
	cmd.Flags().BoolVarP(&config.SkipGithub, "skip-github", "", false, "skip github: skip update github dependencies")
}

// updateDeps executes dependency updates across workspaces with specified configuration
// Handles module information retrieval and orchestrates bulk updates with cleanup
//
// updateDeps 使用指定配置在工作区中执行依赖更新
// 处理模块信息检索并编排批量更新，包括清理操作
func updateDeps(config *worksexec.WorksExec, updateDepsConfig *depbump.UpdateDepsConfig) {
	zaplog.SUG.Debugln(neatjsons.S(updateDepsConfig))

	for _, workspace := range config.GetWorkspaces() {
		for _, projectPath := range workspace.Projects {
			depbump.UpdateDeps(config.GetSubCommand(projectPath), rese.P1(depbump.GetModuleInfo(projectPath)), updateDepsConfig)
			must.Done(GoModTide(config.GetSubCommand(projectPath)))
		}
		if workspace.WorkRoot != "" {
			must.Done(GoWorkSync(config.GetSubCommand(workspace.WorkRoot)))
		}
	}
}

// GoModTide executes go mod tidy with error handling and output logging
// Cleans up module dependencies and ensures consistency
//
// GoModTide 执行 go mod tidy，带有错误处理和输出日志
// 清理模块依赖并确保一致性
func GoModTide(execConfig *osexec.ExecConfig) error {
	output, err := execConfig.Exec("go", "mod", "tidy", "-e")
	if err != nil {
		if len(output) > 0 {
			zaplog.SUG.Warnln(string(output))
		}
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(string(output))
	return nil
}

// GoWorkSync executes go work sync command with error handling and output logging
// Synchronizes workspace configuration and updates dependency relationships
//
// GoWorkSync 执行 go work sync 命令，带错误处理和输出日志
// 同步工作区配置并更新依赖关系
func GoWorkSync(execConfig *osexec.ExecConfig) error {
	output, err := execConfig.Exec("go", "work", "sync")
	if err != nil {
		if len(output) > 0 {
			zaplog.SUG.Warnln(string(output))
		}
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(string(output))
	return nil
}
