package spec

import "fmt"

type ConvertStep struct {
	Format string `yaml:"format"`
}

func (s *ConvertStep) StepType() string {
	return "convert"
}

func (s *ConvertStep) String() string {
	return fmt.Sprintf("[%s] format=%q", s.StepType(), s.Format)
}

func (s *ConvertStep) Validate() error {
	switch s.Format {
	case "json", "jsonl", "csv":
		return nil
	case "":
		return fmt.Errorf("convert format is required: must be one of: json, jsonl, csv")
	default:
		return fmt.Errorf("convert format %q is invalid: must be one of: json, jsonl, csv", s.Format)
	}
}
