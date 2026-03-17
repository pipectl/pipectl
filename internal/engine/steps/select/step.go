package _select

import (
	"fmt"
	"slices"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Fields []string
}

func (s *Step) Name() string {
	return "select"
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
	return s.selectCsv(csvPayload)
}

func (s *Step) selectCsv(csvPayload *payload.CSV) error {
	fmt.Printf("- selecting fields: %v\n", s.Fields)

	headerRow := csvPayload.Rows[0]
	toSelect := make([]bool, len(headerRow))
	for i, header := range headerRow {
		toSelect[i] = slices.Contains(s.Fields, header)
	}

	for i, row := range csvPayload.Rows {
		var selectedRow []string
		for j, value := range row {
			if toSelect[j] {
				selectedRow = append(selectedRow, value)
			}
		}
		csvPayload.Rows[i] = selectedRow
	}

	return nil
}
