package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsUnknownTopLevelField(t *testing.T) {
	content := `id: test
input:
  format: json
output:
  format: json
steps:
  - log: {}
bogus: true
`
	path := writeTempPipeline(t, content)
	_, err := Load(path, nil)
	if err == nil {
		t.Fatal("expected error for unknown top-level field")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Fatalf("expected field name in error, got: %v", err)
	}
}

func TestLoadRejectsMissingID(t *testing.T) {
	content := `input:
  format: json
output:
  format: json
steps:
  - log: {}
`
	path := writeTempPipeline(t, content)
	_, err := Load(path, nil)
	if err == nil {
		t.Fatal("expected error for missing id")
	}
	if err.Error() != "pipeline id must be specified" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsEmptySteps(t *testing.T) {
	content := `id: test
input:
  format: json
output:
  format: json
steps: []
`
	path := writeTempPipeline(t, content)
	_, err := Load(path, nil)
	if err == nil {
		t.Fatal("expected error for empty steps")
	}
	if err.Error() != "pipeline must have at least one step" {
		t.Fatalf("unexpected error: %v", err)
	}
}

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
	_, err := Load(path, nil)
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
	_, err := Load(path, nil)
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
	_, err := Load(path, nil)
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
steps:
  - log: {}
`
	path := writeTempPipeline(t, content)
	p, err := Load(path, nil)
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
	_, err := Load(path, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "line 7") {
		t.Fatalf("expected line number in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "filter requires a condition") {
		t.Fatalf("expected validation message in error, got: %v", err)
	}
}

func TestLoadSubstitutesVars(t *testing.T) {
	content := `id: test
input:
  format: ${INPUT_FORMAT}
output:
  format: json
steps:
  - log: {}
`
	path := writeTempPipeline(t, content)
	p, err := Load(path, map[string]string{"INPUT_FORMAT": "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Input.Format != "json" {
		t.Fatalf("expected input format json, got %q", p.Input.Format)
	}
}

func TestLoadErrorsOnUnresolvedVar(t *testing.T) {
	content := `id: test
input:
  format: ${INPUT_FORMAT}
output:
  format: json
steps:
  - log: {}
`
	path := writeTempPipeline(t, content)
	_, err := Load(path, nil)
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
	if !strings.Contains(err.Error(), "INPUT_FORMAT") {
		t.Fatalf("expected variable name in error, got: %v", err)
	}
}

func TestLoadErrorsOnPartiallyUnresolvedVars(t *testing.T) {
	content := `id: test
input:
  format: ${INPUT_FORMAT}
output:
  format: ${OUTPUT_FORMAT}
steps:
  - log: {}
`
	path := writeTempPipeline(t, content)
	_, err := Load(path, map[string]string{"INPUT_FORMAT": "json"})
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
	if !strings.Contains(err.Error(), "OUTPUT_FORMAT") {
		t.Fatalf("expected unresolved variable name in error, got: %v", err)
	}
}

func TestLoadNoDefaultsBlockPreservesExistingBehavior(t *testing.T) {
	content := `id: test
input:
  format: json
output:
  format: json
steps:
  - http-request:
      url: https://example.com
      method: GET
  - filter:
      field: name
      equals: Bob
`
	path := writeTempPipeline(t, content)
	p, err := Load(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	httpStep := p.Steps[0].Step.(*HTTPRequestStep)
	if httpStep.Proxy != "" {
		t.Fatalf("unexpected proxy: got %q want %q", httpStep.Proxy, "")
	}
	if httpStep.Timeout != 0 {
		t.Fatalf("unexpected timeout: got %d want 0", httpStep.Timeout)
	}

	filterStep := p.Steps[1].Step.(*FilterStep)
	if filterStep.CaseSensitive == nil || !*filterStep.CaseSensitive {
		t.Fatalf("expected case-sensitive to default to true, got %v", filterStep.CaseSensitive)
	}
}

func TestLoadAppliesHTTPDefaults(t *testing.T) {
	content := `id: test
input:
  format: json
defaults:
  http:
    proxy: http://proxy.internal:8080
    timeout: 30
    headers:
      X-Api-Key: shared-key
      X-Common: default-value
steps:
  - http-request:
      url: https://example.com/a
      method: GET
  - http-request:
      url: https://example.com/b
      method: GET
      timeout: 10
      headers:
        X-Common: step-value
output:
  format: json
`
	path := writeTempPipeline(t, content)
	p, err := Load(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	inherited := p.Steps[0].Step.(*HTTPRequestStep)
	if inherited.Proxy != "http://proxy.internal:8080" {
		t.Fatalf("unexpected proxy: got %q", inherited.Proxy)
	}
	if inherited.Timeout != 30 {
		t.Fatalf("unexpected timeout: got %d want 30", inherited.Timeout)
	}
	if inherited.Headers["X-Api-Key"] != "shared-key" {
		t.Fatalf("expected inherited header X-Api-Key, got %+v", inherited.Headers)
	}

	overridden := p.Steps[1].Step.(*HTTPRequestStep)
	if overridden.Proxy != "http://proxy.internal:8080" {
		t.Fatalf("expected proxy to still be inherited: got %q", overridden.Proxy)
	}
	if overridden.Timeout != 10 {
		t.Fatalf("expected step timeout to override default: got %d want 10", overridden.Timeout)
	}
	if overridden.Headers["X-Api-Key"] != "shared-key" {
		t.Fatalf("expected non-conflicting default header to merge in, got %+v", overridden.Headers)
	}
	if overridden.Headers["X-Common"] != "step-value" {
		t.Fatalf("expected step header to win on key conflict, got %+v", overridden.Headers)
	}
}

func TestLoadAppliesTextDefaults(t *testing.T) {
	content := `id: test
input:
  format: json
defaults:
  text:
    case-sensitive: false
steps:
  - filter:
      field: name
      equals: bob
  - assert:
      field-exists: name
      case-sensitive: true
  - dedupe:
      fields: [name]
output:
  format: json
`
	path := writeTempPipeline(t, content)
	p, err := Load(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	filterStep := p.Steps[0].Step.(*FilterStep)
	if filterStep.CaseSensitive == nil || *filterStep.CaseSensitive {
		t.Fatalf("expected filter to inherit case-sensitive: false, got %v", filterStep.CaseSensitive)
	}

	assertStep := p.Steps[1].Step.(*AssertStep)
	if assertStep.CaseSensitive == nil || !*assertStep.CaseSensitive {
		t.Fatalf("expected assert step-level case-sensitive: true to override default, got %v", assertStep.CaseSensitive)
	}

	dedupeStep := p.Steps[2].Step.(*DedupeStep)
	if dedupeStep.CaseSensitive == nil || *dedupeStep.CaseSensitive {
		t.Fatalf("expected dedupe to inherit case-sensitive: false, got %v", dedupeStep.CaseSensitive)
	}
}

func TestLoadRejectsOutOfBoundsHTTPDefaultTimeout(t *testing.T) {
	content := `id: test
input:
  format: json
defaults:
  http:
    timeout: 500
steps:
  - http-request:
      url: https://example.com
      method: GET
output:
  format: json
`
	path := writeTempPipeline(t, content)
	_, err := Load(path, nil)
	if err == nil {
		t.Fatal("expected error for out-of-bounds defaults.http timeout")
	}
	if !strings.Contains(err.Error(), "defaults.http timeout") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsUnknownDefaultsField(t *testing.T) {
	content := `id: test
input:
  format: json
defaults:
  http:
    bogus: true
steps:
  - http-request:
      url: https://example.com
      method: GET
output:
  format: json
`
	path := writeTempPipeline(t, content)
	_, err := Load(path, nil)
	if err == nil {
		t.Fatal("expected error for unknown defaults.http field")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Fatalf("expected field name in error, got: %v", err)
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
