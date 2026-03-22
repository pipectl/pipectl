package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type ConvertStep struct {
	Format string `yaml:"format"`
}

func (s *ConvertStep) StepType() string {
	return "convert"
}

func (s *ConvertStep) String() string {
	return fmt.Sprintf("[%s] format=%q", s.StepType(), s.Format)
}

func (s *ConvertStep) UnmarshalYAML(b []byte) error {
	type rawConvertStep ConvertStep
	var raw rawConvertStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}

	*s = ConvertStep(raw)

	switch s.Format {
	case "json", "jsonl", "csv":
		return nil
	default:
		return fmt.Errorf("convert format must be one of: json, jsonl, csv")
	}
}
