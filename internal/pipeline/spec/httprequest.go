package spec

import (
	"fmt"
	"strings"
)

type HTTPRequestStep struct {
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Proxy   string            `yaml:"proxy"`
	Headers map[string]string `yaml:"headers"`
	Timeout int               `yaml:"timeout"`
}

func (s *HTTPRequestStep) StepType() string {
	return "http-request"
}

func (s *HTTPRequestStep) String() string {
	return fmt.Sprintf("[%s]: %v %v proxy=%v headers=%v timeout=%v", s.StepType(), s.URL, s.Method, s.Proxy, s.Headers, s.Timeout)
}

func (s *HTTPRequestStep) Validate() error {
	s.Method = strings.ToUpper(strings.TrimSpace(s.Method))

	if strings.TrimSpace(s.URL) == "" {
		return fmt.Errorf("http-request url is required")
	}

	if s.Method == "" {
		return fmt.Errorf("http-request method is required")
	}

	validMethod := false
	for _, m := range validHTTPMethods {
		if s.Method == m {
			validMethod = true
			break
		}
	}
	if !validMethod {
		return fmt.Errorf("http-request method must be one of: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
	}

	if s.Timeout < 0 {
		return fmt.Errorf("http-request timeout must be >= 0")
	}

	if s.Timeout > 300 {
		return fmt.Errorf("http-request timeout must be <= 300 seconds")
	}

	return nil
}
