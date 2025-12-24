# Agent Hub Go SDK

Go SDK for [agent-hub](https://github.com/workpi-ai/agent-hub) - A centralized repository for AI agent configurations.

## Features

- 🚀 **Embedded data**: Works offline with embedded agent configurations
- 🔄 **Auto-update**: Automatically checks for updates on startup
- 📦 **Lightweight**: Minimal overhead with embedded markdown files
- 🤖 **Multi-agent**: Supports various specialized agents
- 🏷️ **Type filtering**: Filter agents by type (general, openai, etc.)
- 📂 **Nested structure**: Supports hierarchical agent organization

## Installation

```bash
go get github.com/workpi-ai/agent-hub-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/workpi-ai/agent-hub-go/pkg/hub"
)

func main() {
    home, _ := os.UserHomeDir()
    configDir := filepath.Join(home, ".codev", "agents", "registry")
    
    h, err := hub.New(hub.Options{
        LocalStandardAgentsDir:   filepath.Join(configDir, "agents"),
        LocalStandardCommandsDir: filepath.Join(configDir, "commands"),
        MetadataFile:             filepath.Join(configDir, "metadata.json"),
        AutoUpdate:               false,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer h.Close()
    
    // Get agent by full name (includes path prefix for nested agents)
    agent, _ := h.Agent(hub.AgentEngineering)  // "Engineering"
    fmt.Printf("Agent: %s (Type: %s)\n", agent.Name, agent.Type)
    fmt.Printf("Description: %s\n", agent.Description)
    fmt.Printf("Tools: %v\n", agent.Tools)
    
    // List all agents
    agents := h.Agents()
    fmt.Printf("Total agents: %d\n", len(agents))
    for _, a := range agents {
        fmt.Printf("  - %s (%s): %s\n", a.Name, a.Type, a.Description)
    }
    
    // Filter agents by type
    generalAgents := h.AgentsByType(hub.AgentTypeGeneral)
    fmt.Printf("\nGeneral agents: %d\n", len(generalAgents))
    for _, a := range generalAgents {
        fmt.Printf("  - %s: %s\n", a.Name, a.Description)
    }
}
```

## Agent Naming Convention

Agents use the `name` field from their frontmatter metadata, regardless of file location:

- `Engineering` - Agent in `agents/general/engineering.md` with `name: Engineering`
- `GPT-5 Codex` - Agent in `agents/openai/gpt-5-codex.md` with `name: GPT-5 Codex`
- Agent names must be unique across all directories

Use predefined constants from `hub` package:
- `hub.AgentEngineering` - "Engineering"
- `hub.AgentGPT5Codex` - "GPT-5 Codex"
- etc.

## API Reference

### Hub Methods

- `Agent(name string) (*Agent, error)` - Get agent by full name
- `Agents() []*Agent` - List all agents sorted by name
- `AgentsByType(agentType string) []*Agent` - Filter agents by type
- `Command(name string) (*Command, error)` - Get command by name
- `Commands() []*Command` - List all commands sorted by name
- `ForceUpdate() error` - Manually trigger update from GitHub
- `Close() error` - Stop auto-update loop and cleanup

### Agent Types

- `hub.AgentTypeGeneral` - "general"
- `hub.AgentTypeOpenAI` - "openai"

## Data Priority

1. **Local cache** (`$HOME/.codev/agents/registry/agents/`) - Downloaded from GitHub Release
2. **Embedded data** - Bundled from the agent-hub Go module dependency

## How It Works

1. **Compile Time**: Embeds agent data from the agent-hub Go module
2. **Runtime**: 
   - On startup, checks for updates from GitHub Release (if AutoUpdate is enabled)
   - Downloads new version to local cache if available
   - Loads data with priority: local cache > embedded data

## License

MIT
