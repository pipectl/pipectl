package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/steps/validatejson"
)

func (s *ValidateJSONStep) StepType() string {
	return "validate-json"
}

func (s *ValidateJSONStep) BuildExecutor() (engine.ExecutableStep, error) {
	return &validatejson.Step{
		Schema: s.Schema,
	}, nil
}

func (s *ValidateJSONStep) String() string {
	return fmt.Sprintf("[%s] schema: %v", s.StepType(), s.Schema)
}

type ValidateJSONStep struct {
	Schema string `yaml:"schema"`
}
