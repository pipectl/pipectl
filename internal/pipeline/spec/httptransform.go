package spec

import "fmt"

type HTTPTransformStep struct {
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Proxy   string            `yaml:"proxy"`
	Headers map[string]string `yaml:"headers"`
	Timeout int               `yaml:"timeout"`
}

func (s *HTTPTransformStep) StepType() string {
	return "http-transform"
}

func (s *HTTPTransformStep) String() string {
	return fmt.Sprintf("[%s]: %v %v proxy=%v headers=%v timeout=%v", s.StepType(), s.URL, s.Method, s.Proxy, s.Headers, s.Timeout)
}
