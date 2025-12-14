# Agent Hub Go SDK

Go SDK for [agent-hub](https://github.com/workpi-ai/agent-hub) - A centralized repository for AI agent configurations.

## Features

- 🚀 **Embedded data**: Works offline with embedded agent configurations
- 🔄 **Auto-update**: Automatically checks for updates on startup
- 📦 **Lightweight**: Minimal overhead with embedded markdown files
- 🤖 **Multi-agent**: Supports various specialized agents

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
        ConfigDir:  configDir,
        AutoUpdate: false,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer h.Close()
    
    // Get agent
    agent, _ := h.GetAgent("engineering")
    fmt.Printf("Agent: %s\n", agent.Name)
    fmt.Printf("Description: %s\n", agent.Description)
    fmt.Printf("Tools: %v\n", agent.Tools)
    
    // List all agents
    agents := h.ListAgents()
    fmt.Printf("Total agents: %d\n", len(agents))
}
```

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
