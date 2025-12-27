package hub

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	registry "github.com/workpi-ai/agent-hub"
)

var validate = validator.New()

type Loader struct {
	localStandardAgentsDir   string
	localUserAgentsDirs      []string
	localStandardCommandsDir string
	localUserCommandsDirs    []string
}

func NewLoader(localStandardAgentsDir string, localUserAgentsDirs []string, localStandardCommandsDir string, localUserCommandsDirs []string) *Loader {
	return &Loader{
		localStandardAgentsDir:   localStandardAgentsDir,
		localUserAgentsDirs:      localUserAgentsDirs,
		localStandardCommandsDir: localStandardCommandsDir,
		localUserCommandsDirs:    localUserCommandsDirs,
	}
}

func (l *Loader) Load() (map[string]*Agent, map[string]*Command, error) {
	agents, agentsErr := l.loadAgents()
	commands, commandsErr := l.loadCommands()

	if agentsErr != nil || commandsErr != nil {
		return agents, commands, errors.Join(agentsErr, commandsErr)
	}

	return agents, commands, nil
}

func (l *Loader) loadAgents() (map[string]*Agent, error) {
	agents, err := l.loadStandardAgents()
	if err != nil {
		return nil, fmt.Errorf("load standard agents: %w", err)
	}

	userAgents, err := l.loadUserAgents()
	if err != nil {
		return nil, fmt.Errorf("load user agents: %w", err)
	}

	for name, agent := range userAgents {
		agents[name] = agent
	}

	return agents, nil
}

func (l *Loader) loadStandardAgents() (map[string]*Agent, error) {
	fsys, subPath := l.getFS(l.localStandardAgentsDir, registry.AgentsDir)
	return l.loadAgentsFromFS(fsys, subPath)
}

func (l *Loader) loadUserAgents() (map[string]*Agent, error) {
	agents := make(map[string]*Agent)
	for _, dir := range l.localUserAgentsDirs {
		if !l.dirExists(dir) {
			continue
		}
		dirAgents, err := l.loadAgentsFromFS(os.DirFS(dir), currentDir)
		if err != nil {
			return nil, err
		}
		for name, agent := range dirAgents {
			agents[name] = agent
		}
	}
	return agents, nil
}

func (l *Loader) loadAgentsFromFS(fsys fs.FS, subPath string) (map[string]*Agent, error) {
	agents := make(map[string]*Agent)
	var errs []error

	err := l.walkMarkdownFiles(fsys, subPath, func(path string, data []byte) error {
		agent, err := parseAgent(data)
		if err != nil {
			errs = append(errs, fmt.Errorf("parse %s: %w", path, err))
			return nil
		}

		if agent.Name == "" {
			relPath := path
			if subPath != currentDir {
				relPath = strings.TrimPrefix(path, subPath+"/")
			}
			agent.Name = strings.TrimSuffix(filepath.Base(relPath), markdownExt)
		}

		if err := validate.Struct(agent); err != nil {
			errs = append(errs, fmt.Errorf("validate %s: %w", path, err))
			return nil
		}

		if _, exists := agents[agent.Name]; exists {
			errs = append(errs, fmt.Errorf("duplicate agent name: %s (path: %s)", agent.Name, path))
		}
		agents[agent.Name] = agent
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(errs) > 0 {
		slog.Warn("agent load error", "error", errors.Join(errs...))
	}

	return agents, nil
}

func (l *Loader) loadCommands() (map[string]*Command, error) {
	commands, err := l.loadStandardCommands()
	if err != nil {
		return nil, fmt.Errorf("load standard commands: %w", err)
	}

	userCommands, err := l.loadUserCommands()
	if err != nil {
		return nil, fmt.Errorf("load user commands: %w", err)
	}

	for name, cmd := range userCommands {
		commands[name] = cmd
	}

	return commands, nil
}

func (l *Loader) loadStandardCommands() (map[string]*Command, error) {
	fsys, subPath := l.getFS(l.localStandardCommandsDir, registry.CommandsDir)
	return l.loadCommandsFromFS(fsys, subPath)
}

func (l *Loader) loadUserCommands() (map[string]*Command, error) {
	commands := make(map[string]*Command)
	for _, dir := range l.localUserCommandsDirs {
		if !l.dirExists(dir) {
			continue
		}
		dirCommands, err := l.loadCommandsFromFS(os.DirFS(dir), currentDir)
		if err != nil {
			return nil, err
		}
		for name, cmd := range dirCommands {
			commands[name] = cmd
		}
	}
	return commands, nil
}

func (l *Loader) loadCommandsFromFS(fsys fs.FS, subPath string) (map[string]*Command, error) {
	commands := make(map[string]*Command)
	var errs []error

	err := l.walkMarkdownFiles(fsys, subPath, func(path string, data []byte) error {
		cmd, err := parseCommand(data)
		if err != nil {
			errs = append(errs, fmt.Errorf("parse %s: %w", path, err))
			return nil
		}

		relPath := path
		if subPath != currentDir {
			relPath = strings.TrimPrefix(path, subPath+"/")
		}

		dir := filepath.Dir(relPath)
		if cmd.Name == "" {
			cmd.Name = strings.TrimSuffix(filepath.Base(relPath), markdownExt)
		}

		if dir != currentDir {
			cmd.Name = strings.ReplaceAll(dir, "/", ":") + ":" + cmd.Name
		}

		if err := validate.Struct(cmd); err != nil {
			errs = append(errs, fmt.Errorf("validate %s: %w", path, err))
			return nil
		}

		if _, exists := commands[cmd.Name]; exists {
			errs = append(errs, fmt.Errorf("duplicate command name: %s (path: %s)", cmd.Name, path))
		}
		commands[cmd.Name] = cmd
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(errs) > 0 {
		slog.Warn("command load error", "error", errors.Join(errs...))
	}

	return commands, nil
}

func (l *Loader) getFS(localPath, registryDir string) (fs.FS, string) {
	if l.dirExists(localPath) {
		return os.DirFS(localPath), currentDir
	}
	return registry.All, registryDir
}

func (l *Loader) dirExists(path string) bool {
	if path == "" {
		return false
	}
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func (l *Loader) walkMarkdownFiles(fsys fs.FS, subPath string, fn func(path string, data []byte) error) error {
	return fs.WalkDir(fsys, subPath, func(path string, d fs.DirEntry, err error) error {
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

		return fn(path, data)
	})
}
