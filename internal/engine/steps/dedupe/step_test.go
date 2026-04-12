package dedupe

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "dedupe" {
		t.Fatalf("expected step name %q, got %q", "dedupe", step.Name())
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

func TestExecuteDedupeJSON(t *testing.T) {
	step := &Step{Fields: []string{"department"}}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{"name": "Alice", "department": "Engineering"},
				{"name": "Bob", "department": "Marketing"},
				{"name": "Carol", "department": "Engineering"},
				{"name": "Dave", "department": "HR"},
				{"name": "Eve", "department": "Marketing"},
			},
			Shape: payload.JSONArrayShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 3 {
		t.Fatalf("expected 3 records after dedupe, got %d", got)
	}
	assertContains(t, buf.String(), "removed 2 duplicate records")
}

func TestExecuteDedupeJSONL(t *testing.T) {
	step := &Step{Fields: []string{"country"}}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"name": "Alice", "country": "au"},
				{"name": "Bob", "country": "us"},
				{"name": "Carol", "country": "au"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 2 {
		t.Fatalf("expected 2 records after dedupe, got %d", got)
	}
	assertContains(t, buf.String(), "removed 1 duplicate records")
}

func TestExecuteDedupeCompositeKey(t *testing.T) {
	step := &Step{Fields: []string{"first", "last"}}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"first": "John", "last": "Smith"},
				{"first": "John", "last": "Jones"},
				{"first": "John", "last": "Smith"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 2 {
		t.Fatalf("expected 2 records after dedupe on composite key, got %d", got)
	}
}

func TestExecuteDedupeCSV(t *testing.T) {
	step := &Step{Fields: []string{"plan"}}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.CSV{
			Rows: [][]string{
				{"name", "plan"},
				{"Alice", "starter"},
				{"Bob", "pro"},
				{"Carol", "starter"},
				{"Dave", "enterprise"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 3 {
		t.Fatalf("expected 3 records after dedupe, got %d", got)
	}
	csvPayload := ctx.Payload.(*payload.CSV)
	if csvPayload.Rows[0][0] != "name" {
		t.Fatalf("header row not preserved: %v", csvPayload.Rows[0])
	}
	assertContains(t, buf.String(), "removed 1 duplicate rows")
}

func TestExecuteDedupeNoDuplicates(t *testing.T) {
	step := &Step{Fields: []string{"id"}}
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
		t.Fatalf("expected 3 records unchanged, got %d", got)
	}
	if strings.Contains(buf.String(), "removed") {
		t.Fatal("expected no removed message when there are no duplicates")
	}
}

func TestExecuteDedupeCaseInsensitive(t *testing.T) {
	step := &Step{Fields: []string{"department"}, CaseSensitive: false}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"name": "Alice", "department": "Engineering"},
				{"name": "Bob", "department": "engineering"},
				{"name": "Carol", "department": "ENGINEERING"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 1 {
		t.Fatalf("expected 1 record after case-insensitive dedupe, got %d", got)
	}
	assertContains(t, buf.String(), "removed 2 duplicate records")
}

func TestExecuteDedupeCaseSensitiveByDefault(t *testing.T) {
	step := &Step{Fields: []string{"department"}, CaseSensitive: true}
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&bytes.Buffer{}, true),
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"name": "Alice", "department": "Engineering"},
				{"name": "Bob", "department": "engineering"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if got := ctx.Payload.RecordCount(); got != 2 {
		t.Fatalf("expected 2 records (case-sensitive by default), got %d", got)
	}
}

func TestExecuteDedupeEmptyCSV(t *testing.T) {
	step := &Step{Fields: []string{"plan"}}
	ctx := &engine.ExecutionContext{
		Logger:  engine.NewLoggerWithWriter(&bytes.Buffer{}, true),
		Payload: &payload.CSV{},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error on empty CSV: %v", err)
	}
}

func assertContains(t *testing.T, value, expected string) {
	t.Helper()
	if !strings.Contains(value, expected) {
		t.Fatalf("expected output to contain %q, got %q", expected, value)
	}
}
