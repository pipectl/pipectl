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
	case *payload.CSV:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	csvPayload := context.Payload.(*payload.CSV)
	return s.filterCsv(csvPayload)
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
