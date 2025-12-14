package hub

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	registry "github.com/workpi-ai/agent-hub"
)

type Loader struct {
	configDir string
}

func NewLoader(configDir string) *Loader {
	return &Loader{configDir: configDir}
}

func (l *Loader) Load() (map[string]*Agent, error) {
	var fsys fs.FS
	var subPath string

	localPath := filepath.Join(l.configDir, agentsDir)
	if stat, err := os.Stat(localPath); err == nil && stat.IsDir() {
		fsys = os.DirFS(localPath)
		subPath = "."
	} else {
		fsys = registry.Agents
		subPath = agentsDir
	}

	return l.parseFS(fsys, subPath)
}

func (l *Loader) parseFS(fsys fs.FS, subPath string) (map[string]*Agent, error) {
	agents := make(map[string]*Agent)

	err := fs.WalkDir(fsys, subPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, markdownExt) {
			return nil
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		agent, err := parseAgent(data)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		if _, exists := agents[agent.Name]; exists {
			return fmt.Errorf("duplicate agent name: %s", agent.Name)
		}
		agents[agent.Name] = agent

		return nil
	})

	if err != nil {
		return nil, err
	}

	return agents, nil
}
