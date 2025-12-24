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
	// Agent types
	AgentTypeGeneral = "general"
	AgentTypeOpenAI  = "openai"

	// General agents
	AgentEngineering = "Engineering"
	AgentDesign      = "Design"

	// OpenAI agents
	AgentGPT             = "GPT"
	AgentGPT5Codex       = "GPT-5 Codex"
	AgentGPT51CodexMax   = "GPT-5.1 Codex Max"
	AgentGPT52Codex      = "GPT-5.2 Codex"
	AgentGPT51           = "GPT-5.1"
	AgentGPT52           = "GPT-5.2"
)
