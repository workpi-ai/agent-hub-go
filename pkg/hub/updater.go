package hub

import (
	registry "github.com/workpi-ai/agent-hub"
	"github.com/workpi-ai/go-utils/ghrelease"
)

func NewUpdater(metadataFile, localStandardAgentsDir, localStandardCommandsDir string) (*ghrelease.Updater, error) {
	return ghrelease.NewUpdater(ghrelease.UpdaterConfig{
		RepoOwner:    repoOwner,
		RepoName:     repoName,
		MetadataFile: metadataFile,
		Targets: []ghrelease.ExtractTarget{
			{
				PathTransformer: &ghrelease.SubDirTransformer{SubDir: registry.AgentsDir, Ext: markdownExt},
				DestDir:         localStandardAgentsDir,
			},
			{
				PathTransformer: &ghrelease.SubDirTransformer{SubDir: registry.CommandsDir, Ext: markdownExt},
				DestDir:         localStandardCommandsDir,
			},
		},
	})
}
