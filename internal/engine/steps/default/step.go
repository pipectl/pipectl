package _default

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Fields map[string]interface{}
}

func (s *Step) Name() string {
	return "default"
}

func (s *Step) Supports(p payload.Payload) bool {
	return p.Type() == payload.JSONType || p.Type() == payload.CSVType
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	jsonPayload, jsonOk := context.Payload.(*payload.JSON)
	if jsonOk {
		return s.defaultJSON(jsonPayload)
	}

	csvPayload, csvOk := context.Payload.(*payload.CSV)
	if csvOk {
		return s.defaultCSV(csvPayload)
	}

	return fmt.Errorf("%v received invalid payload type %v", s.Name(), context.Payload.Type())
}

func (s *Step) defaultJSON(jsonPayload *payload.JSON) error {
	if jsonPayload.Data == nil {
		jsonPayload.Data = map[string]interface{}{}
	}

	for key, value := range s.Fields {
		if _, exists := jsonPayload.Data[key]; exists {
			continue
		}

		fmt.Printf("- applying default: %v => %v\n", key, value)
		jsonPayload.Data[key] = value
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
			fmt.Printf("- adding default column: %v => %v\n", key, defaultValue)
			headerRow = append(headerRow, key)
			csvPayload.Rows[0] = headerRow
			fieldIndex = len(headerRow) - 1
			headerIndex[key] = fieldIndex

			for i := 1; i < len(csvPayload.Rows); i++ {
				csvPayload.Rows[i] = append(csvPayload.Rows[i], defaultValue)
			}

			continue
		}

		fmt.Printf("- applying default for CSV field: %v => %v\n", key, defaultValue)
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
