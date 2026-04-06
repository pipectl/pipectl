package _log

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "log" {
		t.Fatalf("expected step name %q, got %q", "log", step.Name())
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

func TestExecutePrintsJSONLSample(t *testing.T) {
	step := &Step{
		Count:  true,
		Sample: 1,
	}

	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, false),
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"id": 1, "name": "alice"},
				{"id": 2, "name": "bob"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	output := buf.String()
	assertContains(t, output, "  records: 2")
	assertContains(t, output, "  sample (1):")
	assertContains(t, output, `"name":"alice"`)
	assertNotContains(t, output, `"name":"bob"`)
}

func TestExecuteDefaultsMessageCountAndSample(t *testing.T) {
	step := &Step{
		Count:  true,
		Sample: 10,
	}

	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, false),
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "name"},
				{"1", "alice"},
				{"2", "bob"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	output := buf.String()
	assertNotContains(t, output, "message:")
	assertContains(t, output, "  records: 2")
	assertContains(t, output, "  sample (2):")
	assertContains(t, output, "id,name")
	assertContains(t, output, "1,alice")
	assertContains(t, output, "2,bob")
}

func TestExecutePrintsMessageAndRespectsCountAndSample(t *testing.T) {
	step := &Step{
		Message: "Payload after step 2",
		Count:   false,
		Sample:  1,
	}

	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, false),
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "name"},
				{"1", "alice"},
				{"2", "bob"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	output := buf.String()
	assertContains(t, output, "  message: Payload after step 2")
	assertNotContains(t, output, "records:")
	assertContains(t, output, "  sample (1):")
	assertContains(t, output, "id,name")
	assertContains(t, output, "1,alice")
	assertNotContains(t, output, "2,bob")
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
