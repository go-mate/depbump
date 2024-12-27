package depbump

import (
	"regexp"

	"github.com/yyle88/erero"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

func UpdateModule(projectPath string, modulePath string) error {
	output, err := osexec.NewOsCommand().WithPath(projectPath).
		WithMatchMore(true).
		WithMatchPipe(func(line string) bool {
			upgradeInfo, matched := MatchUpgrade(line)
			if matched {
				zaplog.SUG.Debugln("match-output-message:", eroticgo.GREEN.Sprint(neatjsons.S(upgradeInfo)))
			}
			return matched
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
