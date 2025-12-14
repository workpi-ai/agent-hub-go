package hub

import "time"

const (
	DefaultCheckInterval = 1 * time.Hour
	repoOwner = "workpi-ai"
	repoName  = "agent-hub"
	agentsDir = "agents"
	metadataFile = "metadata.json"
	markdownExt = ".md"
	defaultDirPerm = 0755
	frontmatterDelimiter = "---"
)

const (
	AgentEngineering = "engineering"
	AgentPlan        = "plan"
	AgentDebug       = "debug"
)
