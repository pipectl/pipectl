package count

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "count" {
		t.Fatalf("expected step name %q, got %q", "count", step.Name())
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

func TestExecutePrintsRawRecordCountWithoutCommas(t *testing.T) {
	rows := make([][]string, 0, 1224)
	rows = append(rows, []string{"id"})
	for i := 1; i <= 1223; i++ {
		rows = append(rows, []string{fmt.Sprintf("%d", i)})
	}

	step := &Step{}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger:  engine.NewLoggerWithWriter(&buf, false),
		Payload: &payload.CSV{Rows: rows},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	output := buf.String()
	assertContains(t, output, "  records: 1223")
	assertNotContains(t, output, "1,223")
	assertNotContains(t, output, "message:")
}

func TestExecutePrintsMessageLikeLogStep(t *testing.T) {
	step := &Step{Message: "Message goes here"}
	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, false),
		Payload: &payload.JSON{
			Items: []map[string]interface{}{
				{
					"name1": "value1",
					"name2": "value2",
					"name3": "value3",
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	output := buf.String()
	assertContains(t, output, "  message: Message goes here")
	assertContains(t, output, "  records: 1")
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
