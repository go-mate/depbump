// Package depbump: Advanced dependency update engine with version control
// Implements intelligent dependency upgrading with toolchain management
// Supports pattern matching for upgrade output and error diagnostics
//
// depbump: 高级依赖更新引擎，带版本控制
// 实现智能依赖升级，包含工具链管理
// 支持升级输出的模式匹配和错误诊断
package depbump

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/must/muststrings"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/tern"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// GetMode defines the strategy for dependency version acquisition
// Latest mode gets the newest available version, Update mode gets compatible upgrades
//
// GetMode 定义依赖版本获取策略
// Latest 模式获取最新可用版本，Update 模式获取兼容升级
type GetMode string

const (
	GetModeLatest GetMode = "LATEST" // Get latest available version // 获取最新可用版本
	GetModeUpdate GetMode = "UPDATE" // Get compatible updates // 获取兼容更新
)

// UpdateConfig specifies parameters for individual module updates
// Controls toolchain version and update strategy for dependency upgrades
//
// UpdateConfig 指定单个模块更新的参数
// 控制工具链版本和依赖升级的更新策略
type UpdateConfig struct {
	Toolchain string  // Go toolchain version to use // 使用的 Go 工具链版本
	Mode      GetMode // Update strategy mode // 更新策略模式
}

// UpdateModule performs dependency update for a specific module path
// Uses specified toolchain and mode to execute go get commands with output monitoring
//
// UpdateModule 为特定模块路径执行依赖更新
// 使用指定的工具链和模式执行 go get 命令，并监控输出
func UpdateModule(execConfig *osexec.ExecConfig, modulePath string, updateConfig *UpdateConfig) error {
	must.Nice(execConfig)
	must.Nice(modulePath)
	must.Nice(updateConfig)
	must.Nice(updateConfig.Toolchain)

	commands := tern.BFF(updateConfig.Mode == GetModeLatest, func() []string {
		modulePathLatest := tern.BVF(strings.HasSuffix(modulePath, "@latest"), modulePath, func() string {
			muststrings.NotContains(modulePath, "@")
			return modulePath + "@latest"
		})

		return []string{"go", "get", modulePathLatest}
	}, func() []string {
		return []string{"go", "get", "-u", modulePath}
	})
	zaplog.LOG.Debug("update-module:", zap.String("module-path", modulePath), zap.Strings("commands", commands))

	output, err := execConfig.NewConfig().
		WithEnvs([]string{"GOTOOLCHAIN=" + updateConfig.Toolchain}). // Use project Go version to suppress dependency Go version requirements // 在升级时需要用项目的go版本号压制住依赖的go版本号
		WithMatchMore(true).
		WithMatchPipe(func(line string) bool {
			if upgradeInfo, matched := MatchUpgrade(line); matched {
				zaplog.SUG.Debugln("match-upgrade-output-message:", eroticgo.GREEN.Sprint(neatjsons.S(upgradeInfo)))
				return true
			}
			if waToolchain, matched := MatchToolchainVersionMismatch(line); matched {
				zaplog.SUG.Debugln("go-toolchain-mismatch-output:", eroticgo.RED.Sprint(neatjsons.S(waToolchain)))
				return true
			}
			return false
		}).ExecInPipe(commands[0], commands[1:]...)
	if err != nil {
		if len(output) > 0 {
			zaplog.SUG.Warnln(string(output))
		}
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(string(output))
	return nil
}

// UpgradeInfo captures successful dependency upgrade information
// Parsed from go get command output to track version changes
//
// UpgradeInfo 捕获成功的依赖升级信息
// 从 go get 命令输出解析，跟踪版本变化
type UpgradeInfo struct {
	Module     string `json:"module"`      // Module path that was upgraded // 已升级的模块路径
	OldVersion string `json:"old_version"` // Previous version // 之前的版本
	NewVersion string `json:"new_version"` // Updated version after upgrade // 升级后的更新版本
}

// MatchUpgrade parses go get output to extract upgrade information
// Returns upgrade details when the output matches the expected pattern
//
// MatchUpgrade 解析 go get 输出以提取升级信息
// 当输出匹配预期模式时返回升级详情
func MatchUpgrade(outputLine string) (*UpgradeInfo, bool) {
	pattern := `go: upgraded ([^\s]+) ([^\s]+) => ([^\s]+)`
	re := regexp.MustCompile(pattern)

	// Match the input string // 匹配输入字符串
	matches := re.FindStringSubmatch(outputLine)
	if len(matches) != 4 {
		return nil, false
	}

	// Extract module, old version, and new version // 提取模块、旧版本和新版本
	return &UpgradeInfo{
		Module:     matches[1],
		OldVersion: matches[2],
		NewVersion: matches[3],
	}, true
}

// ToolchainVersionMismatch represents Go toolchain version compatibility issues
// Contains detailed information about version conflicts during dependency updates
//
// ToolchainVersionMismatch 代表 Go 工具链版本兼容性问题
// 包含依赖更新期间版本冲突的详细信息
type ToolchainVersionMismatch struct {
	ModulePath        string // Module path with version conflict // 存在版本冲突的模块路径
	ModuleVersion     string // Specific module version // 特定模块版本
	RequiredGoVersion string // Minimum required Go version // 所需最低 Go 版本
	RunningGoVersion  string // Active Go version in use // 当前使用的 Go 版本
	Toolchain         string // GOTOOLCHAIN environment value // GOTOOLCHAIN 环境变量值
}

// MatchToolchainVersionMismatch parses toolchain version conflict error messages
// Extracts structured information from go command error output about version mismatches
//
// MatchToolchainVersionMismatch 解析工具链版本冲突错误消息
// 从 go 命令错误输出中提取版本不匹配的结构化信息
func MatchToolchainVersionMismatch(outputLine string) (*ToolchainVersionMismatch, bool) {
	pattern := `^go: ([^\s]+)@([^\s]+) requires go >= ([^\s]+) \(running go ([^\s]+); GOTOOLCHAIN=([^\s]+)\)$`
	re := regexp.MustCompile(pattern)

	// 匹配输入字符串
	matches := re.FindStringSubmatch(outputLine)
	if len(matches) != 6 {
		return nil, false
	}

	// 提取信息并返回
	return &ToolchainVersionMismatch{
		ModulePath:        matches[1],
		ModuleVersion:     matches[2],
		RequiredGoVersion: matches[3],
		RunningGoVersion:  matches[4],
		Toolchain:         matches[5],
	}, true
}

// UpdateDepsConfig provides comprehensive configuration for batch dependency updates
// Supports selective updating based on dependency categories and source filtering
//
// UpdateDepsConfig 为批量依赖更新提供全面的配置
// 支持基于依赖类别和源过滤的选择性更新
type UpdateDepsConfig struct {
	Cate       DepCate // Dependency category filter // 依赖类别过滤器
	Mode       GetMode // Update mode strategy // 更新模式策略
	GitlabOnly bool    // Update just GitLab dependencies // 仅更新 GitLab 依赖
	SkipGitlab bool    // Skip GitLab dependencies // 跳过 GitLab 依赖
	GithubOnly bool    // Update just GitHub dependencies // 仅更新 GitHub 依赖
	SkipGithub bool    // Skip GitHub dependencies // 跳过 GitHub 依赖
}

// UpdateDeps orchestrates batch dependency updates according to configuration
// Processes filtered dependencies with progress tracking and error collection
//
// UpdateDeps 根据配置编排批量依赖更新
// 处理过滤后的依赖，带有进度跟踪和错误收集
func UpdateDeps(execConfig *osexec.CommandConfig, moduleInfo *ModuleInfo, updateDepsConfig *UpdateDepsConfig) {
	must.Nice(execConfig)
	must.Nice(updateDepsConfig)

	toolchainVersion := moduleInfo.GetToolchainVersion()
	must.Nice(toolchainVersion)

	type Warning struct {
		Path string `json:"path"`
		Warn string `json:"warn"`
	}

	var warnings []*Warning
	requires := moduleInfo.GetScopedRequires(updateDepsConfig.Cate)
	for idx, dep := range requires {
		processMessage := fmt.Sprintf("(%d/%d)", idx, len(requires))
		zaplog.LOG.Debug("upgrade:", zap.String("idx", processMessage), zap.String("path", dep.Path), zap.String("from", dep.Version))

		if updateDepsConfig.GitlabOnly && !strings.HasPrefix(dep.Path, "gitlab.") {
			zaplog.LOG.Debug("skip-non-gitlab:", zap.String("path", dep.Path), zap.String("from", dep.Version))
			continue
		}

		if updateDepsConfig.SkipGitlab && strings.HasPrefix(dep.Path, "gitlab.") {
			zaplog.LOG.Debug("skip-gitlab-dep:", zap.String("path", dep.Path), zap.String("from", dep.Version))
			continue
		}

		if updateDepsConfig.GithubOnly && !strings.HasPrefix(dep.Path, "github.com/") {
			zaplog.LOG.Debug("skip-non-github:", zap.String("path", dep.Path), zap.String("from", dep.Version))
			continue
		}

		if updateDepsConfig.SkipGithub && strings.HasPrefix(dep.Path, "github.com/") {
			zaplog.LOG.Debug("skip-github-dep:", zap.String("path", dep.Path), zap.String("from", dep.Version))
			continue
		}

		if err := UpdateModule(execConfig, dep.Path, &UpdateConfig{
			Toolchain: toolchainVersion,
			Mode:      updateDepsConfig.Mode,
		}); err != nil {
			warnings = append(warnings, &Warning{
				Path: dep.Path,
				Warn: err.Error(),
			})
		}
	}

	if len(warnings) > 0 {
		eroticgo.RED.ShowMessage("WARNING>>>")
		for idx, warning := range warnings {
			zaplog.LOG.Debug("warning:", zap.Int("idx", idx), zap.String("path", warning.Path))
			fmt.Println(eroticgo.RED.Sprint(warning.Warn))
		}
		eroticgo.RED.ShowMessage("<<<WARNING")
	} else {
		eroticgo.GREEN.ShowMessage("SUCCESS")
	}
}
