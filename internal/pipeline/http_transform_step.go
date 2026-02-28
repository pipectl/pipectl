package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/steps"
	"github.com/shanebell/pipectl/internal/steps/http_transform"
)

type HTTPTransformStep struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
}

func (s *HTTPTransformStep) StepType() string {
	return "http-transform"
}

func (s *HTTPTransformStep) BuildExecutor() (steps.ExecutableStep, error) {
	return &http_transform.Step{
		URL:    s.URL,
		Method: s.Method,
	}, nil
}

func (s *HTTPTransformStep) String() string {
	return fmt.Sprintf("[%s]: %v %v", s.StepType(), s.URL, s.Method)
}
