package redact

import (
	"fmt"
	"slices"

	"github.com/shanebell/pipectl/internal/steps"
)

type RedactStep struct {
	Fields []string
}

func (s *RedactStep) Name() string {
	return "redact"
}

func (s *RedactStep) Supports(payload steps.Payload) bool {
	return payload.Type() == "json" || payload.Type() == "csv"
}

func (s *RedactStep) redactCsv(csvPayload *steps.CSVPayload) error {
	headerRow := csvPayload.Rows[0]
	toRedact := make([]bool, len(headerRow))
	for i, header := range headerRow {
		toRedact[i] = slices.Contains(s.Fields, header)
	}

	for _, row := range csvPayload.Rows[1:] {
		for i, value := range row {
			if toRedact[i] {
				fmt.Printf("- redacting field: %v, value: %v\n", headerRow[i], value)
				row[i] = "REDACTED"
			}
		}
	}

	return nil
}

// TODO only handles top-level fields, make this recursive
func (s *RedactStep) redactJson(jsonPayload *steps.JSONPayload) error {
	for k, v := range jsonPayload.Data {
		if slices.Contains(s.Fields, k) {
			fmt.Printf("- redacting field: %v, value: %v\n", k, v)
			jsonPayload.Data[k] = "REDACTED"
		}
	}

	return nil
}

func (s *RedactStep) Execute(context *steps.ExecutionContext) error {
	jsonPayload, jsonOk := context.Payload.(*steps.JSONPayload)
	if jsonOk {
		return s.redactJson(jsonPayload)
	}

	csvPayload, csvOk := context.Payload.(*steps.CSVPayload)
	if csvOk {
		return s.redactCsv(csvPayload)
	}

	return fmt.Errorf("%v requires either JSON or CSV payload, got %s", s.Name(), context.Payload.Type())
}
