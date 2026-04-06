package filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

const (
	OpEquals      = "equals"
	OpNotEquals   = "not-equals"
	OpContains    = "contains"
	OpStartsWith  = "starts-with"
	OpEndsWith    = "ends-with"
	OpGreaterThan = "greater-than"
	OpLessThan    = "less-than"
)

type Step struct {
	Field        string
	Op           string
	Value        string
	NumericValue float64
}

func (s *Step) Name() string {
	return "filter"
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
	switch p := context.Payload.(type) {
	case payload.JSONRecordPayload:
		return s.filterJSON(p, context.Logger)
	case *payload.CSV:
		return s.filterCsv(p, context.Logger)
	default:
		return fmt.Errorf("unsupported payload type %T", context.Payload)
	}
}

func (s *Step) filterJSON(p payload.JSONRecordPayload, logger *engine.Logger) error {
	records := p.Records()
	filtered := records[:0]
	excluded := 0

	for _, record := range records {
		value, exists := record[s.Field]
		if !exists {
			excluded++
			continue
		}
		matched, err := s.matches(fmt.Sprintf("%v", value))
		if err != nil {
			return err
		}
		if !matched {
			excluded++
			continue
		}
		filtered = append(filtered, record)
	}

	if excluded > 0 {
		logger.Debug("  excluded %d records (%s %s %v)", excluded, s.Field, s.Op, s.logValue())
	}

	switch p := p.(type) {
	case *payload.JSON:
		p.Items = filtered
	case *payload.JSONL:
		p.Items = filtered
	}

	return nil
}

func (s *Step) filterCsv(csvPayload *payload.CSV, logger *engine.Logger) error {
	headerRow := csvPayload.Rows[0]
	colIndex := -1
	for i, header := range headerRow {
		if s.Field == header {
			colIndex = i
			break
		}
	}

	var filteredRows [][]string
	filteredRows = append(filteredRows, headerRow)
	excluded := 0

	for _, row := range csvPayload.Rows[1:] {
		if colIndex == -1 {
			excluded++
			continue
		}
		matched, err := s.matches(row[colIndex])
		if err != nil {
			return err
		}
		if !matched {
			excluded++
			continue
		}
		filteredRows = append(filteredRows, row)
	}

	if excluded > 0 {
		logger.Debug("  excluded %d rows (%s %s %v)", excluded, s.Field, s.Op, s.logValue())
	}

	csvPayload.Rows = filteredRows
	return nil
}

func (s *Step) equalValues(fieldValue string) bool {
	fField, errField := strconv.ParseFloat(strings.TrimSpace(fieldValue), 64)
	fThreshold, errThreshold := strconv.ParseFloat(strings.TrimSpace(s.Value), 64)
	if errField == nil && errThreshold == nil {
		return fField == fThreshold
	}
	return fieldValue == s.Value
}

func (s *Step) logValue() interface{} {
	switch s.Op {
	case OpGreaterThan, OpLessThan:
		return s.NumericValue
	default:
		return fmt.Sprintf("%q", s.Value)
	}
}

func (s *Step) matches(value string) (bool, error) {
	switch s.Op {
	case OpEquals:
		return s.equalValues(value), nil
	case OpNotEquals:
		return !s.equalValues(value), nil
	case OpContains:
		return strings.Contains(value, s.Value), nil
	case OpStartsWith:
		return strings.HasPrefix(value, s.Value), nil
	case OpEndsWith:
		return strings.HasSuffix(value, s.Value), nil
	case OpGreaterThan:
		f, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return false, fmt.Errorf("filter: field %q value %q is not a number", s.Field, value)
		}
		return f > s.NumericValue, nil
	case OpLessThan:
		f, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return false, fmt.Errorf("filter: field %q value %q is not a number", s.Field, value)
		}
		return f < s.NumericValue, nil
	default:
		return false, nil
	}
}
