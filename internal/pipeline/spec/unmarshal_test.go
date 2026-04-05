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

func TestStepWrapperUnmarshalCastStep(t *testing.T) {
	raw := []byte(`cast:
  fields:
    age:
      type: int
    created_at:
      type: time
      format: "2006-01-02"
    active:
      type: bool
      true_values: ["yes", "1"]
      false_values: ["no", "0"]
`)

	var step StepWrapper
	if err := yaml.Unmarshal(raw, &step); err != nil {
		t.Fatalf("unmarshal returned error: %v", err)
	}

	castStep, ok := step.Step.(*CastStep)
	if !ok {
		t.Fatalf("expected *CastStep, got %T", step.Step)
	}

	if got := castStep.Fields["age"].Type; got != "int" {
		t.Fatalf("unexpected age type: got %q want %q", got, "int")
	}
	if got := castStep.Fields["created_at"].Format; got != "2006-01-02" {
		t.Fatalf("unexpected created_at format: got %q want %q", got, "2006-01-02")
	}
	if got := castStep.Fields["active"].TrueValues; len(got) != 2 || got[0] != "yes" || got[1] != "1" {
		t.Fatalf("unexpected active true_values: %#v", got)
	}
}

func TestStepWrapperUnmarshalCastStepRejectsInvalidType(t *testing.T) {
	raw := []byte(`cast:
  fields:
    age:
      type: decimal
`)

	var step StepWrapper
	err := yaml.Unmarshal(raw, &step)
	if err == nil {
		t.Fatal("expected unmarshal error for invalid cast type")
	}
	if !strings.Contains(err.Error(), `cast field "age" type must be one of: int, float, bool, time, string`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepWrapperUnmarshalCastStepRejectsInvalidOptionCombinations(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		message string
	}{
		{
			name: "format on non-time field",
			raw: `cast:
  fields:
    age:
      type: int
      format: "2006-01-02"
`,
			message: `cast field "age" format is only supported for type time`,
		},
		{
			name: "bool values on non-bool field",
			raw: `cast:
  fields:
    age:
      type: int
      true_values: ["yes"]
`,
			message: `cast field "age" true_values/false_values are only supported for type bool`,
		},
		{
			name: "overlapping bool values",
			raw: `cast:
  fields:
    active:
      type: bool
      true_values: ["yes"]
      false_values: ["yes"]
`,
			message: `cast field "active" bool true_values and false_values must not overlap`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var step StepWrapper
			err := yaml.Unmarshal([]byte(tt.raw), &step)
			if err == nil {
				t.Fatal("expected unmarshal error")
			}
			if !strings.Contains(err.Error(), tt.message) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestStepWrapperUnmarshalConvertStep(t *testing.T) {
	raw := []byte(`convert:
  format: jsonl
`)

	var step StepWrapper
	if err := yaml.Unmarshal(raw, &step); err != nil {
		t.Fatalf("unmarshal returned error: %v", err)
	}

	convertStep, ok := step.Step.(*ConvertStep)
	if !ok {
		t.Fatalf("expected *ConvertStep, got %T", step.Step)
	}

	if convertStep.Format != "jsonl" {
		t.Fatalf("unexpected format: got %q want %q", convertStep.Format, "jsonl")
	}
}

func TestStepWrapperUnmarshalConvertStepRejectsInvalidFormat(t *testing.T) {
	raw := []byte(`convert:
  format: xml
`)

	var step StepWrapper
	err := yaml.Unmarshal(raw, &step)
	if err == nil {
		t.Fatal("expected unmarshal error for invalid convert format")
	}
	if !strings.Contains(err.Error(), "convert format must be one of: json, jsonl, csv") {
		t.Fatalf("unexpected error: %v", err)
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

func TestStepWrapperUnmarshalLimitStep(t *testing.T) {
	raw := []byte(`limit:
  count: 50
`)

	var step StepWrapper
	if err := yaml.Unmarshal(raw, &step); err != nil {
		t.Fatalf("unmarshal returned error: %v", err)
	}

	limitStep, ok := step.Step.(*LimitStep)
	if !ok {
		t.Fatalf("expected *LimitStep, got %T", step.Step)
	}

	if limitStep.Count != 50 {
		t.Fatalf("unexpected count: got %d want 50", limitStep.Count)
	}
}

func TestStepWrapperUnmarshalLimitStepRejectsZero(t *testing.T) {
	raw := []byte(`limit:
  count: 0
`)

	var step StepWrapper
	err := yaml.Unmarshal(raw, &step)
	if err == nil {
		t.Fatal("expected unmarshal error for count of 0")
	}
	if !strings.Contains(err.Error(), "limit count must be at least 1") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepWrapperUnmarshalLimitStepRejectsNegativeCount(t *testing.T) {
	raw := []byte(`limit:
  count: -5
`)

	var step StepWrapper
	err := yaml.Unmarshal(raw, &step)
	if err == nil {
		t.Fatal("expected unmarshal error for negative count")
	}
	if !strings.Contains(err.Error(), "limit count must be at least 1") {
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
