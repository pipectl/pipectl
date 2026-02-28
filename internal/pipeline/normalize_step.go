package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/steps/normalize"
)

type NormalizeStep struct {
	Fields map[string]string `yaml:"fields"`
}

func (s *NormalizeStep) StepType() string {
	return "normalize"
}

func (s *NormalizeStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *NormalizeStep) BuildExecutor() (engine.ExecutableStep, error) {
	return &normalize.Step{
		Fields: s.Fields,
	}, nil
}
