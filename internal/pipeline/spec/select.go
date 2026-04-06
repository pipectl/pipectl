package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type SelectStep struct {
	Fields []string `yaml:"fields"`
}

func (s *SelectStep) StepType() string {
	return "select"
}

func (s *SelectStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *SelectStep) UnmarshalYAML(b []byte) error {
	type rawSelectStep SelectStep
	var raw rawSelectStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}

	*s = SelectStep(raw)

	if len(s.Fields) == 0 {
		return fmt.Errorf("select requires at least one field")
	}

	return nil
}
