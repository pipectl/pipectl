package spec

import "fmt"

type CountStep struct {
	Message string `yaml:"message"`
}

func (s *CountStep) StepType() string {
	return "count"
}

func (s *CountStep) String() string {
	return fmt.Sprintf("[%s] message=%q", s.StepType(), s.Message)
}
