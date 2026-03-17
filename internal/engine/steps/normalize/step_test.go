package normalize

import (
	"reflect"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestNormalizeValueStrategies(t *testing.T) {
	step := &Step{}

	testCases := []struct {
		name     string
		input    string
		strategy string
		expected string
	}{
		{
			name:     "lower",
			input:    " HeLLo ",
			strategy: "lower",
			expected: " hello ",
		},
		{
			name:     "upper",
			input:    " HeLLo ",
			strategy: "upper",
			expected: " HELLO ",
		},
		{
			name:     "trim",
			input:    " \t hello \n ",
			strategy: "trim",
			expected: "hello",
		},
		{
			name:     "trim-left",
			input:    " \t hello \n ",
			strategy: "trim-left",
			expected: "hello \n ",
		},
		{
			name:     "trim-right",
			input:    " \t hello \n ",
			strategy: "trim-right",
			expected: " \t hello",
		},
		{
			name:     "collapse-spaces",
			input:    "  hello   world \t from\npipectl ",
			strategy: "collapse-spaces",
			expected: "hello world from pipectl",
		},
		{
			name:     "capitalize",
			input:    "aLICE",
			strategy: "capitalize",
			expected: "Alice",
		},
		{
			name:     "capitalize-empty",
			input:    "",
			strategy: "capitalize",
			expected: "",
		},
		{
			name:     "unknown-strategy-returns-original",
			input:    "  Keep Me ",
			strategy: "does-not-exist",
			expected: "  Keep Me ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := step.normalizeValue(tc.input, tc.strategy); got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "normalize" {
		t.Fatalf("expected step name %q, got %q", "normalize", step.Name())
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

func TestExecuteNormalizesJSONFields(t *testing.T) {
	step := &Step{
		Fields: map[string]string{
			"name":   "trim",
			"status": "lower",
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Records: []map[string]interface{}{
				{
					"name":   "  Alice  ",
					"status": " ACTIVE ",
					"count":  7,
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
		"name":   "Alice",
		"status": " active ",
		"count":  7,
	}
	if !reflect.DeepEqual(out.Records[0], expected) {
		t.Fatalf("unexpected normalized JSON data:\nexpected: %#v\ngot: %#v", expected, out.Records[0])
	}
}

func TestExecuteNormalizesJSONLFields(t *testing.T) {
	step := &Step{
		Fields: map[string]string{
			"name": "trim",
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Records: []map[string]interface{}{
				{"name": " Alice "},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out := ctx.Payload.(*payload.JSONL)
	if out.Records[0]["name"] != "Alice" {
		t.Fatalf("unexpected normalized JSONL data: %#v", out.Records[0])
	}
}

func TestExecuteNormalizesCSVFields(t *testing.T) {
	step := &Step{
		Fields: map[string]string{
			"name":  "trim",
			"email": "lower",
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"name", "email", "id"},
				{" Alice ", "Alice@Example.Com", "1"},
				{" Bob ", "Bob@Example.Com", "2"},
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
		{"name", "email", "id"},
		{"Alice", "alice@example.com", "1"},
		{"Bob", "bob@example.com", "2"},
	}
	if !reflect.DeepEqual(out.Rows, expected) {
		t.Fatalf("unexpected normalized CSV rows:\nexpected: %#v\ngot: %#v", expected, out.Rows)
	}
}
