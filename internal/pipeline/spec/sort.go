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
	var raw rawSortStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	*s = SortStep(raw)
	if s.Direction == "" {
		s.Direction = "asc"
	}
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
