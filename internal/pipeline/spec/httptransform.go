package spec

import "fmt"

type HTTPTransformStep struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
	Proxy  string `yaml:"proxy"`
}

func (s *HTTPTransformStep) StepType() string {
	return "http-transform"
}

func (s *HTTPTransformStep) String() string {
	return fmt.Sprintf("[%s]: %v %v proxy=%v", s.StepType(), s.URL, s.Method, s.Proxy)
}
