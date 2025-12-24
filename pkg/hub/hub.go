package hub

import (
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/workpi-ai/go-utils/ghrelease"
	"golang.org/x/exp/maps"
)

type Options struct {
	LocalStandardAgentsDir   string
	LocalUserAgentsDirs      []string
	LocalStandardCommandsDir string
	LocalUserCommandsDirs    []string
	MetadataFile             string
	AutoUpdate               bool
	CheckInterval            time.Duration
}

func New(opts Options) (*Hub, error) {
	// Empty paths are allowed - will use embedded registry
	if opts.CheckInterval == 0 {
		opts.CheckInterval = DefaultCheckInterval
	}

	var updater *ghrelease.Updater
	var err error
	
	// Only create updater if local paths are provided
	if opts.LocalStandardAgentsDir != "" && opts.LocalStandardCommandsDir != "" {
		updater, err = NewUpdater(opts.MetadataFile, opts.LocalStandardAgentsDir, opts.LocalStandardCommandsDir)
		if err != nil {
			return nil, err
		}
	}

	hub := &Hub{
		agents:   make(map[string]*Agent),
		commands: make(map[string]*Command),
		loader:   NewLoader(opts.LocalStandardAgentsDir, opts.LocalUserAgentsDirs, opts.LocalStandardCommandsDir, opts.LocalUserCommandsDirs),
		updater:  updater,
		stopChan: make(chan struct{}),
	}

	if err := hub.reload(); err != nil {
		slog.Warn("initial load partially failed", "error", err)
	}

	if opts.AutoUpdate {
		go hub.autoUpdateLoop(opts.CheckInterval)
	}

	return hub, nil
}

func (h *Hub) Agent(name string) (*Agent, error) {
	h.mu.RLock()
	agent, ok := h.agents[name]
	h.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("agent not found: %s", name)
	}
	return agent, nil
}

func (h *Hub) Agents() []*Agent {
	h.mu.RLock()
	agents := maps.Values(h.agents)
	h.mu.RUnlock()

	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Name < agents[j].Name
	})
	return agents
}

func (h *Hub) AgentsByType(agentType string) []*Agent {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var filtered []*Agent
	for _, agent := range h.agents {
		if agent.Type == agentType {
			filtered = append(filtered, agent)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name < filtered[j].Name
	})
	return filtered
}

func (h *Hub) Command(name string) (*Command, error) {
	h.mu.RLock()
	cmd, ok := h.commands[name]
	h.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("command not found: %s", name)
	}
	return cmd, nil
}

func (h *Hub) Commands() []*Command {
	h.mu.RLock()
	commands := maps.Values(h.commands)
	h.mu.RUnlock()

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})
	return commands
}

func (h *Hub) ForceUpdate() error {
	if h.updater == nil {
		return fmt.Errorf("updater not available: using embedded registry")
	}
	if err := h.updater.Update(); err != nil {
		return err
	}
	return h.reload()
}

func (h *Hub) reload() error {
	newAgents, newCommands, err := h.loader.Load()

	h.mu.Lock()
	if newAgents != nil {
		h.agents = newAgents
	}
	if newCommands != nil {
		h.commands = newCommands
	}
	h.mu.Unlock()

	return err
}

func (h *Hub) autoUpdateLoop(interval time.Duration) {
	h.updateAndReload()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.updateAndReload()
		case <-h.stopChan:
			return
		}
	}
}

func (h *Hub) updateAndReload() {
	if h.updater == nil {
		return
	}
	if err := h.updater.Update(); err != nil {
		slog.Error("update failed", "error", err)
	}
	if err := h.reload(); err != nil {
		slog.Error("reload failed", "error", err)
	}
}

func (h *Hub) Close() error {
	h.closeOnce.Do(func() {
		close(h.stopChan)
	})
	return nil
}
