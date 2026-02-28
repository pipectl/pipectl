package spec

import "fmt"

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
