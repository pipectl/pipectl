package sort

import (
	"reflect"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "sort" {
		t.Fatalf("expected step name %q, got %q", "sort", step.Name())
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

func TestExecuteRejectsJSONObjectShape(t *testing.T) {
	step := &Step{Field: "name", Direction: DirectionAsc}
	ctx := &engine.ExecutionContext{Payload: &payload.JSON{
		Shape: payload.JSONObjectShape,
		Items: []map[string]interface{}{{"name": "alice"}},
	}}

	if err := step.Execute(ctx); err == nil {
		t.Fatal("expected error for JSON object shape")
	}
}

func TestExecuteSortsJSONArrayAsc(t *testing.T) {
	step := &Step{Field: "name", Direction: DirectionAsc}
	ctx := &engine.ExecutionContext{Payload: &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"name": "carol"},
			{"name": "alice"},
			{"name": "bob"},
		},
	}}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := ctx.Payload.(*payload.JSON)
	names := []string{out.Items[0]["name"].(string), out.Items[1]["name"].(string), out.Items[2]["name"].(string)}
	expected := []string{"alice", "bob", "carol"}
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("unexpected order: got %v want %v", names, expected)
	}
}

func TestExecuteSortsJSONArrayDesc(t *testing.T) {
	step := &Step{Field: "name", Direction: DirectionDesc}
	ctx := &engine.ExecutionContext{Payload: &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"name": "carol"},
			{"name": "alice"},
			{"name": "bob"},
		},
	}}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := ctx.Payload.(*payload.JSON)
	names := []string{out.Items[0]["name"].(string), out.Items[1]["name"].(string), out.Items[2]["name"].(string)}
	expected := []string{"carol", "bob", "alice"}
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("unexpected order: got %v want %v", names, expected)
	}
}

func TestExecuteSortsJSONNumericFieldAsc(t *testing.T) {
	step := &Step{Field: "age", Direction: DirectionAsc}
	ctx := &engine.ExecutionContext{Payload: &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"name": "carol", "age": float64(30)},
			{"name": "alice", "age": float64(25)},
			{"name": "bob", "age": float64(35)},
		},
	}}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := ctx.Payload.(*payload.JSON)
	ages := []float64{out.Items[0]["age"].(float64), out.Items[1]["age"].(float64), out.Items[2]["age"].(float64)}
	expected := []float64{25, 30, 35}
	if !reflect.DeepEqual(ages, expected) {
		t.Fatalf("unexpected order: got %v want %v", ages, expected)
	}
}

func TestExecuteSortsJSONNullsLast(t *testing.T) {
	step := &Step{Field: "name", Direction: DirectionAsc}
	ctx := &engine.ExecutionContext{Payload: &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"name": "carol"},
			{"id": "2"}, // missing name field
			{"name": "alice"},
		},
	}}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := ctx.Payload.(*payload.JSON)
	if out.Items[0]["name"] != "alice" || out.Items[1]["name"] != "carol" {
		t.Fatalf("expected alice, carol first: got %v", out.Items)
	}
	if _, exists := out.Items[2]["name"]; exists {
		t.Fatalf("expected missing-field record last")
	}
}

func TestExecuteSortsJSONLRecords(t *testing.T) {
	step := &Step{Field: "score", Direction: DirectionDesc}
	ctx := &engine.ExecutionContext{Payload: &payload.JSONL{
		Items: []map[string]interface{}{
			{"id": "1", "score": float64(70)},
			{"id": "2", "score": float64(90)},
			{"id": "3", "score": float64(80)},
		},
	}}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := ctx.Payload.(*payload.JSONL)
	ids := []string{out.Items[0]["id"].(string), out.Items[1]["id"].(string), out.Items[2]["id"].(string)}
	expected := []string{"2", "3", "1"}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("unexpected order: got %v want %v", ids, expected)
	}
}

func TestExecuteSortsCSVRowsAsc(t *testing.T) {
	step := &Step{Field: "name", Direction: DirectionAsc}
	ctx := &engine.ExecutionContext{Payload: &payload.CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1", "carol"},
			{"2", "alice"},
			{"3", "bob"},
		},
	}}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := ctx.Payload.(*payload.CSV)
	if out.Rows[0][1] != "name" {
		t.Fatalf("header row moved: got %v", out.Rows[0])
	}
	names := []string{out.Rows[1][1], out.Rows[2][1], out.Rows[3][1]}
	expected := []string{"alice", "bob", "carol"}
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("unexpected order: got %v want %v", names, expected)
	}
}

func TestExecuteSortsCSVNumericFieldAsc(t *testing.T) {
	step := &Step{Field: "age", Direction: DirectionAsc}
	ctx := &engine.ExecutionContext{Payload: &payload.CSV{
		Rows: [][]string{
			{"name", "age"},
			{"carol", "30"},
			{"alice", "25"},
			{"bob", "9"},
		},
	}}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := ctx.Payload.(*payload.CSV)
	ages := []string{out.Rows[1][1], out.Rows[2][1], out.Rows[3][1]}
	expected := []string{"9", "25", "30"}
	if !reflect.DeepEqual(ages, expected) {
		t.Fatalf("unexpected order: got %v want %v (should sort numerically, not lexically)", ages, expected)
	}
}

func TestExecuteSortsCSVNullsLast(t *testing.T) {
	step := &Step{Field: "name", Direction: DirectionAsc}
	ctx := &engine.ExecutionContext{Payload: &payload.CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1", "carol"},
			{"2", ""},
			{"3", "alice"},
		},
	}}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := ctx.Payload.(*payload.CSV)
	if out.Rows[1][1] != "alice" || out.Rows[2][1] != "carol" {
		t.Fatalf("expected alice, carol first: got %v", out.Rows)
	}
	if out.Rows[3][1] != "" {
		t.Fatalf("expected empty value last: got %q", out.Rows[3][1])
	}
}

func TestExecuteCSVHeaderPreservedWhenEmpty(t *testing.T) {
	step := &Step{Field: "name", Direction: DirectionAsc}
	ctx := &engine.ExecutionContext{Payload: &payload.CSV{
		Rows: [][]string{{"id", "name"}},
	}}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := ctx.Payload.(*payload.CSV)
	if len(out.Rows) != 1 || out.Rows[0][1] != "name" {
		t.Fatalf("unexpected rows: %v", out.Rows)
	}
}
