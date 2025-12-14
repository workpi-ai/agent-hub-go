package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/workpi-ai/agent-hub-go/pkg/hub"
)

func main() {
	h, err := hub.New(hub.Options{
		LocalStandardAgentsDir:   filepath.Join("testdata", "agents", "standard"),
		LocalStandardCommandsDir: filepath.Join("testdata", "commands", "standard"),
		MetadataFile:             filepath.Join("testdata", "metadata.json"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	for _, agent := range h.Agents() {
		fmt.Printf("Agent: %s - %s\n", agent.Name, agent.Description)
	}

	for _, cmd := range h.Commands() {
		fmt.Printf("Command: %s - %s\n", cmd.Name, cmd.Description)
	}
}
