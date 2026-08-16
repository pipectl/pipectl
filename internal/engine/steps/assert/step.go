package assert

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

// FieldCheck holds a field name and comparison value for a per-record value
// assertion. Regex is set only when the check is used as FieldMatches.
type FieldCheck struct {
	Field string
	Value string
	Regex *regexp.Regexp
}

type Step struct {
	payload.JSONCSVSupport
	MinRecords    *int
	MaxRecords    *int
	RecordsEqual  *int
	FieldExists   string
	FieldEquals   *FieldCheck
	FieldContains *FieldCheck
	FieldMatches  *FieldCheck
	CaseSensitive bool
}

func (s *Step) Name() string {
	return "assert"
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	switch context.Payload.(type) {
	case payload.JSONRecordPayload, *payload.CSV:
	default:
		return fmt.Errorf("unsupported payload type %T", context.Payload)
	}

	recordCount := context.Payload.RecordCount()
	context.Logger.Debug("  records: %d", recordCount)

	if s.RecordsEqual != nil {
		context.Logger.Debug("  records-equal: %d", *s.RecordsEqual)
		if recordCount != *s.RecordsEqual {
			return fmt.Errorf("assert failed: records %d is not equal to expected %d", recordCount, *s.RecordsEqual)
		}
	}

	if s.MinRecords != nil {
		context.Logger.Debug("  min-records: >= %d", *s.MinRecords)
		if recordCount < *s.MinRecords {
			return fmt.Errorf("assert failed: records %d is less than minimum %d", recordCount, *s.MinRecords)
		}
	}

	if s.MaxRecords != nil {
		context.Logger.Debug("  max-records: <= %d", *s.MaxRecords)
		if recordCount > *s.MaxRecords {
			return fmt.Errorf("assert failed: records %d is greater than maximum %d", recordCount, *s.MaxRecords)
		}
	}

	if s.FieldExists != "" {
		context.Logger.Debug("  field-exists: %q", s.FieldExists)
		if !s.fieldExists(context.Payload) {
			return fmt.Errorf("assert failed: field %q does not exist", s.FieldExists)
		}
	}

	if s.FieldEquals != nil {
		context.Logger.Debug("  field-equals: %s == %q", s.FieldEquals.Field, s.FieldEquals.Value)
		if err := s.checkFieldEquals(context.Payload); err != nil {
			return err
		}
	}

	if s.FieldContains != nil {
		context.Logger.Debug("  field-contains: %s contains %q", s.FieldContains.Field, s.FieldContains.Value)
		if err := s.checkFieldContains(context.Payload); err != nil {
			return err
		}
	}

	if s.FieldMatches != nil {
		context.Logger.Debug("  field-matches: %s matches %q", s.FieldMatches.Field, s.FieldMatches.Value)
		if err := s.checkFieldMatches(context.Payload); err != nil {
			return err
		}
	}

	return nil
}

func (s *Step) checkFieldEquals(p payload.Payload) error {
	return s.forEachRecord(p, func(idx int, record map[string]interface{}) error {
		value, exists := record[s.FieldEquals.Field]
		if !exists {
			return fmt.Errorf("assert failed: field %q is missing from record %d", s.FieldEquals.Field, idx)
		}
		actual := fmt.Sprintf("%v", value)
		if !s.valuesEqual(actual, s.FieldEquals.Value) {
			return fmt.Errorf("assert failed: field %q in record %d is %q, want %q", s.FieldEquals.Field, idx, actual, s.FieldEquals.Value)
		}
		return nil
	})
}

func (s *Step) checkFieldContains(p payload.Payload) error {
	return s.forEachRecord(p, func(idx int, record map[string]interface{}) error {
		value, exists := record[s.FieldContains.Field]
		if !exists {
			return fmt.Errorf("assert failed: field %q is missing from record %d", s.FieldContains.Field, idx)
		}
		actual := fmt.Sprintf("%v", value)
		haystack, needle := actual, s.FieldContains.Value
		if !s.CaseSensitive {
			haystack, needle = strings.ToLower(haystack), strings.ToLower(needle)
		}
		if !strings.Contains(haystack, needle) {
			return fmt.Errorf("assert failed: field %q in record %d value %q does not contain %q", s.FieldContains.Field, idx, actual, s.FieldContains.Value)
		}
		return nil
	})
}

func (s *Step) checkFieldMatches(p payload.Payload) error {
	return s.forEachRecord(p, func(idx int, record map[string]interface{}) error {
		value, exists := record[s.FieldMatches.Field]
		if !exists {
			return fmt.Errorf("assert failed: field %q is missing from record %d", s.FieldMatches.Field, idx)
		}
		actual := fmt.Sprintf("%v", value)
		if !s.FieldMatches.Regex.MatchString(actual) {
			return fmt.Errorf("assert failed: field %q in record %d value %q does not match pattern %q", s.FieldMatches.Field, idx, actual, s.FieldMatches.Value)
		}
		return nil
	})
}

func (s *Step) valuesEqual(a, b string) bool {
	if s.CaseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
}

// forEachRecord visits every record in the payload, passing a 1-based index.
// CSV rows are converted to records via the header row; the header row
// itself is not visited.
func (s *Step) forEachRecord(p payload.Payload, fn func(idx int, record map[string]interface{}) error) error {
	switch v := p.(type) {
	case *payload.CSV:
		if len(v.Rows) == 0 {
			return nil
		}
		header := v.Rows[0]
		for i, row := range v.Rows[1:] {
			if err := fn(i+1, payload.CSVRowToRecord(header, row)); err != nil {
				return err
			}
		}
		return nil
	case payload.JSONRecordPayload:
		for i, record := range v.Records() {
			if err := fn(i+1, record); err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
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
