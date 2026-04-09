package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type RenameStep struct {
	Fields map[string]string `yaml:"fields"`
}

func (s *RenameStep) StepType() string {
	return "rename"
}

func (s *RenameStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *RenameStep) UnmarshalYAML(b []byte) error {
	type rawRenameStep RenameStep
	var raw rawRenameStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	*s = RenameStep(raw)
	return s.Validate()
}

func (s *RenameStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("rename requires at least one field")
	}
	return nil
}
