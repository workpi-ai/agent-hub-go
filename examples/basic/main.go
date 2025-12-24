package main

import (
	"fmt"
	"log"

	"github.com/workpi-ai/agent-hub-go/pkg/hub"
)

func main() {
	// Use empty paths to load from embedded registry
	h, err := hub.New(hub.Options{
		LocalStandardAgentsDir:   "",
		LocalStandardCommandsDir: "",
		MetadataFile:             "",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	fmt.Println("=== All Agents ===")
	for _, agent := range h.Agents() {
		fmt.Printf("Agent: %s (Type: %s) - %s\n", agent.Name, agent.Type, agent.Description)
	}

	fmt.Println("\n=== General Agents ===")
	for _, agent := range h.AgentsByType(hub.AgentTypeGeneral) {
		fmt.Printf("Agent: %s - %s\n", agent.Name, agent.Description)
	}

	fmt.Println("\n=== OpenAI Agents ===")
	for _, agent := range h.AgentsByType(hub.AgentTypeOpenAI) {
		fmt.Printf("Agent: %s - %s\n", agent.Name, agent.Description)
	}

	fmt.Println("\n=== Get Specific Agent ===")
	agent, err := h.Agent(hub.AgentEngineering)
	if err != nil {
		log.Printf("Warning: %v\n", err)
	} else {
		fmt.Printf("Agent: %s (Type: %s)\n", agent.Name, agent.Type)
		fmt.Printf("Description: %s\n", agent.Description)
		fmt.Printf("Tools: %v\n", agent.Tools)
	}

	fmt.Println("\n=== All Commands ===")
	for _, cmd := range h.Commands() {
		fmt.Printf("Command: %s - %s\n", cmd.Name, cmd.Description)
	}
}
