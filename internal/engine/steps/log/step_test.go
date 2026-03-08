package _log

import (
	"io"
	"os"
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

	if !step.Supports(&payload.CSV{}) {
		t.Fatal("expected step to support CSV payload")
	}

	if !step.Supports(&payload.Text{}) {
		t.Fatal("expected step to support Text payload")
	}
}

func TestExecuteDefaultsMessageCountAndSample(t *testing.T) {
	step := &Step{
		Count:  true,
		Sample: 10,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "name"},
				{"1", "alice"},
				{"2", "bob"},
			},
		},
	}

	output := captureStdout(t, func() {
		if err := step.Execute(ctx); err != nil {
			t.Fatalf("execute returned error: %v", err)
		}
	})

	assertNotContains(t, output, "message:")
	assertContains(t, output, "- records: 2\n")
	assertContains(t, output, "- sample (2):\n")
	assertContains(t, output, "id,name\n")
	assertContains(t, output, "1,alice\n")
	assertContains(t, output, "2,bob\n")
}

func TestExecutePrintsMessageAndRespectsCountAndSample(t *testing.T) {
	step := &Step{
		Message: "Payload after step 2",
		Count:   false,
		Sample:  1,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "name"},
				{"1", "alice"},
				{"2", "bob"},
			},
		},
	}

	output := captureStdout(t, func() {
		if err := step.Execute(ctx); err != nil {
			t.Fatalf("execute returned error: %v", err)
		}
	})

	assertContains(t, output, "- message: Payload after step 2\n")
	assertNotContains(t, output, "records:")
	assertContains(t, output, "- sample (1):\n")
	assertContains(t, output, "id,name\n")
	assertContains(t, output, "1,alice\n")
	assertNotContains(t, output, "2,bob\n")
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
