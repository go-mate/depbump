package depbump

import (
	"github.com/yyle88/erero"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

func UpdateModule(projectPath string, modulePath string) error {
	output, err := osexec.NewOsCommand().WithPath(projectPath).StreamExec("go", "get", "-u", modulePath)
	if err != nil {
		if len(output) > 0 {
			zaplog.SUG.Warnln(string(output))
		}
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(string(output))
	return nil
}
