package _select

import (
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestExecuteSelectsCSVFields(t *testing.T) {
	step := &Step{
		Fields: []string{"id", "email"},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "name", "email", "country"},
				{"1", "Alice", "alice@example.com", "AU"},
				{"2", "Bob", "bob@example.com", "NZ"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out, ok := ctx.Payload.(*payload.CSV)
	if !ok {
		t.Fatalf("expected payload.CSV, got %T", ctx.Payload)
	}

	expected := [][]string{
		{"id", "email"},
		{"1", "alice@example.com"},
		{"2", "bob@example.com"},
	}

	if len(out.Rows) != len(expected) {
		t.Fatalf("unexpected row count: got %d want %d", len(out.Rows), len(expected))
	}

	for i := range expected {
		if len(out.Rows[i]) != len(expected[i]) {
			t.Fatalf("unexpected column count in row %d: got %d want %d", i, len(out.Rows[i]), len(expected[i]))
		}
		for j := range expected[i] {
			if out.Rows[i][j] != expected[i][j] {
				t.Fatalf("unexpected value at row %d col %d: got %q want %q", i, j, out.Rows[i][j], expected[i][j])
			}
		}
	}
}
