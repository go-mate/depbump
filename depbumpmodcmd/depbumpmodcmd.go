// Package depbumpmodcmd: Command-line interface to update Go modules
// Provides Cobra-based CLI commands to handle module updates with go get -u ./...
// Supports workspace operations with recursive module processing
//
// depbumpmodcmd: 更新 Go 模块的命令行接口
// 提供基于 Cobra 的 CLI 命令，使用 go get -u ./... 处理模块更新
// 支持递归处理工作区中的模块
package depbumpmodcmd

import (
	"fmt"

	"github.com/go-mate/depbump"
	"github.com/go-mate/go-work/workspath"
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
// Contains subcommand: recursive (R)
//
// NewModuleCmd 创建更新 Go 模块的命令
// 包含子命令：recursive (R)
func NewModuleCmd(execConfig *osexec.ExecConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module",
		Short: "Update module dependencies",
		Long:  "Update module dependencies using go get -u ./...",
		Run: func(cmd *cobra.Command, args []string) {
			UpdateModules(execConfig)
		},
	}

	// module R: module recursive // module R：模块递归更新
	{
		cmd.AddCommand(&cobra.Command{
			Use:     "recursive",
			Aliases: []string{"R"},
			Short:   "Update module dependencies across workspace",
			Long:    "Update module dependencies across all modules in the workspace.",
			Args:    cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				UpdateModulesRecursive(execConfig)
			},
		})
	}

	return cmd
}

// UpdateModules performs comprehensive module updates
//
// UpdateModules 执行全面的模块更新
func UpdateModules(execConfig *osexec.ExecConfig) {
	projectDIR := osmustexist.ROOT(execConfig.Path)
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
				zaplog.SUG.Debugln("match-upgrade-output-message:", eroticgo.GREEN.Sprint(neatjsons.S(upgradeInfo)))
				return true
			}
			if warnMessage, matched := depbump.MatchToolchainVersionMismatch(line); matched {
				zaplog.SUG.Debugln("go-toolchain-mismatch-output:", eroticgo.RED.Sprint(neatjsons.S(warnMessage)))
				success = false
				return true
			}
			if sdkInfo, matched := depbump.MatchGoDownloadingSdkInfo(line); matched {
				zaplog.SUG.Debugln("go-downloading-sdk-info:", eroticgo.CYAN.Sprint(neatjsons.S(sdkInfo)))
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

// UpdateModulesRecursive executes module updates across workspace modules
//
// UpdateModulesRecursive 在工作区模块中执行模块更新
func UpdateModulesRecursive(execConfig *osexec.ExecConfig) {
	workPath := osmustexist.ROOT(execConfig.Path)

	options := workspath.NewOptions().
		WithIncludeCurrentProject(true).
		WithIncludeCurrentPackage(false).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(false)

	moduleRoots := workspath.GetModulePaths(workPath, options)

	zaplog.SUG.Infoln("Recursive mode: found", eroticgo.CYAN.Sprint(len(moduleRoots)), "modules")

	for idx, modulePath := range moduleRoots {
		zaplog.SUG.Infoln("Module", eroticgo.GREEN.Sprint(fmt.Sprintf("(%d/%d)", idx+1, len(moduleRoots))), "Processing:", eroticgo.CYAN.Sprint(modulePath))

		moduleExecConfig := execConfig.NewConfig().WithPath(modulePath)
		UpdateModules(moduleExecConfig)
	}

	zaplog.SUG.Infoln("✅ Recursive module updates completed!")
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
