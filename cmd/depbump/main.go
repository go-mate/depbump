package main

import (
	"os"

	"github.com/go-mate/depbump/depbumpsubcmd"
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/go-mate/go-work/workspath"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// go run main.go
// go run main.go directs
// go run main.go directs --gitlab-only
// go run main.go directs --skip-gitlab
// go run main.go directs --github-only
// go run main.go directs --skip-github
func main() {
	currentPath := rese.C1(os.Getwd())
	zaplog.LOG.Debug("current:", zap.String("path", currentPath))

	executePath := rese.C1(os.Executable())
	zaplog.LOG.Debug("execute:", zap.String("path", executePath))

	projectPath, _, ok := workspath.GetProjectPath(currentPath)
	must.True(ok)
	zaplog.LOG.Debug("project:", zap.String("path", projectPath))
	must.Nice(projectPath)

	execConfig := osexec.NewCommandConfig().WithBash().WithDebug()
	workspaces := []*workspace.Workspace{
		workspace.NewWorkspace("", []string{projectPath}),
	}
	config := worksexec.NewWorksExec(execConfig, workspaces)

	cmd := depbumpsubcmd.NewUpdateCmd(config)
	must.Done(cmd.Execute())
}
