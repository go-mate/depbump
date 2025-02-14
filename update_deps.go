package depbump

import (
	"fmt"
	"regexp"

	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func UpdateModule(execConfig *osexec.ExecConfig, modulePath string, toolchain string) error {
	output, err := execConfig.ShallowClone().
		WithEnvs([]string{"GOTOOLCHAIN=" + toolchain}). //在升级时需要用项目的go版本号压制住依赖的go版本号
		WithMatchMore(true).
		WithMatchPipe(func(line string) bool {
			upgradeInfo, matched := MatchUpgrade(line)
			if matched {
				zaplog.SUG.Debugln("match-upgrade-output-message:", eroticgo.GREEN.Sprint(neatjsons.S(upgradeInfo)))
				return true
			}
			toolchainVersionMismatch, matched := MatchToolchainVersionMismatch(line)
			if matched {
				zaplog.SUG.Debugln("go-toolchain-mismatch-output:", eroticgo.RED.Sprint(neatjsons.S(toolchainVersionMismatch)))
				return true
			}
			return false
		}).ExecInPipe("go", "get", "-u", modulePath)
	if err != nil {
		if len(output) > 0 {
			zaplog.SUG.Warnln(string(output))
		}
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(string(output))
	return nil
}

type UpgradeInfo struct {
	Module     string `json:"module"`
	OldVersion string `json:"old_version"`
	NewVersion string `json:"new_version"`
}

func MatchUpgrade(outputLine string) (*UpgradeInfo, bool) {
	pattern := `go: upgraded ([^\s]+) ([^\s]+) => ([^\s]+)`
	re := regexp.MustCompile(pattern)

	// Match the input string
	matches := re.FindStringSubmatch(outputLine)
	if len(matches) != 4 {
		return nil, false
	}

	// Extract module, old version, and new version
	return &UpgradeInfo{
		Module:     matches[1],
		OldVersion: matches[2],
		NewVersion: matches[3],
	}, true
}

// ToolchainVersionMismatch 表示 Go 工具链版本不匹配的信息
type ToolchainVersionMismatch struct {
	ModulePath        string // 模块路径
	ModuleVersion     string // 模块版本
	RequiredGoVersion string // 所需的最低 Go 版本
	RunningGoVersion  string // 当前运行的 Go 版本
	Toolchain         string // GOTOOLCHAIN 的值
}

// MatchToolchainVersionMismatch 解析工具链版本不匹配的错误信息
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

func UpdateDirectRequires(execConfig *osexec.CommandConfig, moduleInfo *ModuleInfo) {
	UpdateRequires(execConfig, moduleInfo.GetDirectModules(), moduleInfo.GetToolchainVersion())
}

func UpdateRequires(execConfig *osexec.CommandConfig, requires []*Require, toolchain string) {
	type Warning struct {
		Path string `json:"path"`
		Warn string `json:"warn"`
	}

	var warnings []*Warning
	for idx, dep := range requires {
		processMessage := fmt.Sprintf("(%d/%d)", idx, len(requires))
		zaplog.LOG.Debug("upgrade:", zap.String("idx", processMessage), zap.String("path", dep.Path), zap.String("from", dep.Version))
		if err := UpdateModule(execConfig, dep.Path, toolchain); err != nil {
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
