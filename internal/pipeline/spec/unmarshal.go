package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

var stepRegistry = map[string]func() Step{
	"validate-json":  func() Step { return &ValidateJSONStep{} },
	"normalize":      func() Step { return &NormalizeStep{} },
	"default":        func() Step { return &DefaultStep{} },
	"cast":           func() Step { return &CastStep{} },
	"convert":        func() Step { return &ConvertStep{} },
	"assert":         func() Step { return &AssertStep{} },
	"rename":         func() Step { return &RenameStep{} },
	"redact":         func() Step { return &RedactStep{} },
	"select":         func() Step { return &SelectStep{} },
	"filter":         func() Step { return &FilterStep{} },
	"log":            func() Step { return &LogStep{} },
	"count":          func() Step { return &CountStep{} },
	"http-transform": func() Step { return &HTTPTransformStep{} },
}

func (w *StepWrapper) UnmarshalYAML(b []byte) error {
	var raw map[string]yaml.RawMessage
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}

	if len(raw) != 1 {
		return fmt.Errorf("step must contain exactly one key")
	}

	for key, value := range raw {
		factory, ok := stepRegistry[key]
		if !ok {
			return fmt.Errorf("unknown step type: %s", key)
		}

		step := factory()
		if err := yaml.Unmarshal(value, step); err != nil {
			return err
		}

		w.Step = step
	}

	return nil
}
