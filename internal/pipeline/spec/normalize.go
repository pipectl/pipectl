package spec

import "fmt"

type NormalizeStep struct {
	Fields map[string]string `yaml:"fields"`
}

func (s *NormalizeStep) StepType() string {
	return "normalize"
}

func (s *NormalizeStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}
