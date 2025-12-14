package hub

import (
	"sync"

	"github.com/workpi-ai/go-utils/ghrelease"
)

type Hub struct {
	agents   map[string]*Agent
	commands map[string]*Command

	mu        sync.RWMutex
	loader    *Loader
	updater   *ghrelease.Updater
	stopChan  chan struct{}
	closeOnce sync.Once
}

type Agent struct {
	Name         string   `yaml:"name" validate:"required"`
	Description  string   `yaml:"description" validate:"required"`
	Tools        []string `yaml:"tools"`
	SystemPrompt string   `validate:"required"`
}

type Command struct {
	Name        string `yaml:"name" validate:"required"`
	Description string `yaml:"description" validate:"required"`
	Prompt      string `validate:"required"`
}
