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
		return s.filterJSON(p)
	case *payload.CSV:
		return s.filterCsv(p)
	default:
		return fmt.Errorf("filter: unsupported payload type %T", context.Payload)
	}
}

func (s *Step) filterJSON(p payload.JSONRecordPayload) error {
	records := p.Records()
	filtered := records[:0]
	for _, record := range records {
		value, exists := record[s.Field]
		if !exists {
			fmt.Printf("- excluding record: field %q not found\n", s.Field)
			continue
		}
		str := fmt.Sprintf("%v", value)
		if s.matches(str) {
			filtered = append(filtered, record)
		} else {
			fmt.Printf("- excluding record: %v %v %v\n", s.Field, s.Op, s.Value)
		}
	}

	switch p := p.(type) {
	case *payload.JSON:
		p.Items = filtered
	case *payload.JSONL:
		p.Items = filtered
	}

	return nil
}

func (s *Step) filterCsv(csvPayload *payload.CSV) error {
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

	for _, row := range csvPayload.Rows[1:] {
		if colIndex == -1 || !s.matches(row[colIndex]) {
			fmt.Printf("- excluding row: %v\n", row[0:len(headerRow)])
			continue
		}
		filteredRows = append(filteredRows, row)
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
