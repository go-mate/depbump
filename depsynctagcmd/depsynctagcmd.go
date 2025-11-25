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
	"github.com/go-xlan/gitgo"
	"github.com/spf13/cobra"
	"github.com/yyle88/done"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/zaplog"
)

// SetupSyncCmd creates sync command and adds it to root command
// Provides tag-based synchronization subcommands
//
// SetupSyncCmd 创建同步命令并添加到根命令
// 提供基于标签的同步子命令
func SetupSyncCmd(rootCmd *cobra.Command, execConfig *osexec.ExecConfig) {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "dep sync",
		Long:  "dep sync",
	}
	cmd.AddCommand(SyncTagsCmd(execConfig))
	cmd.AddCommand(SyncSubsCmd(execConfig))

	rootCmd.AddCommand(cmd)
}

// SyncTagsCmd creates command that synchronizes dependencies to the latest Git tags
// Updates dependencies to match corresponding Git tag versions
//
// SyncTagsCmd 创建用于将依赖同步到最新 Git 标签的命令
// 更新依赖以匹配其相应的 Git 标签版本
func SyncTagsCmd(execConfig *osexec.ExecConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "sync tags",
		Long:  "sync tags",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(SyncTags(execConfig, depbump.GetModeUpdate))
		},
	}
	return cmd
}

// SyncSubsCmd creates command to sync dependencies with latest tag fallback
// Uses latest tag when dependencies have no specific tag
//
// SyncSubsCmd 创建用于同步依赖的命令，带有最新标签回退
// 当依赖没有特定标签时使用最新标签
func SyncSubsCmd(execConfig *osexec.ExecConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subs",
		Short: "sync subs",
		Long:  "sync subs",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(SyncTags(execConfig, depbump.GetModeLatest))
		},
	}
	return cmd
}

// SyncTags performs Git tag-based package synchronization
// Compares current package versions with Git tags and updates when different
//
// SyncTags 执行基于 Git 标签的依赖同步
// 比较当前依赖版本与 Git 标签，在不同时进行更新
func SyncTags(execConfig *osexec.ExecConfig, mode depbump.GetMode) error {
	pkgTagsMap := GetPkgTagsMap(execConfig)
	zaplog.SUG.Debugln(neatjsons.S(pkgTagsMap))

	projectDIR := osmustexist.ROOT(execConfig.Path)
	moduleInfo := done.VCE(depbump.GetModuleInfo(projectDIR)).Nice()
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
			if mode == depbump.GetModeLatest {
				pkgTag = "latest" // Use "latest" when no tag version exists // 当没有标签版本时使用 "latest"
			} else {
				continue
			}
		}

		if pkgTag == module.Version {
			zaplog.SUG.Debugln(projectDIR, module.Path, module.Version, "same")
			continue
		}
		zaplog.SUG.Debugln(projectDIR, module.Path, module.Version, "sync", "=>", pkgTag)

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
		output, err := execConfig.Exec("go", "get", module.Path+"@"+pkgTag)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(output))

		zaplog.SUG.Debugln(projectDIR, module.Path, module.Version, "sync", "=>", pkgTag, "done")
	}

	zaplog.SUG.Debugln(neatjsons.S(pkgTagsMap))
	return nil
}

// GetPkgTagsMap retrieves latest Git tags within the module
// Creates a mapping from module paths to matching Git tag versions
//
// GetPkgTagsMap 获取模块的最新 Git 标签
// 创建从模块路径到其相应 Git 标签版本的映射
func GetPkgTagsMap(execConfig *osexec.ExecConfig) map[string]string {
	pkgTagsMap := make(map[string]string)

	projectDIR := osmustexist.ROOT(execConfig.Path)
	moduleInfo := done.VCE(depbump.GetModuleInfo(projectDIR)).Nice()

	tagName, _ := gitgo.NewGcm(projectDIR, execConfig).LatestGitTag()
	if tagName != "" {
		zaplog.SUG.Debugln("pkg:", moduleInfo.Module.Path, "tag:", tagName)
	}

	pkgTagsMap[moduleInfo.Module.Path] = tagName
	return pkgTagsMap
}
