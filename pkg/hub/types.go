package hub

import "github.com/workpi-ai/go-utils/ghrelease"

type Hub struct {
	Agents map[string]*Agent

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
