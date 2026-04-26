package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeValidatePipeline(t *testing.T, yaml string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "pipeline.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name        string
		yaml        string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid pipeline",
			yaml: `id: my-pipeline
input:
  format: json
steps:
  - log: {}
output:
  format: json
`,
		},
		{
			name: "missing id",
			yaml: `input:
  format: json
steps:
  - log: {}
output:
  format: json
`,
			wantErr:     true,
			errContains: "id",
		},
		{
			name: "invalid input format",
			yaml: `id: p
input:
  format: badformat
steps:
  - log: {}
output:
  format: json
`,
			wantErr:     true,
			errContains: "input format",
		},
		{
			name: "unknown step type",
			yaml: `id: p
input:
  format: json
steps:
  - not-a-real-step: {}
output:
  format: json
`,
			wantErr: true,
		},
		{
			name: "invalid step config",
			yaml: `id: p
input:
  format: json
steps:
  - limit:
      count: 0
output:
  format: json
`,
			wantErr:     true,
			errContains: "count",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := writeValidatePipeline(t, tc.yaml)

			err := validateCommand.RunE(validateCommand, []string{path})

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Fatalf("expected error containing %q, got %q", tc.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
