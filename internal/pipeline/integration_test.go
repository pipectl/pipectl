package pipeline_test

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/pipectl/pipectl/internal/pipeline"
)

var update = flag.Bool("update", false, "update golden files")

type testCase struct {
	name, pipeline, input, golden string
	vars                          map[string]string
}

func TestStepPipelines(t *testing.T) {
	run(t, []testCase{
		{name: "filter/operators", pipeline: "step/filter-operators.yaml", input: "people.jsonl", golden: "step/filter-operators.json"},
		{name: "filter/conditions", pipeline: "step/filter-conditions.yaml", input: "people.jsonl", golden: "step/filter-conditions.json"},
		{name: "cast", pipeline: "step/cast.yaml", input: "products.json", golden: "step/cast.json"},
		{name: "cast/csv", pipeline: "step/cast-csv.yaml", input: "products.csv", golden: "step/cast-csv.csv"},
		{name: "normalize", pipeline: "step/normalize.yaml", input: "people.jsonl", golden: "step/normalize.jsonl"},
		{name: "sort", pipeline: "step/sort.yaml", input: "people.jsonl", golden: "step/sort.jsonl"},
		{name: "redact", pipeline: "step/redact.yaml", input: "people.jsonl", golden: "step/redact.jsonl"},
		{name: "redact/partial", pipeline: "step/redact-partial.yaml", input: "people.jsonl", golden: "step/redact-partial.jsonl"},
		{name: "select", pipeline: "step/select.yaml", input: "people.jsonl", golden: "step/select.jsonl"},
		{name: "rename", pipeline: "step/rename.yaml", input: "people.jsonl", golden: "step/rename.jsonl"},
		{name: "default", pipeline: "step/default.yaml", input: "people.jsonl", golden: "step/default.jsonl"},
		{name: "limit", pipeline: "step/limit.yaml", input: "people.jsonl", golden: "step/limit.jsonl"},
		{name: "convert", pipeline: "step/convert.yaml", input: "products.json", golden: "step/convert.json"},
		{name: "validate-json", pipeline: "step/validate-json.yaml", input: "people.jsonl", golden: "step/validate-json.jsonl"},
		{name: "log", pipeline: "step/log.yaml", input: "people.jsonl", golden: "step/log.jsonl"},
		{name: "dedupe", pipeline: "step/dedupe.yaml", input: "people.jsonl", golden: "step/dedupe.jsonl"},
		{name: "dedupe/csv", pipeline: "step/dedupe-csv.yaml", input: "customers.csv", golden: "step/dedupe-csv.csv"},
	})
}

func TestVarPipelines(t *testing.T) {
	run(t, []testCase{
		{name: "vars/basic", pipeline: "vars/basic.yaml", input: "people.jsonl", golden: "vars/basic.jsonl", vars: map[string]string{"LIMIT": "3"}},
	})
}

func TestWorkflowPipelines(t *testing.T) {
	run(t, []testCase{
		{name: "csv-enrichment", pipeline: "workflow/csv-enrichment.yaml", input: "customers.csv", golden: "workflow/csv-enrichment.jsonl"},
		{name: "jsonl-filtering", pipeline: "workflow/jsonl-filtering.yaml", input: "customers.jsonl", golden: "workflow/jsonl-filtering.json"},
		{name: "json-array-to-csv", pipeline: "workflow/json-array-to-csv.yaml", input: "customers-array.json", golden: "workflow/json-array-to-csv.csv"},
	})
}

// TestAssertFailureReturnsError confirms that a failing assert step propagates as a non-nil
// error from pipeline.Run, which causes the test (and CLI) to report failure.
func TestAssertFailureReturnsError(t *testing.T) {
	input, err := os.ReadFile("testdata/input/people.jsonl")
	if err != nil {
		t.Fatalf("read input: %v", err)
	}
	var buf bytes.Buffer
	err = pipeline.Run("testdata/pipelines/step/failing-assert.yaml", input, &buf, false, false, nil)
	if err == nil {
		t.Fatal("expected pipeline with failing assert to return an error, got nil")
	}
}

func run(t *testing.T, cases []testCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input, err := os.ReadFile(filepath.Join("testdata", "input", tc.input))
			if err != nil {
				t.Fatalf("read input: %v", err)
			}

			var buf bytes.Buffer
			if err := pipeline.Run(filepath.Join("testdata", "pipelines", tc.pipeline), input, &buf, false, false, tc.vars); err != nil {
				t.Fatalf("pipeline failed: %v", err)
			}

			goldenPath := filepath.Join("testdata", "golden", tc.golden)
			if *update {
				if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
					t.Fatalf("create golden dir: %v", err)
				}
				if err := os.WriteFile(goldenPath, buf.Bytes(), 0644); err != nil {
					t.Fatalf("write golden: %v", err)
				}
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read golden (run with -update to generate): %v", err)
			}
			if !bytes.Equal(buf.Bytes(), want) {
				t.Errorf("output mismatch\n--- want ---\n%s\n--- got ---\n%s", want, buf.Bytes())
			}
		})
	}
}
