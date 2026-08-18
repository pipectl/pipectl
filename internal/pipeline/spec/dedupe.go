package spec

import (
	"fmt"
)

type DedupeStep struct {
	Fields        []string `yaml:"fields"`
	CaseSensitive *bool    `yaml:"case-sensitive,omitempty"`
}

func (s *DedupeStep) StepType() string {
	return "dedupe"
}

func (s *DedupeStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("dedupe fields is required")
	}
	return nil
}
