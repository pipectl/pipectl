package main

import (
	"strings"
	"testing"
)

func TestParseByteSize(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        int64
		wantErr     bool
		errContains string
	}{
		{name: "megabytes", input: "256MB", want: 256 * megabyte},
		{name: "kilobytes", input: "64KB", want: 64 * kilobyte},
		{name: "gigabytes", input: "1GB", want: 1 * gigabyte},
		{name: "bytes with suffix", input: "100B", want: 100},
		{name: "bare integer is bytes", input: "100", want: 100},
		{name: "zero", input: "0", want: 0},
		{name: "lowercase suffix", input: "256mb", want: 256 * megabyte},
		{name: "mixed case suffix", input: "256Mb", want: 256 * megabyte},
		{name: "whitespace tolerant", input: " 256 MB ", want: 256 * megabyte},
		{name: "fractional value", input: "1.5MB", want: int64(1.5 * megabyte)},
		{
			name:        "empty string",
			input:       "",
			wantErr:     true,
			errContains: "empty value",
		},
		{
			name:        "garbage",
			input:       "abc",
			wantErr:     true,
			errContains: "invalid number",
		},
		{
			name:        "missing number",
			input:       "MB",
			wantErr:     true,
			errContains: "invalid number",
		},
		{
			name:        "negative value",
			input:       "-5MB",
			wantErr:     true,
			errContains: "must not be negative",
		},
		{
			name:        "bad suffix",
			input:       "5XB",
			wantErr:     true,
			errContains: "invalid number",
		},
		{
			name:        "overflow",
			input:       "999999999999999999999GB",
			wantErr:     true,
			errContains: "too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseByteSize(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseByteSize(%q) = %d, nil; want error", tt.input, got)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("parseByteSize(%q) error = %q; want it to contain %q", tt.input, err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseByteSize(%q) returned unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("parseByteSize(%q) = %d; want %d", tt.input, got, tt.want)
			}
		})
	}
}
