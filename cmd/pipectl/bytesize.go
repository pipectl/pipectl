package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	kilobyte = 1024
	megabyte = 1024 * kilobyte
	gigabyte = 1024 * megabyte
)

// parseByteSize parses a human-readable size string such as "256MB", "64KB",
// "1GB", or a bare integer (interpreted as bytes) into a byte count.
func parseByteSize(s string) (int64, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return 0, fmt.Errorf("byte size %q: empty value", s)
	}

	upper := strings.ToUpper(trimmed)
	var unit int64 = 1
	numPart := upper
	switch {
	case strings.HasSuffix(upper, "GB"):
		unit = gigabyte
		numPart = upper[:len(upper)-2]
	case strings.HasSuffix(upper, "MB"):
		unit = megabyte
		numPart = upper[:len(upper)-2]
	case strings.HasSuffix(upper, "KB"):
		unit = kilobyte
		numPart = upper[:len(upper)-2]
	case strings.HasSuffix(upper, "B"):
		unit = 1
		numPart = upper[:len(upper)-1]
	}
	numPart = strings.TrimSpace(numPart)

	if numPart == "" {
		return 0, fmt.Errorf("byte size %q: invalid number %q", s, numPart)
	}

	value, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return 0, fmt.Errorf("byte size %q: invalid number %q", s, numPart)
	}
	if value < 0 {
		return 0, fmt.Errorf("byte size %q: must not be negative", s)
	}

	bytes := value * float64(unit)
	if bytes > float64(math.MaxInt64) {
		return 0, fmt.Errorf("byte size %q: value too large", s)
	}

	return int64(bytes), nil
}
