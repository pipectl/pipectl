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

func (s *HTTPRequestStep) Validate() error {
	if strings.TrimSpace(s.URL) == "" {
		return fmt.Errorf("http-request url is required")
	}

	method, err := validateHTTPMethod(s.StepType(), s.Method)
	if err != nil {
		return err
	}
	s.Method = method

	if s.Timeout < minTimeoutSeconds {
		return fmt.Errorf("http-request timeout must be >= %d", minTimeoutSeconds)
	}

	if s.Timeout > maxTimeoutSeconds {
		return fmt.Errorf("http-request timeout must be <= %d seconds", maxTimeoutSeconds)
	}

	return nil
}
