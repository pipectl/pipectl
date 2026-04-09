package spec

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

type AssertStep struct {
	MinRecords   *int   `yaml:"min-records"`
	MaxRecords   *int   `yaml:"max-records"`
	RecordsEqual *int   `yaml:"records-equal"`
	FieldExists  string `yaml:"field-exists"`
}

func (s *AssertStep) StepType() string {
	return "assert"
}

func (s *AssertStep) String() string {
	return fmt.Sprintf("[%s] min-records=%v max-records=%v records-equal=%v field-exists=%q", s.StepType(), s.MinRecords, s.MaxRecords, s.RecordsEqual, s.FieldExists)
}

func (s *AssertStep) UnmarshalYAML(b []byte) error {
	type rawAssertStep AssertStep
	var raw rawAssertStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	*s = AssertStep(raw)
	return s.Validate()
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

	if s.MinRecords == nil && s.MaxRecords == nil && s.RecordsEqual == nil && strings.TrimSpace(s.FieldExists) == "" {
		return fmt.Errorf("assert requires at least one option: min-records, max-records, records-equal, or field-exists")
	}

	return nil
}
