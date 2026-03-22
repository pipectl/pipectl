package plan

import (
	"testing"

	"github.com/shanebell/pipectl/internal/engine/steps/assert"
	"github.com/shanebell/pipectl/internal/engine/steps/cast"
	"github.com/shanebell/pipectl/internal/engine/steps/convert"
	"github.com/shanebell/pipectl/internal/engine/steps/count"
	"github.com/shanebell/pipectl/internal/engine/steps/default"
	_log "github.com/shanebell/pipectl/internal/engine/steps/log"
	"github.com/shanebell/pipectl/internal/engine/steps/rename"
	"github.com/shanebell/pipectl/internal/pipeline/spec"
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
		t.Fatalf("unexpected active true_values: %#v", got)
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
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.LogStep{},
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
