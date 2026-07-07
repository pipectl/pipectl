package spec

import (
	"fmt"
	"strings"
)

type CastField struct {
	Type        string   `yaml:"type"`
	Format      string   `yaml:"format,omitempty"`
	TrueValues  []string `yaml:"true-values,omitempty"`
	FalseValues []string `yaml:"false-values,omitempty"`
}

type CastStep struct {
	Fields map[string]CastField `yaml:"fields"`
}

func (s *CastStep) StepType() string {
	return "cast"
}

func (s *CastStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *CastStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("cast requires at least one field")
	}

	for fieldName, field := range s.Fields {
		switch field.Type {
		case "int", "float", "bool", "time", "string":
		default:
			return fmt.Errorf("cast field %q type must be one of: int, float, bool, time, string", fieldName)
		}

		if field.Type != "time" && field.Format != "" {
			return fmt.Errorf("cast field %q format is only supported for type time", fieldName)
		}

		if field.Type != "bool" && (len(field.TrueValues) > 0 || len(field.FalseValues) > 0) {
			return fmt.Errorf("cast field %q true-values/false-values are only supported for type bool", fieldName)
		}

		if field.Type == "bool" {
			seen := make(map[string]struct{}, len(field.TrueValues))
			for _, value := range field.TrueValues {
				seen[strings.ToLower(strings.TrimSpace(value))] = struct{}{}
			}
			for _, value := range field.FalseValues {
				normalized := strings.ToLower(strings.TrimSpace(value))
				if _, exists := seen[normalized]; exists {
					return fmt.Errorf("cast field %q bool true-values and false-values must not overlap", fieldName)
				}
			}
		}
	}

	return nil
}
