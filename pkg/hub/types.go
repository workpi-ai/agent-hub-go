package hub

import (
	"sync"

	"github.com/workpi-ai/go-utils/ghrelease"
)

type Hub struct {
	agents map[string]*Agent

	mu        sync.RWMutex
	configDir string
	loader    *Loader
	updater   *ghrelease.Updater
	stopChan  chan struct{}
}

type Agent struct {
	Name         string
	Description  string
	Tools        []string
	SystemPrompt string
}
