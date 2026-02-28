package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/steps"
	"github.com/shanebell/pipectl/internal/steps/filter"
)

type FilterStep struct {
	Field  string `yaml:"field"`
	Equals string `yaml:"equals"`
}

func (s *FilterStep) StepType() string {
	return "filter"
}

func (s *FilterStep) BuildExecutor() (steps.ExecutableStep, error) {
	return &filter.Step{
		Field: s.Field,
		Value: s.Equals,
	}, nil
}

func (s *FilterStep) String() string {
	return fmt.Sprintf("[%s] filter: %v = %v", s.StepType(), s.Field, s.Equals)
}
