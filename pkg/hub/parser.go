package hub

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type frontmatterAgent struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Tools       []string `yaml:"tools"`
}

func parseAgent(path string, data []byte) (*Agent, error) {
	parts := bytes.SplitN(data, []byte(frontmatterDelimiter), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format in %s: expected --- delimiters", path)
	}

	var fm frontmatterAgent
	if err := yaml.Unmarshal(parts[1], &fm); err != nil {
		return nil, fmt.Errorf("parse frontmatter in %s: %w", path, err)
	}

	if fm.Name == "" {
		return nil, fmt.Errorf("agent name is required in %s", path)
	}
	if fm.Description == "" {
		return nil, fmt.Errorf("agent description is required in %s", path)
	}

	content := strings.TrimSpace(string(parts[2]))

	return &Agent{
		Name:         fm.Name,
		Description:  fm.Description,
		Tools:        fm.Tools,
		SystemPrompt: content,
	}, nil
}
