package hub

import (
	"strings"

	"github.com/workpi-ai/go-utils/ghrelease"
)

type agentFilter struct{}

func (f *agentFilter) ShouldExtract(filename string) bool {
	return strings.Contains(filename, "/"+agentsDir+"/") && strings.HasSuffix(filename, markdownExt)
}

func NewUpdater(configDir string) *ghrelease.Updater {
	return ghrelease.NewUpdater(ghrelease.UpdaterConfig{
		RepoOwner:        repoOwner,
		RepoName:         repoName,
		DestDir:          configDir,
		MetadataFilename: metadataFile,
		ExtractFilter:    &agentFilter{},
	})
}
