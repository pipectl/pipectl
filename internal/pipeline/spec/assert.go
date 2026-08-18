package spec

import (
	"fmt"
	"regexp"
	"strings"
)

type AssertFieldCheck struct {
	Field string `yaml:"field"`
	Value string `yaml:"value"`
}

type AssertStep struct {
	MinRecords    *int              `yaml:"min-records"`
	MaxRecords    *int              `yaml:"max-records"`
	RecordsEqual  *int              `yaml:"records-equal"`
	FieldExists   string            `yaml:"field-exists"`
	FieldEquals   *AssertFieldCheck `yaml:"field-equals"`
	FieldContains *AssertFieldCheck `yaml:"field-contains"`
	FieldMatches  *AssertFieldCheck `yaml:"field-matches"`
	CaseSensitive *bool             `yaml:"case-sensitive,omitempty"`
}

func (s *AssertStep) StepType() string {
	return "assert"
}

func (s *AssertStep) Validate() error {
	if s.MinRecords != nil && *s.MinRecords < 0 {
		return fmt.Errorf("assert min-records must be >= 0")
	}

	if s.MaxRecords != nil && *s.MaxRecords < 0 {
		return fmt.Errorf("assert max-records must be >= 0")
	}

	if s.RecordsEqual != nil && *s.RecordsEqual < 0 {
		return fmt.Errorf("assert records-equal must be >= 0")
	}

	if s.MinRecords != nil && s.MaxRecords != nil && *s.MinRecords > *s.MaxRecords {
		return fmt.Errorf("assert min-records must be <= max-records")
	}

	if s.RecordsEqual != nil && s.MinRecords != nil && *s.RecordsEqual < *s.MinRecords {
		return fmt.Errorf("assert records-equal must be >= min-records")
	}

	if s.RecordsEqual != nil && s.MaxRecords != nil && *s.RecordsEqual > *s.MaxRecords {
		return fmt.Errorf("assert records-equal must be <= max-records")
	}

	if s.FieldExists != "" && strings.TrimSpace(s.FieldExists) == "" {
		return fmt.Errorf("assert field-exists must be a non-empty string")
	}

	if err := validateAssertFieldCheck("field-equals", s.FieldEquals); err != nil {
		return err
	}
	if err := validateAssertFieldCheck("field-contains", s.FieldContains); err != nil {
		return err
	}
	if err := validateAssertFieldCheck("field-matches", s.FieldMatches); err != nil {
		return err
	}
	if s.FieldMatches != nil {
		if _, err := regexp.Compile(s.FieldMatches.Value); err != nil {
			return fmt.Errorf("assert field-matches value must be a valid regular expression: %v", err)
		}
	}

	if s.MinRecords == nil && s.MaxRecords == nil && s.RecordsEqual == nil &&
		strings.TrimSpace(s.FieldExists) == "" &&
		s.FieldEquals == nil && s.FieldContains == nil && s.FieldMatches == nil {
		return fmt.Errorf("assert requires at least one option: min-records, max-records, records-equal, field-exists, field-equals, field-contains, or field-matches")
	}

	return nil
}

func validateAssertFieldCheck(option string, check *AssertFieldCheck) error {
	if check == nil {
		return nil
	}
	if strings.TrimSpace(check.Field) == "" {
		return fmt.Errorf("assert %s field must be a non-empty string", option)
	}
	return nil
}
