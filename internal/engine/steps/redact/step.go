package redact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Strategy string
	Fields   []string
}

func (s *Step) Name() string {
	return "redact"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case payload.RecordPayload, *payload.CSV:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	recordPayload, recordOk := context.Payload.(payload.RecordPayload)
	if recordOk {
		return s.redactJson(recordPayload)
	}

	csvPayload, csvOk := context.Payload.(*payload.CSV)
	if csvOk {
		return s.redactCsv(csvPayload)
	}

	return fmt.Errorf("%v received invalid payload type %v", s.Name(), context.Payload.Type())
}

func (s *Step) redactCsv(csvPayload *payload.CSV) error {
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
func (s *Step) redactJson(recordPayload payload.RecordPayload) error {
	for _, record := range recordPayload.GetRecords() {
		if record == nil {
			continue
		}

		for k, v := range record {
			if slices.Contains(s.Fields, k) {
				switch value := v.(type) {

				case string:
					fmt.Printf("- redacting field: %v, value: '%v'\n", k, v)
					record[k] = s.redactSingleValue(value)

				default:
					fmt.Printf("Cannot redact field %v of unsupported type %T\n", k, v)
				}
			}
		}
	}

	return nil
}

func (s *Step) redactSingleValue(value string) string {
	// TODO print something useful when strategy is not defined
	fmt.Printf("- redacting text: %v with strategy %s\n", value, s.Strategy)
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
