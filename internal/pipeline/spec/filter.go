package spec

import "fmt"

type FilterStep struct {
	Field  string `yaml:"field"`
	Equals string `yaml:"equals"`
}

func (s *FilterStep) StepType() string {
	return "filter"
}

func (s *FilterStep) String() string {
	return fmt.Sprintf("[%s] filter: %v = %v", s.StepType(), s.Field, s.Equals)
}
