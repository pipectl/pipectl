package payload

import (
	"fmt"
	"strconv"
	"strings"
)

type pathSegment struct {
	key   string
	index *int
}

// ParseFieldPath validates a field path such as "user.address" or "tags[0]".
func ParseFieldPath(path string) error {
	_, err := parseFieldPath(path)
	return err
}

// ResolveFieldPath looks up path in record and returns its value along with a
// function that replaces the value in place.
func ResolveFieldPath(record map[string]interface{}, path string) (interface{}, func(interface{}) error, error) {
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
			return nil, fmt.Errorf("field %q not found in record", path)
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
