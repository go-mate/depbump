package main

import (
	"os"

	"github.com/go-mate/depbump/depbumpsubcmd"
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// go run main.go
// go run main.go directs
func main() {
	projectPath := rese.C1(os.Getwd())
	zaplog.LOG.Debug("project:", zap.String("path", projectPath))

	executePath := rese.C1(os.Executable())
	zaplog.LOG.Debug("execute:", zap.String("path", executePath))

	wsp := workspace.NewWorkspace("", []string{projectPath})

	execConfig := osexec.NewCommandConfig()
	execConfig.WithBash()
	execConfig.WithDebugMode(true)

	workspaces := []*workspace.Workspace{wsp}

	config := worksexec.NewWorksExec(execConfig, workspaces)

	cmd := depbumpsubcmd.NewUpdateCmd(config)
	must.Done(cmd.Execute())
}
