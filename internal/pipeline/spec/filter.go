package spec

import (
	"fmt"
	"strconv"

	"github.com/goccy/go-yaml"
)

type FilterCondition struct {
	// Leaf rule fields
	Field       string `yaml:"field"`
	Equals      string `yaml:"equals"`
	NotEquals   string `yaml:"not-equals"`
	Contains    string `yaml:"contains"`
	StartsWith  string `yaml:"starts-with"`
	EndsWith    string `yaml:"ends-with"`
	GreaterThan string `yaml:"greater-than"`
	LessThan    string `yaml:"less-than"`
	// Group fields
	All []FilterCondition `yaml:"all"`
	Any []FilterCondition `yaml:"any"`
}

func (c *FilterCondition) UnmarshalYAML(b []byte) error {
	type rawFilterCondition FilterCondition
	var raw rawFilterCondition
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	*c = FilterCondition(raw)

	isGroup := len(c.All) > 0 || len(c.Any) > 0
	isLeaf := c.Field != "" || c.Equals != "" || c.NotEquals != "" || c.Contains != "" ||
		c.StartsWith != "" || c.EndsWith != "" || c.GreaterThan != "" || c.LessThan != ""

	if isGroup && isLeaf {
		return fmt.Errorf("filter condition cannot mix group (all/any) and rule fields")
	}

	if isGroup {
		if len(c.All) > 0 && len(c.Any) > 0 {
			return fmt.Errorf("filter condition cannot specify both all and any")
		}
		return nil
	}

	return validateFilterRule(c.Field, c.Equals, c.NotEquals, c.Contains, c.StartsWith, c.EndsWith, c.GreaterThan, c.LessThan)
}

type FilterStep struct {
	Field       string            `yaml:"field"`
	Equals      string            `yaml:"equals"`
	NotEquals   string            `yaml:"not-equals"`
	Contains    string            `yaml:"contains"`
	StartsWith  string            `yaml:"starts-with"`
	EndsWith    string            `yaml:"ends-with"`
	GreaterThan string            `yaml:"greater-than"`
	LessThan    string            `yaml:"less-than"`
	All         []FilterCondition `yaml:"all"`
	Any         []FilterCondition `yaml:"any"`
}

func (s *FilterStep) StepType() string {
	return "filter"
}

func (s *FilterStep) String() string {
	switch {
	case len(s.All) > 0:
		return fmt.Sprintf("[%s] all: %d conditions", s.StepType(), len(s.All))
	case len(s.Any) > 0:
		return fmt.Sprintf("[%s] any: %d conditions", s.StepType(), len(s.Any))
	case s.Equals != "":
		return fmt.Sprintf("[%s] filter: %v equals %v", s.StepType(), s.Field, s.Equals)
	case s.NotEquals != "":
		return fmt.Sprintf("[%s] filter: %v not-equals %v", s.StepType(), s.Field, s.NotEquals)
	case s.Contains != "":
		return fmt.Sprintf("[%s] filter: %v contains %v", s.StepType(), s.Field, s.Contains)
	case s.StartsWith != "":
		return fmt.Sprintf("[%s] filter: %v starts-with %v", s.StepType(), s.Field, s.StartsWith)
	case s.EndsWith != "":
		return fmt.Sprintf("[%s] filter: %v ends-with %v", s.StepType(), s.Field, s.EndsWith)
	case s.GreaterThan != "":
		return fmt.Sprintf("[%s] filter: %v greater-than %v", s.StepType(), s.Field, s.GreaterThan)
	case s.LessThan != "":
		return fmt.Sprintf("[%s] filter: %v less-than %v", s.StepType(), s.Field, s.LessThan)
	default:
		return fmt.Sprintf("[%s] filter: %v", s.StepType(), s.Field)
	}
}

func (s *FilterStep) Validate() error {
	hasGroup := len(s.All) > 0 || len(s.Any) > 0
	hasFlat := s.Field != "" || s.Equals != "" || s.NotEquals != "" || s.Contains != "" ||
		s.StartsWith != "" || s.EndsWith != "" || s.GreaterThan != "" || s.LessThan != ""

	if hasGroup && hasFlat {
		return fmt.Errorf("filter cannot mix group (all/any) and rule fields")
	}

	if hasGroup {
		if len(s.All) > 0 && len(s.Any) > 0 {
			return fmt.Errorf("filter cannot specify both all and any at the top level")
		}
		return nil
	}

	if !hasFlat {
		return fmt.Errorf("filter requires a condition: specify field with an operator, or use all/any for grouped conditions")
	}

	return validateFilterRule(s.Field, s.Equals, s.NotEquals, s.Contains, s.StartsWith, s.EndsWith, s.GreaterThan, s.LessThan)
}

func validateFilterRule(field, equals, notEquals, contains, startsWith, endsWith, greaterThan, lessThan string) error {
	if field == "" {
		return fmt.Errorf("filter field is required")
	}

	set := 0
	if equals != "" {
		set++
	}
	if notEquals != "" {
		set++
	}
	if contains != "" {
		set++
	}
	if startsWith != "" {
		set++
	}
	if endsWith != "" {
		set++
	}
	if greaterThan != "" {
		if _, err := strconv.ParseFloat(greaterThan, 64); err != nil {
			return fmt.Errorf("filter greater-than must be a number")
		}
		set++
	}
	if lessThan != "" {
		if _, err := strconv.ParseFloat(lessThan, 64); err != nil {
			return fmt.Errorf("filter less-than must be a number")
		}
		set++
	}

	if set == 0 {
		return fmt.Errorf("filter requires exactly one operator: equals, not-equals, contains, starts-with, ends-with, greater-than, or less-than")
	}
	if set > 1 {
		return fmt.Errorf("filter requires exactly one operator: equals, not-equals, contains, starts-with, ends-with, greater-than, or less-than")
	}

	return nil
}
