package spec

import (
	"fmt"
	"strings"
)

type ValidateJSONStep struct {
	Schema string `yaml:"schema"`
}

func (s *ValidateJSONStep) StepType() string {
	return "validate-json"
}

func (s *ValidateJSONStep) Validate() error {
	if strings.TrimSpace(s.Schema) == "" {
		return fmt.Errorf("validate-json schema is required")
	}
	return nil
}
