package assert

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
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
	minRecords := 1
	maxRecords := 3
	equal := 2
	step := &Step{
		MinRecords:   &minRecords,
		MaxRecords:   &maxRecords,
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
	minRecords := 1
	maxRecords := 1
	equal := 1
	step := &Step{
		MinRecords:   &minRecords,
		MaxRecords:   &maxRecords,
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
	minRecords := 1
	maxRecords := 2
	equal := 2
	step := &Step{
		MinRecords:   &minRecords,
		MaxRecords:   &maxRecords,
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
	minRecords := 2
	step := &Step{
		MinRecords: &minRecords,
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
	maxRecords := 1
	step := &Step{
		MaxRecords: &maxRecords,
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

func TestExecuteSucceedsWithFieldEqualsForJSON(t *testing.T) {
	step := &Step{
		FieldEquals:   &FieldCheck{Field: "status", Value: "active"},
		CaseSensitive: true,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"status": "active"},
				{"status": "active"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteSucceedsWithFieldEqualsForCSV(t *testing.T) {
	step := &Step{
		FieldEquals:   &FieldCheck{Field: "status", Value: "active"},
		CaseSensitive: true,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"status"},
				{"active"},
				{"active"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteFailsWhenFieldEqualsMismatches(t *testing.T) {
	step := &Step{
		FieldEquals:   &FieldCheck{Field: "status", Value: "active"},
		CaseSensitive: true,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{
				{"status": "active"},
				{"status": "inactive"},
			},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error when field-equals mismatches")
	}
	if !strings.Contains(err.Error(), `field "status" in record 2 is "inactive", want "active"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFieldEqualsIsCaseInsensitiveWhenConfigured(t *testing.T) {
	step := &Step{
		FieldEquals:   &FieldCheck{Field: "status", Value: "Active"},
		CaseSensitive: false,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{{"status": "active"}},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteFailsWhenFieldEqualsFieldMissing(t *testing.T) {
	step := &Step{
		FieldEquals:   &FieldCheck{Field: "status", Value: "active"},
		CaseSensitive: true,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{{"name": "Alice"}},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error when field-equals field is missing")
	}
	if !strings.Contains(err.Error(), `field "status" is missing from record 1`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteSucceedsWithFieldContainsForCSV(t *testing.T) {
	step := &Step{
		FieldContains: &FieldCheck{Field: "email", Value: "@"},
		CaseSensitive: true,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"email"},
				{"alice@example.com"},
				{"bob@example.com"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteFailsWhenFieldDoesNotContain(t *testing.T) {
	step := &Step{
		FieldContains: &FieldCheck{Field: "email", Value: "@"},
		CaseSensitive: true,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{{"email": "alice-example.com"}},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error when field-contains fails")
	}
	if !strings.Contains(err.Error(), `field "email" in record 1 value "alice-example.com" does not contain "@"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFieldContainsIsCaseInsensitiveWhenConfigured(t *testing.T) {
	step := &Step{
		FieldContains: &FieldCheck{Field: "name", Value: "ALICE"},
		CaseSensitive: false,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{{"name": "alice smith"}},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteSucceedsWithFieldMatchesForJSON(t *testing.T) {
	step := &Step{
		FieldMatches: &FieldCheck{Field: "email", Value: "^[^@]+@[^@]+$", Regex: regexp.MustCompile("^[^@]+@[^@]+$")},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{{"email": "alice@example.com"}},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteFailsWhenFieldDoesNotMatch(t *testing.T) {
	step := &Step{
		FieldMatches: &FieldCheck{Field: "email", Value: "^[^@]+@[^@]+$", Regex: regexp.MustCompile("^[^@]+@[^@]+$")},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{{"email": "not-an-email"}},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error when field-matches fails")
	}
	if !strings.Contains(err.Error(), `field "email" in record 1 value "not-an-email" does not match pattern "^[^@]+@[^@]+$"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteFailsWhenFieldMatchesFieldMissingInCSV(t *testing.T) {
	step := &Step{
		FieldMatches: &FieldCheck{Field: "email", Value: "^[^@]+@[^@]+$", Regex: regexp.MustCompile("^[^@]+@[^@]+$")},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"name"},
				{"Alice"},
			},
		},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error when field-matches field is missing")
	}
	if !strings.Contains(err.Error(), `field "email" is missing from record 1`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteLogsAssertions(t *testing.T) {
	minRecords := 1
	maxRecords := 2
	equal := 2
	step := &Step{
		MinRecords:    &minRecords,
		MaxRecords:    &maxRecords,
		RecordsEqual:  &equal,
		FieldExists:   "email",
		FieldEquals:   &FieldCheck{Field: "email", Value: "alice@example.com"},
		FieldContains: &FieldCheck{Field: "email", Value: "@"},
		FieldMatches:  &FieldCheck{Field: "email", Value: "^[^@]+@[^@]+$", Regex: regexp.MustCompile("^[^@]+@[^@]+$")},
		CaseSensitive: true,
	}

	var buf bytes.Buffer
	ctx := &engine.ExecutionContext{
		Logger: engine.NewLoggerWithWriter(&buf, true),
		Payload: &payload.CSV{
			Rows: [][]string{
				{"name", "email"},
				{"Alice", "alice@example.com"},
				{"Alicia", "alice@example.com"},
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
	assertContains(t, output, `  field-equals: email == "alice@example.com"`)
	assertContains(t, output, `  field-contains: email contains "@"`)
	assertContains(t, output, `  field-matches: email matches "^[^@]+@[^@]+$"`)
}

func assertContains(t *testing.T, value, expected string) {
	t.Helper()
	if !strings.Contains(value, expected) {
		t.Fatalf("expected output to contain %q, got %q", expected, value)
	}
}
