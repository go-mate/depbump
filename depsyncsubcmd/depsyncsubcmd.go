package depsyncsubcmd

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
		Use:   "direct",
		Short: "go module update direct",
		Long:  "go module update direct",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(SyncDeps(config))
		},
	}
	return cmd
}

func SyncDeps(config *workconfig.WorkspacesExecConfig) error {
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
			if pkgTag == module.Version {
				zaplog.SUG.Debugln(module.Path, module.Version, "same")
				continue
			}
			zaplog.SUG.Debugln(module.Path, module.Version, "sync")

			output, err := config.GetSubCommand(projectPath).Exec("go", "get", "-u", module.Path+"@"+pkgTag)
			if err != nil {
				return erero.Wro(err)
			}
			zaplog.SUG.Debugln(string(output))

			zaplog.SUG.Debugln(module.Path, module.Version, "done")
		}
	}

	zaplog.SUG.Debugln(neatjsons.S(pkgTagsMap))
	return nil
}

// GetPkgTagsMap 获得若干个模块的最新tag标签
func GetPkgTagsMap(config *workconfig.WorkspacesExecConfig) map[string]string {
	pkgTagsMap := make(map[string]string)
	must.Done(config.ForeachProjectExec(func(projectPath string, execConfig *osexec.ExecConfig) error {
		tagName, _ := gitgo.NewGcm(projectPath, config.GetNewCommand()).LatestGitTag()
		if tagName != "" {
			moduleInfo := done.VCE(depbump.GetModuleInfo(projectPath)).Nice()
			zaplog.SUG.Debugln("pkg:", moduleInfo.Module.Path, "tag:", tagName)
			pkgTagsMap[moduleInfo.Module.Path] = tagName
		}
		return nil
	}))
	return pkgTagsMap
}
