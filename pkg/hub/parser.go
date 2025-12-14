package hub

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type frontmatterAgent struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Tools       []string `yaml:"tools"`
}

func ParseAgent(data []byte) (*Agent, error) {
	return parseAgent(data)
}

func ParseAgentFromFile(filePath string) (*Agent, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", filePath, err)
	}
	agent, err := parseAgent(data)
	if err != nil {
		return nil, fmt.Errorf("parse file %s: %w", filePath, err)
	}
	return agent, nil
}

func parseAgent(data []byte) (*Agent, error) {
	parts := bytes.SplitN(data, []byte(frontmatterDelimiter), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format: expected --- delimiters")
	}

	var fm frontmatterAgent
	if err := yaml.Unmarshal(parts[1], &fm); err != nil {
		return nil, fmt.Errorf("parse frontmatter: %w", err)
	}

	if fm.Name == "" {
		return nil, fmt.Errorf("agent name is required")
	}
	if fm.Description == "" {
		return nil, fmt.Errorf("agent description is required")
	}

	content := strings.TrimSpace(string(parts[2]))

	return &Agent{
		Name:         fm.Name,
		Description:  fm.Description,
		Tools:        fm.Tools,
		SystemPrompt: content,
	}, nil
}
