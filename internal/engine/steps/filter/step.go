package filter

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Field string
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
		if fmt.Sprintf("%v", value) == s.Value {
			filtered = append(filtered, record)
		} else {
			fmt.Printf("- excluding record: %v = %v\n", s.Field, value)
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
	toFilter := make([]*string, len(headerRow))
	for i, header := range headerRow {
		if s.Field == header {
			toFilter[i] = &s.Value
		} else {
			toFilter[i] = nil
		}

	}

	var filteredRows [][]string
	filteredRows = append(filteredRows, headerRow)

	for _, row := range csvPayload.Rows[1:] {
		match := true
		for i, value := range row {
			filterValue := toFilter[i]
			if filterValue != nil {
				match = value == *filterValue
			}
		}

		if match {
			filteredRows = append(filteredRows, row)
		} else {
			fmt.Printf("- excluding row: %v\n", row[0:len(headerRow)])
		}
	}

	csvPayload.Rows = filteredRows
	return nil
}
