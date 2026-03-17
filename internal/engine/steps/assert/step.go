package assert

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	MinRecords   *int
	MaxRecords   *int
	RecordsEqual *int
	FieldExists  string
}

func (s *Step) Name() string {
	return "assert"
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
	switch context.Payload.(type) {
	case payload.JSONRecordPayload, *payload.CSV:
	default:
		return fmt.Errorf("%v received invalid payload type %v", s.Name(), context.Payload.Type())
	}

	recordCount := context.Payload.RecordCount()
	fmt.Printf("- assert records: actual=%d\n", recordCount)

	if s.RecordsEqual != nil {
		fmt.Printf("- assert records-equal: expected=%d\n", *s.RecordsEqual)
		if recordCount != *s.RecordsEqual {
			return fmt.Errorf("assert failed: records %d is not equal to expected %d", recordCount, *s.RecordsEqual)
		}
	}

	if s.MinRecords != nil {
		fmt.Printf("- assert min-records: expected >= %d\n", *s.MinRecords)
	}
	if s.MinRecords != nil && recordCount < *s.MinRecords {
		return fmt.Errorf("assert failed: records %d is less than minimum %d", recordCount, *s.MinRecords)
	}

	if s.MaxRecords != nil {
		fmt.Printf("- assert max-records: expected <= %d\n", *s.MaxRecords)
	}
	if s.MaxRecords != nil && recordCount > *s.MaxRecords {
		return fmt.Errorf("assert failed: records %d is greater than maximum %d", recordCount, *s.MaxRecords)
	}

	if s.FieldExists != "" {
		fmt.Printf("- assert field-exists: %q\n", s.FieldExists)
	}
	if s.FieldExists != "" && !s.fieldExists(context.Payload) {
		return fmt.Errorf("assert failed: field %q does not exist", s.FieldExists)
	}

	return nil
}

func (s *Step) fieldExists(p payload.Payload) bool {
	switch v := p.(type) {
	case *payload.CSV:
		if len(v.Rows) == 0 {
			return false
		}
		for _, header := range v.Rows[0] {
			if header == s.FieldExists {
				return true
			}
		}
		return false
	case payload.JSONRecordPayload:
		for _, record := range v.Records() {
			if _, exists := record[s.FieldExists]; exists {
				return true
			}
		}
		return false
	default:
		return false
	}
}
