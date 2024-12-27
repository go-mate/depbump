package main

import (
	"fmt"
	"os"

	"github.com/go-mate/depbump"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func main() {
	projectPath := rese.C1(os.Getwd())
	zaplog.LOG.Debug("project:", zap.String("path", projectPath))

	executePath := rese.C1(os.Executable())
	zaplog.LOG.Debug("execute:", zap.String("path", executePath))

	modInfo := rese.P1(depbump.GetModInfo(projectPath))
	zaplog.LOG.Debug("require:", zap.Int("size", len(modInfo.Require)))

	type Warning struct {
		Path string `json:"path"`
		Warn string `json:"warn"`
	}

	var warnings []*Warning
	for idx, dep := range modInfo.Require {
		if !dep.Indirect {
			zaplog.LOG.Debug("upgrade:", zap.Int("idx", idx), zap.String("path", dep.Path), zap.String("from", dep.Version))
			if err := depbump.UpdateModule(projectPath, dep.Path); err != nil {
				warnings = append(warnings, &Warning{
					Path: dep.Path,
					Warn: err.Error(),
				})
			}
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
