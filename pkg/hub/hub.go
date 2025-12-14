package hub

import (
	"fmt"
	"os"
	"sort"
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
		agents:    make(map[string]*Agent),
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
	agents := make([]*Agent, 0, len(h.agents))
	for _, agent := range h.agents {
		agents = append(agents, agent)
	}
	h.mu.RUnlock()
	
	// Sort by name
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Name < agents[j].Name
	})
	return agents
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

	h.mu.Lock()
	h.agents = newAgents
	h.mu.Unlock()
	
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
