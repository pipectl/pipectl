package filter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

const (
	OpEquals      = "equals"
	OpNotEquals   = "not-equals"
	OpContains    = "contains"
	OpStartsWith  = "starts-with"
	OpEndsWith    = "ends-with"
	OpMatches     = "matches"
	OpGreaterThan = "greater-than"
	OpLessThan    = "less-than"
)

const (
	OnMissingExclude = "exclude"
	OnMissingInclude = "include"
	OnMissingError   = "error"
)

type Rule struct {
	Field         string
	Op            string
	Value         string
	NumericValue  float64
	CompiledRegex *regexp.Regexp
	OnMissing     string
	CaseSensitive bool
}

// Condition is either a leaf Rule or an All/Any group (recursive).
type Condition struct {
	Rule *Rule
	All  []*Condition
	Any  []*Condition
}

// evaluate reports whether the condition matches the given record.
// Missing fields on a leaf rule are handled per Rule.OnMissing: "exclude" (the
// default, and the zero value) treats them as non-matching, "include" treats
// them as matching, and "error" fails the step with an error naming the field.
func (c *Condition) evaluate(record map[string]interface{}) (bool, error) {
	if c.Rule != nil {
		value, exists := record[c.Rule.Field]
		if !exists {
			switch c.Rule.OnMissing {
			case OnMissingInclude:
				return true, nil
			case OnMissingError:
				return false, fmt.Errorf("filter: field %q is missing from record", c.Rule.Field)
			case OnMissingExclude, "":
				return false, nil
			default:
				return false, nil
			}
		}
		return c.Rule.matches(fmt.Sprintf("%v", value))
	}
	if len(c.All) > 0 {
		for _, sub := range c.All {
			ok, err := sub.evaluate(record)
			if err != nil || !ok {
				return false, err
			}
		}
		return true, nil
	}
	if len(c.Any) > 0 {
		for _, sub := range c.Any {
			ok, err := sub.evaluate(record)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	}
	return false, nil
}

type Step struct {
	payload.JSONCSVSupport
	Condition *Condition
}

func (s *Step) Name() string {
	return "filter"
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
		matched, err := s.Condition.evaluate(record)
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
		logger.Debug("  excluded %d records", excluded)
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
	if len(csvPayload.Rows) == 0 {
		return nil
	}

	headerRow := csvPayload.Rows[0]
	var filteredRows [][]string
	filteredRows = append(filteredRows, headerRow)
	excluded := 0

	for _, row := range csvPayload.Rows[1:] {
		record := payload.CSVRowToRecord(headerRow, row)
		matched, err := s.Condition.evaluate(record)
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
		logger.Debug("  excluded %d rows", excluded)
	}

	csvPayload.Rows = filteredRows
	return nil
}

func (r *Rule) equalValues(fieldValue string) bool {
	fField, errField := strconv.ParseFloat(strings.TrimSpace(fieldValue), 64)
	fThreshold, errThreshold := strconv.ParseFloat(strings.TrimSpace(r.Value), 64)
	if errField == nil && errThreshold == nil {
		return fField == fThreshold
	}
	if r.CaseSensitive {
		return fieldValue == r.Value
	}
	return strings.EqualFold(fieldValue, r.Value)
}

func (r *Rule) matches(value string) (bool, error) {
	switch r.Op {
	case OpEquals:
		return r.equalValues(value), nil
	case OpNotEquals:
		return !r.equalValues(value), nil
	case OpContains:
		if r.CaseSensitive {
			return strings.Contains(value, r.Value), nil
		}
		return strings.Contains(strings.ToLower(value), strings.ToLower(r.Value)), nil
	case OpStartsWith:
		if r.CaseSensitive {
			return strings.HasPrefix(value, r.Value), nil
		}
		return strings.HasPrefix(strings.ToLower(value), strings.ToLower(r.Value)), nil
	case OpEndsWith:
		if r.CaseSensitive {
			return strings.HasSuffix(value, r.Value), nil
		}
		return strings.HasSuffix(strings.ToLower(value), strings.ToLower(r.Value)), nil
	case OpMatches:
		return r.CompiledRegex.MatchString(value), nil
	case OpGreaterThan:
		f, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return false, fmt.Errorf("filter: field %q value %q is not a number", r.Field, value)
		}
		return f > r.NumericValue, nil
	case OpLessThan:
		f, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return false, fmt.Errorf("filter: field %q value %q is not a number", r.Field, value)
		}
		return f < r.NumericValue, nil
	default:
		return false, nil
	}
}
