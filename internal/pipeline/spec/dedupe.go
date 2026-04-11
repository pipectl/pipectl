package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type DedupeStep struct {
	Fields        []string `yaml:"fields"`
	CaseSensitive bool     `yaml:"case-sensitive,omitempty"`
}

func (s *DedupeStep) StepType() string {
	return "dedupe"
}

func (s *DedupeStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *DedupeStep) UnmarshalYAML(b []byte) error {
	type rawDedupeStep DedupeStep
	raw := rawDedupeStep{CaseSensitive: true}
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	*s = DedupeStep(raw)
	return s.Validate()
}

func (s *DedupeStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("dedupe fields is required")
	}
	return nil
}
