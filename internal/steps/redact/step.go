package redact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/shanebell/pipectl/internal/steps"
)

type Step struct {
	Strategy string
	Fields   []string
}

func (s *Step) Name() string {
	return "redact"
}

func (s *Step) Supports(payload steps.Payload) bool {
	return payload.Type() == "json" || payload.Type() == "csv"
}

func (s *Step) redactCsv(csvPayload *steps.CSVPayload) error {
	headerRow := csvPayload.Rows[0]
	toRedact := make([]bool, len(headerRow))
	for i, header := range headerRow {
		toRedact[i] = slices.Contains(s.Fields, header)
	}

	for _, row := range csvPayload.Rows[1:] {
		for i, value := range row {
			if toRedact[i] {
				fmt.Printf("- redacting field: %v, value: %v\n", headerRow[i], value)
				row[i] = s.redactSingleValue(value)
			}
		}
	}

	return nil
}

// TODO only handles top-level fields, make this recursive
// TODO can types other than strings be redacted? eg: changing from an int to "***" seems wrong and could break schema.
func (s *Step) redactJson(jsonPayload *steps.JSONPayload) error {
	for k, v := range jsonPayload.Data {
		if slices.Contains(s.Fields, k) {
			switch value := v.(type) {

			case string:
				fmt.Printf("- redacting field: %v, value: '%v'\n", k, v)
				jsonPayload.Data[k] = s.redactSingleValue(value)

			default:
				fmt.Printf("Cannot redact field %v of unsupported type %T\n", k, v)
			}
		}
	}

	return nil
}

func (s *Step) redactSingleValue(value string) string {
	fmt.Printf("- redacting text: %v with strategy %v\n", value, s.Strategy)
	switch s.Strategy {
	case "mask":
		return strings.Repeat("*", len(value))

	case "sha256":
		hash := sha256.New()
		hash.Write([]byte(value))
		return hex.EncodeToString(hash.Sum(nil))

	default:
		return "REDACTED"
	}
}

func (s *Step) Execute(context *steps.ExecutionContext) error {
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
