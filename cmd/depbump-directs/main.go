package main

import (
	"os"

	"github.com/go-mate/depbump"
	"github.com/go-mate/depbump/depbumpsubcmd"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func main() {
	projectPath := rese.C1(os.Getwd())
	zaplog.LOG.Debug("project:", zap.String("path", projectPath))

	executePath := rese.C1(os.Executable())
	zaplog.LOG.Debug("execute:", zap.String("path", executePath))

	moduleInfo := rese.P1(depbump.GetModuleInfo(projectPath))
	zaplog.LOG.Debug("require:", zap.Int("size", len(moduleInfo.Require)))

	execConfig := osexec.NewExecConfig().WithDebug().WithPath(projectPath)
	depbump.UpdateDirectRequires(execConfig, moduleInfo)

	must.Done(depbumpsubcmd.RunGoModTide(execConfig))
}
