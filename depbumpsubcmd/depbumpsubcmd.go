// Package depbumpsubcmd: Command-line interface to bump deps
// Provides Cobra-based CLI commands to handle module, direct, and comprehensive dep updates
// Supports workspace operations with configurable filtering and update strategies
//
// depbumpsubcmd: 包升级操作的命令行接口
// 提供基于 Cobra 的 CLI 命令，用于模块、直接和全面的依赖更新
// 支持带有可配置过滤和更新策略的工作区操作
package depbumpsubcmd

import (
	"github.com/go-mate/depbump"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/tern"
	"github.com/yyle88/zaplog"
)

// NewUpdateCmd creates update commands and adds them to root command
// Provides both structured (update subcommand) and direct access commands
//
// NewUpdateCmd 创建更新命令并添加到根命令
// 提供结构化（update 子命令）和直接访问命令
func NewUpdateCmd(rootCmd *cobra.Command, execConfig *osexec.ExecConfig) {
	// Create update subcommand group
	// 创建更新子命令组
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update dependencies",
		Long:  "Update dependencies with various strategies and filtering options.",
	}
	updateCmd.AddCommand(NewUpdateModuleCmd(execConfig, []string{"module", "modules"}))
	updateCmd.AddCommand(NewUpdateDirectCmd(execConfig, []string{"direct", "directs"}))
	updateCmd.AddCommand(NewUpdateEveryoneCmd(execConfig, []string{"everyone", "require", "requires"})) // Use "everyone" to avoid confusion with "each" // 使用 "everyone" 避免与 "each" 混淆

	// Add structured command
	// 添加结构化命令
	rootCmd.AddCommand(updateCmd)

	// Add direct access commands
	// 添加直接访问命令
	rootCmd.AddCommand(NewUpdateModuleCmd(execConfig, []string{"module", "modules"}))
	rootCmd.AddCommand(NewUpdateDirectCmd(execConfig, []string{"direct", "directs"}))
	rootCmd.AddCommand(NewUpdateEveryoneCmd(execConfig, []string{"everyone", "require", "requires"}))
}

// NewUpdateModuleCmd creates a command to update Go modules
// Provides module-specific update function with configurable usage name
//
// NewUpdateModuleCmd 创建用于更新 Go 模块的命令
// 提供特定于模块的更新功能，带可配置的用法名称
func NewUpdateModuleCmd(execConfig *osexec.ExecConfig, usageNames []string) *cobra.Command {
	must.Have(usageNames)

	cmd := &cobra.Command{
		Use:     usageNames[0],
		Aliases: usageNames[1:],
		Short:   "depbump module",
		Long:    "depbump module",
		Run: func(cmd *cobra.Command, args []string) {
			UpdateModules(execConfig)
		},
	}
	return cmd
}

// UpdateModules performs comprehensive module updates
// Handles module info fetch, toolchain detection, and cleanup operations
//
// UpdateModules 执行全面的模块更新
// 处理模块信息检索、工具链检测和清理操作
func UpdateModules(execConfig *osexec.ExecConfig) {
	projectDIR := osmustexist.ROOT(execConfig.Path)
	moduleInfo := rese.P1(depbump.GetModuleInfo(projectDIR))
	updateModule(execConfig, moduleInfo.GetToolchainVersion())
	must.Done(GoModTide(execConfig))
}

// updateModule executes go get -u for a single module with toolchain management
// Handles environment setup and output processing with success logging
//
// updateModule 为单个模块执行 go get -u，带工具链管理
// 处理环境设置和输出处理，带成功日志记录
func updateModule(execConfig *osexec.ExecConfig, toolchain string) {
	var success = true
	output := rese.V1(execConfig.NewConfig().
		WithEnvs([]string{"GOTOOLCHAIN=" + toolchain}). // Use project Go version to suppress package Go version requirements // 在升级时用项目的go版本要求压制包的go版本要求
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
			if sdkInfo, matched := depbump.MatchGoDownloadingSdkInfo(line); matched {
				zaplog.SUG.Debugln("go-downloading-sdk-info:", eroticgo.CYAN.Sprint(neatjsons.S(sdkInfo)))
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

// NewUpdateDirectCmd creates a command to update just direct dependencies
// Filters out indirect dependencies and provides selective update management
//
// NewUpdateDirectCmd 创建仅更新直接依赖的命令
// 过滤掉间接包并提供选择性更新控制
func NewUpdateDirectCmd(execConfig *osexec.ExecConfig, usageNames []string) *cobra.Command {
	usageName := must.Have(usageNames)[0]

	const usageNameLatest = "latest"

	updateDepsConfig := &depbump.UpdateDepsConfig{
		Cate: depbump.DepCateDirect,
		Mode: tern.BVV(usageName == usageNameLatest, depbump.GetModeLatest, depbump.GetModeUpdate),
	}
	cmd := &cobra.Command{
		Use:     usageNames[0],
		Aliases: usageNames[1:],
		Short:   "depbump direct (latest)",
		Long:    "depbump direct (latest)",
		Run: func(cmd *cobra.Command, args []string) {
			updateDeps(execConfig, updateDepsConfig)
		},
	}
	setFlags(cmd, updateDepsConfig)

	if usageName != usageNameLatest {
		cmd.AddCommand(NewUpdateDirectCmd(execConfig, []string{usageNameLatest}))
	}
	return cmd
}

// NewUpdateEveryoneCmd creates a command to update each package
// Updates both direct and indirect dependencies with comprehensive filtering options
//
// NewUpdateEveryoneCmd 创建用于更新所有依赖的命令
// 更新直接和间接包，带全面的过滤选项
func NewUpdateEveryoneCmd(execConfig *osexec.ExecConfig, usageNames []string) *cobra.Command {
	usageName := must.Have(usageNames)[0]

	const usageNameLatest = "latest"

	updateDepsConfig := &depbump.UpdateDepsConfig{
		Cate: depbump.DepCateEveryone,
		Mode: tern.BVV(usageName == usageNameLatest, depbump.GetModeLatest, depbump.GetModeUpdate),
	}
	cmd := &cobra.Command{
		Use:     usageNames[0],
		Aliases: usageNames[1:],
		Short:   "depbump require (latest)",
		Long:    "depbump require (latest)",
		Run: func(cmd *cobra.Command, args []string) {
			updateDeps(execConfig, updateDepsConfig)
		},
	}
	setFlags(cmd, updateDepsConfig)

	if usageName != usageNameLatest {
		cmd.AddCommand(NewUpdateEveryoneCmd(execConfig, []string{usageNameLatest}))
	}
	return cmd
}

// setFlags configures command-line flags that handle package filtering and source management
// Provides flags to handle GitLab/GitHub filtering and skip options
//
// setFlags 为包过滤和源代码控制配置命令行标志
// 提供 GitLab/GitHub 过滤和跳过选项的标志
func setFlags(cmd *cobra.Command, config *depbump.UpdateDepsConfig) {
	cmd.Flags().BoolVarP(&config.GitlabOnly, "gitlab-only", "", false, "gitlab exclusive: update gitlab dependencies")
	cmd.Flags().BoolVarP(&config.SkipGitlab, "skip-gitlab", "", false, "skip gitlab: skip update gitlab dependencies")
	cmd.Flags().BoolVarP(&config.GithubOnly, "github-only", "", false, "github exclusive: update github dependencies")
	cmd.Flags().BoolVarP(&config.SkipGithub, "skip-github", "", false, "skip github: skip update github dependencies")
}

// updateDeps executes package updates with specified configuration
// Handles module information access and orchestrates updates with cleanup
//
// updateDeps 使用指定配置执行包更新
// 处理模块信息检索并编排更新，包括清理操作
func updateDeps(execConfig *osexec.ExecConfig, updateDepsConfig *depbump.UpdateDepsConfig) {
	zaplog.SUG.Debugln(neatjsons.S(updateDepsConfig))

	projectDIR := osmustexist.ROOT(execConfig.Path)
	depbump.UpdateDeps(execConfig, rese.P1(depbump.GetModuleInfo(projectDIR)), updateDepsConfig)
	must.Done(GoModTide(execConfig))
}

// GoModTide executes go mod cleanup with error handling and output logging
// Cleans up module dependencies and ensures stable state
//
// GoModTide 执行 go mod cleanup，带有错误处理和输出日志
// 清理模块包并确保一致性
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
// Synchronizes workspace configuration and updates package relationships
//
// GoWorkSync 执行 go work sync 命令，带错误处理和输出日志
// 同步工作区配置并更新包关系
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
