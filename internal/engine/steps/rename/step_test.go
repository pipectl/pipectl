package rename

import (
	"reflect"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "rename" {
		t.Fatalf("expected step name %q, got %q", "rename", step.Name())
	}
}

func TestSupports(t *testing.T) {
	step := &Step{}

	if !step.Supports(&payload.JSON{}) {
		t.Fatal("expected step to support JSON payload")
	}

	if !step.Supports(&payload.CSV{}) {
		t.Fatal("expected step to support CSV payload")
	}
}

func TestExecuteRenamesJSONFields(t *testing.T) {
	step := &Step{
		Fields: map[string]string{
			"firstName": "first_name",
			"lastName":  "last_name",
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Records: []map[string]interface{}{
				{
					"firstName": "Alice",
					"lastName":  "Lee",
					"age":       29,
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
		"first_name": "Alice",
		"last_name":  "Lee",
		"age":        29,
	}
	if !reflect.DeepEqual(out.Records[0], expected) {
		t.Fatalf("unexpected renamed JSON data:\nexpected: %#v\ngot: %#v", expected, out.Records[0])
	}
}

func TestExecuteRenamesCSVHeaderFields(t *testing.T) {
	step := &Step{
		Fields: map[string]string{
			"firstName": "first_name",
			"lastName":  "last_name",
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"firstName", "lastName", "email"},
				{"Alice", "Lee", "alice@example.com"},
				{"Bob", "Ng", "bob@example.com"},
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
		{"first_name", "last_name", "email"},
		{"Alice", "Lee", "alice@example.com"},
		{"Bob", "Ng", "bob@example.com"},
	}
	if !reflect.DeepEqual(out.Rows, expected) {
		t.Fatalf("unexpected renamed CSV rows:\nexpected: %#v\ngot: %#v", expected, out.Rows)
	}
}
