package spec

import (
	"fmt"
	"regexp"
	"strings"

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
		if err := yaml.UnmarshalWithOptions(value, step, yaml.DisallowUnknownField()); err != nil {
			return wrapWithLine(node, err)
		}
		if err := step.Validate(); err != nil {
			return wrapWithLine(node, err)
		}

		w.Step = step
	}

	return nil
}

// goccyPosition matches the [line:col] prefix goccy prepends to error messages.
var goccyPosition = regexp.MustCompile(`^\[\d+:\d+] `)

// goTypePath matches goccy's internal Go type language in type mismatch errors.
var goTypePath = regexp.MustCompile(`Go (?:struct field \S+ of type|value of type) `)

// typeMismatch matches the cleaned "cannot unmarshal X into Y" form.
var typeMismatch = regexp.MustCompile(`cannot unmarshal \S+ into (\S+)`)

// friendlyTypeNames maps Go type names to plain English descriptions.
var friendlyTypeNames = map[string]string{
	"bool":    "true or false",
	"int":     "a number",
	"int8":    "a number",
	"int16":   "a number",
	"int32":   "a number",
	"int64":   "a number",
	"uint":    "a number",
	"uint8":   "a number",
	"uint16":  "a number",
	"uint32":  "a number",
	"uint64":  "a number",
	"float32": "a number",
	"float64": "a number",
	"string":  "text",
}

func wrapWithLine(node ast.Node, err error) error {
	msg := friendlyError(err.Error())
	tk := node.GetToken()
	if tk != nil && tk.Position != nil {
		return fmt.Errorf("line %d: %s", tk.Position.Line, msg)
	}
	return fmt.Errorf("%s", msg)
}

func friendlyError(msg string) string {
	// Strip the [line:col] prefix goccy adds to its own errors — we add our own.
	msg = goccyPosition.ReplaceAllString(msg, "")

	// Strip goccy's internal Go type path to get "cannot unmarshal X into Y".
	msg = goTypePath.ReplaceAllString(msg, "")

	// Replace the remaining technical form with plain English.
	if m := typeMismatch.FindStringSubmatch(msg); m != nil {
		typeName := strings.TrimPrefix(m[1], "*") // handle pointer types e.g. *bool
		if friendly, ok := friendlyTypeNames[typeName]; ok {
			return "expected " + friendly
		}
	}

	return msg
}
