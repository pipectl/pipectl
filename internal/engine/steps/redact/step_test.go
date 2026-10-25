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
			name:     "partial-last non-default N",
			strategy: "partial-last:6",
			input:    "1234-5678-9012-3456",
			expected: "*************2-3456",
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
			name:     "partial-first non-default N",
			strategy: "partial-first:6",
			input:    "1234-5678-9012-3456",
			expected: "1234-5*************",
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

func TestExecuteRedactsNestedJSONFields(t *testing.T) {
	newRecord := func() map[string]interface{} {
		return map[string]interface{}{
			"user": map[string]interface{}{
				"email":  "alice@example.com",
				"salary": float64(85000),
				"active": true,
				"bank":   map[string]interface{}{"bsb": "062-000"},
				"tags":   []interface{}{"vip"},
				"notes":  nil,
			},
			"contacts": []interface{}{
				map[string]interface{}{"phone": "0400123456"},
			},
		}
	}

	testCases := []struct {
		name     string
		strategy string
		field    string
		expected interface{}
		err      string
	}{
		{name: "mask nested object field", strategy: "mask", field: "user.email", expected: "*****************"},
		{name: "partial-last array index field", strategy: "partial-last:4", field: "contacts[0].phone", expected: "******3456"},
		{name: "mask number errors", strategy: "mask", field: "user.salary", err: `field "user.salary" in record 1 is a number; strategy "mask" only applies to strings — cast it to string first, or use strategy: nullify`},
		{name: "default strategy on bool errors", strategy: "", field: "user.active", err: "is a bool; strategy REDACTED (default) only applies to strings"},
		{name: "sha256 object errors", strategy: "sha256", field: "user.bank", err: "is an object"},
		{name: "mask array errors", strategy: "mask", field: "user.tags", err: "is an array"},
		{name: "mask existing null left unchanged", strategy: "mask", field: "user.notes", expected: nil},
		{name: "nullify string", strategy: "nullify", field: "user.email", expected: nil},
		{name: "nullify number", strategy: "nullify", field: "user.salary", expected: nil},
		{name: "nullify object", strategy: "nullify", field: "user.bank", expected: nil},
		{name: "missing nested field errors", strategy: "mask", field: "user.phone", err: "not found in record"},
		{name: "out of range index errors", strategy: "mask", field: "contacts[3].phone", err: "out of range"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			record := newRecord()
			step := &Step{Strategy: tc.strategy, Fields: []string{tc.field}}
			ctx := &engine.ExecutionContext{
				Payload: &payload.JSONL{Items: []map[string]interface{}{record}},
			}

			err := step.Execute(ctx)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("expected error containing %q, got %v", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("execute returned error: %v", err)
			}

			got, _, err := payload.ResolveFieldPath(record, tc.field)
			if err != nil {
				t.Fatalf("resolve field after redact: %v", err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("expected %#v, got %#v", tc.expected, got)
			}
		})
	}
}

func TestExecuteInvalidFieldPathLeavesRecordsUnchanged(t *testing.T) {
	step := &Step{Strategy: "mask", Fields: []string{"email", "user..name"}}
	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{{"email": "alice@example.com"}},
		},
	}

	err := step.Execute(ctx)
	if err == nil || !strings.Contains(err.Error(), "invalid field path") {
		t.Fatalf("expected invalid field path error, got %v", err)
	}
	if got := ctx.Payload.(*payload.JSONL).Items[0]["email"]; got != "alice@example.com" {
		t.Fatalf("expected record to be unchanged, got %#v", got)
	}
}

func TestExecuteNullifiesCSVFields(t *testing.T) {
	step := &Step{Strategy: "nullify", Fields: []string{"email"}}
	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "email"},
				{"1", "alice@example.com"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	expected := [][]string{{"id", "email"}, {"1", ""}}
	if got := ctx.Payload.(*payload.CSV).Rows; !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %#v, got %#v", expected, got)
	}
}
