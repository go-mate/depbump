// Package main: depbump command-line application main package
// Provides automatic package upgrade and management tools
// Supports workspace operations and configurable update strategies
//
// main: depbump 命令行工具入口点
// 提供 Go 模块的自动包升级和管理功能
// 支持工作区操作和可配置的更新策略
package main

import (
	"os"

	"github.com/go-mate/depbump/depbumpkitcmd"
	"github.com/go-mate/depbump/depbumpmodcmd"
	"github.com/go-mate/depbump/depbumpsubcmd"
	"github.com/go-mate/depbump/depsynctagcmd"
	"github.com/go-mate/go-work/workspath"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// main initializes and executes the depbump command with workspace configuration
// Sets up project path detection, workspace management, and command execution
// Commands: module, update (D/E/R), sync, bump
//
// main 初始化并执行 depbump 命令，配置工作区
// 设置项目路径检测、工作区管理和命令执行
// 命令：module、update (D/E/R)、sync、bump
func main() {
	// Get current working DIR
	// 获取当前工作 DIR
	currentPath := rese.C1(os.Getwd())
	zaplog.LOG.Debug("current:", zap.String("path", currentPath))

	// Get executable path
	// 获取可执行文件路径
	executePath := rese.C1(os.Executable())
	zaplog.LOG.Debug("execute:", zap.String("path", executePath))

	// Detect project path from current DIR
	// 从当前 DIR 检测项目路径
	pathInfo, ok := workspath.GetProjectPath(currentPath)
	must.True(ok)
	projectPath := must.Nice(pathInfo.ProjectPath)
	zaplog.LOG.Debug("project:", zap.String("path", projectPath))
	must.Nice(projectPath)

	// Initialize execution configuration with project path
	// 用项目路径初始化执行配置
	execConfig := osexec.NewCommandConfig().WithBash().WithDebug().WithPath(projectPath)

	// Create root command with default module update action
	// 创建根命令，默认执行模块更新操作
	rootCmd := &cobra.Command{
		Use:   "depbump",
		Short: "Go package management assistant",
		Long:  "Check and upgrade outdated dependencies in Go modules, with version bumping.",
		Run: func(cmd *cobra.Command, args []string) {
			depbumpmodcmd.UpdateModules(execConfig)
		},
	}

	// Add subcommands to root
	// 添加子命令到根命令
	rootCmd.AddCommand(depbumpmodcmd.NewModuleCmd(execConfig))
	rootCmd.AddCommand(depbumpsubcmd.NewUpdateCmd(execConfig))
	rootCmd.AddCommand(depsynctagcmd.NewSyncCmd(execConfig))
	rootCmd.AddCommand(depbumpkitcmd.NewBumpCmd(execConfig))

	// Execute CLI application
	// 执行 CLI 应用程序
	must.Done(rootCmd.Execute())
}
