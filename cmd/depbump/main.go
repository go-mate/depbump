package main

import (
	"os"

	"github.com/go-mate/depbump/depbumpsubcmd"
	"github.com/go-mate/go-work/workconfig"
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

	workspace := workconfig.NewWorkspace("", []string{projectPath})
	workspace.MustCheck()

	workspaces := workconfig.NewWorkspaces(workspace)
	workspaces.MustCheck()

	commandConfig := osexec.NewCommandConfig()
	commandConfig.WithBash()
	commandConfig.WithDebugMode(true)

	config := workconfig.NewWorkspacesExecConfig(workspaces, commandConfig)

	cmd := depbumpsubcmd.NewUpdateCmd(config)
	must.Done(cmd.Execute())
}
