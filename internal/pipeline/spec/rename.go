package spec

import "fmt"

type RenameStep struct {
	Fields map[string]string `yaml:"fields"`
}

func (s *RenameStep) StepType() string {
	return "rename"
}

func (s *RenameStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *RenameStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("rename requires at least one field")
	}
	return nil
}
