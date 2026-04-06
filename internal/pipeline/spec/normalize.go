package spec

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

var validNormalizeStrategies = []string{"lower", "upper", "trim", "trim-left", "trim-right", "collapse-spaces", "capitalize"}

type NormalizeStep struct {
	Fields map[string]string `yaml:"fields"`
}

func (s *NormalizeStep) StepType() string {
	return "normalize"
}

func (s *NormalizeStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *NormalizeStep) UnmarshalYAML(b []byte) error {
	type rawNormalizeStep NormalizeStep
	var raw rawNormalizeStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}

	*s = NormalizeStep(raw)

	if len(s.Fields) == 0 {
		return fmt.Errorf("normalize requires at least one field")
	}

	for field, strategy := range s.Fields {
		valid := false
		for _, v := range validNormalizeStrategies {
			if strategy == v {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("normalize field %q has unknown strategy %q: must be one of: %s", field, strategy, strings.Join(validNormalizeStrategies, ", "))
		}
	}

	return nil
}
