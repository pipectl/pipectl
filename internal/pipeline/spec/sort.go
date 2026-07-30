package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type SortStep struct {
	Field     string `yaml:"field"`
	Direction string `yaml:"direction,omitempty"`
}

func (s *SortStep) StepType() string {
	return "sort"
}

func (s *SortStep) String() string {
	return fmt.Sprintf("[%s] field: %v direction: %v", s.StepType(), s.Field, s.Direction)
}

func (s *SortStep) UnmarshalYAML(b []byte) error {
	type rawSortStep SortStep
	raw := rawSortStep{Direction: "asc"}
	if err := yaml.UnmarshalWithOptions(b, &raw, yaml.DisallowUnknownField()); err != nil {
		return err
	}
	*s = SortStep(raw)
	return s.Validate()
}

func (s *SortStep) Validate() error {
	if s.Field == "" {
		return fmt.Errorf("sort field is required")
	}

	if s.Direction != "asc" && s.Direction != "desc" {
		return fmt.Errorf("sort direction must be asc or desc")
	}

	return nil
}
