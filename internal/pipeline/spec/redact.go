package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type RedactStep struct {
	Strategy string   `yaml:"strategy"`
	Fields   []string `yaml:"fields"`
}

func (s *RedactStep) StepType() string {
	return "redact"
}

func (s *RedactStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *RedactStep) UnmarshalYAML(b []byte) error {
	type rawRedactStep RedactStep
	var raw rawRedactStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	*s = RedactStep(raw)
	return s.Validate()
}

func (s *RedactStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("redact requires at least one field")
	}

	switch s.Strategy {
	case "", "mask", "sha256":
	default:
		return fmt.Errorf("redact strategy must be one of: mask, sha256")
	}

	return nil
}
