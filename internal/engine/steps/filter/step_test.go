package filter

import (
	"reflect"
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

	if !step.Supports(&payload.JSON{}) {
		t.Fatal("expected step to support JSON payload")
	}

	if !step.Supports(&payload.JSONL{}) {
		t.Fatal("expected step to support JSONL payload")
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

func TestExecuteFiltersJSONRecords(t *testing.T) {
	step := &Step{
		Field: "status",
		Value: "active",
	}

	jsonPayload := &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"id": "1", "status": "active"},
			{"id": "2", "status": "inactive"},
			{"id": "3", "status": "active"},
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
		{"id": "1", "status": "active"},
		{"id": "3", "status": "active"},
	}

	if !reflect.DeepEqual(out.Items, expected) {
		t.Fatalf("unexpected filtered items:\nexpected: %#v\ngot: %#v", expected, out.Items)
	}
}

func TestExecuteFiltersJSONLRecords(t *testing.T) {
	step := &Step{
		Field: "status",
		Value: "active",
	}

	jsonlPayload := &payload.JSONL{
		Items: []map[string]interface{}{
			{"id": "1", "status": "active"},
			{"id": "2", "status": "inactive"},
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
		{"id": "1", "status": "active"},
	}

	if !reflect.DeepEqual(out.Items, expected) {
		t.Fatalf("unexpected filtered items:\nexpected: %#v\ngot: %#v", expected, out.Items)
	}
}

func TestExecuteFiltersJSONRecordsWithNumericField(t *testing.T) {
	step := &Step{
		Field: "count",
		Value: "5",
	}

	jsonPayload := &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"id": "1", "count": float64(5)},
			{"id": "2", "count": float64(3)},
		},
	}

	ctx := &engine.ExecutionContext{Payload: jsonPayload}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out := ctx.Payload.(*payload.JSON)

	expected := []map[string]interface{}{
		{"id": "1", "count": float64(5)},
	}

	if !reflect.DeepEqual(out.Items, expected) {
		t.Fatalf("unexpected filtered items:\nexpected: %#v\ngot: %#v", expected, out.Items)
	}
}

func TestExecuteFiltersJSONRecordsMissingField(t *testing.T) {
	step := &Step{
		Field: "missing",
		Value: "x",
	}

	jsonPayload := &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"id": "1", "status": "active"},
		},
	}

	ctx := &engine.ExecutionContext{Payload: jsonPayload}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out := ctx.Payload.(*payload.JSON)

	if len(out.Items) != 0 {
		t.Fatalf("expected no items, got %d", len(out.Items))
	}
}
