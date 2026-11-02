package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type SortStep struct {
	Field     string    `yaml:"field,omitempty"`
	Direction string    `yaml:"direction,omitempty"`
	Nulls     string    `yaml:"nulls,omitempty"`
	Keys      []SortKey `yaml:"keys,omitempty"`
}

// SortKey is one entry in a compound sort. Keys are compared in order; later
// keys only break ties left by earlier ones.
type SortKey struct {
	Field     string `yaml:"field"`
	Direction string `yaml:"direction,omitempty"`
	Nulls     string `yaml:"nulls,omitempty"`
}

func (s *SortStep) StepType() string {
	return "sort"
}

func (s *SortStep) UnmarshalYAML(b []byte) error {
	type rawSortStep SortStep
	// Direction is not pre-defaulted so that setting it alongside keys can be
	// detected and rejected.
	raw := rawSortStep{Nulls: "last"}
	if err := yaml.UnmarshalWithOptions(b, &raw, yaml.DisallowUnknownField()); err != nil {
		return err
	}
	*s = SortStep(raw)

	if len(s.Keys) == 0 {
		if s.Direction == "" {
			s.Direction = "asc"
		}
	} else {
		for i := range s.Keys {
			if s.Keys[i].Direction == "" {
				s.Keys[i].Direction = "asc"
			}
			if s.Keys[i].Nulls == "" {
				s.Keys[i].Nulls = s.Nulls
			}
		}
	}

	return s.Validate()
}

func (s *SortStep) Validate() error {
	// Checked first: keys inherit this value, so an invalid one would otherwise
	// be reported against keys[0].
	if s.Nulls != "first" && s.Nulls != "last" {
		return fmt.Errorf("sort nulls must be first or last")
	}

	if len(s.Keys) > 0 {
		if s.Field != "" || s.Direction != "" {
			return fmt.Errorf("sort keys cannot be combined with field or direction")
		}
		for i, k := range s.Keys {
			if k.Field == "" {
				return fmt.Errorf("sort keys[%d] field is required", i)
			}
			if k.Direction != "asc" && k.Direction != "desc" {
				return fmt.Errorf("sort keys[%d] direction must be asc or desc", i)
			}
			if k.Nulls != "first" && k.Nulls != "last" {
				return fmt.Errorf("sort keys[%d] nulls must be first or last", i)
			}
		}
	} else {
		if s.Field == "" {
			return fmt.Errorf("sort field or keys is required")
		}

		if s.Direction != "asc" && s.Direction != "desc" {
			return fmt.Errorf("sort direction must be asc or desc")
		}
	}

	return nil
}

// SortKeys returns the step's sort keys, normalising the single-field form
// into a one-element list.
func (s *SortStep) SortKeys() []SortKey {
	if len(s.Keys) > 0 {
		return s.Keys
	}
	return []SortKey{{Field: s.Field, Direction: s.Direction, Nulls: s.Nulls}}
}
