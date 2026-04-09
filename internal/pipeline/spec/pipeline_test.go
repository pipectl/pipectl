package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsMultiCharDelimiter(t *testing.T) {
	content := `id: test
input:
  format: csv
  delimiter: "||"
output:
  format: json
steps: []
`
	path := writeTempPipeline(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for multi-char delimiter")
	}
	if err.Error() != "input delimiter must be a single character" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsInvalidInputFormat(t *testing.T) {
	content := `id: test
input:
  format: xml
output:
  format: json
steps: []
`
	path := writeTempPipeline(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid input format")
	}
	if err.Error() != "input format must be one of: json, jsonl, csv" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsInvalidOutputFormat(t *testing.T) {
	content := `id: test
input:
  format: json
output:
  format: xml
steps: []
`
	path := writeTempPipeline(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid output format")
	}
	if err.Error() != "output format must be one of: json, jsonl, csv" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadAcceptsSingleCharDelimiter(t *testing.T) {
	content := `id: test
input:
  format: csv
  delimiter: "|"
output:
  format: json
steps: []
`
	path := writeTempPipeline(t, content)
	p, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Input.Delimiter != "|" {
		t.Fatalf("unexpected delimiter: got %q want %q", p.Input.Delimiter, "|")
	}
}

func TestLoadIncludesLineNumberInStepError(t *testing.T) {
	content := `id: test
input:
  format: json
output:
  format: json
steps:
  - filter:
`
	path := writeTempPipeline(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "line 7") {
		t.Fatalf("expected line number in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "filter field is required") {
		t.Fatalf("expected validation message in error, got: %v", err)
	}
}

func writeTempPipeline(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "pipeline.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp pipeline: %v", err)
	}
	return path
}
