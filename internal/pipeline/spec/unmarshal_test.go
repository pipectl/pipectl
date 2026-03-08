package spec

import (
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
