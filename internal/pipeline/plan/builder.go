package plan

import (
	"fmt"
	"strconv"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/steps/assert"
	"github.com/shanebell/pipectl/internal/engine/steps/cast"
	"github.com/shanebell/pipectl/internal/engine/steps/convert"
	"github.com/shanebell/pipectl/internal/engine/steps/count"
	"github.com/shanebell/pipectl/internal/engine/steps/default"
	"github.com/shanebell/pipectl/internal/engine/steps/filter"
	"github.com/shanebell/pipectl/internal/engine/steps/httptransform"
	"github.com/shanebell/pipectl/internal/engine/steps/limit"
	_log "github.com/shanebell/pipectl/internal/engine/steps/log"
	"github.com/shanebell/pipectl/internal/engine/steps/normalize"
	"github.com/shanebell/pipectl/internal/engine/steps/redact"
	"github.com/shanebell/pipectl/internal/engine/steps/rename"
	"github.com/shanebell/pipectl/internal/engine/steps/select"
	_sort "github.com/shanebell/pipectl/internal/engine/steps/sort"
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
	case *spec.CastStep:
		fields := make(map[string]cast.Field, len(s.Fields))
		for name, field := range s.Fields {
			fields[name] = cast.Field{
				Type:        field.Type,
				Format:      field.Format,
				TrueValues:  field.TrueValues,
				FalseValues: field.FalseValues,
			}
		}
		return &cast.Step{
			Fields: fields,
		}, nil
	case *spec.ConvertStep:
		return &convert.Step{
			Format: s.Format,
		}, nil
	case *spec.AssertStep:
		return &assert.Step{
			MinRecords:   s.MinRecords,
			MaxRecords:   s.MaxRecords,
			RecordsEqual: s.RecordsEqual,
			FieldExists:  s.FieldExists,
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
		var op, value string
		var numericValue float64
		switch {
		case s.Equals != "":
			op, value = filter.OpEquals, s.Equals
		case s.NotEquals != "":
			op, value = filter.OpNotEquals, s.NotEquals
		case s.Contains != "":
			op, value = filter.OpContains, s.Contains
		case s.StartsWith != "":
			op, value = filter.OpStartsWith, s.StartsWith
		case s.GreaterThan != "":
			op = filter.OpGreaterThan
			numericValue, _ = strconv.ParseFloat(s.GreaterThan, 64)
		case s.LessThan != "":
			op = filter.OpLessThan
			numericValue, _ = strconv.ParseFloat(s.LessThan, 64)
		}
		return &filter.Step{
			Field:        s.Field,
			Op:           op,
			Value:        value,
			NumericValue: numericValue,
		}, nil
	case *spec.LogStep:
		recordCount := true
		if s.Count != nil {
			recordCount = *s.Count
		}

		sample := 10
		if s.Sample != nil {
			sample = *s.Sample
		}

		return &_log.Step{
			Message: s.Message,
			Count:   recordCount,
			Sample:  sample,
		}, nil
	case *spec.LimitStep:
		return &limit.Step{
			Count: s.Count,
		}, nil
	case *spec.CountStep:
		return &count.Step{
			Message: s.Message,
		}, nil
	case *spec.SortStep:
		return &_sort.Step{
			Field:     s.Field,
			Direction: s.Direction,
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
