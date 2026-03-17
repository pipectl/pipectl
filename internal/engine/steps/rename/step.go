package rename

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Fields map[string]string
}

func (s *Step) Name() string {
	return "rename"
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
	jsonPayload, jsonOk := context.Payload.(payload.JSONRecordPayload)
	if jsonOk {
		return s.renameJSON(jsonPayload)
	}

	csvPayload, csvOk := context.Payload.(*payload.CSV)
	if csvOk {
		return s.renameCSV(csvPayload)
	}

	return fmt.Errorf("%v received invalid payload type %v", s.Name(), context.Payload.Type())
}

func (s *Step) renameJSON(jsonPayload payload.JSONRecordPayload) error {
	for _, record := range jsonPayload.Records() {
		if record == nil {
			continue
		}

		for from, to := range s.Fields {
			value, ok := record[from]
			if !ok {
				continue
			}

			fmt.Printf("- renaming field: %v => %v\n", from, to)
			record[to] = value
			delete(record, from)
		}
	}

	return nil
}

func (s *Step) renameCSV(csvPayload *payload.CSV) error {
	if len(csvPayload.Rows) == 0 {
		return nil
	}

	headerRow := csvPayload.Rows[0]
	for i, header := range headerRow {
		to, ok := s.Fields[header]
		if !ok {
			continue
		}

		fmt.Printf("- renaming field: %v => %v\n", header, to)
		headerRow[i] = to
	}

	return nil
}
