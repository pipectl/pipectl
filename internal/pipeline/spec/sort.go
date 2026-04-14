package spec

import "fmt"

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

func (s *SortStep) Validate() error {
	if s.Field == "" {
		return fmt.Errorf("sort field is required")
	}

	if s.Direction == "" {
		s.Direction = "asc"
	}

	if s.Direction != "asc" && s.Direction != "desc" {
		return fmt.Errorf("sort direction must be asc or desc")
	}

	return nil
}
