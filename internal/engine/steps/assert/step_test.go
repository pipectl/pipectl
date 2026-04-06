package assert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "assert" {
		t.Fatalf("expected step name %q, got %q", "assert", step.Name())
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

func TestExecuteSucceedsForCSV(t *testing.T) {
	min := 1
	max := 3
	equal := 2
	step := &Step{
		MinRecords:   &min,
		MaxRecords:   &max,
		RecordsEqual: &equal,
		FieldExists:  "email",
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"name", "email"},
				{"Alice", "alice@example.com"},
				{"Bob", "bob@example.com"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteSucceedsForJSON(t *testing.T) {
	min := 1
	max := 1
	equal := 1
	step := &Step{
		MinRecords:   &min,
		MaxRecords:   &max,
		RecordsEqual: &equal,
		FieldExists:  "email",
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"email": "alice@example.com",
					"name":  "Alice",
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteSucceedsForJSONL(t *testing.T) {
	min := 1
	max := 2
	equal := 2
	step := &Step{
		MinRecords:   &min,
		MaxRecords:   &max,
		RecordsEqual: &equal,
		FieldExists:  "email",
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"email": "alice@example.com"},
				{"email": "bob@example.com"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteFailsWhenRecordCountBelowMinimum(t *testing.T) {
	min := 2
	step := &Step{
		MinRecords: &min,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Items: []map[string]interface{}{{"email": "alice@example.com"}},
			Shape: payload.JSONObjectShape,
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error when record count is below minimum")
	}
	if !strings.Contains(err.Error(), "less than minimum 2") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFailsWhenRecordCountAboveMaximum(t *testing.T) {
	max := 1
	step := &Step{
		MaxRecords: &max,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id"},
				{"1"},
				{"2"},
			},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error when record count is above maximum")
	}
	if !strings.Contains(err.Error(), "greater than maximum 1") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFailsWhenFieldDoesNotExist(t *testing.T) {
	step := &Step{
		FieldExists: "email",
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"name", "country"},
				{"Alice", "AU"},
			},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error when required field does not exist")
	}
	if !strings.Contains(err.Error(), `field "email" does not exist`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFailsWhenRecordCountDoesNotEqualExpected(t *testing.T) {
	equal := 3
	step := &Step{
		RecordsEqual: &equal,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id"},
				{"1"},
				{"2"},
			},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error when record count does not equal expected")
	}
	if !strings.Contains(err.Error(), "is not equal to expected 3") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteLogsAssertions(t *testing.T) {
	min := 1
	max := 2
	equal := 2
	step := &Step{
		MinRecords:   &min,
		MaxRecords:   &max,
		RecordsEqual: &equal,
		FieldExists:  "email",
	}

	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.CSV{
			Rows: [][]string{
				{"name", "email"},
				{"Alice", "alice@example.com"},
				{"Bob", "bob@example.com"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	output := buf.String()
	assertContains(t, output, "  records: 2")
	assertContains(t, output, "  records-equal: 2")
	assertContains(t, output, "  min-records: >= 1")
	assertContains(t, output, "  max-records: <= 2")
	assertContains(t, output, `  field-exists: "email"`)
}

func assertContains(t *testing.T, value, expected string) {
	t.Helper()
	if !strings.Contains(value, expected) {
		t.Fatalf("expected output to contain %q, got %q", expected, value)
	}
}
