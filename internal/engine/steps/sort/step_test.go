package sort

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
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
	step := &Step{Keys: []Key{{Field: "name", Direction: DirectionAsc}}}
	ctx := &engine.ExecutionContext{Payload: &payload.JSON{
		Shape: payload.JSONObjectShape,
		Items: []map[string]interface{}{{"name": "alice"}},
	}}

	if err := step.Execute(ctx); err == nil {
		t.Fatal("expected error for JSON object shape")
	}
}

func TestExecuteSortsJSONArrayAsc(t *testing.T) {
	step := &Step{Keys: []Key{{Field: "name", Direction: DirectionAsc}}}
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
	step := &Step{Keys: []Key{{Field: "name", Direction: DirectionDesc}}}
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
	step := &Step{Keys: []Key{{Field: "age", Direction: DirectionAsc}}}
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
	step := &Step{Keys: []Key{{Field: "name", Direction: DirectionAsc}}}
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
	step := &Step{Keys: []Key{{Field: "score", Direction: DirectionDesc}}}
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
	step := &Step{Keys: []Key{{Field: "name", Direction: DirectionAsc}}}
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
	step := &Step{Keys: []Key{{Field: "age", Direction: DirectionAsc}}}
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
	step := &Step{Keys: []Key{{Field: "name", Direction: DirectionAsc}}}
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

func TestExecuteErrorsOnMissingCSVField(t *testing.T) {
	step := &Step{Keys: []Key{{Field: "missing", Direction: DirectionAsc}}}
	ctx := &engine.ExecutionContext{Payload: &payload.CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1", "Alice"},
		},
	}}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for missing CSV field, got nil")
	}
	if !strings.Contains(err.Error(), "not found in CSV headers") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestExecuteCSVHeaderPreservedWhenEmpty(t *testing.T) {
	step := &Step{Keys: []Key{{Field: "name", Direction: DirectionAsc}}}
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

func TestExecuteNullPlacementJSON(t *testing.T) {
	tests := []struct {
		name      string
		direction string
		nulls     string
		expected  []string
	}{
		{name: "asc default nulls last", direction: DirectionAsc, nulls: "", expected: []string{"alice", "carol", "missing", "null"}},
		{name: "asc nulls last", direction: DirectionAsc, nulls: NullsLast, expected: []string{"alice", "carol", "missing", "null"}},
		{name: "asc nulls first", direction: DirectionAsc, nulls: NullsFirst, expected: []string{"missing", "null", "alice", "carol"}},
		{name: "desc nulls last", direction: DirectionDesc, nulls: NullsLast, expected: []string{"carol", "alice", "missing", "null"}},
		{name: "desc nulls first", direction: DirectionDesc, nulls: NullsFirst, expected: []string{"missing", "null", "carol", "alice"}},
	}

	for _, tt := range tests {
		for _, format := range []string{"json", "jsonl"} {
			t.Run(tt.name+"/"+format, func(t *testing.T) {
				items := []map[string]interface{}{
					{"id": "carol", "name": "carol"},
					{"id": "missing"},
					{"id": "alice", "name": "alice"},
					{"id": "null", "name": nil},
				}
				var p payload.Payload = &payload.JSON{Shape: payload.JSONArrayShape, Items: items}
				if format == "jsonl" {
					p = &payload.JSONL{Items: items}
				}

				step := &Step{Keys: []Key{{Field: "name", Direction: tt.direction, Nulls: tt.nulls}}}
				ctx := &engine.ExecutionContext{Payload: p}
				if err := step.Execute(ctx); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				ids := make([]string, len(items))
				for i, item := range items {
					ids[i] = item["id"].(string)
				}
				if !reflect.DeepEqual(ids, tt.expected) {
					t.Fatalf("unexpected order: got %v want %v", ids, tt.expected)
				}
			})
		}
	}
}

func TestExecuteNullPlacementCSV(t *testing.T) {
	tests := []struct {
		name      string
		direction string
		nulls     string
		expected  []string
	}{
		{name: "asc default nulls last", direction: DirectionAsc, nulls: "", expected: []string{"3", "1", "2", "4"}},
		{name: "asc nulls first", direction: DirectionAsc, nulls: NullsFirst, expected: []string{"2", "4", "3", "1"}},
		{name: "desc nulls last", direction: DirectionDesc, nulls: NullsLast, expected: []string{"1", "3", "2", "4"}},
		{name: "desc nulls first", direction: DirectionDesc, nulls: NullsFirst, expected: []string{"2", "4", "1", "3"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := &Step{Keys: []Key{{Field: "name", Direction: tt.direction, Nulls: tt.nulls}}}
			ctx := &engine.ExecutionContext{Payload: &payload.CSV{
				Rows: [][]string{
					{"id", "name"},
					{"1", "carol"},
					{"2", ""},
					{"3", "alice"},
					{"4", ""},
				},
			}}

			if err := step.Execute(ctx); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			out := ctx.Payload.(*payload.CSV)
			if out.Rows[0][0] != "id" {
				t.Fatalf("header row moved: got %v", out.Rows[0])
			}
			ids := []string{out.Rows[1][0], out.Rows[2][0], out.Rows[3][0], out.Rows[4][0]}
			if !reflect.DeepEqual(ids, tt.expected) {
				t.Fatalf("unexpected order: got %v want %v", ids, tt.expected)
			}
		})
	}
}

func TestExecuteMultiKeyJSON(t *testing.T) {
	tests := []struct {
		name     string
		keys     []Key
		expected []string
	}{
		{
			name:     "secondary key breaks ties",
			keys:     []Key{{Field: "country", Direction: DirectionAsc}, {Field: "age", Direction: DirectionAsc}},
			expected: []string{"eve", "alice", "carol", "dave", "bob", "frank"},
		},
		{
			name:     "mixed directions",
			keys:     []Key{{Field: "country", Direction: DirectionAsc}, {Field: "age", Direction: DirectionDesc}},
			expected: []string{"carol", "alice", "eve", "dave", "bob", "frank"},
		},
		{
			name:     "per-key nulls first on secondary key",
			keys:     []Key{{Field: "country", Direction: DirectionDesc}, {Field: "age", Direction: DirectionAsc, Nulls: NullsFirst}},
			expected: []string{"frank", "bob", "dave", "eve", "alice", "carol"},
		},
		{
			name:     "full tie keeps input order",
			keys:     []Key{{Field: "country", Direction: DirectionAsc}, {Field: "team", Direction: DirectionAsc}},
			expected: []string{"alice", "carol", "eve", "dave", "bob", "frank"},
		},
	}

	for _, tt := range tests {
		for _, format := range []string{"json", "jsonl"} {
			t.Run(tt.name+"/"+format, func(t *testing.T) {
				items := []map[string]interface{}{
					{"id": "alice", "country": "au", "age": float64(25), "team": "x"},
					{"id": "bob", "country": "us", "age": float64(17), "team": "x"},
					{"id": "carol", "country": "au", "age": float64(32), "team": "x"},
					{"id": "dave", "country": "nz", "age": float64(45), "team": "x"},
					{"id": "eve", "country": "au", "age": float64(19), "team": "x"},
					{"id": "frank", "country": "us", "team": "x"},
				}
				var p payload.Payload = &payload.JSON{Shape: payload.JSONArrayShape, Items: items}
				if format == "jsonl" {
					p = &payload.JSONL{Items: items}
				}

				step := &Step{Keys: tt.keys}
				ctx := &engine.ExecutionContext{Payload: p}
				if err := step.Execute(ctx); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				ids := make([]string, len(items))
				for i, item := range items {
					ids[i] = item["id"].(string)
				}
				if !reflect.DeepEqual(ids, tt.expected) {
					t.Fatalf("unexpected order: got %v want %v", ids, tt.expected)
				}
			})
		}
	}
}

func TestExecuteMultiKeyCSV(t *testing.T) {
	tests := []struct {
		name     string
		keys     []Key
		expected []string
	}{
		{
			name:     "secondary key breaks ties",
			keys:     []Key{{Field: "plan", Direction: DirectionAsc}, {Field: "name", Direction: DirectionAsc}},
			expected: []string{"4", "1", "3", "2"},
		},
		{
			name:     "mixed directions with secondary nulls first",
			keys:     []Key{{Field: "plan", Direction: DirectionDesc}, {Field: "name", Direction: DirectionDesc, Nulls: NullsFirst}},
			expected: []string{"2", "3", "1", "4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := &Step{Keys: tt.keys}
			ctx := &engine.ExecutionContext{Payload: &payload.CSV{
				Rows: [][]string{
					{"id", "plan", "name"},
					{"1", "pro", "carol"},
					{"2", "starter", ""},
					{"3", "pro", "dave"},
					{"4", "enterprise", "bob"},
				},
			}}

			if err := step.Execute(ctx); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			out := ctx.Payload.(*payload.CSV)
			if out.Rows[0][0] != "id" {
				t.Fatalf("header row moved: got %v", out.Rows[0])
			}
			ids := []string{out.Rows[1][0], out.Rows[2][0], out.Rows[3][0], out.Rows[4][0]}
			if !reflect.DeepEqual(ids, tt.expected) {
				t.Fatalf("unexpected order: got %v want %v", ids, tt.expected)
			}
		})
	}
}

func TestExecuteMultiKeyErrorsOnMissingSecondCSVField(t *testing.T) {
	step := &Step{Keys: []Key{{Field: "name", Direction: DirectionAsc}, {Field: "missing", Direction: DirectionAsc}}}
	ctx := &engine.ExecutionContext{Payload: &payload.CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1", "Alice"},
			{"2", "Bob"},
		},
	}}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for missing CSV field, got nil")
	}
	if !strings.Contains(err.Error(), "not found in CSV headers") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
