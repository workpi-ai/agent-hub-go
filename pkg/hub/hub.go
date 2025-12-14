package hub

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Options struct {
	ConfigDir     string
	AutoUpdate    bool
	CheckInterval time.Duration
}

func New(opts Options) (*Hub, error) {
	if opts.ConfigDir == "" {
		return nil, fmt.Errorf("ConfigDir is required")
	}
	if opts.CheckInterval == 0 {
		opts.CheckInterval = DefaultCheckInterval
	}

	if err := os.MkdirAll(opts.ConfigDir, defaultDirPerm); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	hub := &Hub{
		Agents:    make(map[string]*Agent),
		configDir: opts.ConfigDir,
		loader:    NewLoader(opts.ConfigDir),
		updater:   NewUpdater(opts.ConfigDir),
		stopChan:  make(chan struct{}),
	}

	if err := hub.reload(); err != nil {
		return nil, err
	}

	if opts.AutoUpdate {
		go hub.autoUpdateLoop(opts.CheckInterval)
	}

	return hub, nil
}

func (h *Hub) GetAgent(name string) (*Agent, error) {
	agentKey := strings.ToLower(name)
	agent, ok := h.Agents[agentKey]
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", name)
	}
	return agent, nil
}

func (h *Hub) ListAgents() []string {
	names := make([]string, 0, len(h.Agents))
	for _, agent := range h.Agents {
		names = append(names, agent.Name)
	}
	sort.Strings(names)
	return names
}

func (h *Hub) ForceUpdate() error {
	if err := h.updater.Update(); err != nil {
		return err
	}
	return h.reload()
}

func (h *Hub) Close() error {
	close(h.stopChan)
	return nil
}

func (h *Hub) reload() error {
	newAgents, err := h.loader.Load()
	if err != nil {
		return err
	}

	h.Agents = newAgents
	return nil
}

func (h *Hub) autoUpdateLoop(interval time.Duration) {
	if err := h.updater.Update(); err == nil {
		_ = h.reload()
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := h.updater.Update(); err != nil {
				continue
			}

			if err := h.reload(); err != nil {
				continue
			}
		case <-h.stopChan:
			return
		}
	}
}
