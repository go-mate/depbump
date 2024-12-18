package main

import (
	"os"

	"github.com/go-mate/depbump"
	"github.com/yyle88/done"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func main() {
	projectPath := rese.C1(os.Getwd())
	zaplog.LOG.Debug("project:", zap.String("path", projectPath))

	modInfo := rese.P1(depbump.GetModInfo(projectPath))
	zaplog.LOG.Debug("require:", zap.Int("size", len(modInfo.Require)))

	for idx, dep := range modInfo.Require {
		if !dep.Indirect {
			zaplog.LOG.Debug("upgrade:", zap.Int("idx", idx), zap.String("path", dep.Path), zap.String("from", dep.Version))
			done.Soft(depbump.UpdateModule(projectPath, dep.Path))
		}
	}
}
