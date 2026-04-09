package spec

import "fmt"

type LogStep struct {
	Message string `yaml:"message"`
	Count   *bool  `yaml:"count"`
	Sample  *int   `yaml:"sample"`
}

func (s *LogStep) StepType() string {
	return "log"
}

func (s *LogStep) String() string {
	return fmt.Sprintf("[%s] message=%q count=%v sample=%v", s.StepType(), s.Message, s.Count, s.Sample)
}

func (s *LogStep) Validate() error { return nil }
