package normalize

import (
	"fmt"
	"strings"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Fields map[string]string
}

func (s *Step) Name() string {
	return "normalize"
}

func (s *Step) Supports(p payload.Payload) bool {
	return p.Type() == payload.JSONType || p.Type() == payload.CSVType
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	jsonPayload, jsonOk := context.Payload.(*payload.JSON)
	if jsonOk {
		return s.normalizeJSON(jsonPayload)
	}

	csvPayload, csvOk := context.Payload.(*payload.CSV)
	if csvOk {
		return s.normalizeCsv(csvPayload)
	}

	return fmt.Errorf("%v received invalid payload type %v", s.Name(), context.Payload.Type())
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

func (s *Step) normalizeJSON(jsonPayload *payload.JSON) error {
	for key, _ := range jsonPayload.Data {
		if strategy, needsNormalizing := s.Fields[key]; needsNormalizing {
			if currentValue, ok := jsonPayload.Data[key].(string); ok {
				jsonPayload.Data[key] = s.normalizeValue(currentValue, strategy)
			}
		}
	}

	return nil
}

func (s *Step) normalizeCsv(csvPayload *payload.CSV) error {
	headerRow := csvPayload.Rows[0]
	normalizeFunctions := map[int]func(string) string{}
	strategyIndex := map[int]string{}
	for i, header := range headerRow {
		strategy, ok := s.Fields[header]
		if ok {
			normalizeFunctions[i] = strategies[strategy]
			strategyIndex[i] = strategy
		}
	}

	for _, row := range csvPayload.Rows[1:] {
		for i, value := range row {
			if strategy, ok := strategyIndex[i]; ok {
				row[i] = s.normalizeValue(value, strategy)
			}
		}
	}

	return nil
}
