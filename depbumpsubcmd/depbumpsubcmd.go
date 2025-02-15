package depbumpsubcmd

import (
	"github.com/go-mate/depbump"
	"github.com/go-mate/go-work/workconfig"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func NewUpdateCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "go module update -->>",
		Long:  "go module update -->>",
		Run: func(cmd *cobra.Command, args []string) {
			UpdateDeps(config)
		},
	}
	cmd.AddCommand(NewUpdateDirectCmd(config))
	cmd.AddCommand(NewUpdateModuleCmd(config))
	return cmd
}

func UpdateDeps(config *workconfig.WorkspacesExecConfig) {
	for _, workspace := range config.GetWorkspaces() {
		for _, projectPath := range workspace.Projects {
			moduleInfo := rese.P1(depbump.GetModuleInfo(projectPath))
			success, err := updateDeps(config.GetSubCommand(projectPath), projectPath, moduleInfo.GetToolchainVersion())
			if err != nil {
				panic(erero.Wro(err))
			}
			if success {
				zaplog.SUG.Infoln("success", eroticgo.RED.Sprint("success"))
			} else {
				zaplog.SUG.Warnln("warning", eroticgo.RED.Sprint("warning"))
			}
			must.Done(RunGoModTide(config.GetSubCommand(projectPath)))
		}
		if workspace.WorkRoot != "" {
			must.Done(RunGoWorkSync(config.GetSubCommand(workspace.WorkRoot)))
		}
	}
}

func updateDeps(execConfig *osexec.ExecConfig, projectPath string, toolchain string) (bool, error) {
	var success = true
	output, err := execConfig.ShallowClone().
		WithEnvs([]string{"GOTOOLCHAIN=" + toolchain}). //在升级时需要用项目的go版本号压制住依赖的go版本号
		WithPath(projectPath).
		WithMatchMore(true).
		WithMatchPipe(func(line string) bool {
			upgradeInfo, matched := depbump.MatchUpgrade(line)
			if matched {
				zaplog.SUG.Debugln("match-upgrade-output-message:", eroticgo.GREEN.Sprint(neatjsons.S(upgradeInfo)))
				return true
			}
			toolchainVersionMismatch, matched := depbump.MatchToolchainVersionMismatch(line)
			if matched {
				zaplog.SUG.Debugln("go-toolchain-mismatch-output:", eroticgo.RED.Sprint(neatjsons.S(toolchainVersionMismatch)))
				success = false
				return true
			}
			return false
		}).ExecInPipe("go", "get", "-u", "./...")
	if err != nil {
		if len(output) > 0 {
			zaplog.SUG.Warnln(string(output))
		}
		return false, erero.Wro(err)
	}
	if success {
		zaplog.SUG.Debugln(string(output))
	} else {
		zaplog.SUG.Warnln(string(output))
	}
	return success, nil
}

func NewUpdateDirectCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "direct",
		Short: "go module update direct",
		Long:  "go module update direct",
		Run: func(cmd *cobra.Command, args []string) {
			UpdateDirectDeps(config)
		},
	}
}

func UpdateDirectDeps(config *workconfig.WorkspacesExecConfig) {
	for _, workspace := range config.GetWorkspaces() {
		for _, projectPath := range workspace.Projects {
			depbump.UpdateDirectRequires(config.GetSubCommand(projectPath), rese.P1(depbump.GetModuleInfo(projectPath)))
			must.Done(RunGoModTide(config.GetSubCommand(projectPath)))
		}
		if workspace.WorkRoot != "" {
			must.Done(RunGoWorkSync(config.GetSubCommand(workspace.WorkRoot)))
		}
	}
}

func NewUpdateModuleCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "module",
		Short: "go module update module",
		Long:  "go module update module",
		Run: func(cmd *cobra.Command, args []string) {
			UpdateDepModules(config)
		},
	}
}

func UpdateDepModules(config *workconfig.WorkspacesExecConfig) {
	for _, workspace := range config.GetWorkspaces() {
		for _, projectPath := range workspace.Projects {
			moduleInfo := rese.P1(depbump.GetModuleInfo(projectPath))
			must.Done(updateDepModules(config.GetSubCommand(projectPath), projectPath, moduleInfo.GetToolchainVersion()))
			must.Done(RunGoModTide(config.GetSubCommand(projectPath)))
		}
		if workspace.WorkRoot != "" {
			must.Done(RunGoWorkSync(config.GetSubCommand(workspace.WorkRoot)))
		}
	}
}

func updateDepModules(execConfig *osexec.ExecConfig, projectPath string, toolchain string) error {
	success, err := updateDeps(execConfig, projectPath, toolchain)
	if err != nil {
		return erero.Wro(err)
	}
	if !success {
		depbump.UpdateDirectRequires(execConfig.ShallowClone().WithPath(projectPath), rese.P1(depbump.GetModuleInfo(projectPath)))
	}
	return nil
}

func RunGoModTide(execConfig *osexec.ExecConfig) error {
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

func RunGoWorkSync(execConfig *osexec.ExecConfig) error {
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
