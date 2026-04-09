package spec

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

var validHTTPMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

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

func (s *HTTPTransformStep) UnmarshalYAML(b []byte) error {
	type rawHTTPTransformStep HTTPTransformStep
	var raw rawHTTPTransformStep
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}
	*s = HTTPTransformStep(raw)
	s.Method = strings.ToUpper(strings.TrimSpace(s.Method))
	return s.Validate()
}

func (s *HTTPTransformStep) Validate() error {
	if strings.TrimSpace(s.URL) == "" {
		return fmt.Errorf("http-transform url is required")
	}

	if s.Method == "" {
		return fmt.Errorf("http-transform method is required")
	}

	validMethod := false
	for _, m := range validHTTPMethods {
		if s.Method == m {
			validMethod = true
			break
		}
	}
	if !validMethod {
		return fmt.Errorf("http-transform method must be one of: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
	}

	if s.Timeout < 0 {
		return fmt.Errorf("http-transform timeout must be >= 0")
	}

	if s.Timeout > 300 {
		return fmt.Errorf("http-transform timeout must be <= 300 seconds")
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
