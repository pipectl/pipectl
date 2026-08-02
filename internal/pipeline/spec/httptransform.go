package spec

import (
	"fmt"
	"strings"
)

type HTTPTransformStep struct {
	URL          string            `yaml:"url"`
	Method       string            `yaml:"method"`
	Proxy        string            `yaml:"proxy"`
	Headers      map[string]string `yaml:"headers"`
	Timeout      int               `yaml:"timeout"`
	ExpectFormat string            `yaml:"expect-format"`
}

func (s *HTTPTransformStep) StepType() string {
	return "http-transform"
}

func (s *HTTPTransformStep) Validate() error {
	if strings.TrimSpace(s.URL) == "" {
		return fmt.Errorf("http-transform url is required")
	}

	method, err := validateHTTPMethod(s.StepType(), s.Method)
	if err != nil {
		return err
	}
	s.Method = method

	if s.Timeout < minTimeoutSeconds {
		return fmt.Errorf("http-transform timeout must be >= %d", minTimeoutSeconds)
	}

	if s.Timeout > maxTimeoutSeconds {
		return fmt.Errorf("http-transform timeout must be <= %d seconds", maxTimeoutSeconds)
	}

	if s.ExpectFormat != "" {
		switch strings.ToLower(strings.TrimSpace(s.ExpectFormat)) {
		case "json", "jsonl", "csv":
		default:
			return fmt.Errorf("http-transform expect-format must be one of: json, jsonl, csv")
		}
	}

	return nil
}
