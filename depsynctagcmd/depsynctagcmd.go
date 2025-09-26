// Package depsynctagcmd: Git tag synchronization helping package management
// Provides commands that sync package versions with Git tags across workspace
// Supports latest tag resolution and selective package synchronization
//
// depsynctagcmd: 用于依赖管理的 Git 标签同步
// 提供在工作区中同步依赖版本与 Git 标签的命令
// 支持最新标签解析和选择性依赖同步
package depsynctagcmd

import (
	"github.com/go-mate/depbump"
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/go-xlan/gitgo"
	"github.com/spf13/cobra"
	"github.com/yyle88/done"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

// SetupSyncCmd creates sync command and adds it to root command
// Provides basic go work sync features and tag-based synchronization subcommands
//
// SetupSyncCmd 创建同步命令并添加到根命令
// 提供基本的 go work sync 功能和基于标签的同步子命令
func SetupSyncCmd(rootCmd *cobra.Command, config *worksexec.WorksExec) {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "go workspace sync",
		Long:  "go workspace sync",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(config.ForeachWorkRun(func(execConfig *osexec.ExecConfig, wsp *workspace.Workspace) error {
				output, err := execConfig.Exec("go", "work", "sync")
				if err != nil {
					return erero.Wro(err)
				}
				zaplog.SUG.Debugln(string(output))
				return nil
			}))
		},
	}
	cmd.AddCommand(SyncTagsCmd(config))
	cmd.AddCommand(SyncSubsCmd(config))

	rootCmd.AddCommand(cmd)
}

// SyncTagsCmd creates command that synchronizes dependencies to the latest Git tags
// Updates dependencies to match corresponding Git tag versions
//
// SyncTagsCmd 创建用于将依赖同步到最新 Git 标签的命令
// 更新依赖以匹配其相应的 Git 标签版本
func SyncTagsCmd(config *worksexec.WorksExec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "go workspace sync tags",
		Long:  "go workspace sync tags",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(SyncTags(config, false))
		},
	}
	return cmd
}

// SyncSubsCmd creates command for syncing dependencies with latest tag fallback
// Uses latest tag when no specific tag is available for dependencies
//
// SyncSubsCmd 创建用于同步依赖的命令，带有最新标签回退
// 当依赖没有特定标签时使用最新标签
func SyncSubsCmd(config *worksexec.WorksExec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subs",
		Short: "go workspace sync subs",
		Long:  "go workspace sync subs",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(SyncTags(config, true))
		},
	}
	return cmd
}

// SyncTags performs Git tag-based package synchronization across workspace
// Compares current package versions with Git tags and updates when different
//
// SyncTags 在工作区中执行基于 Git 标签的依赖同步
// 比较当前依赖版本与 Git 标签，在不同时进行更新
func SyncTags(config *worksexec.WorksExec, useLatest bool) error {
	pkgTagsMap := GetPkgTagsMap(config)
	zaplog.SUG.Debugln(neatjsons.S(pkgTagsMap))

	for _, projectPath := range config.Subprojects() {
		moduleInfo := done.VCE(depbump.GetModuleInfo(projectPath)).Nice()
		zaplog.SUG.Debugln(moduleInfo.Module.Path)

		for _, module := range moduleInfo.Require {
			if module.Indirect {
				continue
			}

			pkgTag, ok := pkgTagsMap[module.Path]
			if !ok {
				continue
			}
			if pkgTag == "" {
				if useLatest {
					pkgTag = "latest" // Use "latest" when no tag version exists // 当没有标签版本时使用 "latest"
				} else {
					continue
				}
			}

			if pkgTag == module.Version {
				zaplog.SUG.Debugln(projectPath, module.Path, module.Version, "same")
				continue
			}
			zaplog.SUG.Debugln(projectPath, module.Path, module.Version, "sync", "=>", pkgTag)

			// Example command execution patterns:
			// GOTOOLCHAIN=go1.22.8 go get -u github.com/yyle88/syntaxgo@v0.0.45
			// go: golang.org/x/exp@v0.0.0-20250218142911-aa4b98e5adaa requires go >= 1.23.0 (running go 1.22.8; GOTOOLCHAIN=go1.22.8)
			// Correct approach: omit -u option to avoid version conflicts
			// GOTOOLCHAIN=go1.22.8 go get github.com/yyle88/syntaxgo@v0.0.45
			// go: upgraded github.com/yyle88/syntaxgo v0.0.44 => v0.0.45
			//
			// 命令执行模式示例：
			// GOTOOLCHAIN=go1.22.8 go get -u github.com/yyle88/syntaxgo@v0.0.45
			// go: golang.org/x/exp@v0.0.0-20250218142911-aa4b98e5adaa requires go >= 1.23.0 (running go 1.22.8; GOTOOLCHAIN=go1.22.8)
			// 正确的做法：省略 -u 选项以避免版本冲突
			// GOTOOLCHAIN=go1.22.8 go get github.com/yyle88/syntaxgo@v0.0.45
			// go: upgraded github.com/yyle88/syntaxgo v0.0.44 => v0.0.45
			output, err := config.GetSubCommand(projectPath).Exec("go", "get", module.Path+"@"+pkgTag)
			if err != nil {
				return erero.Wro(err)
			}
			zaplog.SUG.Debugln(string(output))

			zaplog.SUG.Debugln(projectPath, module.Path, module.Version, "sync", "=>", pkgTag, "done")
		}
	}

	zaplog.SUG.Debugln(neatjsons.S(pkgTagsMap))
	return nil
}

// GetPkgTagsMap retrieves latest Git tags for all modules in the workspace
// Creates a mapping from module paths to matching Git tag versions
//
// GetPkgTagsMap 获取工作区中所有模块的最新 Git 标签
// 创建从模块路径到其相应 Git 标签版本的映射
func GetPkgTagsMap(config *worksexec.WorksExec) map[string]string {
	pkgTagsMap := make(map[string]string)
	must.Done(config.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		moduleInfo := done.VCE(depbump.GetModuleInfo(projectPath)).Nice()

		tagName, _ := gitgo.NewGcm(projectPath, config.GetNewCommand()).LatestGitTag()
		if tagName != "" {
			zaplog.SUG.Debugln("pkg:", moduleInfo.Module.Path, "tag:", tagName)
		}

		pkgTagsMap[moduleInfo.Module.Path] = tagName
		return nil
	}))
	return pkgTagsMap
}
