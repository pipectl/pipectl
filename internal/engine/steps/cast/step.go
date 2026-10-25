package cast

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

var (
	defaultTrueValues  = []string{"true", "t", "1", "yes", "y", "on"}
	defaultFalseValues = []string{"false", "f", "0", "no", "n", "off"}
)

type Field struct {
	Type        string
	Format      string
	TrueValues  []string
	FalseValues []string
}

type Step struct {
	payload.JSONCSVSupport
	Fields map[string]Field
}

func (s *Step) Name() string {
	return "cast"
}

func (s *Step) Execute(ctx *engine.ExecutionContext) error {
	for fieldPath, config := range s.Fields {
		ctx.Logger.Debug("  %s → %s", fieldPath, config.Type)
	}

	switch p := ctx.Payload.(type) {
	case payload.JSONRecordPayload:
		return s.executeJSON(p)
	case *payload.CSV:
		return s.executeCSV(p)
	default:
		return fmt.Errorf("unsupported payload type %T", ctx.Payload)
	}
}

func (s *Step) executeJSON(jsonPayload payload.JSONRecordPayload) error {
	for recordIndex, record := range jsonPayload.Records() {
		if record == nil {
			continue
		}

		for fieldPath, config := range s.Fields {
			value, assign, err := payload.ResolveFieldPath(record, fieldPath)
			if err != nil {
				return fmt.Errorf("cast field %q in record %d: %w", fieldPath, recordIndex+1, err)
			}

			casted, err := castValue(value, config)
			if err != nil {
				return fmt.Errorf("cast field %q in record %d: %w", fieldPath, recordIndex+1, err)
			}

			if err := assign(casted); err != nil {
				return fmt.Errorf("cast field %q in record %d: %w", fieldPath, recordIndex+1, err)
			}
		}
	}

	return nil
}

func (s *Step) executeCSV(csvPayload *payload.CSV) error {
	if len(csvPayload.Rows) == 0 {
		return nil
	}

	headerRow := csvPayload.Rows[0]
	fields := make([]string, 0, len(s.Fields))
	for field := range s.Fields {
		fields = append(fields, field)
	}
	colIndex, err := payload.FindColumnIndices(headerRow, fields)
	if err != nil {
		return fmt.Errorf("cast: %w", err)
	}

	for rowIndex, row := range csvPayload.Rows[1:] {
		for field, config := range s.Fields {
			i := colIndex[field]
			casted, err := castScalarValue(row[i], config)
			if err != nil {
				return fmt.Errorf("cast field %q in record %d: %w", field, rowIndex+1, err)
			}
			str, err := castToString(casted)
			if err != nil {
				return fmt.Errorf("cast field %q in record %d: %w", field, rowIndex+1, err)
			}
			row[i] = str
		}
	}

	return nil
}

func castValue(value interface{}, field Field) (interface{}, error) {
	if values, ok := value.([]interface{}); ok {
		casted := make([]interface{}, len(values))
		for i, item := range values {
			itemValue, err := castScalarValue(item, field)
			if err != nil {
				return nil, fmt.Errorf("array index %d: %w", i, err)
			}
			casted[i] = itemValue
		}
		return casted, nil
	}

	return castScalarValue(value, field)
}

func castScalarValue(value interface{}, field Field) (interface{}, error) {
	switch field.Type {
	case "int":
		return castToInt(value)
	case "float":
		return castToFloat(value)
	case "bool":
		return castToBool(value, field)
	case "time":
		return castToTime(value, field)
	case "string":
		return castToString(value)
	default:
		return nil, fmt.Errorf("unsupported target type %q", field.Type)
	}
}

func castToInt(value interface{}) (int, error) {
	switch typed := value.(type) {
	case int:
		return typed, nil
	case float64:
		return truncateFloat(typed)
	case string:
		trimmed := strings.TrimSpace(typed)
		parsed, err := strconv.Atoi(trimmed)
		if err == nil {
			return parsed, nil
		}

		floatValue, floatErr := strconv.ParseFloat(trimmed, 64)
		if floatErr != nil {
			return 0, fmt.Errorf("cannot cast %q to int: %w", typed, err)
		}
		return truncateFloat(floatValue)
	default:
		return 0, fmt.Errorf("cannot cast %T to int", value)
	}
}

func castToFloat(value interface{}) (float64, error) {
	switch typed := value.(type) {
	case int:
		return float64(typed), nil
	case float64:
		return typed, nil
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return 0, fmt.Errorf("cannot cast %q to float: %w", typed, err)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("cannot cast %T to float", value)
	}
}

func castToBool(value interface{}, field Field) (bool, error) {
	switch typed := value.(type) {
	case bool:
		return typed, nil
	case int:
		switch typed {
		case 1:
			return true, nil
		case 0:
			return false, nil
		default:
			return false, fmt.Errorf("cannot cast %d to bool", typed)
		}
	case float64:
		switch typed {
		case 1:
			return true, nil
		case 0:
			return false, nil
		default:
			return false, fmt.Errorf("cannot cast %v to bool", typed)
		}
	case string:
		normalized := strings.ToLower(strings.TrimSpace(typed))
		for _, candidate := range boolTrueValues(field) {
			if normalized == candidate {
				return true, nil
			}
		}
		for _, candidate := range boolFalseValues(field) {
			if normalized == candidate {
				return false, nil
			}
		}
		return false, fmt.Errorf("cannot cast %q to bool", typed)
	default:
		return false, fmt.Errorf("cannot cast %T to bool", value)
	}
}

func castToTime(value interface{}, field Field) (time.Time, error) {
	switch typed := value.(type) {
	case time.Time:
		return typed, nil
	case string:
		format := field.Format
		if format == "" {
			format = time.RFC3339
		}

		parsed, err := time.Parse(format, strings.TrimSpace(typed))
		if err != nil {
			return time.Time{}, fmt.Errorf("cannot cast %q to time with format %q: %w", typed, format, err)
		}
		return parsed, nil
	default:
		return time.Time{}, fmt.Errorf("cannot cast %T to time", value)
	}
}

func castToString(value interface{}) (string, error) {
	switch typed := value.(type) {
	case string:
		return typed, nil
	case bool:
		return strconv.FormatBool(typed), nil
	case int:
		return strconv.Itoa(typed), nil
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64), nil
	case time.Time:
		return typed.UTC().Format(time.RFC3339), nil
	default:
		return "", fmt.Errorf("cannot cast %T to string", value)
	}
}

func boolTrueValues(field Field) []string {
	if len(field.TrueValues) == 0 {
		return defaultTrueValues
	}
	return normalizeBoolValues(field.TrueValues)
}

func boolFalseValues(field Field) []string {
	if len(field.FalseValues) == 0 {
		return defaultFalseValues
	}
	return normalizeBoolValues(field.FalseValues)
}

func normalizeBoolValues(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		normalized = append(normalized, strings.ToLower(strings.TrimSpace(value)))
	}
	return normalized
}

func truncateFloat(value float64) (int, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("float value %v is not finite", value)
	}
	if value < float64(minInt()) || value > float64(maxInt()) {
		return 0, fmt.Errorf("float value %v overflows int", value)
	}
	return int(value), nil
}

func maxInt() int {
	return int(^uint(0) >> 1)
}

func minInt() int {
	return -maxInt() - 1
}
