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

func (s *DefaultStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("default requires at least one field")
	}
	return nil
}
