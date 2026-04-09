package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type LimitStep struct {
	Count int `yaml:"count"`
}

func (s *LimitStep) StepType() string {
	return "limit"
}

func (s *LimitStep) String() string {
	return fmt.Sprintf("[%s] count=%d", s.StepType(), s.Count)
}

func (s *LimitStep) UnmarshalYAML(b []byte) error {
	type rawLimitStep LimitStep
	var raw rawLimitStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	*s = LimitStep(raw)
	return s.Validate()
}

func (s *LimitStep) Validate() error {
	if s.Count < 1 {
		return fmt.Errorf("limit count must be at least 1")
	}
	return nil
}
