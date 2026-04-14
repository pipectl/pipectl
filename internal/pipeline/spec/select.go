package spec

import "fmt"

type SelectStep struct {
	Fields []string `yaml:"fields"`
}

func (s *SelectStep) StepType() string {
	return "select"
}

func (s *SelectStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *SelectStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("select requires at least one field")
	}
	return nil
}
