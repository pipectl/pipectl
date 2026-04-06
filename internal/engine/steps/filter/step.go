package filter

import (
	"fmt"
	"strings"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

const (
	OpEquals     = "equals"
	OpNotEquals  = "not-equals"
	OpContains   = "contains"
	OpStartsWith = "starts-with"
)

type Step struct {
	Field string
	Op    string
	Value string
}

func (s *Step) Name() string {
	return "filter"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case *payload.CSV, payload.JSONRecordPayload:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	switch p := context.Payload.(type) {
	case payload.JSONRecordPayload:
		return s.filterJSON(p, context.Logger)
	case *payload.CSV:
		return s.filterCsv(p, context.Logger)
	default:
		return fmt.Errorf("filter: unsupported payload type %T", context.Payload)
	}
}

func (s *Step) filterJSON(p payload.JSONRecordPayload, logger *engine.Logger) error {
	records := p.Records()
	filtered := records[:0]
	excluded := 0

	for _, record := range records {
		value, exists := record[s.Field]
		if !exists || !s.matches(fmt.Sprintf("%v", value)) {
			excluded++
			continue
		}
		filtered = append(filtered, record)
	}

	if excluded > 0 {
		logger.Debug("  excluded %d records (%s %s %q)", excluded, s.Field, s.Op, s.Value)
	}

	switch p := p.(type) {
	case *payload.JSON:
		p.Items = filtered
	case *payload.JSONL:
		p.Items = filtered
	}

	return nil
}

func (s *Step) filterCsv(csvPayload *payload.CSV, logger *engine.Logger) error {
	headerRow := csvPayload.Rows[0]
	colIndex := -1
	for i, header := range headerRow {
		if s.Field == header {
			colIndex = i
			break
		}
	}

	var filteredRows [][]string
	filteredRows = append(filteredRows, headerRow)
	excluded := 0

	for _, row := range csvPayload.Rows[1:] {
		if colIndex == -1 || !s.matches(row[colIndex]) {
			excluded++
			continue
		}
		filteredRows = append(filteredRows, row)
	}

	if excluded > 0 {
		logger.Debug("  excluded %d rows (%s %s %q)", excluded, s.Field, s.Op, s.Value)
	}

	csvPayload.Rows = filteredRows
	return nil
}

func (s *Step) matches(value string) bool {
	switch s.Op {
	case OpEquals:
		return value == s.Value
	case OpNotEquals:
		return value != s.Value
	case OpContains:
		return strings.Contains(value, s.Value)
	case OpStartsWith:
		return strings.HasPrefix(value, s.Value)
	default:
		return false
	}
}
