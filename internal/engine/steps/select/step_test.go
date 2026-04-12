package _select

import (
	"reflect"
	"testing"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

func TestSupports(t *testing.T) {
	step := &Step{}

	if !step.Supports(&payload.CSV{}) {
		t.Fatal("expected step to support CSV payload")
	}

	if !step.Supports(&payload.JSON{}) {
		t.Fatal("expected step to support JSON payload")
	}

	if !step.Supports(&payload.JSONL{}) {
		t.Fatal("expected step to support JSONL payload")
	}
}

func TestExecuteSelectsJSONFields(t *testing.T) {
	step := &Step{
		Fields: []string{"id", "email"},
	}

	jsonPayload := &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"id": "1", "name": "Alice", "email": "alice@example.com", "country": "AU"},
			{"id": "2", "name": "Bob", "email": "bob@example.com", "country": "NZ"},
		},
	}

	ctx := &engine.ExecutionContext{Payload: jsonPayload}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out, ok := ctx.Payload.(*payload.JSON)
	if !ok {
		t.Fatalf("expected payload.JSON, got %T", ctx.Payload)
	}

	expected := []map[string]interface{}{
		{"id": "1", "email": "alice@example.com"},
		{"id": "2", "email": "bob@example.com"},
	}

	if !reflect.DeepEqual(out.Items, expected) {
		t.Fatalf("unexpected items:\nexpected: %#v\ngot: %#v", expected, out.Items)
	}
}

func TestExecuteSelectsJSONLFields(t *testing.T) {
	step := &Step{
		Fields: []string{"id", "email"},
	}

	jsonlPayload := &payload.JSONL{
		Items: []map[string]interface{}{
			{"id": "1", "name": "Alice", "email": "alice@example.com"},
			{"id": "2", "name": "Bob", "email": "bob@example.com"},
		},
	}

	ctx := &engine.ExecutionContext{Payload: jsonlPayload}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out, ok := ctx.Payload.(*payload.JSONL)
	if !ok {
		t.Fatalf("expected payload.JSONL, got %T", ctx.Payload)
	}

	expected := []map[string]interface{}{
		{"id": "1", "email": "alice@example.com"},
		{"id": "2", "email": "bob@example.com"},
	}

	if !reflect.DeepEqual(out.Items, expected) {
		t.Fatalf("unexpected items:\nexpected: %#v\ngot: %#v", expected, out.Items)
	}
}

func TestExecuteSelectsJSONFieldsMissingFieldOmitted(t *testing.T) {
	step := &Step{
		Fields: []string{"id", "email"},
	}

	jsonPayload := &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"id": "1", "name": "Alice"},
		},
	}

	ctx := &engine.ExecutionContext{Payload: jsonPayload}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out := ctx.Payload.(*payload.JSON)

	expected := []map[string]interface{}{
		{"id": "1"},
	}

	if !reflect.DeepEqual(out.Items, expected) {
		t.Fatalf("unexpected items:\nexpected: %#v\ngot: %#v", expected, out.Items)
	}
}

func TestExecuteSelectsCSVFields(t *testing.T) {
	step := &Step{
		Fields: []string{"id", "email"},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "name", "email", "country"},
				{"1", "Alice", "alice@example.com", "AU"},
				{"2", "Bob", "bob@example.com", "NZ"},
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
		{"id", "email"},
		{"1", "alice@example.com"},
		{"2", "bob@example.com"},
	}

	if len(out.Rows) != len(expected) {
		t.Fatalf("unexpected row count: got %d want %d", len(out.Rows), len(expected))
	}

	for i := range expected {
		if len(out.Rows[i]) != len(expected[i]) {
			t.Fatalf("unexpected column count in row %d: got %d want %d", i, len(out.Rows[i]), len(expected[i]))
		}
		for j := range expected[i] {
			if out.Rows[i][j] != expected[i][j] {
				t.Fatalf("unexpected value at row %d col %d: got %q want %q", i, j, out.Rows[i][j], expected[i][j])
			}
		}
	}
}
