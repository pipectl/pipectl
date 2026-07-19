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

func (s *HTTPTransformStep) String() string {
	return fmt.Sprintf("[%s]: %v %v proxy=%v headers=%v timeout=%v expect-format=%v", s.StepType(), s.URL, s.Method, s.Proxy, s.Headers, s.Timeout, s.ExpectFormat)
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
