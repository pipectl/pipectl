package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/steps/httptransform"
)

type HTTPTransformStep struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
}

func (s *HTTPTransformStep) StepType() string {
	return "http-transform"
}

func (s *HTTPTransformStep) BuildExecutor() (engine.ExecutableStep, error) {
	return &httptransform.Step{
		URL:    s.URL,
		Method: s.Method,
	}, nil
}

func (s *HTTPTransformStep) String() string {
	return fmt.Sprintf("[%s]: %v %v", s.StepType(), s.URL, s.Method)
}
