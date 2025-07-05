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

func NewUpdateCmd(config *worksexec.WorksExec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "depbump -->>",
		Long:  "depbump -->>",
		Run: func(cmd *cobra.Command, args []string) {
			updateModules(config)
		},
	}
	cmd.AddCommand(NewUpdateModuleCmd(config, "module"))
	cmd.AddCommand(NewUpdateDirectCmd(config, "direct"))
	cmd.AddCommand(NewUpdateEveryoneCmd(config, "everyone")) //这不用"all"避免和"all"混淆
	return cmd
}

func NewUpdateModuleCmd(config *worksexec.WorksExec, usageName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     usageName,
		Aliases: aliasesMap[usageName],
		Short:   "depbump module",
		Long:    "depbump module",
		Run: func(cmd *cobra.Command, args []string) {
			updateModules(config)
		},
	}
	return cmd
}

func updateModules(config *worksexec.WorksExec) {
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

func updateModule(execConfig *osexec.ExecConfig, projectPath string, toolchain string) {
	var success = true
	output := rese.V1(execConfig.NewConfig().
		WithEnvs([]string{"GOTOOLCHAIN=" + toolchain}). //在升级时需要用项目的go版本号压制住依赖的go版本号
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

func NewUpdateDirectCmd(config *worksexec.WorksExec, usageName string) *cobra.Command {
	const usageNameLatest = "latest"
	mode := tern.BVV(usageName == usageNameLatest, depbump.GetModeLatest, depbump.GetModeUpdate)

	cmd := &cobra.Command{
		Use:     usageName,
		Aliases: aliasesMap[usageName],
		Short:   "depbump direct (latest)",
		Long:    "depbump direct (latest)",
		Run: func(cmd *cobra.Command, args []string) {
			updateDepsConfig := &depbump.UpdateDepsConfig{
				Cate: depbump.DepCateDirect,
				Mode: mode,
			}
			updateDeps(config, updateDepsConfig)
		},
	}
	if usageName != usageNameLatest {
		cmd.AddCommand(NewUpdateDirectCmd(config, usageNameLatest))
	}
	return cmd
}

func NewUpdateEveryoneCmd(config *worksexec.WorksExec, usageName string) *cobra.Command {
	const usageNameLatest = "latest"
	mode := tern.BVV(usageName == usageNameLatest, depbump.GetModeLatest, depbump.GetModeUpdate)

	cmd := &cobra.Command{
		Use:     usageName,
		Aliases: aliasesMap[usageName],
		Short:   "depbump require (latest)",
		Long:    "depbump require (latest)",
		Run: func(cmd *cobra.Command, args []string) {
			updateDepsConfig := &depbump.UpdateDepsConfig{
				Cate: depbump.DepCateEveryone,
				Mode: mode,
			}
			updateDeps(config, updateDepsConfig)
		},
	}
	if usageName != usageNameLatest {
		cmd.AddCommand(NewUpdateEveryoneCmd(config, usageNameLatest))
	}
	return cmd
}

func updateDeps(config *worksexec.WorksExec, updateDepsConfig *depbump.UpdateDepsConfig) {
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
