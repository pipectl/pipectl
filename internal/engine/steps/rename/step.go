package rename

import (
	"fmt"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
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
	for from, to := range s.Fields {
		context.Logger.Debug("  %s → %s", from, to)
	}

	switch p := context.Payload.(type) {
	case payload.JSONRecordPayload:
		return s.renameJSON(p)
	case *payload.CSV:
		return s.renameCSV(p)
	default:
		return fmt.Errorf("unsupported payload type %T", context.Payload)
	}
}

func (s *Step) renameJSON(jsonPayload payload.JSONRecordPayload) error {
	for _, record := range jsonPayload.Records() {
		if record == nil {
			continue
		}

		for from, to := range s.Fields {
			if _, exists := record[from]; !exists {
				return fmt.Errorf("rename: field %q not found in record", from)
			}
			if _, exists := record[to]; exists {
				return fmt.Errorf("rename: target field %q already exists in record", to)
			}
		}

		for from, to := range s.Fields {
			record[to] = record[from]
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
	headerIndex := payload.HeaderIndex(headerRow)

	for from, to := range s.Fields {
		if _, exists := headerIndex[from]; !exists {
			return fmt.Errorf("rename: field %q not found in CSV headers", from)
		}
		if _, exists := headerIndex[to]; exists {
			return fmt.Errorf("rename: target field %q already exists in CSV headers", to)
		}
	}

	for from, to := range s.Fields {
		headerRow[headerIndex[from]] = to
	}

	return nil
}
