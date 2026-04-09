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
		var condition *filter.Condition
		if len(s.All) > 0 {
			subs, err := buildFilterConditions(s.All)
			if err != nil {
				return nil, err
			}
			condition = &filter.Condition{All: subs}
		} else if len(s.Any) > 0 {
			subs, err := buildFilterConditions(s.Any)
			if err != nil {
				return nil, err
			}
			condition = &filter.Condition{Any: subs}
		} else {
			rule := buildFilterRule(s.Field, s.Equals, s.NotEquals, s.Contains, s.StartsWith, s.EndsWith, s.GreaterThan, s.LessThan)
			condition = &filter.Condition{Rule: rule}
		}
		return &filter.Step{Condition: condition}, nil
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

func buildFilterConditions(conditions []spec.FilterCondition) ([]*filter.Condition, error) {
	result := make([]*filter.Condition, 0, len(conditions))
	for _, c := range conditions {
		built, err := buildFilterCondition(c)
		if err != nil {
			return nil, err
		}
		result = append(result, built)
	}
	return result, nil
}

func buildFilterCondition(c spec.FilterCondition) (*filter.Condition, error) {
	if len(c.All) > 0 {
		subs, err := buildFilterConditions(c.All)
		if err != nil {
			return nil, err
		}
		return &filter.Condition{All: subs}, nil
	}
	if len(c.Any) > 0 {
		subs, err := buildFilterConditions(c.Any)
		if err != nil {
			return nil, err
		}
		return &filter.Condition{Any: subs}, nil
	}
	rule := buildFilterRule(c.Field, c.Equals, c.NotEquals, c.Contains, c.StartsWith, c.EndsWith, c.GreaterThan, c.LessThan)
	return &filter.Condition{Rule: rule}, nil
}

func buildFilterRule(field, equals, notEquals, contains, startsWith, endsWith, greaterThan, lessThan string) *filter.Rule {
	rule := &filter.Rule{Field: field}
	switch {
	case equals != "":
		rule.Op, rule.Value = filter.OpEquals, equals
	case notEquals != "":
		rule.Op, rule.Value = filter.OpNotEquals, notEquals
	case contains != "":
		rule.Op, rule.Value = filter.OpContains, contains
	case startsWith != "":
		rule.Op, rule.Value = filter.OpStartsWith, startsWith
	case endsWith != "":
		rule.Op, rule.Value = filter.OpEndsWith, endsWith
	case greaterThan != "":
		rule.Op = filter.OpGreaterThan
		rule.NumericValue, _ = strconv.ParseFloat(greaterThan, 64) // already validated as a number in spec
	case lessThan != "":
		rule.Op = filter.OpLessThan
		rule.NumericValue, _ = strconv.ParseFloat(lessThan, 64) // already validated as a number in spec
	}
	return rule
}
