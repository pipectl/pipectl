package spec

import (
	"fmt"
	"strconv"

	"github.com/goccy/go-yaml"
)

type FilterStep struct {
	Field       string `yaml:"field"`
	Equals      string `yaml:"equals"`
	NotEquals   string `yaml:"not-equals"`
	Contains    string `yaml:"contains"`
	StartsWith  string `yaml:"starts-with"`
	GreaterThan string `yaml:"greater-than"`
	LessThan    string `yaml:"less-than"`
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
	case s.GreaterThan != "":
		return fmt.Sprintf("[%s] filter: %v greater-than %v", s.StepType(), s.Field, s.GreaterThan)
	case s.LessThan != "":
		return fmt.Sprintf("[%s] filter: %v less-than %v", s.StepType(), s.Field, s.LessThan)
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
	if s.GreaterThan != "" {
		if _, err := strconv.ParseFloat(s.GreaterThan, 64); err != nil {
			return fmt.Errorf("filter greater-than must be a number")
		}
		set++
	}
	if s.LessThan != "" {
		if _, err := strconv.ParseFloat(s.LessThan, 64); err != nil {
			return fmt.Errorf("filter less-than must be a number")
		}
		set++
	}

	if set == 0 {
		return fmt.Errorf("filter requires exactly one operator: equals, not-equals, contains, starts-with, greater-than, or less-than")
	}
	if set > 1 {
		return fmt.Errorf("filter requires exactly one operator: equals, not-equals, contains, starts-with, greater-than, or less-than")
	}

	return nil
}
