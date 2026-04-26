package redact

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "redact" {
		t.Fatalf("expected step name %q, got %q", "redact", step.Name())
	}
}

func TestSupports(t *testing.T) {
	step := &Step{}

	if !step.Supports(&payload.JSON{}) {
		t.Fatal("expected step to support JSON payload")
	}
	if !step.Supports(&payload.JSONL{}) {
		t.Fatal("expected step to support JSONL payload")
	}

	if !step.Supports(&payload.CSV{}) {
		t.Fatal("expected step to support CSV payload")
	}
}

func TestRedactSingleValueStrategies(t *testing.T) {
	testCases := []struct {
		name     string
		strategy string
		input    string
		expected string
	}{
		{
			name:     "mask",
			strategy: "mask",
			input:    "abc123",
			expected: "******",
		},
		{
			name:     "sha256",
			strategy: "sha256",
			input:    "secret",
			expected: "2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b",
		},
		{
			name:     "default",
			strategy: "unknown",
			input:    "secret",
			expected: "REDACTED",
		},
		{
			name:     "partial-last explicit N",
			strategy: "partial-last:4",
			input:    "1234-5678-9012-3456",
			expected: "***************3456",
		},
		{
			name:     "partial-last bare defaults to 4",
			strategy: "partial-last",
			input:    "1234-5678-9012-3456",
			expected: "***************3456",
		},
		{
			name:     "partial-last N >= len returns value unchanged",
			strategy: "partial-last:4",
			input:    "abc",
			expected: "abc",
		},
		{
			name:     "partial-first explicit N",
			strategy: "partial-first:4",
			input:    "1234-5678-9012-3456",
			expected: "1234***************",
		},
		{
			name:     "partial-first bare defaults to 4",
			strategy: "partial-first",
			input:    "1234-5678-9012-3456",
			expected: "1234***************",
		},
		{
			name:     "partial-first N >= len returns value unchanged",
			strategy: "partial-first:4",
			input:    "ab",
			expected: "ab",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			step := &Step{Strategy: tc.strategy}
			if got := step.redactSingleValue(tc.input); got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestExecuteRedactsJSONFields(t *testing.T) {
	step := &Step{
		Strategy: "mask",
		Fields:   []string{"email", "ssn"},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"name":  "Alice",
					"email": "alice@example.com",
					"ssn":   "123-45-6789",
					"age":   42,
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out, ok := ctx.Payload.(*payload.JSON)
	if !ok {
		t.Fatalf("expected payload.JSON, got %T", ctx.Payload)
	}

	expected := map[string]interface{}{
		"name":  "Alice",
		"email": "*****************",
		"ssn":   "***********",
		"age":   42,
	}
	if !reflect.DeepEqual(out.Items[0], expected) {
		t.Fatalf("unexpected redacted JSON data:\nexpected: %#v\ngot: %#v", expected, out.Items[0])
	}
}

func TestExecuteRedactsJSONLFields(t *testing.T) {
	step := &Step{
		Strategy: "mask",
		Fields:   []string{"email"},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"email": "alice@example.com"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out := ctx.Payload.(*payload.JSONL)
	if out.Items[0]["email"] != "*****************" {
		t.Fatalf("unexpected redacted JSONL data: %#v", out.Items[0])
	}
}

func TestExecuteErrorsOnMissingJSONField(t *testing.T) {
	step := &Step{
		Strategy: "mask",
		Fields:   []string{"email", "missing"},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{"email": "alice@example.com"},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for missing field, got nil")
	}
	if !strings.Contains(err.Error(), "not found in record") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestExecuteErrorsOnMissingCSVField(t *testing.T) {
	step := &Step{
		Strategy: "mask",
		Fields:   []string{"email", "missing"},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "email"},
				{"1", "alice@example.com"},
			},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for missing CSV field, got nil")
	}
	if !strings.Contains(err.Error(), "not found in CSV headers") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestExecuteRedactsCSVFields(t *testing.T) {
	step := &Step{
		Strategy: "mask",
		Fields:   []string{"email", "token"},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "email", "token"},
				{"1", "alice@example.com", "abc123"},
				{"2", "bob@example.com", "xyz789"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out, ok := ctx.Payload.(*payload.CSV)
	if !ok {
		t.Fatalf("expected payload.CSV, got %T", ctx.Payload)
	}

	expected := [][]string{
		{"id", "email", "token"},
		{"1", "*****************", "******"},
		{"2", "***************", "******"},
	}
	if !reflect.DeepEqual(out.Rows, expected) {
		t.Fatalf("unexpected redacted CSV rows:\nexpected: %#v\ngot: %#v", expected, out.Rows)
	}
}
