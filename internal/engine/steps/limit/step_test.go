package limit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "limit" {
		t.Fatalf("expected step name %q, got %q", "limit", step.Name())
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

func TestExecuteLimitsJSONRecords(t *testing.T) {
	step := &Step{Count: 2}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{"id": 1},
				{"id": 2},
				{"id": 3},
				{"id": 4},
			},
			Shape: payload.JSONArrayShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 2 {
		t.Fatalf("expected 2 records after limit, got %d", got)
	}
	assertContains(t, buf.String(), "limited 4 records to 2")
}

func TestExecuteLimitsJSONLRecords(t *testing.T) {
	step := &Step{Count: 1}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"id": 1},
				{"id": 2},
				{"id": 3},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 1 {
		t.Fatalf("expected 1 record after limit, got %d", got)
	}
	assertContains(t, buf.String(), "limited 3 records to 1")
}

func TestExecuteLimitsCSVRecords(t *testing.T) {
	step := &Step{Count: 2}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "name"},
				{"1", "Alice"},
				{"2", "Bob"},
				{"3", "Carol"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 2 {
		t.Fatalf("expected 2 records after limit, got %d", got)
	}
	// Header row must be preserved
	csvPayload := ctx.Payload.(*payload.CSV)
	if len(csvPayload.Rows[0]) != 2 || csvPayload.Rows[0][0] != "id" {
		t.Fatalf("header row not preserved: %v", csvPayload.Rows[0])
	}
	assertContains(t, buf.String(), "limited 3 records to 2")
}

func TestExecuteDoesNotTruncateWhenUnderLimit(t *testing.T) {
	step := &Step{Count: 100}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{"id": 1},
				{"id": 2},
			},
			Shape: payload.JSONArrayShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 2 {
		t.Fatalf("expected 2 records (unchanged), got %d", got)
	}
	output := buf.String()
	assertContains(t, output, "limit of 100 not reached")
	assertNotContains(t, output, "limited")
}

func TestExecuteExactlyAtLimit(t *testing.T) {
	step := &Step{Count: 3}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"id": 1},
				{"id": 2},
				{"id": 3},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 3 {
		t.Fatalf("expected 3 records (unchanged), got %d", got)
	}
	assertContains(t, buf.String(), "limit of 3 not reached")
}

func assertContains(t *testing.T, value, expected string) {
	t.Helper()
	if !strings.Contains(value, expected) {
		t.Fatalf("expected output to contain %q, got %q", expected, value)
	}
}

func assertNotContains(t *testing.T, value, expected string) {
	t.Helper()
	if strings.Contains(value, expected) {
		t.Fatalf("did not expect output to contain %q, got %q", expected, value)
	}
}
