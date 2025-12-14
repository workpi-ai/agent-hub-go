package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/workpi-ai/agent-hub-go/pkg/hub"
)

func main() {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".codev", "agents", "registry")

	h, err := hub.New(hub.Options{
		ConfigDir:     configDir,
		AutoUpdate:    false,
		CheckInterval: 1 * time.Hour,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	// List all available agents
	agents := h.ListAgents()
	fmt.Printf("Available agents (%d):\n", len(agents))
	for _, name := range agents {
		fmt.Printf("  - %s\n", name)
	}
	fmt.Println()

	// Get engineering agent details
	agent, err := h.GetAgent("engineering")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Agent: %s\n", agent.Name)
	fmt.Printf("Description: %s\n", agent.Description)
	fmt.Printf("Tools: %v\n", agent.Tools)
	fmt.Printf("SystemPrompt length: %d bytes\n", len(agent.SystemPrompt))
}
