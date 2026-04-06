package filter

import (
	"reflect"
	"strconv"
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

	if !step.Supports(&payload.JSON{}) {
		t.Fatal("expected step to support JSON payload")
	}

	if !step.Supports(&payload.JSONL{}) {
		t.Fatal("expected step to support JSONL payload")
	}
}

func TestExecuteFiltersCSVRows(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		value    string
		rows     [][]string
		expected [][]string
	}{
		{
			name:  "equals",
			op:    OpEquals,
			value: "active",
			rows: [][]string{
				{"id", "status"},
				{"1", "active"},
				{"2", "inactive"},
				{"3", "active"},
			},
			expected: [][]string{
				{"id", "status"},
				{"1", "active"},
				{"3", "active"},
			},
		},
		{
			name:  "not-equals",
			op:    OpNotEquals,
			value: "inactive",
			rows: [][]string{
				{"id", "status"},
				{"1", "active"},
				{"2", "inactive"},
				{"3", "active"},
			},
			expected: [][]string{
				{"id", "status"},
				{"1", "active"},
				{"3", "active"},
			},
		},
		{
			name:  "contains",
			op:    OpContains,
			value: "example",
			rows: [][]string{
				{"id", "email"},
				{"1", "alice@example.com"},
				{"2", "bob@other.org"},
				{"3", "carol@example.com"},
			},
			expected: [][]string{
				{"id", "email"},
				{"1", "alice@example.com"},
				{"3", "carol@example.com"},
			},
		},
		{
			name:  "starts-with",
			op:    OpStartsWith,
			value: "alice",
			rows: [][]string{
				{"id", "email"},
				{"1", "alice@example.com"},
				{"2", "bob@example.com"},
			},
			expected: [][]string{
				{"id", "email"},
				{"1", "alice@example.com"},
			},
		},
		{
			name:  "greater-than",
			op:    OpGreaterThan,
			value: "30",
			rows: [][]string{
				{"id", "age"},
				{"1", "25"},
				{"2", "35"},
				{"3", "30"},
			},
			expected: [][]string{
				{"id", "age"},
				{"2", "35"},
			},
		},
		{
			name:  "less-than",
			op:    OpLessThan,
			value: "30",
			rows: [][]string{
				{"id", "age"},
				{"1", "25"},
				{"2", "35"},
				{"3", "30"},
			},
			expected: [][]string{
				{"id", "age"},
				{"1", "25"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var numericValue float64
			if tt.op == OpGreaterThan || tt.op == OpLessThan {
				numericValue, _ = strconv.ParseFloat(tt.value, 64)
			}
			step := &Step{Field: "status", Op: tt.op, Value: tt.value, NumericValue: numericValue}
			switch tt.op {
			case OpContains, OpStartsWith:
				step.Field = "email"
			case OpGreaterThan, OpLessThan:
				step.Field = "age"
			}

			ctx := &engine.ExecutionContext{Payload: &payload.CSV{Rows: tt.rows}}

			if err := step.Execute(ctx); err != nil {
				t.Fatalf("execute returned error: %v", err)
			}

			out := ctx.Payload.(*payload.CSV)
			if !reflect.DeepEqual(out.Rows, tt.expected) {
				t.Fatalf("unexpected rows:\nexpected: %#v\ngot: %#v", tt.expected, out.Rows)
			}
		})
	}
}

func TestExecuteFiltersJSONRecords(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		field    string
		value    string
		items    []map[string]interface{}
		expected []map[string]interface{}
	}{
		{
			name:  "equals",
			op:    OpEquals,
			field: "status",
			value: "active",
			items: []map[string]interface{}{
				{"id": "1", "status": "active"},
				{"id": "2", "status": "inactive"},
				{"id": "3", "status": "active"},
			},
			expected: []map[string]interface{}{
				{"id": "1", "status": "active"},
				{"id": "3", "status": "active"},
			},
		},
		{
			name:  "not-equals",
			op:    OpNotEquals,
			field: "status",
			value: "inactive",
			items: []map[string]interface{}{
				{"id": "1", "status": "active"},
				{"id": "2", "status": "inactive"},
				{"id": "3", "status": "active"},
			},
			expected: []map[string]interface{}{
				{"id": "1", "status": "active"},
				{"id": "3", "status": "active"},
			},
		},
		{
			name:  "contains",
			op:    OpContains,
			field: "email",
			value: "example",
			items: []map[string]interface{}{
				{"id": "1", "email": "alice@example.com"},
				{"id": "2", "email": "bob@other.org"},
				{"id": "3", "email": "carol@example.com"},
			},
			expected: []map[string]interface{}{
				{"id": "1", "email": "alice@example.com"},
				{"id": "3", "email": "carol@example.com"},
			},
		},
		{
			name:  "starts-with",
			op:    OpStartsWith,
			field: "email",
			value: "alice",
			items: []map[string]interface{}{
				{"id": "1", "email": "alice@example.com"},
				{"id": "2", "email": "bob@example.com"},
			},
			expected: []map[string]interface{}{
				{"id": "1", "email": "alice@example.com"},
			},
		},
		{
			name:  "numeric field with equals integer",
			op:    OpEquals,
			field: "count",
			value: "5",
			items: []map[string]interface{}{
				{"id": "1", "count": float64(5)},
				{"id": "2", "count": float64(3)},
			},
			expected: []map[string]interface{}{
				{"id": "1", "count": float64(5)},
			},
		},
		{
			name:  "numeric field with equals decimal notation matches integer",
			op:    OpEquals,
			field: "count",
			value: "5.0",
			items: []map[string]interface{}{
				{"id": "1", "count": float64(5)},
				{"id": "2", "count": float64(3)},
			},
			expected: []map[string]interface{}{
				{"id": "1", "count": float64(5)},
			},
		},
		{
			name:  "numeric field with equals scientific notation",
			op:    OpEquals,
			field: "count",
			value: "1e6",
			items: []map[string]interface{}{
				{"id": "1", "count": float64(1000000)},
				{"id": "2", "count": float64(3)},
			},
			expected: []map[string]interface{}{
				{"id": "1", "count": float64(1000000)},
			},
		},
		{
			name:  "missing field excluded",
			op:    OpEquals,
			field: "missing",
			value: "x",
			items: []map[string]interface{}{
				{"id": "1", "status": "active"},
			},
			expected: []map[string]interface{}{},
		},
		{
			name:  "greater-than with float64 field",
			op:    OpGreaterThan,
			field: "age",
			value: "30",
			items: []map[string]interface{}{
				{"id": "1", "age": float64(25)},
				{"id": "2", "age": float64(35)},
				{"id": "3", "age": float64(30)},
			},
			expected: []map[string]interface{}{
				{"id": "2", "age": float64(35)},
			},
		},
		{
			name:  "less-than with float64 field",
			op:    OpLessThan,
			field: "age",
			value: "30",
			items: []map[string]interface{}{
				{"id": "1", "age": float64(25)},
				{"id": "2", "age": float64(35)},
				{"id": "3", "age": float64(30)},
			},
			expected: []map[string]interface{}{
				{"id": "1", "age": float64(25)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var numericValue float64
			if tt.op == OpGreaterThan || tt.op == OpLessThan {
				numericValue, _ = strconv.ParseFloat(tt.value, 64)
			}
			step := &Step{Field: tt.field, Op: tt.op, Value: tt.value, NumericValue: numericValue}

			jsonPayload := &payload.JSON{
				Shape: payload.JSONArrayShape,
				Items: tt.items,
			}

			ctx := &engine.ExecutionContext{Payload: jsonPayload}

			if err := step.Execute(ctx); err != nil {
				t.Fatalf("execute returned error: %v", err)
			}

			out := ctx.Payload.(*payload.JSON)

			if tt.expected == nil {
				tt.expected = []map[string]interface{}{}
			}
			if len(out.Items) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(out.Items, tt.expected) {
				t.Fatalf("unexpected items:\nexpected: %#v\ngot: %#v", tt.expected, out.Items)
			}
		})
	}
}

func TestExecuteReturnsErrorForNonNumericFieldWithGreaterThan(t *testing.T) {
	step := &Step{Field: "status", Op: OpGreaterThan, NumericValue: 10}

	ctx := &engine.ExecutionContext{Payload: &payload.JSON{
		Shape: payload.JSONArrayShape,
		Items: []map[string]interface{}{
			{"status": "active"},
		},
	}}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for non-numeric field value")
	}
	if !strings.Contains(err.Error(), `field "status"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFiltersJSONLRecords(t *testing.T) {
	step := &Step{
		Field: "status",
		Op:    OpEquals,
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
		t.Fatalf("unexpected items:\nexpected: %#v\ngot: %#v", expected, out.Items)
	}
}
