package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/steps/validate_json"
)

func (s *ValidateJSONStep) StepType() string {
	return "validate-json"
}

func (s *ValidateJSONStep) BuildExecutor() (engine.ExecutableStep, error) {
	return &validate_json.Step{
		Schema: s.Schema,
	}, nil
}

func (s *ValidateJSONStep) String() string {
	return fmt.Sprintf("[%s] schema: %v", s.StepType(), s.Schema)
}

type ValidateJSONStep struct {
	Schema string `yaml:"schema"`
}
