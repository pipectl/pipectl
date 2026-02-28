package spec

import "fmt"

type ValidateJSONStep struct {
	Schema string `yaml:"schema"`
}

func (s *ValidateJSONStep) StepType() string {
	return "validate-json"
}

func (s *ValidateJSONStep) String() string {
	return fmt.Sprintf("[%s] schema: %v", s.StepType(), s.Schema)
}
