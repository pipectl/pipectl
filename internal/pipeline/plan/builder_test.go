package plan

import (
	"testing"

	"github.com/pipectl/pipectl/internal/engine/steps/assert"
	"github.com/pipectl/pipectl/internal/engine/steps/cast"
	"github.com/pipectl/pipectl/internal/engine/steps/convert"
	"github.com/pipectl/pipectl/internal/engine/steps/count"
	"github.com/pipectl/pipectl/internal/engine/steps/default"
	"github.com/pipectl/pipectl/internal/engine/steps/filter"
	"github.com/pipectl/pipectl/internal/engine/steps/limit"
	_log "github.com/pipectl/pipectl/internal/engine/steps/log"
	"github.com/pipectl/pipectl/internal/engine/steps/rename"
	"github.com/pipectl/pipectl/internal/pipeline/spec"
)

func TestBuildRenameStep(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.RenameStep{
					Fields: map[string]string{
						"firstName": "first_name",
						"lastName":  "last_name",
					},
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	renameStep, ok := executableSteps[0].(*rename.Step)
	if !ok {
		t.Fatalf("expected *rename.Step, got %T", executableSteps[0])
	}

	expected := map[string]string{
		"firstName": "first_name",
		"lastName":  "last_name",
	}
	if len(renameStep.Fields) != len(expected) {
		t.Fatalf("unexpected fields count: got %d want %d", len(renameStep.Fields), len(expected))
	}

	for from, to := range expected {
		if got := renameStep.Fields[from]; got != to {
			t.Fatalf("unexpected field mapping for %q: got %q want %q", from, got, to)
		}
	}
}

func TestBuildDefaultStep(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.DefaultStep{
					Fields: map[string]interface{}{
						"country":  "AU",
						"password": "Passw0rd",
					},
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	defaultStep, ok := executableSteps[0].(*_default.Step)
	if !ok {
		t.Fatalf("expected *_default.Step, got %T", executableSteps[0])
	}

	expected := map[string]interface{}{
		"country":  "AU",
		"password": "Passw0rd",
	}
	if len(defaultStep.Fields) != len(expected) {
		t.Fatalf("unexpected fields count: got %d want %d", len(defaultStep.Fields), len(expected))
	}

	for key, value := range expected {
		if got := defaultStep.Fields[key]; got != value {
			t.Fatalf("unexpected default value for %q: got %v want %v", key, got, value)
		}
	}
}

func TestBuildCastStep(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.CastStep{
					Fields: map[string]spec.CastField{
						"age": {
							Type: "int",
						},
						"active": {
							Type:        "bool",
							TrueValues:  []string{"yes"},
							FalseValues: []string{"no"},
						},
					},
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	castStep, ok := executableSteps[0].(*cast.Step)
	if !ok {
		t.Fatalf("expected *cast.Step, got %T", executableSteps[0])
	}

	if got := castStep.Fields["age"].Type; got != "int" {
		t.Fatalf("unexpected age type: got %q want %q", got, "int")
	}
	if got := castStep.Fields["active"].TrueValues; len(got) != 1 || got[0] != "yes" {
		t.Fatalf("unexpected active true-values: %#v", got)
	}
}

func TestBuildConvertStep(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.ConvertStep{
					Format: "jsonl",
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	convertStep, ok := executableSteps[0].(*convert.Step)
	if !ok {
		t.Fatalf("expected *convert.Step, got %T", executableSteps[0])
	}

	if convertStep.Format != "jsonl" {
		t.Fatalf("unexpected format: got %q want %q", convertStep.Format, "jsonl")
	}
}

func TestBuildLogStepDefaults(t *testing.T) {
	// Sample defaulting now happens in spec.LogStep.UnmarshalYAML (see
	// spec/unmarshal_test.go); here we supply it explicitly since Build
	// trusts the spec layer to have already populated it.
	sample := 10
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.LogStep{Sample: &sample},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	logStep, ok := executableSteps[0].(*_log.Step)
	if !ok {
		t.Fatalf("expected *_log.Step, got %T", executableSteps[0])
	}

	if logStep.Message != "" {
		t.Fatalf("unexpected message: got %q want empty", logStep.Message)
	}
	if !logStep.Count {
		t.Fatalf("unexpected count default: got %v want true", logStep.Count)
	}
	if logStep.Sample != 10 {
		t.Fatalf("unexpected sample default: got %d want %d", logStep.Sample, 10)
	}
}

func TestBuildLogStepCustomValues(t *testing.T) {
	countRecords := false
	sample := 3
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.LogStep{
					Message: "after transform",
					Count:   &countRecords,
					Sample:  &sample,
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	logStep, ok := executableSteps[0].(*_log.Step)
	if !ok {
		t.Fatalf("expected *_log.Step, got %T", executableSteps[0])
	}

	if logStep.Message != "after transform" {
		t.Fatalf("unexpected message: got %q want %q", logStep.Message, "after transform")
	}
	if logStep.Count {
		t.Fatalf("unexpected count: got %v want false", logStep.Count)
	}
	if logStep.Sample != 3 {
		t.Fatalf("unexpected sample: got %d want %d", logStep.Sample, 3)
	}
}

func TestBuildCountStep(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.CountStep{
					Message: "records before output",
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	countStep, ok := executableSteps[0].(*count.Step)
	if !ok {
		t.Fatalf("expected *count.Step, got %T", executableSteps[0])
	}

	if countStep.Message != "records before output" {
		t.Fatalf("unexpected message: got %q want %q", countStep.Message, "records before output")
	}
}

func TestBuildLimitStep(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.LimitStep{
					Count: 25,
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	limitStep, ok := executableSteps[0].(*limit.Step)
	if !ok {
		t.Fatalf("expected *limit.Step, got %T", executableSteps[0])
	}

	if limitStep.Count != 25 {
		t.Fatalf("unexpected count: got %d want 25", limitStep.Count)
	}
}

func TestBuildAssertStep(t *testing.T) {
	minRecords := 10
	maxRecords := 1000
	equal := 100
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.AssertStep{
					MinRecords:   &minRecords,
					MaxRecords:   &maxRecords,
					RecordsEqual: &equal,
					FieldExists:  "email",
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	assertStep, ok := executableSteps[0].(*assert.Step)
	if !ok {
		t.Fatalf("expected *assert.Step, got %T", executableSteps[0])
	}

	if assertStep.MinRecords == nil || *assertStep.MinRecords != 10 {
		t.Fatalf("unexpected min-records: got %v want 10", assertStep.MinRecords)
	}
	if assertStep.MaxRecords == nil || *assertStep.MaxRecords != 1000 {
		t.Fatalf("unexpected max-records: got %v want 1000", assertStep.MaxRecords)
	}
	if assertStep.RecordsEqual == nil || *assertStep.RecordsEqual != 100 {
		t.Fatalf("unexpected records-equal: got %v want 100", assertStep.RecordsEqual)
	}
	if assertStep.FieldExists != "email" {
		t.Fatalf("unexpected field-exists: got %q want %q", assertStep.FieldExists, "email")
	}
}

func TestBuildAssertStepFieldChecks(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.AssertStep{
					FieldEquals:   &spec.AssertFieldCheck{Field: "status", Value: "active"},
					FieldContains: &spec.AssertFieldCheck{Field: "email", Value: "@"},
					FieldMatches:  &spec.AssertFieldCheck{Field: "status", Value: "^active$"},
					CaseSensitive: false,
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	assertStep, ok := executableSteps[0].(*assert.Step)
	if !ok {
		t.Fatalf("expected *assert.Step, got %T", executableSteps[0])
	}

	if assertStep.CaseSensitive {
		t.Fatalf("unexpected case-sensitive: got true want false")
	}

	if assertStep.FieldEquals == nil || assertStep.FieldEquals.Field != "status" || assertStep.FieldEquals.Value != "active" {
		t.Fatalf("unexpected field-equals: got %+v", assertStep.FieldEquals)
	}
	if assertStep.FieldContains == nil || assertStep.FieldContains.Field != "email" || assertStep.FieldContains.Value != "@" {
		t.Fatalf("unexpected field-contains: got %+v", assertStep.FieldContains)
	}
	if assertStep.FieldMatches == nil || assertStep.FieldMatches.Field != "status" || assertStep.FieldMatches.Value != "^active$" {
		t.Fatalf("unexpected field-matches: got %+v", assertStep.FieldMatches)
	}
	if assertStep.FieldMatches.Regex == nil {
		t.Fatal("expected field-matches regex to be compiled")
	}
	if !assertStep.FieldMatches.Regex.MatchString("active") {
		t.Fatal("expected compiled field-matches regex to match")
	}
	// case-sensitive: false should fold the regex to case-insensitive.
	if !assertStep.FieldMatches.Regex.MatchString("ACTIVE") {
		t.Fatal("expected compiled field-matches regex to be case-insensitive")
	}
}

func TestBuildFilterStep(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.FilterStep{
					Field:     "status",
					Equals:    "active",
					OnMissing: "include",
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	filterStep, ok := executableSteps[0].(*filter.Step)
	if !ok {
		t.Fatalf("expected *filter.Step, got %T", executableSteps[0])
	}

	if filterStep.Condition.Rule == nil {
		t.Fatal("expected a leaf rule")
	}
	if filterStep.Condition.Rule.Field != "status" {
		t.Fatalf("unexpected field: got %q want %q", filterStep.Condition.Rule.Field, "status")
	}
	if filterStep.Condition.Rule.OnMissing != "include" {
		t.Fatalf("unexpected on-missing: got %q want %q", filterStep.Condition.Rule.OnMissing, "include")
	}
}

func TestBuildFilterStepGreaterThanNumericOrText(t *testing.T) {
	tests := []struct {
		name        string
		greaterThan string
		wantNumeric bool
		wantValue   float64
		wantText    string
	}{
		{name: "numeric threshold", greaterThan: "18", wantNumeric: true, wantValue: 18},
		{name: "text threshold", greaterThan: "banana", wantNumeric: false, wantText: "banana"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipeline := spec.Pipeline{
				Steps: []spec.StepWrapper{
					{Step: &spec.FilterStep{Field: "name", GreaterThan: tt.greaterThan}},
				},
			}

			executableSteps, err := Build(pipeline)
			if err != nil {
				t.Fatalf("build returned error: %v", err)
			}

			rule := executableSteps[0].(*filter.Step).Condition.Rule
			if rule.Numeric != tt.wantNumeric {
				t.Fatalf("unexpected Numeric: got %v want %v", rule.Numeric, tt.wantNumeric)
			}
			if tt.wantNumeric && rule.NumericValue != tt.wantValue {
				t.Fatalf("unexpected NumericValue: got %v want %v", rule.NumericValue, tt.wantValue)
			}
			if !tt.wantNumeric && rule.Value != tt.wantText {
				t.Fatalf("unexpected Value: got %q want %q", rule.Value, tt.wantText)
			}
		})
	}
}

func TestBuildFilterStepAllAnyThreadsOnMissing(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.FilterStep{
					OnMissing: "error",
					All: []spec.FilterCondition{
						{Field: "status", Equals: "active"},
						{
							Any: []spec.FilterCondition{
								{Field: "age", GreaterThan: "18"},
								{Field: "department", Equals: "HR"},
							},
						},
					},
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	filterStep, ok := executableSteps[0].(*filter.Step)
	if !ok {
		t.Fatalf("expected *filter.Step, got %T", executableSteps[0])
	}

	if len(filterStep.Condition.All) != 2 {
		t.Fatalf("unexpected all count: got %d want %d", len(filterStep.Condition.All), 2)
	}
	if got := filterStep.Condition.All[0].Rule.OnMissing; got != "error" {
		t.Fatalf("unexpected on-missing on flat rule: got %q want %q", got, "error")
	}

	nested := filterStep.Condition.All[1]
	if len(nested.Any) != 2 {
		t.Fatalf("unexpected any count: got %d want %d", len(nested.Any), 2)
	}
	for _, sub := range nested.Any {
		if got := sub.Rule.OnMissing; got != "error" {
			t.Fatalf("unexpected on-missing on nested rule: got %q want %q", got, "error")
		}
	}
}

func TestBuildFilterStepCaseSensitivity(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.FilterStep{
					CaseSensitive: false,
					All: []spec.FilterCondition{
						{Field: "status", Equals: "active"},
						{
							Any: []spec.FilterCondition{
								{Field: "age", GreaterThan: "18"},
								{Field: "department", Equals: "HR"},
							},
						},
					},
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	filterStep, ok := executableSteps[0].(*filter.Step)
	if !ok {
		t.Fatalf("expected *filter.Step, got %T", executableSteps[0])
	}

	if got := filterStep.Condition.All[0].Rule.CaseSensitive; got != false {
		t.Fatalf("unexpected case-sensitive on flat rule: got %v want %v", got, false)
	}

	nested := filterStep.Condition.All[1]
	for _, sub := range nested.Any {
		if got := sub.Rule.CaseSensitive; got != false {
			t.Fatalf("unexpected case-sensitive on nested rule: got %v want %v", got, false)
		}
	}
}
