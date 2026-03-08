package spec

import "fmt"

type DefaultStep struct {
	Fields map[string]interface{} `yaml:"fields"`
}

func (s *DefaultStep) StepType() string {
	return "default"
}

func (s *DefaultStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}
