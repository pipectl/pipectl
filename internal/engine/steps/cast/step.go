package cast

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
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
	Fields map[string]Field
}

type pathSegment struct {
	key   string
	index *int
}

func (s *Step) Name() string {
	return "cast"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case payload.JSONRecordPayload:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(ctx *engine.ExecutionContext) error {
	jsonPayload, ok := ctx.Payload.(payload.JSONRecordPayload)
	if !ok {
		return fmt.Errorf("%v received invalid payload type %v", s.Name(), ctx.Payload.Type())
	}

	for fieldPath, config := range s.Fields {
		ctx.Logger.Debug("  %s → %s", fieldPath, config.Type)
	}

	for recordIndex, record := range jsonPayload.Records() {
		if record == nil {
			continue
		}

		for fieldPath, config := range s.Fields {
			value, assign, err := resolveFieldPath(record, fieldPath)
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

func resolveFieldPath(record map[string]interface{}, path string) (interface{}, func(interface{}) error, error) {
	segments, err := parseFieldPath(path)
	if err != nil {
		return nil, nil, err
	}

	var current interface{} = record
	for i := 0; i < len(segments)-1; i++ {
		current, err = getSegmentValue(current, segments[i], path)
		if err != nil {
			return nil, nil, err
		}
	}

	last := segments[len(segments)-1]
	value, err := getSegmentValue(current, last, path)
	if err != nil {
		return nil, nil, err
	}

	assign := func(newValue interface{}) error {
		return setSegmentValue(current, last, newValue, path)
	}

	return value, assign, nil
}

func parseFieldPath(path string) ([]pathSegment, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("field path must not be empty")
	}

	var segments []pathSegment
	i := 0
	for i < len(path) {
		if path[i] == '.' {
			return nil, fmt.Errorf("invalid field path %q", path)
		}

		start := i
		for i < len(path) && path[i] != '.' && path[i] != '[' {
			i++
		}
		if start != i {
			segments = append(segments, pathSegment{key: path[start:i]})
		}

		for i < len(path) && path[i] == '[' {
			i++
			indexStart := i
			for i < len(path) && path[i] != ']' {
				i++
			}
			if indexStart == i || i >= len(path) || path[i] != ']' {
				return nil, fmt.Errorf("invalid field path %q", path)
			}

			index, err := strconv.Atoi(path[indexStart:i])
			if err != nil || index < 0 {
				return nil, fmt.Errorf("invalid field path %q", path)
			}

			segments = append(segments, pathSegment{index: &index})
			i++
		}

		if i < len(path) {
			if path[i] != '.' {
				return nil, fmt.Errorf("invalid field path %q", path)
			}
			i++
			if i >= len(path) {
				return nil, fmt.Errorf("invalid field path %q", path)
			}
		}
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("invalid field path %q", path)
	}

	return segments, nil
}

func getSegmentValue(current interface{}, segment pathSegment, path string) (interface{}, error) {
	if segment.index == nil {
		object, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("path %q expected object before key %q, got %T", path, segment.key, current)
		}

		value, exists := object[segment.key]
		if !exists {
			return nil, fmt.Errorf("path %q missing key %q", path, segment.key)
		}

		return value, nil
	}

	array, ok := current.([]interface{})
	if !ok {
		return nil, fmt.Errorf("path %q expected array before index %d, got %T", path, *segment.index, current)
	}
	if *segment.index >= len(array) {
		return nil, fmt.Errorf("path %q index %d out of range", path, *segment.index)
	}

	return array[*segment.index], nil
}

func setSegmentValue(current interface{}, segment pathSegment, newValue interface{}, path string) error {
	if segment.index == nil {
		object, ok := current.(map[string]interface{})
		if !ok {
			return fmt.Errorf("path %q expected object before key %q, got %T", path, segment.key, current)
		}
		object[segment.key] = newValue
		return nil
	}

	array, ok := current.([]interface{})
	if !ok {
		return fmt.Errorf("path %q expected array before index %d, got %T", path, *segment.index, current)
	}
	if *segment.index >= len(array) {
		return fmt.Errorf("path %q index %d out of range", path, *segment.index)
	}

	array[*segment.index] = newValue
	return nil
}
