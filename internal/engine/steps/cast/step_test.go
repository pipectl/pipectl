package cast

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "cast" {
		t.Fatalf("expected step name %q, got %q", "cast", step.Name())
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
	if step.Supports(&payload.CSV{}) {
		t.Fatal("expected step not to support CSV payload")
	}
}

func TestExecuteCastsJSONFields(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"age": {
				Type: "int",
			},
			"pi": {
				Type: "int",
			},
			"score": {
				Type: "float",
			},
			"attempts": {
				Type: "float",
			},
			"active": {
				Type: "bool",
			},
			"enabled": {
				Type: "bool",
			},
			"disabled": {
				Type: "bool",
			},
			"created_at": {
				Type:   "time",
				Format: "2006-01-02",
			},
			"visits": {
				Type: "string",
			},
			"verified": {
				Type: "string",
			},
			"ratio": {
				Type: "int",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"age":        "42",
					"pi":         "3.14",
					"score":      "98.5",
					"attempts":   7,
					"active":     "YES",
					"enabled":    1,
					"disabled":   0.0,
					"created_at": "2026-03-22",
					"visits":     12.0,
					"verified":   true,
					"ratio":      9.9,
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

	record := out.Items[0]
	if got := record["age"]; got != 42 {
		t.Fatalf("unexpected int cast: %#v", got)
	}
	if got := record["pi"]; got != 3 {
		t.Fatalf("unexpected decimal string to int cast: %#v", got)
	}
	if got := record["score"]; got != 98.5 {
		t.Fatalf("unexpected float cast: %#v", got)
	}
	if got := record["attempts"]; got != 7.0 {
		t.Fatalf("unexpected int to float cast: %#v", got)
	}
	if got := record["active"]; got != true {
		t.Fatalf("unexpected bool cast: %#v", got)
	}
	if got := record["enabled"]; got != true {
		t.Fatalf("unexpected numeric true bool cast: %#v", got)
	}
	if got := record["disabled"]; got != false {
		t.Fatalf("unexpected numeric false bool cast: %#v", got)
	}

	createdAt, ok := record["created_at"].(time.Time)
	if !ok {
		t.Fatalf("expected created_at to be time.Time, got %T", record["created_at"])
	}
	expectedTime := time.Date(2026, time.March, 22, 0, 0, 0, 0, time.UTC)
	if !createdAt.Equal(expectedTime) {
		t.Fatalf("unexpected time cast: got %v want %v", createdAt, expectedTime)
	}

	if got := record["visits"]; got != "12" {
		t.Fatalf("unexpected number to string cast: %#v", got)
	}
	if got := record["verified"]; got != "true" {
		t.Fatalf("unexpected bool to string cast: %#v", got)
	}
	if got := record["ratio"]; got != 9 {
		t.Fatalf("unexpected float to int cast: %#v", got)
	}
}

func TestExecuteCastsNestedJSONFields(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"headers.Content-Length": {
				Type: "int",
			},
			"items": {
				Type: "int",
			},
			"flags[1].active": {
				Type: "bool",
			},
			"scores[0]": {
				Type: "float",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"headers": map[string]interface{}{
						"Content-Length": "232",
					},
					"items":  []interface{}{"1", "0"},
					"flags":  []interface{}{map[string]interface{}{"active": "no"}, map[string]interface{}{"active": "yes"}},
					"scores": []interface{}{"10.5", "11.5"},
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	record := ctx.Payload.(*payload.JSON).Items[0]
	headers := record["headers"].(map[string]interface{})
	if got := headers["Content-Length"]; got != 232 {
		t.Fatalf("unexpected nested int cast: %#v", got)
	}

	items := record["items"].([]interface{})
	if got := items[0]; got != 1 {
		t.Fatalf("unexpected array cast first value: %#v", got)
	}
	if got := items[1]; got != 0 {
		t.Fatalf("unexpected array cast second value: %#v", got)
	}
	flags := record["flags"].([]interface{})
	secondFlag := flags[1].(map[string]interface{})
	if got := secondFlag["active"]; got != true {
		t.Fatalf("unexpected nested bool cast: %#v", got)
	}
	scores := record["scores"].([]interface{})
	if got := scores[0]; got != 10.5 {
		t.Fatalf("unexpected indexed array cast: %#v", got)
	}
}

func TestExecuteCastsNestedArrayField(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"payload.metrics": {
				Type: "float",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"payload": map[string]interface{}{
						"metrics": []interface{}{"1.25", "2.5", "3"},
					},
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	record := ctx.Payload.(*payload.JSON).Items[0]
	payloadObject := record["payload"].(map[string]interface{})
	metrics := payloadObject["metrics"].([]interface{})
	expected := []interface{}{1.25, 2.5, 3.0}
	if !reflect.DeepEqual(metrics, expected) {
		t.Fatalf("unexpected nested array cast:\nexpected: %#v\ngot: %#v", expected, metrics)
	}
}

func TestExecuteFailsWhenWholeArrayCastContainsInvalidValue(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"items": {
				Type: "int",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"items": []interface{}{"1", "oops"},
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for invalid array cast value")
	}
	if !strings.Contains(err.Error(), `cast field "items" in record 1: array index 1: cannot cast "oops" to int`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteKeepsIndexedPathCasting(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"items[1].active": {
				Type: "bool",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"items": []interface{}{
						map[string]interface{}{"active": "no"},
						map[string]interface{}{"active": "yes"},
					},
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	record := ctx.Payload.(*payload.JSON).Items[0]
	items := record["items"].([]interface{})
	second := items[1].(map[string]interface{})
	if got := second["active"]; got != true {
		t.Fatalf("unexpected nested bool cast: %#v", got)
	}
}

func TestExecuteCastsJSONLFieldsWithCustomBoolValues(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"active": {
				Type:        "bool",
				TrueValues:  []string{"enabled"},
				FalseValues: []string{"disabled"},
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"active": "enabled"},
				{"active": "disabled"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out := ctx.Payload.(*payload.JSONL)
	expected := []map[string]interface{}{
		{"active": true},
		{"active": false},
	}
	if !reflect.DeepEqual(out.Items, expected) {
		t.Fatalf("unexpected JSONL data:\nexpected: %#v\ngot: %#v", expected, out.Items)
	}
}

func TestExecuteFailsWhenFieldMissing(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"age": {
				Type: "int",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{"name": "Alice"},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for missing field")
	}
	if !strings.Contains(err.Error(), `cast field "age" in record 1: path "age" missing key "age"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFailsWhenNestedPathMissing(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"headers.Content-Length": {
				Type: "int",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"headers": map[string]interface{}{},
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for missing nested path")
	}
	if !strings.Contains(err.Error(), `cast field "headers.Content-Length" in record 1: path "headers.Content-Length" missing key "Content-Length"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFailsWhenArrayIndexOutOfRange(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"items[2].active": {
				Type: "bool",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"items": []interface{}{
						map[string]interface{}{"active": "yes"},
					},
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for out of range array index")
	}
	if !strings.Contains(err.Error(), `cast field "items[2].active" in record 1: path "items[2].active" index 2 out of range`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFailsWhenPathSyntaxInvalid(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"items[].active": {
				Type: "bool",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"items": []interface{}{
						map[string]interface{}{"active": "yes"},
					},
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for invalid field path")
	}
	if !strings.Contains(err.Error(), `cast field "items[].active" in record 1: invalid field path "items[].active"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFailsWhenCastIsInvalid(t *testing.T) {
	step := &Step{
		Fields: map[string]Field{
			"active": {
				Type: "bool",
			},
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{"active": "maybe"},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for invalid bool cast")
	}
	if !strings.Contains(err.Error(), `cast field "active" in record 1: cannot cast "maybe" to bool`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFailsForUnsupportedPayload(t *testing.T) {
	step := &Step{}
	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{{"age"}, {"42"}},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for unsupported payload")
	}
	if !strings.Contains(err.Error(), "unsupported payload type") {
		t.Fatalf("unexpected error: %v", err)
	}
}
