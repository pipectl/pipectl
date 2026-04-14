package normalize

import (
	"fmt"
	"strings"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

type Step struct {
	payload.JSONCSVSupport
	Fields map[string]string
}

func (s *Step) Name() string {
	return "normalize"
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	for field, strategy := range s.Fields {
		context.Logger.Debug("  %s: %s", field, strategy)
	}

	switch p := context.Payload.(type) {
	case payload.JSONRecordPayload:
		return s.normalizeJSON(p)
	case *payload.CSV:
		return s.normalizeCsv(p)
	default:
		return fmt.Errorf("unsupported payload type %T", context.Payload)
	}
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
	// TODO handle multiple strategies pipe separated eg: "trim|lower"
	fn, ok := strategies[strategy]
	if !ok {
		return value
	}
	return fn(value)
}

func (s *Step) normalizeJSON(jsonPayload payload.JSONRecordPayload) error {
	for _, record := range jsonPayload.Records() {
		if record == nil {
			continue
		}

		for field, strategy := range s.Fields {
			val, exists := record[field]
			if !exists {
				return fmt.Errorf("normalize: field %q not found in record", field)
			}
			if currentValue, ok := val.(string); ok {
				record[field] = s.normalizeValue(currentValue, strategy)
			}
		}
	}

	return nil
}

func (s *Step) normalizeCsv(csvPayload *payload.CSV) error {
	headerRow := csvPayload.Rows[0]
	matched := make(map[string]bool)
	normalizeFunctions := map[int]func(string) string{}
	strategyIndex := map[int]string{}
	for i, header := range headerRow {
		if strategy, ok := s.Fields[header]; ok {
			normalizeFunctions[i] = strategies[strategy]
			strategyIndex[i] = strategy
			matched[header] = true
		}
	}

	for field := range s.Fields {
		if !matched[field] {
			return fmt.Errorf("normalize: field %q not found in CSV headers", field)
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
