package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

const defaultLogSample = 10

type LogStep struct {
	Message string `yaml:"message"`
	Count   *bool  `yaml:"count"`
	Sample  *int   `yaml:"sample"`
}

func (s *LogStep) StepType() string {
	return "log"
}

func (s *LogStep) String() string {
	return fmt.Sprintf("[%s] message=%q count=%v sample=%v", s.StepType(), s.Message, s.Count, s.Sample)
}

func (s *LogStep) UnmarshalYAML(b []byte) error {
	type rawLogStep LogStep
	sample := defaultLogSample
	raw := rawLogStep{Sample: &sample}
	if err := yaml.UnmarshalWithOptions(b, &raw, yaml.DisallowUnknownField()); err != nil {
		return err
	}
	*s = LogStep(raw)
	return s.Validate()
}

func (s *LogStep) Validate() error { return nil }
