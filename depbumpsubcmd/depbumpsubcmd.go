// Package depbumpsubcmd: Command-line interface to update deps
// Provides update command with -D/-E/-L/-R flags
// Supports workspace operations with configurable filtering and update strategies
//
// depbumpsubcmd: 更新依赖的命令行接口
// 提供带有 -D/-E/-L/-R 标志的 update 命令
// 支持带有可配置过滤和更新策略的工作区操作
package depbumpsubcmd

import (
	"github.com/go-mate/depbump"
	"github.com/go-mate/depbump/internal/utils"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must/mustboolean"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/tern"
	"github.com/yyle88/zaplog"
)

// NewUpdateCmd creates update command with -D/-E/-L/-R flags
//
// NewUpdateCmd 创建 update 命令，使用 -D/-E/-L/-R 标志
func NewUpdateCmd(execConfig *osexec.ExecConfig) *cobra.Command {
	var (
		upDirectXX bool
		upEveryone bool
		upToLatest bool
		recurseXqt bool
	)

	config := &depbump.UpdateDepsConfig{
		Cate: depbump.DepCateDirect,
		Mode: depbump.GetModeUpdate,
	}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update dependencies",
		Long:  "Update dependencies with various strategies and filtering options.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure direct and everyone flags cannot be combined
			// 确保 direct 和 everyone 标志不能同时使用
			mustboolean.Conflict(upDirectXX, upEveryone)

			config.Cate = tern.BVV(upEveryone, depbump.DepCateEveryone, depbump.DepCateDirect)
			config.Mode = tern.BVV(upToLatest, depbump.GetModeLatest, depbump.GetModeUpdate)

			if recurseXqt {
				updateDepsRecursive(execConfig, config)
			} else {
				updateDeps(execConfig, config)
			}
		},
	}

	// Add flags to update command
	// 给 update 命令添加标志
	cmd.Flags().BoolVarP(&upDirectXX, "D", "D", false, "Update direct dependencies (default)")
	cmd.Flags().BoolVarP(&upEveryone, "E", "E", false, "Update each dependencies (direct + indirect)")
	cmd.Flags().BoolVarP(&upToLatest, "L", "L", false, "Use latest versions (including prerelease)")
	cmd.Flags().BoolVarP(&recurseXqt, "R", "R", false, "Process dependencies across workspace modules")
	cmd.Flags().BoolVarP(&config.GitlabOnly, "gitlab-only", "", false, "Update gitlab dependencies")
	cmd.Flags().BoolVarP(&config.SkipGitlab, "skip-gitlab", "", false, "Skip gitlab dependencies")
	cmd.Flags().BoolVarP(&config.GithubOnly, "github-only", "", false, "Update github dependencies")
	cmd.Flags().BoolVarP(&config.SkipGithub, "skip-github", "", false, "Skip github dependencies")

	return cmd
}

// updateDeps executes package updates with specified configuration
//
// updateDeps 使用指定配置执行包更新
func updateDeps(execConfig *osexec.ExecConfig, config *depbump.UpdateDepsConfig) {
	projectDIR := osmustexist.ROOT(execConfig.Path)
	zaplog.SUG.Infoln("Starting", string(config.Cate), "update:", eroticgo.CYAN.Sprint(projectDIR))
	zaplog.SUG.Debugln("Update config:", neatjsons.S(config))

	depbump.UpdateDeps(execConfig, rese.P1(depbump.GetModuleInfo(projectDIR)), config)
	rese.V1(execConfig.Exec("go", "mod", "tidy", "-e"))
}

// updateDepsRecursive executes package updates across workspace modules
//
// updateDepsRecursive 在工作区模块中执行包更新
func updateDepsRecursive(execConfig *osexec.ExecConfig, config *depbump.UpdateDepsConfig) {
	utils.ForeachModule(execConfig, func(moduleExecConfig *osexec.ExecConfig) {
		updateDeps(moduleExecConfig, config)
	})
}
