package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/steps/filter"
)

type FilterStep struct {
	Field  string `yaml:"field"`
	Equals string `yaml:"equals"`
}

func (s *FilterStep) StepType() string {
	return "filter"
}

func (s *FilterStep) BuildExecutor() (engine.ExecutableStep, error) {
	return &filter.Step{
		Field: s.Field,
		Value: s.Equals,
	}, nil
}

func (s *FilterStep) String() string {
	return fmt.Sprintf("[%s] filter: %v = %v", s.StepType(), s.Field, s.Equals)
}
