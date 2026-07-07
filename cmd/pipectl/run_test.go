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

func writeRunTestPipeline(t *testing.T, dir string) string {
	t.Helper()
	pipelinePath := filepath.Join(dir, "pipeline.yaml")
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
	return pipelinePath
}

func TestRunCommandRejectsOversizedStdinInput(t *testing.T) {
	tempDir := t.TempDir()
	pipelinePath := writeRunTestPipeline(t, tempDir)

	inputPath := filepath.Join(tempDir, "input.json")
	if err := os.WriteFile(inputPath, []byte(`{"aaaaaaa":1}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer inputFile.Close()

	originalStdin := os.Stdin
	originalMaxInputSizeStr := maxInputSizeStr
	os.Stdin = inputFile
	maxInputSizeStr = "10B"
	defer func() {
		os.Stdin = originalStdin
		maxInputSizeStr = originalMaxInputSizeStr
	}()

	err = runCommand.RunE(runCommand, []string{pipelinePath})
	if err == nil {
		t.Fatalf("RunE returned nil error, want oversized stdin error")
	}
	if !strings.Contains(err.Error(), "stdin input exceeds maximum input size") {
		t.Fatalf("RunE error = %q; want it to contain %q", err.Error(), "stdin input exceeds maximum input size")
	}
}

func TestRunCommandRejectsOversizedInputFile(t *testing.T) {
	tempDir := t.TempDir()
	pipelinePath := writeRunTestPipeline(t, tempDir)

	inputFilePath := filepath.Join(tempDir, "input.json")
	if err := os.WriteFile(inputFilePath, []byte(`{"aaaaaaa":1}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	originalInputPath := inputPath
	originalMaxInputSizeStr := maxInputSizeStr
	inputPath = inputFilePath
	maxInputSizeStr = "10B"
	defer func() {
		inputPath = originalInputPath
		maxInputSizeStr = originalMaxInputSizeStr
	}()

	err := runCommand.RunE(runCommand, []string{pipelinePath})
	if err == nil {
		t.Fatalf("RunE returned nil error, want oversized input file error")
	}
	if !strings.Contains(err.Error(), "exceeds maximum input size") {
		t.Fatalf("RunE error = %q; want it to contain %q", err.Error(), "exceeds maximum input size")
	}
	if !strings.Contains(err.Error(), inputFilePath) {
		t.Fatalf("RunE error = %q; want it to contain the file path %q", err.Error(), inputFilePath)
	}
}

func TestRunCommandAcceptsInputAtExactMaxSize(t *testing.T) {
	const exactlyTenBytes = `{"a":1234}`
	if len(exactlyTenBytes) != 10 {
		t.Fatalf("test fixture must be exactly 10 bytes, got %d", len(exactlyTenBytes))
	}

	t.Run("stdin", func(t *testing.T) {
		tempDir := t.TempDir()
		pipelinePath := writeRunTestPipeline(t, tempDir)

		inputPath := filepath.Join(tempDir, "input.json")
		if err := os.WriteFile(inputPath, []byte(exactlyTenBytes), 0o644); err != nil {
			t.Fatalf("WriteFile returned error: %v", err)
		}
		inputFile, err := os.Open(inputPath)
		if err != nil {
			t.Fatalf("Open returned error: %v", err)
		}
		defer inputFile.Close()

		originalStdin := os.Stdin
		originalMaxInputSizeStr := maxInputSizeStr
		os.Stdin = inputFile
		maxInputSizeStr = "10B"
		defer func() {
			os.Stdin = originalStdin
			maxInputSizeStr = originalMaxInputSizeStr
		}()

		if err := runCommand.RunE(runCommand, []string{pipelinePath}); err != nil {
			t.Fatalf("RunE returned error: %v", err)
		}
	})

	t.Run("file", func(t *testing.T) {
		tempDir := t.TempDir()
		pipelinePath := writeRunTestPipeline(t, tempDir)

		inputFilePath := filepath.Join(tempDir, "input.json")
		if err := os.WriteFile(inputFilePath, []byte(exactlyTenBytes), 0o644); err != nil {
			t.Fatalf("WriteFile returned error: %v", err)
		}

		originalInputPath := inputPath
		originalMaxInputSizeStr := maxInputSizeStr
		inputPath = inputFilePath
		maxInputSizeStr = "10B"
		defer func() {
			inputPath = originalInputPath
			maxInputSizeStr = originalMaxInputSizeStr
		}()

		if err := runCommand.RunE(runCommand, []string{pipelinePath}); err != nil {
			t.Fatalf("RunE returned error: %v", err)
		}
	})
}

func TestRunCommandRejectsInvalidMaxInputSizeFlag(t *testing.T) {
	tempDir := t.TempDir()
	pipelinePath := writeRunTestPipeline(t, tempDir)

	originalMaxInputSizeStr := maxInputSizeStr
	maxInputSizeStr = "not-a-size"
	defer func() {
		maxInputSizeStr = originalMaxInputSizeStr
	}()

	err := runCommand.RunE(runCommand, []string{pipelinePath})
	if err == nil {
		t.Fatalf("RunE returned nil error, want invalid --max-input-size error")
	}
	if !strings.Contains(err.Error(), "--max-input-size") {
		t.Fatalf("RunE error = %q; want it to contain %q", err.Error(), "--max-input-size")
	}
}
