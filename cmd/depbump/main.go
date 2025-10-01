// Package main: depbump command-line application main package
// Provides automatic package upgrade and management tools
// Supports workspace operations and configurable update strategies
//
// main: depbump 命令行工具入口点
// 为 Go 模块提供自动包升级和管理功能
// 支持工作区操作和可配置的更新策略
package main

import (
	"os"

	"github.com/go-mate/depbump/depbumpkitcmd"
	"github.com/go-mate/depbump/depbumpsubcmd"
	"github.com/go-mate/depbump/depsynctagcmd"
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
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
// Supports various update modes: basic, direct, comprehensive, intelligent with filtering options
//
// main 初始化并执行 depbump 命令，配置工作区
// 设置项目路径检测、工作区管理和命令执行
// 支持各种更新模式：基础、直接、全部、智能，带过滤选项
//
// Usage examples:
// go run main.go
// go run main.go direct
// go run main.go direct --gitlab-only
// go run main.go direct --skip-gitlab
// go run main.go direct --github-only
// go run main.go direct --skip-github
// go run main.go sync
// go run main.go sync tags
// go run main.go sync subs
// go run main.go bump
// go run main.go bump direct
// go run main.go bump direct latest
// go run main.go bump everyone
// go run main.go bump everyone latest
func main() {
	currentPath := rese.C1(os.Getwd())
	zaplog.LOG.Debug("current:", zap.String("path", currentPath))

	executePath := rese.C1(os.Executable())
	zaplog.LOG.Debug("execute:", zap.String("path", executePath))

	projectPath, _, ok := workspath.GetProjectPath(currentPath)
	must.True(ok)
	zaplog.LOG.Debug("project:", zap.String("path", projectPath))
	must.Nice(projectPath)

	execConfig := osexec.NewCommandConfig().WithBash().WithDebug()
	workspaces := []*workspace.Workspace{
		workspace.NewWorkspace("", []string{projectPath}),
	}
	config := worksexec.NewWorksExec(execConfig, workspaces)

	rootCmd := &cobra.Command{
		Use:   "depbump",
		Short: "Go package management assistant",
		Long:  "Check and upgrade outdated dependencies in Go modules, with version bumping.",
		Run: func(cmd *cobra.Command, args []string) {
			depbumpsubcmd.UpdateModules(config)
		},
	}

	depbumpsubcmd.NewUpdateCmd(rootCmd, config)
	depsynctagcmd.SetupSyncCmd(rootCmd, config)
	depbumpkitcmd.SetupBumpCmd(rootCmd, config)

	must.Done(rootCmd.Execute())
}
