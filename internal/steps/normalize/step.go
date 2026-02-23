package normalize

import (
	"fmt"
	"strings"

	"github.com/shanebell/pipectl/internal/steps"
)

type Step struct {
	Fields map[string]string
}

func (s *Step) Name() string {
	return "normalize"
}

func (s *Step) Supports(payload steps.Payload) bool {
	return payload.Type() == "json" || payload.Type() == "csv"
}

var strategies = map[string]func(string) string{
	"lower": strings.ToLower,
	"upper": strings.ToUpper,
	"trim":  strings.TrimSpace,
	"trim-left": func(s string) string {
		return strings.TrimLeft(s, " \t\n\r")
	},
	"trim-right": func(s string) string {
		return strings.TrimRight(s, " \t\n\r")
	},
	"collapse-spaces": func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	},
	"capitalize": func(s string) string {
		if len(s) == 0 {
			return s
		}
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	},
}

func (s *Step) normalizeValue(value string, strategy string) string {
	fmt.Printf("Normalizing value: '%v' with strategy: %v\n", value, strategy)

	// TODO handle multiple strategies pipe separated eg: "trim|lower"
	fn, ok := strategies[strategy]
	if !ok {
		return value
	}
	return fn(value)
}

func (s *Step) Execute(context *steps.ExecutionContext) error {
	jsonPayload, ok := context.Payload.(*steps.JSONPayload)
	if !ok {
		return fmt.Errorf("%v requires JSON payload, got %s", s.Name(), context.Payload.Type())
	}

	for key, _ := range jsonPayload.Data {
		if strategy, needsNormalizing := s.Fields[key]; needsNormalizing {
			if currentValue, ok := jsonPayload.Data[key].(string); ok {
				jsonPayload.Data[key] = s.normalizeValue(currentValue, strategy)
			}
		}
	}

	// TODO handle CSV

	return nil
}
