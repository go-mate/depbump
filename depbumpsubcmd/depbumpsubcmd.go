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
			UpdateModules(config)
		},
	}
	cmd.AddCommand(NewUpdateDirectsCmd(config))
	return cmd
}

func NewUpdateDirectsCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "directs",
		Short: "go module update directs",
		Long:  "go module update directs",
		Run: func(cmd *cobra.Command, args []string) {
			UpdateDirects(config)
		},
	}
}

func UpdateDirects(config *workconfig.WorkspacesExecConfig) {
	for _, workspace := range config.GetWorkspaces() {
		for _, projectPath := range workspace.Projects {
			depbump.UpdateDirectRequires(config.GetSubCommand(projectPath), rese.P1(depbump.GetModuleInfo(projectPath)))
			must.Done(RunGoModuleTide(config.GetSubCommand(projectPath)))
		}
		if workspace.WorkRoot != "" {
			must.Done(RunGoWorkSync(config.GetSubCommand(workspace.WorkRoot)))
		}
	}
}

func RunGoModuleTide(execConfig *osexec.ExecConfig) error {
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

func UpdateModules(config *workconfig.WorkspacesExecConfig) {
	for _, workspace := range config.GetWorkspaces() {
		for _, projectPath := range workspace.Projects {
			must.Done(updateModules(config.GetSubCommand(projectPath), projectPath))
			must.Done(RunGoModuleTide(config.GetSubCommand(projectPath)))
		}
		if workspace.WorkRoot != "" {
			must.Done(RunGoWorkSync(config.GetSubCommand(workspace.WorkRoot)))
		}
	}
}

func updateModules(execConfig *osexec.ExecConfig, projectPath string) error {
	var matchedOnce bool
	output, err := execConfig.ShallowClone().
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
				matchedOnce = true
				return true
			}
			return false
		}).ExecInPipe("go", "get", "-u", "./...")
	if err != nil {
		if len(output) > 0 {
			zaplog.SUG.Warnln(string(output))
		}
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(string(output))

	if matchedOnce {
		depbump.UpdateDirectRequires(execConfig.ShallowClone().WithPath(projectPath), rese.P1(depbump.GetModuleInfo(projectPath)))
	}
	return nil
}
