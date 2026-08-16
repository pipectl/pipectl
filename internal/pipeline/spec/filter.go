package spec

import (
	"fmt"
	"regexp"

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
	Matches     string `yaml:"matches"`
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
	return c.Validate()
}

func (c *FilterCondition) Validate() error {
	return validateFilterCondition(c.All, c.Any, c.Field, c.Equals, c.NotEquals, c.Contains, c.StartsWith, c.EndsWith, c.Matches, c.GreaterThan, c.LessThan)
}

type FilterStep struct {
	Field         string            `yaml:"field"`
	Equals        string            `yaml:"equals"`
	NotEquals     string            `yaml:"not-equals"`
	Contains      string            `yaml:"contains"`
	StartsWith    string            `yaml:"starts-with"`
	EndsWith      string            `yaml:"ends-with"`
	Matches       string            `yaml:"matches"`
	GreaterThan   string            `yaml:"greater-than"`
	LessThan      string            `yaml:"less-than"`
	All           []FilterCondition `yaml:"all"`
	Any           []FilterCondition `yaml:"any"`
	OnMissing     string            `yaml:"on-missing,omitempty"`
	CaseSensitive bool              `yaml:"case-sensitive,omitempty"`
}

func (s *FilterStep) StepType() string {
	return "filter"
}

func (s *FilterStep) UnmarshalYAML(b []byte) error {
	type rawFilterStep FilterStep
	raw := rawFilterStep{OnMissing: "exclude", CaseSensitive: true}
	if err := yaml.UnmarshalWithOptions(b, &raw, yaml.DisallowUnknownField()); err != nil {
		return err
	}
	// An empty step body (e.g. "filter:\n") unmarshals as a no-op that resets
	// the whole raw struct to its zero value, wiping the pre-set default.
	if raw.OnMissing == "" {
		raw.OnMissing = "exclude"
	}
	*s = FilterStep(raw)
	return s.Validate()
}

func (s *FilterStep) Validate() error {
	// "" is accepted alongside "exclude" because a fully empty step body (e.g.
	// "filter:\n") never reaches UnmarshalYAML's default-seeding, so Validate
	// can be called directly against a zero-valued FilterStep; the engine
	// treats an unset OnMissing the same as "exclude".
	switch s.OnMissing {
	case "", "exclude", "include", "error":
	default:
		return fmt.Errorf("filter on-missing must be exclude, include, or error")
	}
	return validateFilterCondition(s.All, s.Any, s.Field, s.Equals, s.NotEquals, s.Contains, s.StartsWith, s.EndsWith, s.Matches, s.GreaterThan, s.LessThan)
}

func validateFilterCondition(all, any []FilterCondition, field, equals, notEquals, contains, startsWith, endsWith, matches, greaterThan, lessThan string) error {
	isGroup := len(all) > 0 || len(any) > 0
	isLeaf := field != "" || equals != "" || notEquals != "" || contains != "" ||
		startsWith != "" || endsWith != "" || matches != "" || greaterThan != "" || lessThan != ""

	if isGroup && isLeaf {
		return fmt.Errorf("filter cannot mix group (all/any) and rule fields")
	}

	if isGroup {
		if len(all) > 0 && len(any) > 0 {
			return fmt.Errorf("filter cannot specify both all and any")
		}
		return nil
	}

	if !isLeaf {
		return fmt.Errorf("filter requires a condition: specify field with an operator, or use all/any for grouped conditions")
	}

	return validateFilterRule(field, equals, notEquals, contains, startsWith, endsWith, matches, greaterThan, lessThan)
}

func validateFilterRule(field, equals, notEquals, contains, startsWith, endsWith, matches, greaterThan, lessThan string) error {
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
	if matches != "" {
		if _, err := regexp.Compile(matches); err != nil {
			return fmt.Errorf("filter matches must be a valid regular expression: %v", err)
		}
		set++
	}
	if greaterThan != "" {
		set++
	}
	if lessThan != "" {
		set++
	}

	if set == 0 {
		return fmt.Errorf("filter requires exactly one operator: equals, not-equals, contains, starts-with, ends-with, matches, greater-than, or less-than")
	}
	if set > 1 {
		return fmt.Errorf("filter requires exactly one operator: equals, not-equals, contains, starts-with, ends-with, matches, greater-than, or less-than")
	}

	return nil
}
