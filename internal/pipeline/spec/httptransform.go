package spec

import "fmt"

type HTTPTransformStep struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
}

func (s *HTTPTransformStep) StepType() string {
	return "http-transform"
}

func (s *HTTPTransformStep) String() string {
	return fmt.Sprintf("[%s]: %v %v", s.StepType(), s.URL, s.Method)
}
