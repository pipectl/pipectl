package count

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
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
	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: rows,
		},
	}

	output := captureStdout(t, func() {
		if err := step.Execute(ctx); err != nil {
			t.Fatalf("execute returned error: %v", err)
		}
	})

	assertContains(t, output, "- records: 1223\n")
	assertNotContains(t, output, "- records: 1,223")
	assertNotContains(t, output, "message:")
}

func TestExecutePrintsMessageLikeLogStep(t *testing.T) {
	step := &Step{Message: "Message goes here"}
	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Data: map[string]interface{}{
				"name1": "value1",
				"name2": "value2",
				"name3": "value3",
			},
		},
	}

	output := captureStdout(t, func() {
		if err := step.Execute(ctx); err != nil {
			t.Fatalf("execute returned error: %v", err)
		}
	})

	assertContains(t, output, "- message: Message goes here\n")
	assertContains(t, output, "- records: 1\n")
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe returned error: %v", err)
	}
	defer reader.Close()

	os.Stdout = writer
	defer func() {
		os.Stdout = original
	}()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("closing writer returned error: %v", err)
	}

	out, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("reading stdout returned error: %v", err)
	}

	return string(out)
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
