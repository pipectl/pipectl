package validatejson

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "validate-json" {
		t.Fatalf("expected step name %q, got %q", "validate-json", step.Name())
	}
}

func TestSupports(t *testing.T) {
	step := &Step{}

	if !step.Supports(&payload.JSON{}) {
		t.Fatal("expected step to support JSON payload")
	}

	if step.Supports(&payload.CSV{}) {
		t.Fatal("did not expect step to support CSV payload")
	}
}

func TestExecuteValidJSONAgainstSchemaFile(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "schema.json")
	schema := `{"type":"object","required":["email"],"properties":{"email":{"type":"string"}}}`
	if err := os.WriteFile(schemaPath, []byte(schema), 0o644); err != nil {
		t.Fatalf("failed to write schema file: %v", err)
	}

	step := &Step{Schema: schemaPath}
	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{Data: map[string]interface{}{"email": "alice@example.com"}},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
}

func TestExecuteReturnsValidationError(t *testing.T) {
	step := &Step{Schema: `{"type":"object","required":["email"],"properties":{"email":{"type":"string"}}}`}
	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{Data: map[string]interface{}{"id": 123}},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected schema validation error")
	}
	if !strings.Contains(err.Error(), "JSON schema validation failed") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "email is required") {
		t.Fatalf("unexpected validation details: %v", err)
	}
}

func TestExecuteReturnsErrorWhenSchemaIsMissing(t *testing.T) {
	step := &Step{Schema: " "}
	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{Data: map[string]interface{}{"email": "alice@example.com"}},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error for missing schema")
	}
	if !strings.Contains(err.Error(), "schema is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}
