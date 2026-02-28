package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/steps/redact"
)

type RedactStep struct {
	Strategy string   `yaml:"strategy"`
	Fields   []string `yaml:"fields"`
}

func (s *RedactStep) StepType() string {
	return "redact"
}

func (s *RedactStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *RedactStep) BuildExecutor() (engine.ExecutableStep, error) {
	return &redact.Step{
		Fields:   s.Fields,
		Strategy: s.Strategy,
	}, nil
}
