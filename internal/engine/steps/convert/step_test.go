package convert

import (
	"testing"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "convert" {
		t.Fatalf("expected step name %q, got %q", "convert", step.Name())
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

func TestExecuteConvertsJSONToCSV(t *testing.T) {
	step := &Step{Format: payload.CSVType}
	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{"id": float64(1), "name": "alice"},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	csvPayload, ok := ctx.Payload.(*payload.CSV)
	if !ok {
		t.Fatalf("expected CSV payload, got %T", ctx.Payload)
	}

	if len(csvPayload.Rows) != 2 {
		t.Fatalf("unexpected row count: got %d want %d", len(csvPayload.Rows), 2)
	}
	if csvPayload.Rows[0][0] != "id" || csvPayload.Rows[0][1] != "name" {
		t.Fatalf("unexpected header row: %#v", csvPayload.Rows[0])
	}
	if csvPayload.Rows[1][0] != "1" || csvPayload.Rows[1][1] != "alice" {
		t.Fatalf("unexpected data row: %#v", csvPayload.Rows[1])
	}
}

func TestExecuteConvertsCSVToJSONL(t *testing.T) {
	step := &Step{Format: payload.JSONLType}
	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "nested.name", "values"},
				{"1", "alice", `["quick","brown"]`},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	jsonlPayload, ok := ctx.Payload.(*payload.JSONL)
	if !ok {
		t.Fatalf("expected JSONL payload, got %T", ctx.Payload)
	}

	if len(jsonlPayload.Items) != 1 {
		t.Fatalf("unexpected record count: got %d want %d", len(jsonlPayload.Items), 1)
	}
	if jsonlPayload.Items[0]["id"] != "1" {
		t.Fatalf("unexpected record: %#v", jsonlPayload.Items[0])
	}

	nested, ok := jsonlPayload.Items[0]["nested"].(map[string]interface{})
	if !ok || nested["name"] != "alice" {
		t.Fatalf("unexpected nested record: %#v", jsonlPayload.Items[0]["nested"])
	}

	values, ok := jsonlPayload.Items[0]["values"].([]interface{})
	if !ok || len(values) != 2 || values[0] != "quick" || values[1] != "brown" {
		t.Fatalf("unexpected array value: %#v", jsonlPayload.Items[0]["values"])
	}
}

func TestExecuteReturnsErrorForInvalidCSV(t *testing.T) {
	step := &Step{Format: payload.JSONType}
	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "name"},
				{"1"},
			},
		},
	}

	if err := step.Execute(ctx); err == nil {
		t.Fatal("expected execute to return error")
	}
}
