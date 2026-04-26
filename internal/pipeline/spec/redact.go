package spec

import (
	"fmt"
	"strconv"
	"strings"
)

type RedactStep struct {
	Strategy string   `yaml:"strategy"`
	Fields   []string `yaml:"fields"`
}

func (s *RedactStep) StepType() string {
	return "redact"
}

func (s *RedactStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *RedactStep) Validate() error {
	if len(s.Fields) == 0 {
		return fmt.Errorf("redact requires at least one field")
	}

	switch {
	case s.Strategy == "", s.Strategy == "mask", s.Strategy == "sha256":
	case s.Strategy == "partial-first", s.Strategy == "partial-last":
	case strings.HasPrefix(s.Strategy, "partial-first:"),
		strings.HasPrefix(s.Strategy, "partial-last:"):
		suffix := s.Strategy[strings.Index(s.Strategy, ":")+1:]
		n, err := strconv.Atoi(suffix)
		if err != nil || n <= 0 {
			return fmt.Errorf("redact partial strategy requires a positive integer suffix (e.g. partial-last:4)")
		}
	default:
		return fmt.Errorf("redact strategy must be one of: mask, sha256, partial-first[:N], partial-last[:N]")
	}

	return nil
}
