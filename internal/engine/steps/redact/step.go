package redact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

type Step struct {
	payload.JSONCSVSupport
	Strategy string
	Fields   []string
}

func (s *Step) Name() string {
	return "redact"
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	strategy := s.Strategy
	if strategy == "" {
		strategy = "remove"
	}
	context.Logger.Debug("  fields: [%s] (%s)", strings.Join(s.Fields, ", "), strategy)

	switch p := context.Payload.(type) {
	case payload.JSONRecordPayload:
		return s.redactJson(p)
	case *payload.CSV:
		return s.redactCsv(p)
	default:
		return fmt.Errorf("unsupported payload type %T", context.Payload)
	}
}

func (s *Step) redactCsv(csvPayload *payload.CSV) error {
	headerRow := csvPayload.Rows[0]
	colIndex, err := payload.FindColumnIndices(headerRow, s.Fields)
	if err != nil {
		return fmt.Errorf("redact: %w", err)
	}

	toRedact := make([]bool, len(headerRow))
	for _, i := range colIndex {
		toRedact[i] = true
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

func (s *Step) redactJson(jsonPayload payload.JSONRecordPayload) error {
	for _, field := range s.Fields {
		if err := payload.ParseFieldPath(field); err != nil {
			return fmt.Errorf("redact: %w", err)
		}
	}

	for recordIndex, record := range jsonPayload.Records() {
		if record == nil {
			continue
		}

		for _, field := range s.Fields {
			v, assign, err := payload.ResolveFieldPath(record, field)
			if err != nil {
				return fmt.Errorf("redact: field %q in record %d: %w", field, recordIndex+1, err)
			}

			var redacted interface{}
			switch value := v.(type) {
			case nil:
				continue
			case string:
				if s.Strategy == "nullify" {
					redacted = nil
				} else {
					redacted = s.redactSingleValue(value)
				}
			default:
				if s.Strategy != "nullify" {
					return fmt.Errorf(
						"redact: field %q in record %d is %s; strategy %s only applies to strings — cast it to string first, or use strategy: nullify",
						field, recordIndex+1, describeType(value), s.strategyLabel(),
					)
				}
			}

			if err := assign(redacted); err != nil {
				return fmt.Errorf("redact: field %q in record %d: %w", field, recordIndex+1, err)
			}
		}
	}

	return nil
}

func (s *Step) strategyLabel() string {
	if s.Strategy == "" {
		return "REDACTED (default)"
	}
	return fmt.Sprintf("%q", s.Strategy)
}

func describeType(value interface{}) string {
	switch value.(type) {
	case bool:
		return "a bool"
	case map[string]interface{}:
		return "an object"
	case []interface{}:
		return "an array"
	case int, int64, float64:
		return "a number"
	default:
		return fmt.Sprintf("of type %T", value)
	}
}

func (s *Step) redactSingleValue(value string) string {
	switch {
	case s.Strategy == "nullify":
		return ""

	case s.Strategy == "mask":
		return strings.Repeat("*", len(value))

	case s.Strategy == "sha256":
		hash := sha256.New()
		hash.Write([]byte(value))
		return hex.EncodeToString(hash.Sum(nil))

	case strings.HasPrefix(s.Strategy, "partial-last"):
		n := partialN(s.Strategy)
		if n >= len(value) {
			return value
		}
		return strings.Repeat("*", len(value)-n) + value[len(value)-n:]

	case strings.HasPrefix(s.Strategy, "partial-first"):
		n := partialN(s.Strategy)
		if n >= len(value) {
			return value
		}
		return value[:n] + strings.Repeat("*", len(value)-n)

	default:
		return "REDACTED"
	}
}

// partialN extracts N from "partial-last:4" or "partial-first:4".
// Returns 4 when no suffix is present (bare "partial-last" / "partial-first").
func partialN(strategy string) int {
	if i := strings.Index(strategy, ":"); i >= 0 {
		n, _ := strconv.Atoi(strategy[i+1:])
		return n
	}
	return 4
}
