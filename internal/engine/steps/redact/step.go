package redact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Strategy string
	Fields   []string
}

func (s *Step) Name() string {
	return "redact"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case payload.JSONRecordPayload, *payload.CSV:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	strategy := s.Strategy
	if strategy == "" {
		strategy = "remove"
	}
	context.Logger.Debug("  fields: [%s] (%s)", strings.Join(s.Fields, ", "), strategy)

	jsonPayload, jsonOk := context.Payload.(payload.JSONRecordPayload)
	if jsonOk {
		return s.redactJson(jsonPayload)
	}

	csvPayload, csvOk := context.Payload.(*payload.CSV)
	if csvOk {
		return s.redactCsv(csvPayload)
	}

	return fmt.Errorf("unsupported payload type %T", context.Payload)
}

func (s *Step) redactCsv(csvPayload *payload.CSV) error {
	headerRow := csvPayload.Rows[0]
	toRedact := make([]bool, len(headerRow))
	matched := make(map[string]bool)
	for i, header := range headerRow {
		if slices.Contains(s.Fields, header) {
			toRedact[i] = true
			matched[header] = true
		}
	}

	for _, field := range s.Fields {
		if !matched[field] {
			return fmt.Errorf("redact: field %q not found in CSV headers", field)
		}
	}

	for _, row := range csvPayload.Rows[1:] {
		for i, value := range row {
			if toRedact[i] {
				row[i] = s.redactSingleValue(value)
			}
		}
	}

	return nil
}

// TODO only handles top-level fields, make this recursive
// TODO can types other than strings be redacted? eg: changing from an int to "***" seems wrong and could break schema.
func (s *Step) redactJson(jsonPayload payload.JSONRecordPayload) error {
	for _, record := range jsonPayload.Records() {
		if record == nil {
			continue
		}

		for _, field := range s.Fields {
			v, exists := record[field]
			if !exists {
				return fmt.Errorf("redact: field %q not found in record", field)
			}
			if value, ok := v.(string); ok {
				record[field] = s.redactSingleValue(value)
			}
		}
	}

	return nil
}

func (s *Step) redactSingleValue(value string) string {
	switch s.Strategy {
	case "mask":
		return strings.Repeat("*", len(value))

	case "sha256":
		hash := sha256.New()
		hash.Write([]byte(value))
		return hex.EncodeToString(hash.Sum(nil))

	default:
		return "REDACTED"
	}
}
