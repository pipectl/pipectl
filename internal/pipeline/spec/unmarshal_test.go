package spec

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestStepWrapperUnmarshalRenameStep(t *testing.T) {
	raw := []byte(`rename:
  fields:
    firstName: first_name
    lastName: last_name
`)

	var step StepWrapper
	if err := yaml.Unmarshal(raw, &step); err != nil {
		t.Fatalf("unmarshal returned error: %v", err)
	}

	renameStep, ok := step.Step.(*RenameStep)
	if !ok {
		t.Fatalf("expected *RenameStep, got %T", step.Step)
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

func TestStepWrapperUnmarshalDefaultStep(t *testing.T) {
	raw := []byte(`default:
  fields:
    country: AU
    password: Passw0rd
`)

	var step StepWrapper
	if err := yaml.Unmarshal(raw, &step); err != nil {
		t.Fatalf("unmarshal returned error: %v", err)
	}

	defaultStep, ok := step.Step.(*DefaultStep)
	if !ok {
		t.Fatalf("expected *DefaultStep, got %T", step.Step)
	}

	expected := map[string]interface{}{
		"country":  "AU",
		"password": "Passw0rd",
	}
	if len(defaultStep.Fields) != len(expected) {
		t.Fatalf("unexpected fields count: got %d want %d", len(defaultStep.Fields), len(expected))
	}

	for key, expectedValue := range expected {
		got, exists := defaultStep.Fields[key]
		if !exists {
			t.Fatalf("expected field %q to exist", key)
		}
		if got != expectedValue {
			t.Fatalf("unexpected default value for %q: got %v want %v", key, got, expectedValue)
		}
	}
}

func TestStepWrapperUnmarshalLogStep(t *testing.T) {
	raw := []byte(`log:
  message: Payload after step 2
  count: true
  sample: 10
`)

	var step StepWrapper
	if err := yaml.Unmarshal(raw, &step); err != nil {
		t.Fatalf("unmarshal returned error: %v", err)
	}

	logStep, ok := step.Step.(*LogStep)
	if !ok {
		t.Fatalf("expected *LogStep, got %T", step.Step)
	}

	if logStep.Message != "Payload after step 2" {
		t.Fatalf("unexpected message: got %q want %q", logStep.Message, "Payload after step 2")
	}

	if logStep.Count == nil || !*logStep.Count {
		t.Fatalf("unexpected count: got %v want true", logStep.Count)
	}

	if logStep.Sample == nil || *logStep.Sample != 10 {
		t.Fatalf("unexpected sample: got %v want 10", logStep.Sample)
	}
}

func TestStepWrapperUnmarshalCountStep(t *testing.T) {
	raw := []byte(`count:
  message: Payload before output
`)

	var step StepWrapper
	if err := yaml.Unmarshal(raw, &step); err != nil {
		t.Fatalf("unmarshal returned error: %v", err)
	}

	countStep, ok := step.Step.(*CountStep)
	if !ok {
		t.Fatalf("expected *CountStep, got %T", step.Step)
	}

	if countStep.Message != "Payload before output" {
		t.Fatalf("unexpected message: got %q want %q", countStep.Message, "Payload before output")
	}
}

func TestStepWrapperUnmarshalAssertStep(t *testing.T) {
	raw := []byte(`assert:
  min-records: 10
  max-records: 1000
  records-equal: 100
  field-exists: email
`)

	var step StepWrapper
	if err := yaml.Unmarshal(raw, &step); err != nil {
		t.Fatalf("unmarshal returned error: %v", err)
	}

	assertStep, ok := step.Step.(*AssertStep)
	if !ok {
		t.Fatalf("expected *AssertStep, got %T", step.Step)
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

func TestStepWrapperUnmarshalAssertStepRequiresAtLeastOneOption(t *testing.T) {
	raw := []byte(`assert: {}`)

	var step StepWrapper
	err := yaml.Unmarshal(raw, &step)
	if err == nil {
		t.Fatal("expected unmarshal error for assert with no options")
	}
	if !strings.Contains(err.Error(), "assert requires at least one option") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepWrapperUnmarshalAssertStepValidatesBounds(t *testing.T) {
	raw := []byte(`assert:
  min-records: 100
  max-records: 10
`)

	var step StepWrapper
	err := yaml.Unmarshal(raw, &step)
	if err == nil {
		t.Fatal("expected unmarshal error for invalid min/max bounds")
	}
	if !strings.Contains(err.Error(), "assert min-records must be <= max-records") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepWrapperUnmarshalAssertStepValidatesRecordsEqualBounds(t *testing.T) {
	raw := []byte(`assert:
  min-records: 10
  records-equal: 9
`)

	var step StepWrapper
	err := yaml.Unmarshal(raw, &step)
	if err == nil {
		t.Fatal("expected unmarshal error for inconsistent records-equal/min-records bounds")
	}
	if !strings.Contains(err.Error(), "assert records-equal must be >= min-records") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepWrapperUnmarshalAssertStepWithRecordsEqualOnly(t *testing.T) {
	raw := []byte(`assert:
  records-equal: 16
`)

	var step StepWrapper
	if err := yaml.Unmarshal(raw, &step); err != nil {
		t.Fatalf("unmarshal returned error: %v", err)
	}

	assertStep, ok := step.Step.(*AssertStep)
	if !ok {
		t.Fatalf("expected *AssertStep, got %T", step.Step)
	}

	if assertStep.RecordsEqual == nil || *assertStep.RecordsEqual != 16 {
		t.Fatalf("unexpected records-equal: got %v want 16", assertStep.RecordsEqual)
	}
}
