package hub

import "time"

const (
	DefaultCheckInterval = 1 * time.Hour
	repoOwner            = "workpi-ai"
	repoName             = "agent-hub"
	markdownExt          = ".md"
	frontmatterDelimiter = "---"
	currentDir           = "."
)

const (
	AgentEngineering = "Engineering"
	AgentPlan        = "Plan"
	AgentDebug       = "Debug"
)
