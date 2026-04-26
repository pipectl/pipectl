package redact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
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
	switch {
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
