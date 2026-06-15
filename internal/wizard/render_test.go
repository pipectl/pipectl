package wizard

import (
	"strings"
	"testing"
)

func TestRender_basic(t *testing.T) {
	r := Result{
		ID:           "test-pipeline",
		InputFormat:  "json",
		Steps:        []string{"select", "rename"},
		OutputFormat: "csv",
	}
	got := Render(r)

	if !strings.Contains(got, "id: test-pipeline") {
		t.Error("missing pipeline id")
	}
	if !strings.Contains(got, "  format: json") {
		t.Error("missing input format")
	}
	if !strings.Contains(got, "  format: csv") {
		t.Error("missing output format")
	}
	if !strings.Contains(got, "- select:") {
		t.Error("missing select step")
	}
	if !strings.Contains(got, "- rename:") {
		t.Error("missing rename step")
	}
}

func TestRender_stepOrder(t *testing.T) {
	r := Result{
		ID:           "order-test",
		InputFormat:  "csv",
		Steps:        []string{"limit", "filter", "select"},
		OutputFormat: "jsonl",
	}
	got := Render(r)

	limitIdx := strings.Index(got, "- limit:")
	filterIdx := strings.Index(got, "- filter:")
	selectIdx := strings.Index(got, "- select:")

	if limitIdx == -1 || filterIdx == -1 || selectIdx == -1 {
		t.Fatal("one or more steps missing from output")
	}
	if !(limitIdx < filterIdx && filterIdx < selectIdx) {
		t.Error("steps rendered out of order")
	}
}

func TestRender_noSteps(t *testing.T) {
	r := Result{
		ID:           "empty",
		InputFormat:  "jsonl",
		Steps:        nil,
		OutputFormat: "json",
	}
	got := Render(r)
	if !strings.Contains(got, "steps:\n") {
		t.Error("steps section missing")
	}
	if strings.Contains(got, "- ") {
		t.Error("unexpected step content with empty steps")
	}
}

func TestRender_allStepsHaveTemplates(t *testing.T) {
	allSteps := []string{
		"select", "rename", "cast", "default", "normalize", "redact", "convert",
		"filter", "dedupe", "sort", "limit",
		"validate-json", "assert",
		"count", "log",
		"http-request", "http-transform",
	}
	for _, step := range allSteps {
		if _, ok := stepTemplates[step]; !ok {
			t.Errorf("no template for step %q", step)
		}
	}
}

func TestRender_outputToStdout(t *testing.T) {
	r := Result{
		ID:           "stdout-test",
		InputFormat:  "json",
		Steps:        []string{"select"},
		OutputFormat: "json",
		OutputFile:   "",
	}
	got := Render(r)
	if got == "" {
		t.Error("expected non-empty output")
	}
}
