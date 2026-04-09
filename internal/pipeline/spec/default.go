package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type DefaultStep struct {
	Fields map[string]interface{} `yaml:"fields"`
}

func (s *DefaultStep) StepType() string {
	return "default"
}

func (s *DefaultStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *DefaultStep) UnmarshalYAML(b []byte) error {
	type rawDefaultStep DefaultStep
	var raw rawDefaultStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	*s = DefaultStep(raw)
	return s.Validate()
}

func (s *DefaultStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("default requires at least one field")
	}
	return nil
}
