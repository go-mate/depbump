// Package depbumpmodcmd: Command-line interface to update Go modules
// Provides module command with -R flag using go get -u ./...
// Supports workspace operations with recursive module processing
//
// depbumpmodcmd: 更新 Go 模块的命令行接口
// 提供带有 -R 标志的 module 命令，使用 go get -u ./... 处理模块更新
// 支持递归处理工作区中的模块
package depbumpmodcmd

import (
	"github.com/go-mate/depbump"
	"github.com/go-mate/depbump/internal/utils"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

// NewModuleCmd creates command to update Go modules
// Uses -R flag to enable recursive mode
//
// NewModuleCmd 创建更新 Go 模块的命令
// 使用 -R 标志启用递归模式
func NewModuleCmd(execConfig *osexec.ExecConfig) *cobra.Command {
	var recurseXqt bool

	cmd := &cobra.Command{
		Use:   "module",
		Short: "Update module dependencies",
		Long:  "Update module dependencies using go get -u ./...",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if recurseXqt {
				UpdateModulesRecursive(execConfig)
			} else {
				UpdateModules(execConfig)
			}
		},
	}

	// Add flags to module command
	// 给 module 命令添加标志
	cmd.Flags().BoolVarP(&recurseXqt, "R", "R", false, "Process modules across workspace")

	return cmd
}

// UpdateModules performs comprehensive module updates
//
// UpdateModules 执行全面的模块更新
func UpdateModules(execConfig *osexec.ExecConfig) {
	projectDIR := osmustexist.ROOT(execConfig.Path)
	zaplog.SUG.Infoln("Starting module update:", eroticgo.CYAN.Sprint(projectDIR))
	moduleInfo := rese.P1(depbump.GetModuleInfo(projectDIR))
	updateModule(execConfig, moduleInfo.GetToolchainVersion())
	must.Done(GoModTide(execConfig))
}

// updateModule executes go get -u on a single module with toolchain management
//
// updateModule 在单个模块上执行 go get -u，带工具链管理
func updateModule(execConfig *osexec.ExecConfig, toolchain string) {
	var success = true
	output := rese.V1(execConfig.NewConfig().
		WithEnvs([]string{"GOTOOLCHAIN=" + toolchain}).
		WithMatchMore(true).
		WithMatchPipe(func(line string) bool {
			if upgradeInfo, matched := depbump.MatchUpgrade(line); matched {
				zaplog.SUG.Debugln("Upgrade detected:", eroticgo.GREEN.Sprint(neatjsons.S(upgradeInfo)))
				return true
			}
			if warnMessage, matched := depbump.MatchToolchainVersionMismatch(line); matched {
				zaplog.SUG.Debugln("Toolchain mismatch:", eroticgo.RED.Sprint(neatjsons.S(warnMessage)))
				success = false
				return true
			}
			if sdkInfo, matched := depbump.MatchGoDownloadingSdkInfo(line); matched {
				zaplog.SUG.Debugln("Downloading SDK:", eroticgo.CYAN.Sprint(neatjsons.S(sdkInfo)))
				return true
			}
			return false
		}).ExecInPipe("go", "get", "-u", "./..."))
	if success {
		zaplog.SUG.Debugln(string(output))
		zaplog.SUG.Infoln("Module update", eroticgo.GREEN.Sprint("success"))
	} else {
		zaplog.SUG.Warnln(string(output))
		zaplog.SUG.Warnln("Module update", eroticgo.RED.Sprint("has warnings"))
	}
}

// UpdateModulesRecursive executes module updates across workspace modules
//
// UpdateModulesRecursive 在工作区模块中执行模块更新
func UpdateModulesRecursive(execConfig *osexec.ExecConfig) {
	utils.ForeachModule(execConfig, func(moduleExecConfig *osexec.ExecConfig) {
		UpdateModules(moduleExecConfig)
	})
}

// GoModTide executes go mod tidy
//
// GoModTide 执行 go mod tidy
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
