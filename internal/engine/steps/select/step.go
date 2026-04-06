package _select

import (
	"fmt"
	"slices"
	"strings"

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
	case *payload.CSV, payload.JSONRecordPayload:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	context.Logger.Debug("  fields: [%s]", strings.Join(s.Fields, ", "))

	switch p := context.Payload.(type) {
	case payload.JSONRecordPayload:
		return s.selectJSON(p)
	case *payload.CSV:
		return s.selectCsv(p)
	default:
		return fmt.Errorf("select: unsupported payload type %T", context.Payload)
	}
}

func (s *Step) selectJSON(p payload.JSONRecordPayload) error {
	records := p.Records()
	for i, record := range records {
		selected := make(map[string]interface{}, len(s.Fields))
		for _, field := range s.Fields {
			if value, exists := record[field]; exists {
				selected[field] = value
			}
		}
		records[i] = selected
	}

	switch typed := p.(type) {
	case *payload.JSON:
		typed.Items = records
	case *payload.JSONL:
		typed.Items = records
	}

	return nil
}

func (s *Step) selectCsv(csvPayload *payload.CSV) error {
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
