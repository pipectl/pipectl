package filter

import (
	"reflect"
	"strings"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "filter" {
		t.Fatalf("expected step name %q, got %q", "filter", step.Name())
	}
}

func TestSupports(t *testing.T) {
	step := &Step{}

	if !step.Supports(&payload.CSV{}) {
		t.Fatal("expected step to support CSV payload")
	}

	if step.Supports(&payload.JSON{}) {
		t.Fatal("did not expect step to support JSON payload")
	}
}

func TestExecuteFiltersCSVRows(t *testing.T) {
	step := &Step{
		Field: "status",
		Value: "active",
	}

	csvPayload := &payload.CSV{
		Rows: [][]string{
			{"id", "status"},
			{"1", "active"},
			{"2", "inactive"},
			{"3", "active"},
		},
	}

	ctx := &engine.ExecutionContext{Payload: csvPayload}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out, ok := ctx.Payload.(*payload.CSV)
	if !ok {
		t.Fatalf("expected payload.CSV, got %T", ctx.Payload)
	}

	expected := [][]string{
		{"id", "status"},
		{"1", "active"},
		{"3", "active"},
	}

	if !reflect.DeepEqual(out.Rows, expected) {
		t.Fatalf("unexpected filtered rows:\nexpected: %#v\ngot: %#v", expected, out.Rows)
	}
}

func TestExecuteReturnsErrorForNonCSVPayload(t *testing.T) {
	step := &Step{
		Field: "status",
		Value: "active",
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{Data: map[string]interface{}{"status": "active"}},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error for non-CSV payload")
	}
	if !strings.Contains(err.Error(), "requires CSV payload") {
		t.Fatalf("unexpected error: %v", err)
	}
}
