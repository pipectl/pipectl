package filter

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
)

type Step struct {
	Field string
	Value string
}

func (s *Step) Name() string {
	return "filter"
}

func (s *Step) Supports(p payload.Payload) bool {
	return p.Type() == payload.CSVType
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
			fmt.Printf("Excluding row: %v\n", row[0:len(headerRow)])
		}
	}

	csvPayload.Rows = filteredRows
	return nil
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	csvPayload, ok := context.Payload.(*payload.CSV)
	if !ok {
		return fmt.Errorf("%v requires CSV payload, got %s", s.Name(), context.Payload.Type())
	}

	s.filterCsv(csvPayload)

	return nil
}
