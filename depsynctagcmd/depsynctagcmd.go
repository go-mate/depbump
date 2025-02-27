package depsynctagcmd

import (
	"github.com/go-mate/depbump"
	"github.com/go-mate/go-work/workconfig"
	"github.com/go-xlan/gitgo"
	"github.com/spf13/cobra"
	"github.com/yyle88/done"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

func SyncDepsCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "go workspace sync",
		Long:  "go workspace sync",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(config.ForeachWorkRootRun(func(workspace *workconfig.Workspace, execConfig *osexec.ExecConfig) error {
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
	return cmd
}

func SyncTagsCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
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

func SyncSubsCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
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

func SyncTags(config *workconfig.WorkspacesExecConfig, useLatest bool) error {
	pkgTagsMap := GetPkgTagsMap(config)
	zaplog.SUG.Debugln(neatjsons.S(pkgTagsMap))

	for _, projectPath := range config.CollectSubprojectPaths() {
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
					pkgTag = "latest" // 假如没有 tag 版本号，则默认为 latest 的
				} else {
					continue
				}
			}

			if pkgTag == module.Version {
				zaplog.SUG.Debugln(projectPath, module.Path, module.Version, "same")
				continue
			}
			zaplog.SUG.Debugln(projectPath, module.Path, module.Version, "sync", "=>", pkgTag)

			// GOTOOLCHAIN=go1.22.8 go get -u github.com/yyle88/syntaxgo@v0.0.45
			// go: golang.org/x/exp@v0.0.0-20250218142911-aa4b98e5adaa requires go >= 1.23.0 (running go 1.22.8; GOTOOLCHAIN=go1.22.8)
			// 因此正确的做法是不带 -u 选项
			// GOTOOLCHAIN=go1.22.8 go get github.com/yyle88/syntaxgo@v0.0.45
			// go: upgraded github.com/yyle88/syntaxgo v0.0.44 => v0.0.45
			// 这里需要注意
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

// GetPkgTagsMap 获得若干个模块的最新tag标签
func GetPkgTagsMap(config *workconfig.WorkspacesExecConfig) map[string]string {
	pkgTagsMap := make(map[string]string)
	must.Done(config.ForeachProjectExec(func(projectPath string, execConfig *osexec.ExecConfig) error {
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
