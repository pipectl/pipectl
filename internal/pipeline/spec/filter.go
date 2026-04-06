package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type FilterStep struct {
	Field      string `yaml:"field"`
	Equals     string `yaml:"equals"`
	NotEquals  string `yaml:"not-equals"`
	Contains   string `yaml:"contains"`
	StartsWith string `yaml:"starts-with"`
}

func (s *FilterStep) StepType() string {
	return "filter"
}

func (s *FilterStep) String() string {
	switch {
	case s.Equals != "":
		return fmt.Sprintf("[%s] filter: %v equals %v", s.StepType(), s.Field, s.Equals)
	case s.NotEquals != "":
		return fmt.Sprintf("[%s] filter: %v not-equals %v", s.StepType(), s.Field, s.NotEquals)
	case s.Contains != "":
		return fmt.Sprintf("[%s] filter: %v contains %v", s.StepType(), s.Field, s.Contains)
	case s.StartsWith != "":
		return fmt.Sprintf("[%s] filter: %v starts-with %v", s.StepType(), s.Field, s.StartsWith)
	default:
		return fmt.Sprintf("[%s] filter: %v", s.StepType(), s.Field)
	}
}

func (s *FilterStep) UnmarshalYAML(b []byte) error {
	type rawFilterStep FilterStep
	var raw rawFilterStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}

	*s = FilterStep(raw)

	if s.Field == "" {
		return fmt.Errorf("filter field is required")
	}

	set := 0
	if s.Equals != "" {
		set++
	}
	if s.NotEquals != "" {
		set++
	}
	if s.Contains != "" {
		set++
	}
	if s.StartsWith != "" {
		set++
	}

	if set == 0 {
		return fmt.Errorf("filter requires exactly one operator: equals, not-equals, contains, or starts-with")
	}
	if set > 1 {
		return fmt.Errorf("filter requires exactly one operator: equals, not-equals, contains, or starts-with")
	}

	return nil
}
