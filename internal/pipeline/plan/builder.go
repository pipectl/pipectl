package plan

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/steps/default"
	"github.com/shanebell/pipectl/internal/engine/steps/filter"
	"github.com/shanebell/pipectl/internal/engine/steps/httptransform"
	"github.com/shanebell/pipectl/internal/engine/steps/normalize"
	"github.com/shanebell/pipectl/internal/engine/steps/redact"
	"github.com/shanebell/pipectl/internal/engine/steps/rename"
	"github.com/shanebell/pipectl/internal/engine/steps/select"
	"github.com/shanebell/pipectl/internal/engine/steps/validatejson"
	"github.com/shanebell/pipectl/internal/pipeline/spec"
)

func Build(p spec.Pipeline) ([]engine.ExecutableStep, error) {
	executableSteps := make([]engine.ExecutableStep, 0, len(p.Steps))
	for _, stepWrapper := range p.Steps {
		executable, err := buildStep(stepWrapper.Step)
		if err != nil {
			return nil, err
		}
		executableSteps = append(executableSteps, executable)
	}

	return executableSteps, nil
}

func buildStep(step spec.Step) (engine.ExecutableStep, error) {
	switch s := step.(type) {
	case *spec.ValidateJSONStep:
		return &validatejson.Step{
			Schema: s.Schema,
		}, nil
	case *spec.NormalizeStep:
		return &normalize.Step{
			Fields: s.Fields,
		}, nil
	case *spec.DefaultStep:
		return &_default.Step{
			Fields: s.Fields,
		}, nil
	case *spec.RenameStep:
		return &rename.Step{
			Fields: s.Fields,
		}, nil
	case *spec.RedactStep:
		return &redact.Step{
			Fields:   s.Fields,
			Strategy: s.Strategy,
		}, nil
	case *spec.SelectStep:
		return &_select.Step{
			Fields: s.Fields,
		}, nil
	case *spec.FilterStep:
		return &filter.Step{
			Field: s.Field,
			Value: s.Equals,
		}, nil
	case *spec.HTTPTransformStep:
		return &httptransform.Step{
			URL:          s.URL,
			Method:       s.Method,
			Proxy:        s.Proxy,
			Headers:      s.Headers,
			Timeout:      s.Timeout,
			ExpectFormat: s.ExpectFormat,
		}, nil
	default:
		return nil, fmt.Errorf("invalid step type %T", step)
	}
}
