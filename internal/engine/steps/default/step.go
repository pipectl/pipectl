package _default

import (
	"fmt"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

type Step struct {
	Fields map[string]interface{}
}

func (s *Step) Name() string {
	return "default"
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
	for key, value := range s.Fields {
		context.Logger.Debug("  %s: %v", key, value)
	}

	jsonPayload, jsonOk := context.Payload.(payload.JSONRecordPayload)
	if jsonOk {
		return s.defaultJSON(jsonPayload)
	}

	csvPayload, csvOk := context.Payload.(*payload.CSV)
	if csvOk {
		return s.defaultCSV(csvPayload)
	}

	return fmt.Errorf("unsupported payload type %T", context.Payload)
}

func (s *Step) defaultJSON(jsonPayload payload.JSONRecordPayload) error {
	records := jsonPayload.Records()
	if len(records) == 0 {
		return fmt.Errorf("%s requires at least one JSON record", s.Name())
	}

	for _, record := range records {
		if record == nil {
			continue
		}

		for key, value := range s.Fields {
			if _, exists := record[key]; exists {
				continue
			}
			record[key] = value
		}
	}

	return nil
}

func (s *Step) defaultCSV(csvPayload *payload.CSV) error {
	if len(csvPayload.Rows) == 0 {
		return nil
	}

	headerRow := csvPayload.Rows[0]
	headerIndex := make(map[string]int, len(headerRow))
	for i, header := range headerRow {
		headerIndex[header] = i
	}

	for key, value := range s.Fields {
		defaultValue := stringify(value)
		fieldIndex, exists := headerIndex[key]
		if !exists {
			headerRow = append(headerRow, key)
			csvPayload.Rows[0] = headerRow
			fieldIndex = len(headerRow) - 1
			headerIndex[key] = fieldIndex

			for i := 1; i < len(csvPayload.Rows); i++ {
				csvPayload.Rows[i] = append(csvPayload.Rows[i], defaultValue)
			}

			continue
		}

		for i := 1; i < len(csvPayload.Rows); i++ {
			row := csvPayload.Rows[i]
			if len(row) <= fieldIndex {
				row = append(row, make([]string, fieldIndex-len(row)+1)...)
				csvPayload.Rows[i] = row
			}

			if row[fieldIndex] == "" {
				row[fieldIndex] = defaultValue
			}
		}
	}

	return nil
}

func stringify(value interface{}) string {
	if value == nil {
		return ""
	}

	return fmt.Sprint(value)
}
