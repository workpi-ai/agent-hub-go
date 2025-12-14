package hub

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func validateName(name string) error {
	if name != "" && !nameRegex.MatchString(name) {
		return fmt.Errorf("name must contain only letters and numbers: %s", name)
	}
	return nil
}

func parseAgent(data []byte) (*Agent, error) {
	agent, content, err := parseFrontmatter[Agent](data)
	if err != nil {
		return nil, err
	}

	if err := validateName(agent.Name); err != nil {
		return nil, err
	}

	agent.SystemPrompt = content
	return agent, nil
}

func parseCommand(data []byte) (*Command, error) {
	cmd, content, err := parseFrontmatter[Command](data)
	if err != nil {
		return nil, err
	}

	if err := validateName(cmd.Name); err != nil {
		return nil, err
	}

	cmd.Prompt = content
	return cmd, nil
}

func parseFrontmatter[T any](data []byte) (*T, string, error) {
	parts := bytes.SplitN(data, []byte(frontmatterDelimiter), 3)
	if len(parts) < 3 {
		var zero T
		return &zero, strings.TrimSpace(string(data)), nil
	}

	var fm T
	if err := yaml.Unmarshal(parts[1], &fm); err != nil {
		return nil, "", fmt.Errorf("parse frontmatter: %w", err)
	}

	content := strings.TrimSpace(string(parts[2]))

	return &fm, content, nil
}
