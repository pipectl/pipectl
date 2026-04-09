package spec

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
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
	"limit":          func() Step { return &LimitStep{} },
	"log":            func() Step { return &LogStep{} },
	"count":          func() Step { return &CountStep{} },
	"http-transform": func() Step { return &HTTPTransformStep{} },
	"sort":           func() Step { return &SortStep{} },
}

func (w *StepWrapper) UnmarshalYAML(node ast.Node) error {
	var raw map[string]yaml.RawMessage
	if err := yaml.NodeToValue(node, &raw); err != nil {
		return err
	}

	if len(raw) != 1 {
		return wrapWithLine(node, fmt.Errorf("step must contain exactly one key"))
	}

	for key, value := range raw {
		factory, ok := stepRegistry[key]
		if !ok {
			return wrapWithLine(node, fmt.Errorf("unknown step type: %s", key))
		}

		step := factory()
		if err := yaml.Unmarshal(value, step); err != nil {
			return wrapWithLine(node, err)
		}
		if err := step.Validate(); err != nil {
			return wrapWithLine(node, err)
		}

		w.Step = step
	}

	return nil
}

func wrapWithLine(node ast.Node, err error) error {
	tk := node.GetToken()
	if tk != nil && tk.Position != nil {
		return fmt.Errorf("line %d: %w", tk.Position.Line, err)
	}
	return err
}
