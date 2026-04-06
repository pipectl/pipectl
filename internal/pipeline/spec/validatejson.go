package spec

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

type ValidateJSONStep struct {
	Schema string `yaml:"schema"`
}

func (s *ValidateJSONStep) StepType() string {
	return "validate-json"
}

func (s *ValidateJSONStep) String() string {
	return fmt.Sprintf("[%s] schema: %v", s.StepType(), s.Schema)
}

func (s *ValidateJSONStep) UnmarshalYAML(b []byte) error {
	type rawValidateJSONStep ValidateJSONStep
	var raw rawValidateJSONStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}

	*s = ValidateJSONStep(raw)

	if strings.TrimSpace(s.Schema) == "" {
		return fmt.Errorf("validate-json schema is required")
	}

	return nil
}
