package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCommandWritesPipelineOutputToFile(t *testing.T) {
	tempDir := t.TempDir()

	pipelinePath := filepath.Join(tempDir, "pipeline.yaml")
	pipelineYAML := `id: test-pipeline
input:
  format: json
steps:
  - log: {}
output:
  format: json
`
	if err := os.WriteFile(pipelinePath, []byte(pipelineYAML), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	inputPath := filepath.Join(tempDir, "input.json")
	if err := os.WriteFile(inputPath, []byte(`{"name":"alice"}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	outputFilePath := filepath.Join(tempDir, "output.json")

	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer inputFile.Close()

	originalStdin := os.Stdin
	originalOutputPath := outputPath
	os.Stdin = inputFile
	outputPath = outputFilePath
	defer func() {
		os.Stdin = originalStdin
		outputPath = originalOutputPath
	}()

	if err := runCommand.RunE(runCommand, []string{pipelinePath}); err != nil {
		t.Fatalf("RunE returned error: %v", err)
	}

	output, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	outputText := string(output)
	if !strings.Contains(outputText, `"name": "alice"`) {
		t.Fatalf("expected output file to contain serialized payload, got %q", outputText)
	}
}
