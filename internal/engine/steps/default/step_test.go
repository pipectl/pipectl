package _default

import (
	"reflect"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

func TestName(t *testing.T) {
	step := &Step{}
	if step.Name() != "default" {
		t.Fatalf("expected step name %q, got %q", "default", step.Name())
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

func TestExecuteAppliesDefaultsToMissingJSONFields(t *testing.T) {
	step := &Step{
		Fields: map[string]interface{}{
			"country":  "AU",
			"password": "Passw0rd",
			"age":      18,
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{
			Records: []map[string]interface{}{
				{
					"name":     "Alice",
					"country":  "NZ",
					"password": "super-secret",
				},
			},
			Shape: payload.JSONObjectShape,
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out, ok := ctx.Payload.(*payload.JSON)
	if !ok {
		t.Fatalf("expected payload.JSON, got %T", ctx.Payload)
	}

	expected := map[string]interface{}{
		"name":     "Alice",
		"country":  "NZ",
		"password": "super-secret",
		"age":      18,
	}
	if !reflect.DeepEqual(out.Records[0], expected) {
		t.Fatalf("unexpected defaulted JSON data:\nexpected: %#v\ngot: %#v", expected, out.Records[0])
	}
}

func TestExecuteInitializesNilJSONMap(t *testing.T) {
	step := &Step{
		Fields: map[string]interface{}{
			"country": "AU",
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out := ctx.Payload.(*payload.JSON)
	if out.Records[0]["country"] != "AU" {
		t.Fatalf("expected default field to be set, got %#v", out.Records)
	}
}

func TestExecuteAppliesDefaultsToCSVRows(t *testing.T) {
	step := &Step{
		Fields: map[string]interface{}{
			"country":  "AU",
			"password": "Passw0rd",
			"age":      18,
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"name", "country", "age"},
				{"Alice", "", ""},
				{"Bob", "NZ", "24"},
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
		{"name", "country", "age", "password"},
		{"Alice", "AU", "18", "Passw0rd"},
		{"Bob", "NZ", "24", "Passw0rd"},
	}
	if !reflect.DeepEqual(out.Rows, expected) {
		t.Fatalf("unexpected defaulted CSV rows:\nexpected: %#v\ngot: %#v", expected, out.Rows)
	}
}

func TestExecuteAddsMissingColumnsForCSVWithNoDataRows(t *testing.T) {
	step := &Step{
		Fields: map[string]interface{}{
			"country": "AU",
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"name"},
			},
		},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out := ctx.Payload.(*payload.CSV)
	expected := [][]string{
		{"name", "country"},
	}
	if !reflect.DeepEqual(out.Rows, expected) {
		t.Fatalf("unexpected CSV rows:\nexpected: %#v\ngot: %#v", expected, out.Rows)
	}
}
